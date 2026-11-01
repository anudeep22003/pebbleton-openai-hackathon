package ai

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"pebbleton/backend/internal/game"
)

func TestOpenAIContractsAndFailures(t *testing.T) {
	for _, failure := range []string{"", "refusal", "missing_probability", "newspaper", "http"} {
		t.Run(failure, func(t *testing.T) {
			calls := 0
			handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				calls++
				if r.Method != "POST" || r.Header.Get("Authorization") != "Bearer test-key" {
					t.Error("bad transport contract")
				}
				if failure == "http" {
					w.WriteHeader(429)
					return
				}
				var body map[string]json.RawMessage
				if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
					t.Fatal(err)
				}
				w.Header().Set("Content-Type", "application/json")
				switch calls {
				case 1, 4:
					if r.URL.Path != "/responses" {
						t.Error("wrong generation endpoint")
					}
					var config struct {
						Format struct {
							Type   string `json:"type"`
							Strict bool   `json:"strict"`
						} `json:"format"`
					}
					json.Unmarshal(body["text"], &config)
					if config.Format.Type != "json_schema" || !config.Format.Strict {
						t.Error("missing strict structured outputs")
					}
					if calls == 4 && failure == "newspaper" {
						fmt.Fprint(w, `{"status":"incomplete","output":[]}`)
						return
					}
					text := `{"headline":"One coin, much ado","article":"A factual newspaper.","image_prompt":"Town square."}`
					if calls == 1 {
						text = `{"candidates":[{"title":"A","reason":"A reason","from_faction":"farmers","to_faction":"merchants"},{"title":"B","reason":"B reason","from_faction":"merchants","to_faction":"workers"},{"title":"C","reason":"C reason","from_faction":"clergy","to_faction":"farmers"}]}`
					}
					json.NewEncoder(w).Encode(map[string]any{"status": "completed", "output": []any{map[string]any{"type": "message", "content": []any{map[string]any{"type": "output_text", "text": text}}}}})
				case 2:
					if r.URL.Path != "/decisions" {
						t.Error("wrong choice endpoint")
					}
					var questions []decisionQuestion
					json.Unmarshal(body["questions"], &questions)
					if len(questions) != 1 || len(questions[0].Choices) != 3 {
						t.Error("bad choice question")
					}
					if failure == "refusal" {
						fmt.Fprint(w, `{"answers":[{"type":"refusal","name":"proposal_selection"}]}`)
						return
					}
					fmt.Fprint(w, `{"answers":[{"type":"choice","name":"proposal_selection","choice":"proposal-1-1","confidence":1,"probabilities":[{"value":"proposal-1-1","probability":0},{"value":"proposal-1-2","probability":1},{"value":"proposal-1-3","probability":0}]}]}`)
				case 3:
					if r.URL.Path != "/decisions" {
						t.Error("wrong predicate endpoint")
					}
					var questions []decisionQuestion
					json.Unmarshal(body["questions"], &questions)
					if len(questions) != 6 {
						t.Error("expected six batched predicates")
					}
					answers := []map[string]any{}
					for _, q := range questions {
						answers = append(answers, map[string]any{"type": "predicate", "name": q.Name, "probability": 1})
					}
					if failure == "missing_probability" {
						delete(answers[0], "probability")
					}
					json.NewEncoder(w).Encode(map[string]any{"answers": answers})
				}
			})
			provider := New("test-key", "text", "gpt-6-luna", "image")
			provider.baseURL = "https://fixture.test"
			provider.client.Transport = roundTripFunc(func(r *http.Request) (*http.Response, error) {
				response := httptest.NewRecorder()
				handler.ServeHTTP(response, r)
				return response.Result(), nil
			})
			next, err := game.Advance(context.Background(), game.NewWorld(42, "test"), provider)
			if failure == "" || failure == "newspaper" {
				if err != nil {
					t.Fatal(err)
				}
				if calls != 4 || next.History[0].SelectedProposal.ID != "proposal-1-2" || next.History[0].TotalYes != 40 || next.History[0].Article == "" {
					t.Fatal("incorrect live-provider flow")
				}
			} else if err == nil {
				t.Fatal("malformed or refused answer advanced the day")
			}
		})
	}
}

type roundTripFunc func(*http.Request) (*http.Response, error)

func (f roundTripFunc) RoundTrip(r *http.Request) (*http.Response, error) { return f(r) }
