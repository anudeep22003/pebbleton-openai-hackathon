package ai

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

type OpenAI struct {
	key            string
	textModel      string
	decisionsModel string
	imageModel     string
	baseURL        string
	client         *http.Client
}

func New(key, textModel, decisionsModel, imageModel string) *OpenAI {
	return &OpenAI{key: key, textModel: textModel, decisionsModel: decisionsModel, imageModel: imageModel, baseURL: "https://api.openai.com/v1", client: &http.Client{Timeout: 3 * time.Minute}}
}

func (a *OpenAI) post(ctx context.Context, endpoint string, payload, result any) error {
	body, err := json.Marshal(payload)
	if err != nil {
		return err
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, a.baseURL+endpoint, bytes.NewReader(body))
	if err != nil {
		return err
	}
	req.Header.Set("Authorization", "Bearer "+a.key)
	req.Header.Set("Content-Type", "application/json")
	response, err := a.client.Do(req)
	if err != nil {
		return fmt.Errorf("OpenAI %s: %w", endpoint, err)
	}
	defer response.Body.Close()
	if response.StatusCode < 200 || response.StatusCode >= 300 {
		// Avoid logging request bodies or credentials, including in upstream error text.
		return fmt.Errorf("OpenAI %s returned HTTP %d (request %s); check model access, key and quota", endpoint, response.StatusCode, response.Header.Get("x-request-id"))
	}
	if err := json.NewDecoder(io.LimitReader(response.Body, 32<<20)).Decode(result); err != nil {
		return fmt.Errorf("OpenAI %s returned malformed JSON: %w", endpoint, err)
	}
	return nil
}

func objectSchema(properties map[string]any, required ...string) map[string]any {
	return map[string]any{"type": "object", "properties": properties, "required": required, "additionalProperties": false}
}

func (a *OpenAI) structured(ctx context.Context, name, instructions string, input any, schema map[string]any, target any) error {
	encoded, err := json.Marshal(input)
	if err != nil {
		return err
	}
	payload := map[string]any{
		"model": a.textModel, "instructions": instructions, "input": string(encoded), "max_output_tokens": 2000,
		"text": map[string]any{"format": map[string]any{"type": "json_schema", "name": name, "strict": true, "schema": schema}},
	}
	var response struct {
		Status string `json:"status"`
		Output []struct {
			Type    string `json:"type"`
			Content []struct {
				Type string `json:"type"`
				Text string `json:"text"`
			} `json:"content"`
		} `json:"output"`
	}
	if err := a.post(ctx, "/responses", payload, &response); err != nil {
		return err
	}
	if response.Status != "completed" {
		return fmt.Errorf("Responses %s was not completed: %s", name, response.Status)
	}
	var text strings.Builder
	for _, output := range response.Output {
		for _, content := range output.Content {
			if content.Type == "refusal" {
				return fmt.Errorf("Responses refused %s", name)
			}
			if output.Type == "message" && content.Type == "output_text" {
				text.WriteString(content.Text)
			}
		}
	}
	decoder := json.NewDecoder(strings.NewReader(text.String()))
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(target); err != nil {
		return fmt.Errorf("invalid %s output: %w", name, err)
	}
	return nil
}
