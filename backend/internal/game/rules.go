package game

import (
	"fmt"
	"math"
	"math/rand"
	"strings"
)

func policyEffect(factionID string, policy *Proposal) int {
	if policy == nil {
		return 0
	}
	if factionID == policy.FromFaction {
		return -1
	}
	if factionID == policy.ToFaction {
		return 1
	}
	return 0
}

func RefreshIncome(w *World) {
	for i := range w.Factions {
		f := &w.Factions[i]
		f.EffectiveIncome = f.BaseIncome + policyEffect(f.ID, w.ActivePolicy)
	}
}

func applyEconomy(w *World) {
	RefreshIncome(w)
	for i := range w.Factions {
		f := &w.Factions[i]
		netIncome := f.EffectiveIncome - f.Upkeep
		f.Money += netIncome
		for j := range w.Protagonists {
			if w.Protagonists[j].OccupationID == f.ID {
				w.Protagonists[j].Money += netIncome
			}
		}
	}
}

func ValidateCandidates(w World, actor Protagonist, candidates []Proposal) error {
	if len(candidates) != 3 {
		return fmt.Errorf("expected exactly three proposals")
	}
	authorExists := false
	for _, p := range w.Protagonists {
		authorExists = authorExists || p.ID == actor.ID
	}
	if !authorExists {
		return fmt.Errorf("unknown author")
	}
	factions := map[string]Faction{}
	for _, f := range w.Factions {
		factions[f.ID] = f
	}
	seen := map[string]bool{}
	for _, p := range candidates {
		from, fromOK := factions[p.FromFaction]
		_, toOK := factions[p.ToFaction]
		pair := p.FromFaction + ":" + p.ToFaction
		if p.AuthorID != actor.ID || !fromOK || !toOK || p.FromFaction == p.ToFaction || from.BaseIncome < 1 || seen[pair] || strings.TrimSpace(p.Title) == "" || strings.TrimSpace(p.Reason) == "" {
			return fmt.Errorf("invalid or duplicate proposal %q: use distinct valid factions, a valid author, a title and a reason", p.Title)
		}
		seen[pair] = true
	}
	return nil
}

func validProbability(p float64) bool { return !math.IsNaN(p) && !math.IsInf(p, 0) && p >= 0 && p <= 1 }

func sampleChoice(rng *rand.Rand, candidates []Proposal, probabilities map[string]float64) (Proposal, error) {
	if len(probabilities) != len(candidates) {
		return Proposal{}, fmt.Errorf("incomplete choice distribution")
	}
	total := 0.0
	for _, p := range candidates {
		probability, ok := probabilities[p.ID]
		if !ok || !validProbability(probability) {
			return Proposal{}, fmt.Errorf("invalid choice probability for %s", p.ID)
		}
		total += probability
	}
	if math.Abs(total-1) > 0.01 {
		return Proposal{}, fmt.Errorf("choice probabilities must sum to one")
	}
	threshold := rng.Float64() * total
	for _, p := range candidates {
		threshold -= probabilities[p.ID]
		if threshold < 0 {
			return p, nil
		}
	}
	return candidates[len(candidates)-1], nil
}

func tally(w World, record *DayRecord, probabilities map[string]float64, rng *rand.Rand) error {
	if len(probabilities) != 6 {
		return fmt.Errorf("expected six reaction probabilities")
	}
	for _, p := range w.Protagonists {
		if p.ID == record.AuthorID {
			continue
		}
		probability, ok := probabilities[p.ID]
		if !ok || !validProbability(probability) {
			return fmt.Errorf("invalid reaction probability for %s", p.ID)
		}
		record.Reactions = append(record.Reactions, Reaction{p.ID, probability, rng.Float64() < probability})
	}
	for _, f := range w.Factions {
		probability, ok := probabilities[f.ID]
		if !ok || !validProbability(probability) {
			return fmt.Errorf("invalid vote probability for %s", f.ID)
		}
		vote := FactionVote{FactionID: f.ID, SupportProbability: probability}
		for citizen := 0; citizen < f.Population; citizen++ {
			if rng.Float64() < probability {
				vote.YesVotes++
			} else {
				vote.NoVotes++
			}
		}
		record.FactionVotes = append(record.FactionVotes, vote)
		record.TotalYes += vote.YesVotes
		record.TotalNo += vote.NoVotes
	}
	record.Passed = record.TotalYes > record.TotalNo
	return nil
}

func applyConsequences(w *World, record DayRecord) {
	proposal := record.SelectedProposal
	loyaltyChange := 5
	if record.Passed {
		w.ActivePolicy = &proposal
		loyaltyChange += 3
	}
	for i := range w.Factions {
		f := &w.Factions[i]
		f.Loyalty[record.AuthorID] = max(0, min(100, f.Loyalty[record.AuthorID]+policyEffect(f.ID, &proposal)*loyaltyChange))
	}
	for _, reaction := range record.Reactions {
		for i := range w.Protagonists {
			p := &w.Protagonists[i]
			if p.ID != reaction.ProtagonistID {
				continue
			}
			change := -2
			if reaction.Supports {
				change = 2
			}
			change += 5 * policyEffect(p.OccupationID, &proposal)
			p.Feelings[record.AuthorID] = max(-100, min(100, p.Feelings[record.AuthorID]+change))
		}
	}
	RefreshIncome(w)
}
