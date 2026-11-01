package game

import (
	"context"
	"encoding/json"
	"errors"
	"math/rand"
	"reflect"
	"testing"
)

func TestThirtyDaysReproducibleAndImmutable(t *testing.T) {
	a, b := NewWorld(42, "test"), NewWorld(42, "test")
	var first Snapshot
	for day := 1; day <= MaxDays; day++ {
		var err error
		a, err = Advance(context.Background(), a, Mock{})
		if err != nil {
			t.Fatal(err)
		}
		b, err = Advance(context.Background(), Clone(b), Mock{})
		if err != nil {
			t.Fatal(err)
		}
		if !reflect.DeepEqual(a, b) {
			t.Fatalf("non-reproducible day %d", day)
		}
		d := a.History[day-1]
		if a.Day != day || len(a.History) != day || d.TotalYes+d.TotalNo != 40 || d.Headline == "" || d.Article == "" {
			t.Fatalf("bad record on day %d", day)
		}
		if d.AuthorID != a.Protagonists[(day-1)%3].ID {
			t.Fatal("wrong author")
		}
		for _, vote := range d.FactionVotes {
			if vote.YesVotes+vote.NoVotes != 10 {
				t.Fatal("lost voter")
			}
		}
		for _, f := range a.Factions {
			if f.EffectiveIncome < 0 {
				t.Fatal("negative income")
			}
			for _, loyalty := range f.Loyalty {
				if loyalty < 0 || loyalty > 100 {
					t.Fatal("invalid loyalty")
				}
			}
		}
		for _, p := range a.Protagonists {
			for _, feeling := range p.Feelings {
				if feeling < -100 || feeling > 100 {
					t.Fatal("invalid feelings")
				}
			}
		}
		if day == 1 {
			first = Clone(d.Snapshot)
		}
		if !reflect.DeepEqual(first, a.History[0].Snapshot) {
			t.Fatal("history mutated")
		}
	}
	if _, err := Advance(context.Background(), a, Mock{}); !errors.Is(err, ErrFinished) {
		t.Fatal("day 31 allowed")
	}
}

func TestEconomyPolicyTimingReplacementAndDebt(t *testing.T) {
	w := NewWorld(42, "test")
	applyEconomy(&w)
	if w.Factions[0].Money != 11 || w.Protagonists[0].Money != 18 {
		t.Fatal("incorrect base income")
	}
	p := Proposal{AuthorID: "pip", FromFaction: "farmers", ToFaction: "merchants"}
	applyConsequences(&w, DayRecord{AuthorID: "pip", SelectedProposal: p, Passed: true})
	if w.Factions[0].Money != 11 {
		t.Fatal("policy changed today's pay")
	}
	applyEconomy(&w)
	if w.Factions[0].Money != 11 || w.Protagonists[0].Money != 22 {
		t.Fatal("next-day policy missing")
	}
	w.ActivePolicy = &Proposal{FromFaction: "workers", ToFaction: "farmers"}
	w.Factions[2].Money = 0
	applyEconomy(&w)
	if w.Factions[0].EffectiveIncome != 4 || w.Factions[1].EffectiveIncome != 5 || w.Factions[2].Money != -1 {
		t.Fatal("policy stacked or debt blocked")
	}
}

func TestValidationAndSampling(t *testing.T) {
	w := NewWorld(1, "test")
	w.Day = 1
	actor := w.Protagonists[0]
	candidates, _ := (Mock{}).GenerateProposals(context.Background(), w, actor, "")
	for i := range candidates {
		candidates[i].ID = string(rune('a' + i))
	}
	for _, mutate := range []func([]Proposal){
		func(p []Proposal) { p[0].ToFaction = p[0].FromFaction },
		func(p []Proposal) { p[0].FromFaction = "unknown" },
		func(p []Proposal) { p[0].AuthorID = "unknown" },
		func(p []Proposal) { p[0].Title = "" },
		func(p []Proposal) { p[1] = p[0] },
	} {
		bad := Clone(candidates)
		mutate(bad)
		if ValidateCandidates(w, actor, bad) == nil {
			t.Fatal("accepted invalid candidates")
		}
	}
	w.Factions[0].BaseIncome = 0
	if ValidateCandidates(w, actor, candidates) == nil {
		t.Fatal("accepted negative income")
	}
	rng := rand.New(rand.NewSource(1))
	choice, err := sampleChoice(rng, candidates, map[string]float64{"a": 0, "b": 1, "c": 0})
	if err != nil || choice.ID != "b" {
		t.Fatal("wrong weighted choice")
	}
	for _, distribution := range []map[string]float64{{"a": 1}, {"a": 0, "b": 0, "c": 0}, {"a": -1, "b": 1, "c": 1}} {
		if _, err := sampleChoice(rng, candidates, distribution); err == nil {
			t.Fatal("accepted malformed distribution")
		}
	}
	for _, allYes := range []bool{false, true} {
		probabilities := map[string]float64{"momo": 0, "zig": 1, "farmers": 0, "merchants": 0, "workers": 1, "clergy": 1}
		if allYes {
			probabilities["farmers"], probabilities["merchants"] = 1, 1
		}
		d := DayRecord{AuthorID: "pip"}
		if err := tally(w, &d, probabilities, rng); err != nil {
			t.Fatal(err)
		}
		if d.Passed != allYes || (!allYes && d.TotalYes != 20) || (allYes && d.TotalYes != 40) {
			t.Fatal("wrong majority/tie result")
		}
	}
}

func TestDirectionalRelationshipsAndClamps(t *testing.T) {
	w := NewWorld(42, "test")
	w.Factions[0].Loyalty["pip"] = 3
	w.Factions[1].Loyalty["pip"] = 98
	w.Protagonists[1].Feelings["pip"] = -98
	d := DayRecord{AuthorID: "pip", SelectedProposal: Proposal{FromFaction: "farmers", ToFaction: "merchants"}, Reactions: []Reaction{{ProtagonistID: "momo", Supports: false}, {ProtagonistID: "zig", Supports: true}}}
	applyConsequences(&w, d)
	if w.Factions[0].Loyalty["pip"] != 0 || w.Factions[1].Loyalty["pip"] != 100 || w.Protagonists[1].Feelings["pip"] != -100 || w.Protagonists[2].Feelings["pip"] != 2 || w.Protagonists[0].Feelings["momo"] != 0 {
		t.Fatal("incorrect directional/clamped effects")
	}
	if w.ActivePolicy != nil {
		t.Fatal("enacted rejected policy")
	}
}

type failingProvider struct {
	Mock
	stage string
	calls *int
}

func (p failingProvider) GenerateProposals(ctx context.Context, w World, actor Protagonist, correction string) ([]Proposal, error) {
	*p.calls++
	if p.stage == "invalid" {
		return []Proposal{}, nil
	}
	if p.stage == "generate" {
		return nil, errors.New("offline")
	}
	return p.Mock.GenerateProposals(ctx, w, actor, correction)
}
func (p failingProvider) EvaluateReactions(ctx context.Context, w World, proposal Proposal) (map[string]float64, error) {
	if p.stage == "reactions" {
		return nil, errors.New("offline")
	}
	return p.Mock.EvaluateReactions(ctx, w, proposal)
}
func (p failingProvider) GenerateNewspaper(context.Context, World, DayRecord) (Newspaper, error) {
	return Newspaper{}, errors.New("offline")
}

func TestAtomicFailuresAndNewspaperFallback(t *testing.T) {
	w := NewWorld(42, "test")
	before, _ := json.Marshal(w)
	for _, stage := range []string{"generate", "invalid", "reactions", "newspaper"} {
		calls := 0
		next, err := Advance(context.Background(), w, failingProvider{stage: stage, calls: &calls})
		if stage == "newspaper" {
			if err != nil || next.Day != 1 || next.History[0].Article == "" {
				t.Fatal("newspaper fallback failed")
			}
		} else if err == nil {
			t.Fatal("core failure succeeded")
		}
		if stage == "invalid" && calls != 2 {
			t.Fatal("invalid proposals must retry exactly once")
		}
		after, _ := json.Marshal(w)
		if string(before) != string(after) {
			t.Fatal("input world partially mutated")
		}
	}
}
