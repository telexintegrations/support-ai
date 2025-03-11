package aicom

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"time"

	genai "github.com/google/generative-ai-go/genai"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"google.golang.org/api/option"
)

type AIService interface{
	GetAIResponse(message string) (string, error)
	GetGeminiEmbedding(text string) ([]float32, error)
	RaggingService(query string) (string, error)
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

func FormatResponse(data interface{}) ([]byte, error) {
	formattedJSON, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return nil, err
	}
	return formattedJSON, nil
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

func (a *AIServiceImpl)RaggingService(query string) (string, error){

	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
	defer cancel()
	
	//TODO VECTOR SEARCH SERVICE AND RESULTS GOES HERE, SHOULD RETURN A CURSOR INTO RESULTS
	var results *mongo.Cursor

	var response []bson.M
	var db_response string
    // Iterate over results
	if results != nil{
		for results.Next(ctx) {
			var result bson.M
			if err := results.Decode(&result); err != nil {
				fmt.Println("Error decoding result:", err)
				continue
			}
			fmt.Printf("Movie Name: %s,\nMovie Plot: %s\n\n", result["title"], result["plot"])
			response = append(response, result)
		}
		
		jsonData, err := json.Marshal(response)
		 if err != nil {
			  panic(err)
		 }
	
		db_response = string(jsonData)
	}
    

	startPrompt := fmt.Sprintf("You are a customer support officer, based on the stringed json below, %s respond to this query: %s, and make your response as humanoid as possible.", db_response, query)
	ai_response, err := a.GetAIResponse(startPrompt)
	if err != nil {
		fmt.Println("Failed to process file: ", err)
		return "", err
	}

	return ai_response, err

	
}