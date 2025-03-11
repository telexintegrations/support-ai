package aicom

import (
	"context"
	"strings"
	"fmt"
	"os"

	genai "github.com/google/generative-ai-go/genai"
	"google.golang.org/api/option"
)

func InitializeGeminiClient() (*genai.Client, error) {
	apiKey := os.Getenv("GEMINI_API_KEY")
	if apiKey == "" {
		os.Exit(1)
		return nil, fmt.Errorf("GEMINI_API_KEY is missing")
	}

	ctx := context.Background()
	client, err := genai.NewClient(ctx, option.WithAPIKey(apiKey))
	if err != nil {
		os.Exit(1)
		return nil, fmt.Errorf("Failed to create Gemini client: %s", err)
	}

	fmt.Println("Listing available Gemini models...")
	modelIterator := client.ListModels(ctx)

	// Iterate over available models
	for {
		model, err := modelIterator.Next()
		if err != nil {
			break // Exit loop when there are no more models
		}
		fmt.Printf("Available Model: %s\n", model.Name)
	}

	return client, nil
}

func GenerateGeminiResponse(client *genai.Client, prompt string) (string, error) {
	// Construct a context-aware prompt
	ctx := context.Background()
	model := client.GenerativeModel("models/gemini-1.5-flash-latest")
	resp, err := model.GenerateContent(ctx, genai.Text(prompt))

	if err != nil {
		return "", err
	}

	if len(resp.Candidates) == 0 || len(resp.Candidates[0].Content.Parts) == 0 {
		return "", fmt.Errorf("no response from Gemini")
	}

	// Convert first response part to text
	if textPart, ok := resp.Candidates[0].Content.Parts[0].(genai.Text); ok {
		return string(textPart), nil
	}

	return "", fmt.Errorf("unexpected response format")
}

func ConstructPrompt(query string, documents []string) string {
	var sb strings.Builder

	sb.WriteString(fmt.Sprintf("User Query: %s\n\n", query))

	sb.WriteString("Relevant Information: \n")
	for i, doc := range documents {
		sb.WriteString(fmt.Sprintf("%d. %s\n", i+1, doc))
	}

	sb.WriteString("\nGenerate a human-like response based on the provided context.")

	return sb.String()
}
