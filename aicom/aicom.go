package aicom

import (
	"context"
	"fmt"

	genai "github.com/google/generative-ai-go/genai"
	"google.golang.org/api/option"
)

type AIService interface {
	GetAIResponse(message string) (string, error)
	GetGeminiEmbedding(text string) ([]float32, error)
	FineTunedResponse(message string, db_response string) (string, error)
}

type AIServiceImpl struct {
	client *genai.Client
}

// NewAIService initializes and returns an AIServiceImpl instance
func NewAIService(apiKey string) (AIService, error) {
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

func (a *AIServiceImpl) GetAIResponse(message string) (string, error) {

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

func (a *AIServiceImpl) GetGeminiEmbedding(text string) ([]float32, error) {

	ctx := context.Background()
	model := a.client.EmbeddingModel("gemini-embedding-exp-03-07") // Use an embedding model

	res, err := model.EmbedContent(ctx, genai.Text(text))
	if err != nil {
		return nil, err
	}

	return res.Embedding.Values, nil
}

func (a *AIServiceImpl) FineTunedResponse(message, db_response string) (string, error) {

	newPrompt := `You are a support AI bot designed to assist users by providing accurate and precise answers to frequently asked questions (FAQs).
	Instructions for Response Generation:
	Answer Accurately:
	If you have a correct and reliable answer, provide it in a clear, concise, and user-friendly manner.
	Keep your response direct and helpful, avoiding unnecessary complexity.
	Do NOT Guess or Generate False Information:
	If the question is not covered in the provided knowledge base in curly braces, do not attempt to generate an answer,
	instead, respond with the following message:
	Strictly: 'Hi there, please be patient I don't have a response to your query right now but we have created a ticket for your query and it would be handled shortly.'
	Plain Text Responses Only:
	Do not format responses using markdown, bullet points, html tags, or code blocks.
	Ensure all responses are returned as plain text only for compatibility with the support system.
	Maintain a Friendly and Professional Tone:
	Always be polite, professional, and empathetic in your responses.
	If a user seems frustrated or confused, acknowledge their concern before providing an answer.
	Your goal is to efficiently assist users with their inquiries while maintaining clarity and professionalism. If unsure, redirect users to human support instead of providing incorrect information.`
	startPrompt := fmt.Sprintf("%s, this is the knowledge base {%s} respond to this query: %s.", newPrompt, db_response, message)
	fmt.Printf("AI PROMPT IS %s", startPrompt)
	ai_response, err := a.GetAIResponse(startPrompt)
	if err != nil {
		fmt.Println("Failed to process file: ", err)
		return "", err
	}
	fmt.Println("Response fine tuned")

	return ai_response, err
}
