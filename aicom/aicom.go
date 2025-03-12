package aicom

import (
	"context"
	"fmt"

	genai "github.com/google/generative-ai-go/genai"
	"google.golang.org/api/option"
)

type AIService interface{
	GetAIResponse(message string) (string, error)
	GetGeminiEmbedding(text string) ([]float32, error)
	FineTunedResponse(message string, db_response string) (string, error)
}

type AIServiceImpl struct {
	client *genai.Client
}

// NewAIService initializes and returns an AIServiceImpl instance
func NewAIService(apiKey string) (AIService, error) {
	fmt.Printf("API KEY is %s", apiKey)
	if apiKey == "" {
		return nil, fmt.Errorf("GEMINI_API_KEY is missing")
	}

	ctx := context.Background()
	client, err := genai.NewClient(ctx, option.WithAPIKey(apiKey))
	if err != nil {
		return nil, fmt.Errorf("failed to create Gemini client: %w", err)
	}

	return &AIServiceImpl{client: client}, nil
}

func (a *AIServiceImpl)GetAIResponse(message string) (string, error) {

	ctx := context.Background()
	model := a.client.GenerativeModel("models/gemini-1.5-flash-latest")

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


func (a *AIServiceImpl)GetGeminiEmbedding(text string) ([]float32, error) {

    ctx := context.Background()
    model := a.client.EmbeddingModel("gemini-embedding-exp-03-07") // Use an embedding model

    res, err := model.EmbedContent(ctx, genai.Text(text))
    if err != nil {
        return nil, err
    }

    return res.Embedding.Values, nil
}


func (a *AIServiceImpl)FineTunedResponse(message, db_response string) (string, error) {

	startPrompt := fmt.Sprintf("You are a customer support officer, based on only the context below, %s respond to this query: %s, and make your response as humanoid as possible.", db_response, message)
	ai_response, err := a.GetAIResponse(startPrompt)
	if err != nil {
		fmt.Println("Failed to process file: ", err)
		return "", err
	}
	fmt.Println("Response fine tuned")

	return ai_response, err
}