package httpapi

import (
	"context"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"

	"pebbleton/backend/internal/game"
)

type ImageGenerator interface {
	GenerateIllustration(context.Context, string) ([]byte, error)
}

// Configure once at startup, before serving requests.
func (s *Server) EnableIllustrations(generator ImageGenerator) { s.images = generator }

func (s *Server) illustrate(worldID string, day game.DayRecord) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Minute)
	defer cancel()
	// Include authoritative facts even if the optional newspaper prompt was empty.
	prompt := fmt.Sprintf("Day %d. %s proposes %s: one daily income coin from %s to %s. Passed: %t, %d yes and %d no. %s", day.Day, day.AuthorID, day.SelectedProposal.Title, day.SelectedProposal.FromFaction, day.SelectedProposal.ToFaction, day.Passed, day.TotalYes, day.TotalNo, day.ImagePrompt)
	data, imageErr := s.images.GenerateIllustration(ctx, prompt)
	s.mu.Lock()
	defer s.mu.Unlock()
	// A reset can complete while the image model is still working.
	if s.world.ID != worldID || len(s.world.History) < day.Day || s.world.History[day.Day-1].SelectedProposal.ID != day.SelectedProposal.ID {
		return
	}
	next := game.Clone(s.world)
	record := &next.History[day.Day-1]
	record.IllustrationStatus = "failed"
	if imageErr == nil {
		filename := fmt.Sprintf("%s-day-%02d.png", worldID, day.Day)
		imageErr = os.WriteFile(filepath.Join(s.directory, "images", filename), data, 0644)
		if imageErr == nil {
			record.IllustrationStatus, record.IllustrationURL = "ready", "/images/"+filename
		}
	}
	if imageErr != nil {
		log.Printf("day %d illustration failed: %v", day.Day, imageErr)
	}
	if err := saveWorld(s.directory, next); err != nil {
		log.Printf("day %d illustration persistence failed: %v", day.Day, err)
		// Stop polling in this process; startup also clears interrupted pending work.
		s.world.History[day.Day-1].IllustrationStatus = "failed"
		return
	}
	s.world = next
}

// An image may finish during a slow advance. Preserve its presentation fields
// when publishing the working copy, which was cloned before the model calls.
func mergeIllustrations(next *game.World, current game.World) {
	for i, previous := range current.History {
		next.History[i].IllustrationStatus = previous.IllustrationStatus
		next.History[i].IllustrationURL = previous.IllustrationURL
	}
}
