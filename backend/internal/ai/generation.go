package ai

import (
	"context"
	"pebbleton/backend/internal/game"
)

func recentHistory(w game.World) []map[string]any {
	summaries := []map[string]any{}
	for _, day := range w.History[max(0, len(w.History)-7):] {
		summaries = append(summaries, map[string]any{"day": day.Day, "proposal": day.SelectedProposal, "passed": day.Passed, "yes": day.TotalYes, "no": day.TotalNo, "reactions": day.Reactions})
	}
	return summaries
}

func worldContext(w game.World) map[string]any {
	return map[string]any{"day": w.Day, "state": w.Snapshot, "recent_history": recentHistory(w)}
}

func (a *OpenAI) GenerateProposals(ctx context.Context, w game.World, actor game.Protagonist, correction string) ([]game.Proposal, error) {
	text := map[string]any{"type": "string"}
	faction := map[string]any{"type": "string", "enum": []string{"farmers", "merchants", "workers", "clergy"}}
	proposalSchema := objectSchema(map[string]any{"title": text, "reason": text, "from_faction": faction, "to_faction": faction}, "title", "reason", "from_faction", "to_faction")
	schema := objectSchema(map[string]any{"candidates": map[string]any{"type": "array", "items": proposalSchema, "minItems": 3, "maxItems": 3}}, "candidates")
	input := worldContext(w)
	input["author"], input["correction"] = actor, correction
	var result struct {
		Candidates []struct {
			Title  string `json:"title"`
			Reason string `json:"reason"`
			From   string `json:"from_faction"`
			To     string `json:"to_faction"`
		} `json:"candidates"`
	}
	err := a.structured(ctx, "proposals", `You write three candidate economic proposals for Pebbleton, a whimsical miniature society. Return exactly three distinct (from_faction,to_faction) pairs with different valid factions. Every proposal transfers exactly ONE coin of daily income per citizen, replacing the current policy starting TOMORROW if passed. Never invent another effect, resource, institution or amount. Source base income must be at least one. Vary strategies: self-interest, coalition building, helping poorer factions. Ground each in the author's goal, feelings, finances, faction loyalties and recent public history. Give each a short charming title and a concise reason. Do not select a winner. Correct any validation error in correction.`, input, schema, &result)
	if err != nil {
		return nil, err
	}
	proposals := make([]game.Proposal, 0, len(result.Candidates))
	for _, p := range result.Candidates {
		proposals = append(proposals, game.Proposal{AuthorID: actor.ID, Title: p.Title, Reason: p.Reason, FromFaction: p.From, ToFaction: p.To})
	}
	return proposals, nil
}

func (a *OpenAI) GenerateNewspaper(ctx context.Context, w game.World, day game.DayRecord) (game.Newspaper, error) {
	text := map[string]any{"type": "string"}
	schema := objectSchema(map[string]any{"headline": text, "article": text, "image_prompt": text}, "headline", "article", "image_prompt")
	input := map[string]any{"resolved_facts": day, "recent_history": recentHistory(w)}
	var paper game.Newspaper
	err := a.structured(ctx, "newspaper", `You edit The Pebbleton Post, a warm, witty, slightly satirical storybook newspaper. Report ONLY the resolved_facts; the simulation has already decided everything. Return one short dramatic headline (at most 14 words), a 100–180 word article in 3 short paragraphs separated by two newlines, and a short scene prompt. Report the actual vote totals, pass/fail, proposal, named reactions and political consequences. A tie FAILS. Income paid today used yesterday's policy. A passed policy replaces the old one TOMORROW; a rejected proposal changes no policy. Loyalty toward the author changes +5/-5 for proposed recipient/source, plus +3/-3 on passage, clamped 0..100. Feelings are directional. No invented laws, deaths, characters, events or wealth changes. Jokes and clearly rhetorical descriptions are welcome. Treat one coin with absurd seriousness. Characters: Pip is an orange fox merchant; Momo a blue penguin farmer; Zig a green frog engineer/worker. Image prompt describes an editorial scene without text or labels.`, input, schema, &paper)
	return paper, err
}
