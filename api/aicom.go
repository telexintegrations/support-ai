package api

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/telexintegrations/support-ai/aicom"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func (s *Server) RaggedResponse(c *gin.Context){

	query := "She gets a lot of money after her uncle dies"
	model := aicom.InitGeminiClient()
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
    

	startPrompt := fmt.Sprintf("Based on the stringed json below, %s recommend these movies based on the original search query: %s, and make your response as humanoid as possible", db_response, query)
	ai_response, err := aicom.GetAIResponse(model, startPrompt)
	if err != nil {
		fmt.Println("Failed to process file: ", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to process file", "detail": err.Error()})
		return
	}

	var parsedResponse map[string]interface{}
	parseErr := json.Unmarshal([]byte(ai_response), &parsedResponse)
	if parseErr != nil {
		// If unmarshalling fails, wrap it in a response struct
		parsedResponse = map[string]interface{}{"response": ai_response}
	}

	formattedAnswer, err := aicom.FormatResponse(parsedResponse)
	if err != nil {
		fmt.Println("Failed to format response: ", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to format response", "detail": err.Error()})
		return
	}
	c.Data(http.StatusOK, "application/json", formattedAnswer)
}