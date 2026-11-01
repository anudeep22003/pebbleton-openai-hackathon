package game

import (
	"context"
	"fmt"
	"strings"
)

type Mock struct{}

func (Mock) GenerateProposals(_ context.Context, w World, actor Protagonist, _ string) ([]Proposal, error) {
	pairs := [][2]string{{"farmers", "merchants"}, {"clergy", "workers"}, {"merchants", "farmers"}, {"merchants", "workers"}, {"workers", "clergy"}, {"farmers", "clergy"}, {"clergy", "farmers"}, {"workers", "merchants"}, {"merchants", "clergy"}, {"farmers", "workers"}, {"workers", "farmers"}, {"clergy", "merchants"}}
	titles := []string{"The Very Necessary Prosperity Act", "The Everybody Will Thank Me Plan", "The One-Coin Solution"}
	var proposals []Proposal
	for i, title := range titles {
		pair := pairs[((w.Day-1)*3+i)%len(pairs)]
		proposals = append(proposals, Proposal{AuthorID: actor.ID, Title: title, Reason: fmt.Sprintf("%s argues that helping %s is a small step toward a much better Pebbleton. The %s have been volunteered to pay.", actor.Name, pair[1], pair[0]), FromFaction: pair[0], ToFaction: pair[1]})
	}
	return proposals, nil
}

func (Mock) ChooseProposal(_ context.Context, w World, actor Protagonist, candidates []Proposal) (map[string]float64, error) {
	weights := map[string]float64{}
	total := 0.0
	for _, p := range candidates {
		weight := 1.0
		if p.ToFaction == actor.OccupationID {
			weight += 2
		}
		for _, f := range w.Factions {
			if f.ID == p.ToFaction {
				weight += float64(100-f.Loyalty[actor.ID]) / 100
				if actor.ID != "pip" && f.Money < 15 {
					weight += 2
				}
			}
		}
		weights[p.ID] = weight
		total += weight
	}
	for id := range weights {
		weights[id] /= total
	}
	return weights, nil
}

func (Mock) EvaluateReactions(_ context.Context, w World, proposal Proposal) (map[string]float64, error) {
	probabilities := map[string]float64{}
	for _, p := range w.Protagonists {
		if p.ID != proposal.AuthorID {
			probabilities[p.ID] = max(0.05, min(0.95, 0.5+float64(p.Feelings[proposal.AuthorID])/250+float64(policyEffect(p.OccupationID, &proposal))*0.3))
		}
	}
	for _, f := range w.Factions {
		economicEffect := policyEffect(f.ID, &proposal) - policyEffect(f.ID, w.ActivePolicy)
		probabilities[f.ID] = max(0.05, min(0.95, 0.5+float64(f.Loyalty[proposal.AuthorID]-50)/200+float64(economicEffect)*0.3))
	}
	return probabilities, nil
}

func (Mock) GenerateNewspaper(_ context.Context, w World, record DayRecord) (Newspaper, error) {
	return FallbackNewspaper(w, record), nil
}

func FallbackNewspaper(w World, record DayRecord) Newspaper {
	p := record.SelectedProposal
	author := p.AuthorID
	for _, character := range w.Protagonists {
		if character.ID == p.AuthorID {
			author = character.Name
		}
	}
	verdict, consequence := "REJECTED", "The existing income rules remain in place."
	rejectedHeadlines := []string{"%s's one-coin dream fails to make change", "%s discovers that no is a complete sentence", "%s's grand plan meets forty small opinions"}
	headline := fmt.Sprintf(rejectedHeadlines[(record.Day-1)%len(rejectedHeadlines)], author)
	if record.TotalYes == record.TotalNo {
		headline = fmt.Sprintf("%s splits the town. The coin stays put.", author)
	}
	if record.Passed {
		verdict = "PASSED"
		passedHeadlines := []string{"%s moves one coin. Pebbleton calls it history.", "%s wins the day. The coin has a new address.", "%s's little plan causes a very big fuss"}
		headline = fmt.Sprintf(passedHeadlines[(record.Day-1)%len(passedHeadlines)], author)
		consequence = fmt.Sprintf("Starting tomorrow, %s receive one extra coin of daily income per citizen, while %s give up one. This replaces any previous policy; today's pay is unchanged by the vote.", p.ToFaction, p.FromFaction)
	}
	reactions := []string{}
	for _, reaction := range record.Reactions {
		stance := "opposed"
		if reaction.Supports {
			stance = "supported"
		}
		for _, character := range w.Protagonists {
			if character.ID == reaction.ProtagonistID {
				reactions = append(reactions, character.Name+" "+stance+" the proposal")
			}
		}
	}
	article := fmt.Sprintf("PEBBLETON, DAY %d — %s brought ‘%s’ before the town today, proposing to redirect one coin of daily income from %s to %s. In a society this small, a single coin is apparently enough to require everybody's attention.\n\nThe measure was %s, with %d votes for and %d against. %s. %s\n\nThe %s warmed to %s for trying to help, while the %s took the suggestion personally. Loyalties shifted, feelings were recorded, and all forty citizens remain responsible for tomorrow's opinions. Our economic correspondent confirms that one coin is still one coin, however dramatically it is announced.", record.Day, author, p.Title, p.FromFaction, p.ToFaction, verdict, record.TotalYes, record.TotalNo, strings.Join(reactions, "; "), consequence, p.ToFaction, author, p.FromFaction)
	return Newspaper{Headline: headline, Article: article, ImagePrompt: fmt.Sprintf("%s presents a one-coin proposal in Pebbleton's town square; the vote is %s.", author, verdict)}
}
