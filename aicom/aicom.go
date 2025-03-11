package aicom

import (
	"context"
	"encoding/json"
	"fmt"
	"os"

	genai "github.com/google/generative-ai-go/genai"
	"google.golang.org/api/option"
)

func InitGeminiClient() *genai.Client {
	apiKey := os.Getenv("GEMINI_API_KEY")
	if apiKey == "" {
		fmt.Println("GEMINI_API_KEY is missing")
		os.Exit(1)
	}

	ctx := context.Background()
	client, err := genai.NewClient(ctx, option.WithAPIKey(apiKey))
	if err != nil {
		fmt.Println("Failed to create Gemini client:", err)
		os.Exit(1)
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

	return client
}

func GetAIResponse(client *genai.Client, message string) (string, error) {

	ctx := context.Background()
	model := client.GenerativeModel("models/gemini-1.5-flash-latest")

	resp, err := model.GenerateContent(ctx, genai.Text(message))
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

func FormatResponse(data interface{}) ([]byte, error) {
	formattedJSON, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return nil, err
	}
	return formattedJSON, nil
}

func GetGeminiEmbedding(client *genai.Client, text string) ([]float32, error) {

    ctx := context.Background()
    model := client.EmbeddingModel("gemini-embedding-exp-03-07") // Use an embedding model

    res, err := model.EmbedContent(ctx, genai.Text(text))
    if err != nil {
        return nil, err
    }

    return res.Embedding.Values, nil
}