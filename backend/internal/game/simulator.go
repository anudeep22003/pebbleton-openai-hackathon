package game

import (
	"context"
	"errors"
	"fmt"
	"log"
	"math/rand"
)

var ErrFinished = errors.New("Pebbleton's 30-day story is complete; reset to begin again")

// Advance works on a detached copy. The caller publishes it only after persistence succeeds.
func Advance(ctx context.Context, current World, provider Provider) (World, error) {
	if current.Day >= MaxDays {
		return World{}, ErrFinished
	}
	w := Clone(current)
	w.Day++
	applyEconomy(&w)
	actor := w.Protagonists[(w.Day-1)%len(w.Protagonists)]
	record := DayRecord{Day: w.Day, AuthorID: actor.ID, IllustrationStatus: "disabled"}
	correction := ""
	for attempt := 0; attempt < 2; attempt++ {
		candidates, err := provider.GenerateProposals(ctx, w, actor, correction)
		if err != nil {
			return World{}, fmt.Errorf("generate proposals: %w", err)
		}
		for i := range candidates {
			candidates[i].ID = fmt.Sprintf("proposal-%d-%d", w.Day, i+1)
		}
		if err := ValidateCandidates(w, actor, candidates); err != nil {
			log.Printf("day %d: proposal validation failed: %v", w.Day, err)
			correction = err.Error()
			if attempt == 1 {
				return World{}, err
			}
			continue
		}
		record.Candidates = candidates
		break
	}
	probabilities, err := provider.ChooseProposal(ctx, w, actor, record.Candidates)
	if err != nil {
		return World{}, fmt.Errorf("choose proposal: %w", err)
	}
	// Each day has its own stream: retries and JSON reloads cannot consume future randomness.
	rng := rand.New(rand.NewSource(w.Seed + int64(w.Day)*104729))
	record.SelectedProposal, err = sampleChoice(rng, record.Candidates, probabilities)
	if err != nil {
		return World{}, err
	}
	record.SelectionProbabilities = probabilities
	reactions, err := provider.EvaluateReactions(ctx, w, record.SelectedProposal)
	if err != nil {
		return World{}, fmt.Errorf("evaluate reactions: %w", err)
	}
	if err := tally(w, &record, reactions, rng); err != nil {
		return World{}, err
	}
	applyConsequences(&w, record)
	record.Snapshot = Clone(w.Snapshot)
	record.Newspaper, err = provider.GenerateNewspaper(ctx, w, record)
	if err != nil || record.Headline == "" || record.Article == "" {
		log.Printf("day %d: newspaper unavailable (%v); using factual fallback", w.Day, err)
		record.Newspaper = FallbackNewspaper(w, record)
	}
	if err := ctx.Err(); err != nil {
		return World{}, err
	}
	w.History = append(w.History, record)
	return w, nil
}
