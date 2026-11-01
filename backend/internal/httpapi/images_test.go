package httpapi

import (
	"context"
	"errors"
	"os"
	"path/filepath"
	"testing"

	"pebbleton/backend/internal/game"
)

type imageStub struct {
	entered chan struct{}
	release chan struct{}
	err     error
}

func (i imageStub) GenerateIllustration(context.Context, string) ([]byte, error) {
	if i.entered != nil {
		close(i.entered)
		<-i.release
	}
	return []byte("image-fixture"), i.err
}

func TestImageLifecycleAndAdvanceMerge(t *testing.T) {
	for _, failed := range []bool{false, true} {
		s, err := New(game.Mock{}, t.TempDir(), 42, "mock")
		if err != nil {
			t.Fatal(err)
		}
		s.world, err = game.Advance(context.Background(), s.world, game.Mock{})
		if err != nil {
			t.Fatal(err)
		}
		s.world.History[0].IllustrationStatus = "pending"
		workingCopy := game.Clone(s.world)
		generator := imageStub{}
		if failed {
			generator.err = errors.New("image service unavailable")
		}
		s.EnableIllustrations(generator)
		s.illustrate(s.world.ID, s.world.History[0])
		status := "ready"
		if failed {
			status = "failed"
		}
		if s.State().World.History[0].IllustrationStatus != status {
			t.Fatal("wrong illustration state")
		}
		mergeIllustrations(&workingCopy, s.State().World)
		if workingCopy.History[0].IllustrationStatus != status {
			t.Fatal("advance lost a completed image")
		}
		if !failed {
			response := request(s.Handler(t.TempDir()), "GET", workingCopy.History[0].IllustrationURL)
			if response.Code != 200 || response.Body.String() != "image-fixture" {
				t.Fatal("image not served")
			}
		}
	}
}

func TestResetDiscardsAnOldImage(t *testing.T) {
	directory := t.TempDir()
	s, err := New(game.Mock{}, directory, 42, "mock")
	if err != nil {
		t.Fatal(err)
	}
	s.world, err = game.Advance(context.Background(), s.world, game.Mock{})
	if err != nil {
		t.Fatal(err)
	}
	generator := imageStub{entered: make(chan struct{}), release: make(chan struct{})}
	s.EnableIllustrations(generator)
	worldID, day := s.world.ID, s.world.History[0]
	done := make(chan struct{})
	go func() { s.illustrate(worldID, day); close(done) }()
	<-generator.entered
	if request(s.Handler(t.TempDir()), "POST", "/api/reset").Code != 200 {
		t.Fatal("reset failed")
	}
	close(generator.release)
	<-done
	files, _ := os.ReadDir(filepath.Join(directory, "images"))
	if s.State().World.Day != 0 || len(files) != 0 {
		t.Fatal("old illustration escaped the reset guard")
	}
}

func TestRestartClearsInterruptedIllustrations(t *testing.T) {
	directory := t.TempDir()
	w, err := game.Advance(context.Background(), game.NewWorld(42, "test"), game.Mock{})
	if err != nil {
		t.Fatal(err)
	}
	w.History[0].IllustrationStatus = "pending"
	if err := saveWorld(directory, w); err != nil {
		t.Fatal(err)
	}
	s, err := New(game.Mock{}, directory, 42, "mock")
	if err != nil {
		t.Fatal(err)
	}
	if s.State().World.History[0].IllustrationStatus != "failed" {
		t.Fatal("interrupted image left pending forever")
	}
}
