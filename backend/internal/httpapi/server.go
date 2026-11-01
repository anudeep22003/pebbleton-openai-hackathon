package httpapi

import (
	"crypto/rand"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"sync"

	"pebbleton/backend/internal/game"
)

type Server struct {
	mu        sync.Mutex
	world     game.World
	running   bool
	provider  game.Provider
	directory string
	seed      int64
	mode      string
	images    ImageGenerator
}

type Response struct {
	World   game.World `json:"world"`
	Running bool       `json:"running"`
	MaxDays int        `json:"max_days"`
	Mode    string     `json:"mode"`
}

func New(provider game.Provider, directory string, seed int64, mode string) (*Server, error) {
	if err := os.MkdirAll(filepath.Join(directory, "images"), 0755); err != nil {
		return nil, err
	}
	s := &Server{provider: provider, directory: directory, seed: seed, mode: mode}
	s.world = loadWorld(directory, game.NewWorld(seed, newID()))
	interrupted := false
	for i := range s.world.History {
		if s.world.History[i].IllustrationStatus == "pending" {
			s.world.History[i].IllustrationStatus = "failed"
			interrupted = true
		}
	}
	if interrupted {
		if err := saveWorld(directory, s.world); err != nil {
			return nil, err
		}
	}
	game.RefreshIncome(&s.world)
	return s, nil
}

func newID() string { return fmt.Sprintf("pebbleton-%x", rand.Text()) }

func (s *Server) Handler(frontendDirectory string) http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /api/world", func(w http.ResponseWriter, r *http.Request) { writeJSON(w, http.StatusOK, s.State()) })
	mux.HandleFunc("POST /api/advance", s.advance)
	mux.HandleFunc("POST /api/reset", s.reset)
	mux.Handle("GET /images/", http.StripPrefix("/images/", http.FileServer(http.Dir(filepath.Join(s.directory, "images")))))
	mux.Handle("GET /", http.FileServer(http.Dir(frontendDirectory)))
	return mux
}

func (s *Server) State() Response {
	s.mu.Lock()
	defer s.mu.Unlock()
	return Response{game.Clone(s.world), s.running, game.MaxDays, s.mode}
}

func (s *Server) advance(w http.ResponseWriter, r *http.Request) {
	s.mu.Lock()
	if s.running {
		s.mu.Unlock()
		apiError(w, 409, "The town is already deliberating.")
		return
	}
	if s.world.Day >= game.MaxDays {
		s.mu.Unlock()
		apiError(w, 422, game.ErrFinished.Error())
		return
	}
	s.running = true
	current := game.Clone(s.world)
	s.mu.Unlock()

	next, err := game.Advance(r.Context(), current, s.provider)
	s.mu.Lock()
	if err == nil {
		mergeIllustrations(&next, s.world)
		if s.images != nil {
			next.History[next.Day-1].IllustrationStatus = "pending"
		}
		err = saveWorld(s.directory, next)
	}
	if err == nil {
		s.world = next
	}
	s.running = false
	response := Response{game.Clone(s.world), false, game.MaxDays, s.mode}
	s.mu.Unlock()
	if err != nil {
		log.Printf("advance failed: %v", err)
		apiError(w, http.StatusBadGateway, "The day's deliberations could not finish. No changes were saved. Please retry; see the server log for details.")
		return
	}
	if s.images != nil {
		go s.illustrate(next.ID, next.History[next.Day-1])
	}
	writeJSON(w, http.StatusOK, response)
}

func (s *Server) reset(w http.ResponseWriter, r *http.Request) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.running {
		apiError(w, 409, "Wait for the town to finish deliberating before resetting.")
		return
	}
	initial := game.NewWorld(s.seed, newID())
	if err := saveWorld(s.directory, initial); err != nil {
		log.Printf("reset failed: %v", err)
		apiError(w, 500, "Could not save the reset. Your world is unchanged.")
		return
	}
	s.world = initial
	writeJSON(w, http.StatusOK, Response{s.world, false, game.MaxDays, s.mode})
}

func writeJSON(w http.ResponseWriter, status int, value any) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Cache-Control", "no-store")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(value); err != nil && !errors.Is(err, http.ErrHandlerTimeout) {
		log.Printf("write response: %v", err)
	}
}

func apiError(w http.ResponseWriter, status int, message string) {
	writeJSON(w, status, map[string]string{"error": message})
}
