package main

import (
	"context"
	"fmt"
	"log"

	"pebbleton/backend/internal/game"
)

func main() {
	w := game.NewWorld(42, "demo")
	for range game.MaxDays {
		next, err := game.Advance(context.Background(), w, game.Mock{})
		if err != nil {
			log.Fatal(err)
		}
		w = next
		d := w.History[w.Day-1]
		fmt.Printf("Day %02d · %s · %d–%d · passed=%t · %s\n", w.Day, d.AuthorID, d.TotalYes, d.TotalNo, d.Passed, d.Headline)
	}
}
