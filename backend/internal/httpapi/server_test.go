package httpapi

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"reflect"
	"testing"

	"pebbleton/backend/internal/game"
)

func request(handler http.Handler, method, path string) *httptest.ResponseRecorder {
	response := httptest.NewRecorder()
	handler.ServeHTTP(response, httptest.NewRequest(method, path, nil))
	return response
}

func TestHTTPStoryPersistenceAndReset(t *testing.T) {
	directory := t.TempDir()
	s, err := New(game.Mock{}, directory, 42, "mock")
	if err != nil {
		t.Fatal(err)
	}
	handler := s.Handler(t.TempDir())
	for day := 1; day <= game.MaxDays; day++ {
		response := request(handler, "POST", "/api/advance")
		if response.Code != 200 {
			t.Fatal(response.Body.String())
		}
		var state Response
		if err := json.Unmarshal(response.Body.Bytes(), &state); err != nil {
			t.Fatal(err)
		}
		if state.World.Day != day || state.Running {
			t.Fatal("incorrect API state")
		}
	}
	if request(handler, "POST", "/api/advance").Code != 422 {
		t.Fatal("allowed day 31")
	}
	reloaded, err := New(game.Mock{}, directory, 42, "mock")
	if err != nil || !reflect.DeepEqual(s.State(), reloaded.State()) {
		t.Fatal("reload changed world")
	}
	oldID := s.State().World.ID
	if request(handler, "POST", "/api/reset").Code != 200 {
		t.Fatal("reset failed")
	}
	state := s.State()
	if state.World.Day != 0 || len(state.World.History) != 0 || state.World.ID == oldID || state.World.Seed != 42 || state.World.ActivePolicy != nil {
		t.Fatal("incorrect reset")
	}
	if request(handler, "POST", "/api/advance").Code != 200 {
		t.Fatal("cannot advance after reset")
	}
	if !reflect.DeepEqual(s.State().World.History[0], reloaded.State().World.History[0]) {
		t.Fatal("reset did not reproduce first day")
	}
}

type blockedProvider struct {
	game.Mock
	entered chan struct{}
	release chan struct{}
}

func (p blockedProvider) GenerateProposals(ctx context.Context, world game.World, actor game.Protagonist, correction string) ([]game.Proposal, error) {
	close(p.entered)
	<-p.release
	return p.Mock.GenerateProposals(ctx, world, actor, correction)
}

func TestConcurrentAdvanceAndResetAreRejected(t *testing.T) {
	provider := blockedProvider{entered: make(chan struct{}), release: make(chan struct{})}
	s, err := New(provider, t.TempDir(), 42, "mock")
	if err != nil {
		t.Fatal(err)
	}
	handler := s.Handler(t.TempDir())
	done := make(chan *httptest.ResponseRecorder, 1)
	go func() { done <- request(handler, "POST", "/api/advance") }()
	<-provider.entered
	if !s.State().Running {
		t.Fatal("missing running status")
	}
	for _, path := range []string{"/api/advance", "/api/reset"} {
		if request(handler, "POST", path).Code != 409 {
			t.Fatalf("concurrent %s succeeded", path)
		}
	}
	if request(handler, "GET", "/api/world").Code != 200 {
		t.Fatal("reads blocked by model call")
	}
	close(provider.release)
	if response := <-done; response.Code != 200 {
		t.Fatal(response.Body.String())
	}
	if s.State().World.Day != 1 {
		t.Fatal("advanced twice")
	}
}

func TestPersistenceFailureDoesNotPublish(t *testing.T) {
	directory := t.TempDir()
	s, err := New(game.Mock{}, directory, 42, "mock")
	if err != nil {
		t.Fatal(err)
	}
	before := s.State()
	// A directory in place of the snapshot forces the final atomic rename to fail.
	if err := os.Mkdir(filepath.Join(directory, "world.json"), 0755); err != nil {
		t.Fatal(err)
	}
	if request(s.Handler(t.TempDir()), "POST", "/api/advance").Code != 502 {
		t.Fatal("save failure succeeded")
	}
	if !reflect.DeepEqual(before, s.State()) {
		t.Fatal("save failure published a partial day")
	}
}
