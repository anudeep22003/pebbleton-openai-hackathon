package httpapi

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"pebbleton/backend/internal/game"
)

func saveWorld(directory string, world game.World) error {
	data, err := json.MarshalIndent(world, "", "  ")
	if err != nil {
		return err
	}
	file, err := os.CreateTemp(directory, ".world-*.json")
	if err != nil {
		return err
	}
	defer os.Remove(file.Name())
	_, writeErr := file.Write(data)
	closeErr := file.Close()
	if writeErr != nil {
		return writeErr
	}
	if closeErr != nil {
		return closeErr
	}
	return os.Rename(file.Name(), filepath.Join(directory, "world.json"))
}

func loadWorld(directory string, initial game.World) game.World {
	data, err := os.ReadFile(filepath.Join(directory, "world.json"))
	if os.IsNotExist(err) {
		return initial
	}
	var world game.World
	if err == nil {
		err = json.Unmarshal(data, &world)
	}
	if err == nil {
		err = validateSnapshot(world)
	}
	if err != nil {
		log.Printf("snapshot unavailable or invalid; starting at Day 0: %v", err)
		return initial
	}
	return world
}

func validateSnapshot(world game.World) error {
	if world.ID == "" || world.Day < 0 || world.Day > game.MaxDays || len(world.History) != world.Day || len(world.Protagonists) != 3 || len(world.Factions) != 4 {
		return fmt.Errorf("invalid world shape")
	}
	initial := game.NewWorld(0, "validation")
	for i, p := range world.Protagonists {
		if p.ID != initial.Protagonists[i].ID || p.OccupationID != initial.Protagonists[i].OccupationID || len(p.Feelings) != 2 {
			return fmt.Errorf("invalid protagonist")
		}
		for id := range initial.Protagonists[i].Feelings {
			feeling, ok := p.Feelings[id]
			if !ok || feeling < -100 || feeling > 100 {
				return fmt.Errorf("invalid feelings")
			}
		}
	}
	for i, f := range world.Factions {
		if f.ID != initial.Factions[i].ID || f.Population != 10 || f.BaseIncome != initial.Factions[i].BaseIncome || f.Upkeep != 2 || len(f.Loyalty) != 3 {
			return fmt.Errorf("invalid faction")
		}
		for id := range initial.Factions[i].Loyalty {
			loyalty, ok := f.Loyalty[id]
			if !ok || loyalty < 0 || loyalty > 100 {
				return fmt.Errorf("invalid loyalty")
			}
		}
	}
	for i, d := range world.History {
		if d.Day != i+1 || d.TotalYes+d.TotalNo != 40 || d.Passed != (d.TotalYes > d.TotalNo) || len(d.Snapshot.Factions) != 4 || len(d.Snapshot.Protagonists) != 3 {
			return fmt.Errorf("invalid day record")
		}
		if err := game.ValidateCandidates(world, world.Protagonists[i%3], d.Candidates); err != nil {
			return err
		}
	}
	return nil
}
