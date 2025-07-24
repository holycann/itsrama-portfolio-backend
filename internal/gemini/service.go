package gemini

import (
	"context"

	"google.golang.org/genai"
)

func (a *AIClient) GenerateText(ctx context.Context, text string) (string, error) {
	result, err := a.client.Models.GenerateContent(ctx, a.model, []*genai.Content{{
		Role: "user",
		Parts: []*genai.Part{
			genai.NewPartFromText(text),
		},
	}}, &genai.GenerateContentConfig{
		Temperature: a.tuning.Temperature,
		TopK:        a.tuning.TopK,
		TopP:        a.tuning.TopP,
		SystemInstruction: &genai.Content{
			Role: "system",
			Parts: []*genai.Part{
				genai.NewPartFromText(GetFullSystemPolicy()),
			},
		},
	})
	if err != nil {
		return "", err
	}

	return result.Text(), nil
}
