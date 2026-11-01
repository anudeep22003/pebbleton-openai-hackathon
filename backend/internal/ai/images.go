package ai

import (
	"bytes"
	"context"
	"encoding/base64"
	"fmt"
	"image/png"
)

func (a *OpenAI) GenerateIllustration(ctx context.Context, prompt string) ([]byte, error) {
	var response struct {
		Data []struct {
			Base64 string `json:"b64_json"`
		} `json:"data"`
	}
	payload := map[string]any{
		"model": a.imageModel, "n": 1, "size": "1536x1024", "quality": "medium", "output_format": "png",
		"prompt": "Create a landscape editorial illustration for The Pebbleton Post, a whimsical tiny society newspaper. Consistent art direction: textured gouache, ink outlines, warm cream paper, sage-green hills, terracotta roofs, a cozy European storybook town. Recurring characters: Pip is an ambitious orange fox in a burgundy merchant waistcoat; Momo is a friendly BLUE penguin in practical ochre farmer overalls and a straw hat; Zig is a clever GREEN frog in teal engineer overalls with brass goggles. Show only these protagonists, plus anonymous animal citizens if useful. Gentle political satire, expressive faces, one coin treated with excessive importance. Do not draw any text, labels, letters, headlines or numbers. Illustrate the supplied facts without inventing laws or events. Scene: " + prompt,
	}
	if err := a.post(ctx, "/images/generations", payload, &response); err != nil {
		return nil, err
	}
	if len(response.Data) != 1 || response.Data[0].Base64 == "" {
		return nil, fmt.Errorf("image generation returned no image")
	}
	data, err := base64.StdEncoding.DecodeString(response.Data[0].Base64)
	if err != nil {
		return nil, fmt.Errorf("invalid image encoding: %w", err)
	}
	if _, err := png.DecodeConfig(bytes.NewReader(data)); err != nil {
		return nil, fmt.Errorf("invalid PNG illustration: %w", err)
	}
	return data, nil
}
