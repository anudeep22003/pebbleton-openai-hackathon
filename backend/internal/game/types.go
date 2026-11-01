package game

import (
	"context"
	"encoding/json"
)

const MaxDays = 30

type Protagonist struct {
	ID           string         `json:"id"`
	Name         string         `json:"name"`
	OccupationID string         `json:"occupation_id"`
	Goal         string         `json:"goal"`
	Money        int            `json:"money"`
	Feelings     map[string]int `json:"feelings"`
}

type Faction struct {
	ID              string         `json:"id"`
	Name            string         `json:"name"`
	Population      int            `json:"population"`
	Money           int            `json:"money"`
	BaseIncome      int            `json:"base_income"`
	EffectiveIncome int            `json:"effective_income"`
	Upkeep          int            `json:"upkeep"`
	Loyalty         map[string]int `json:"loyalty"`
}

type Proposal struct {
	ID          string `json:"id"`
	AuthorID    string `json:"author_id"`
	Title       string `json:"title"`
	Reason      string `json:"reason"`
	FromFaction string `json:"from_faction"`
	ToFaction   string `json:"to_faction"`
}

type FactionVote struct {
	FactionID          string  `json:"faction_id"`
	SupportProbability float64 `json:"support_probability"`
	YesVotes           int     `json:"yes_votes"`
	NoVotes            int     `json:"no_votes"`
}

type Reaction struct {
	ProtagonistID      string  `json:"protagonist_id"`
	SupportProbability float64 `json:"support_probability"`
	Supports           bool    `json:"supports"`
}

type Newspaper struct {
	Headline    string `json:"headline"`
	Article     string `json:"article"`
	ImagePrompt string `json:"image_prompt"`
}

type DayRecord struct {
	Day                    int                `json:"day"`
	AuthorID               string             `json:"author_id"`
	Candidates             []Proposal         `json:"candidates"`
	SelectionProbabilities map[string]float64 `json:"selection_probabilities"`
	SelectedProposal       Proposal           `json:"selected_proposal"`
	FactionVotes           []FactionVote      `json:"faction_votes"`
	Reactions              []Reaction         `json:"reactions"`
	TotalYes               int                `json:"total_yes"`
	TotalNo                int                `json:"total_no"`
	Passed                 bool               `json:"passed"`
	Newspaper
	IllustrationStatus string   `json:"illustration_status"`
	IllustrationURL    string   `json:"illustration_url"`
	Snapshot           Snapshot `json:"snapshot"`
}

type Snapshot struct {
	Protagonists []Protagonist `json:"protagonists"`
	Factions     []Faction     `json:"factions"`
	ActivePolicy *Proposal     `json:"active_policy"`
}

type World struct {
	ID   string `json:"id"`
	Day  int    `json:"day"`
	Seed int64  `json:"seed"`
	Snapshot
	History []DayRecord `json:"history"`
}

// Providers supply language and probabilities; only the simulation changes facts.
type Provider interface {
	GenerateProposals(context.Context, World, Protagonist, string) ([]Proposal, error)
	ChooseProposal(context.Context, World, Protagonist, []Proposal) (map[string]float64, error)
	EvaluateReactions(context.Context, World, Proposal) (map[string]float64, error)
	GenerateNewspaper(context.Context, World, DayRecord) (Newspaper, error)
}

// All world fields are JSON values. A round trip also detaches historical maps.
func Clone[T any](value T) T {
	data, err := json.Marshal(value)
	if err != nil {
		panic(err)
	}
	var copy T
	if err := json.Unmarshal(data, &copy); err != nil {
		panic(err)
	}
	return copy
}

func NewWorld(seed int64, id string) World {
	w := World{ID: id, Seed: seed, History: []DayRecord{}}
	w.Protagonists = []Protagonist{
		{"pip", "Pip", "merchants", "Become the wealthiest and most influential person in Pebbleton", 15, map[string]int{"momo": 0, "zig": 0}},
		{"momo", "Momo", "farmers", "Build a fair society where ordinary citizens prosper", 15, map[string]int{"pip": 0, "zig": 0}},
		{"zig", "Zig", "workers", "Maintain stability, security and order", 15, map[string]int{"pip": 0, "momo": 0}},
	}
	w.Factions = []Faction{
		{"farmers", "Farmers", 10, 10, 3, 3, 2, map[string]int{"pip": 40, "momo": 70, "zig": 50}},
		{"merchants", "Merchants", 10, 10, 5, 5, 2, map[string]int{"pip": 70, "momo": 40, "zig": 50}},
		{"workers", "Workers", 10, 10, 2, 2, 2, map[string]int{"pip": 40, "momo": 60, "zig": 60}},
		{"clergy", "Clergy", 10, 10, 2, 2, 2, map[string]int{"pip": 50, "momo": 50, "zig": 60}},
	}
	return w
}
