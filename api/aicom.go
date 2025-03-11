package api

import (
	"encoding/json"
	"net/http"
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/telexintegrations/support-ai/aicom"
)

func (s *Server) RaggedResponse(c *gin.Context){

	query := "She gets a lot of money after her uncle dies" // should be dynamic to receive payload from telex
	
	ai_response, err := aicom.RaggingService(query)
	if err != nil {
		fmt.Println("Failed to get AI response: ", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to format response", "detail": err.Error()})
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

	//TODO POST DATA BACK TO TELEX SHOULD HAPPEN HERE ON THE FORMATTED RESPONSE

	c.Data(http.StatusOK, "application/json", formattedAnswer)
}


func (s *Server) TestConstructPrompt(c *gin.Context) {
	var req aicom.AIRequest
	err := c.ShouldBindJSON(&req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request format"})
		return
	}

	client, err := aicom.InitializeGeminiClient()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to initialise Gemini"})
		return
	}

	aiPrompt:= aicom.ConstructPrompt(req.Query, req.RetrievedDocs)

	response, genErr := aicom.GenerateGeminiResponse(client, aiPrompt)
	if genErr != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate response"})
		return
	}
	
	c.JSON(http.StatusOK, gin.H{"response": response})
}