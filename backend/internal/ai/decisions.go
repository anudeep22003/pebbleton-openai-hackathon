package ai

import (
	"context"
	"encoding/json"
	"fmt"

	"pebbleton/backend/internal/game"
)

type decisionQuestion struct {
	Type         string           `json:"type"`
	Name         string           `json:"name"`
	Instructions string           `json:"instructions"`
	Choices      []decisionChoice `json:"choices,omitempty"`
}
type decisionChoice struct {
	Value       string `json:"value"`
	Description string `json:"description"`
}
type decisionAnswer struct {
	Type          string   `json:"type"`
	Name          string   `json:"name"`
	Probability   *float64 `json:"probability"`
	Probabilities []struct {
		Value       string   `json:"value"`
		Probability *float64 `json:"probability"`
	} `json:"probabilities"`
}

func (a *OpenAI) decide(ctx context.Context, input any, questions []decisionQuestion) ([]decisionAnswer, error) {
	encoded, err := json.Marshal(input)
	if err != nil {
		return nil, err
	}
	var response struct {
		Answers []decisionAnswer `json:"answers"`
	}
	if err := a.post(ctx, "/decisions", map[string]any{"model": a.decisionsModel, "input": string(encoded), "questions": questions}, &response); err != nil {
		return nil, err
	}
	if len(response.Answers) != len(questions) {
		return nil, fmt.Errorf("Decisions returned missing answers")
	}
	for i, answer := range response.Answers {
		if answer.Type == "refusal" {
			return nil, fmt.Errorf("Decisions refused %s", questions[i].Name)
		}
		if answer.Type != questions[i].Type || answer.Name != questions[i].Name {
			return nil, fmt.Errorf("Decisions returned unexpected answer for %s", questions[i].Name)
		}
	}
	return response.Answers, nil
}

func (a *OpenAI) ChooseProposal(ctx context.Context, w game.World, actor game.Protagonist, proposals []game.Proposal) (map[string]float64, error) {
	question := decisionQuestion{Type: "choice", Name: "proposal_selection", Instructions: "Which supplied proposal would this protagonist most plausibly introduce, given their goal, occupation, money, feelings, faction loyalties and recent public history? Consider self-interest and plausible coalition building. The proposal replaces the current income policy starting tomorrow if passed."}
	for _, p := range proposals {
		question.Choices = append(question.Choices, decisionChoice{p.ID, fmt.Sprintf("%s: %s Transfer one income coin per citizen from %s to %s.", p.Title, p.Reason, p.FromFaction, p.ToFaction)})
	}
	input := worldContext(w)
	input["author"] = actor
	answers, err := a.decide(ctx, input, []decisionQuestion{question})
	if err != nil {
		return nil, err
	}
	probabilities := map[string]float64{}
	for _, item := range answers[0].Probabilities {
		if _, duplicate := probabilities[item.Value]; duplicate || item.Probability == nil {
			return nil, fmt.Errorf("Decisions returned duplicate or missing choice probability")
		}
		probabilities[item.Value] = *item.Probability
	}
	// The engine validates the full distribution and samples it, ignoring choice/confidence.
	return probabilities, nil
}

func (a *OpenAI) EvaluateReactions(ctx context.Context, w game.World, proposal game.Proposal) (map[string]float64, error) {
	questions := []decisionQuestion{}
	for _, p := range w.Protagonists {
		if p.ID != proposal.AuthorID {
			questions = append(questions, decisionQuestion{Type: "predicate", Name: p.ID, Instructions: fmt.Sprintf("%s publicly supports this proposal, given their goal, occupation, money, feelings toward the author and recent public history. Judge likely support, not moral correctness.", p.Name)})
		}
	}
	for _, f := range w.Factions {
		questions = append(questions, decisionQuestion{Type: "predicate", Name: f.ID, Instructions: fmt.Sprintf("A randomly selected citizen of %s votes YES, given their money per person, income, loyalty to the author and recent history. Compare the proposed replacement policy with the existing policy; include loss or recovery of previous benefits. Estimate an individual citizen's support, not the likelihood that a majority passes it.", f.Name)})
	}
	input := worldContext(w)
	input["proposal"] = proposal
	answers, err := a.decide(ctx, input, questions)
	if err != nil {
		return nil, err
	}
	probabilities := map[string]float64{}
	for _, answer := range answers {
		if answer.Probability == nil {
			return nil, fmt.Errorf("Decisions omitted probability for %s", answer.Name)
		}
		probabilities[answer.Name] = *answer.Probability
	}
	return probabilities, nil
}
