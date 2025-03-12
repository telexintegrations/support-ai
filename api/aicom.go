package api

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/telexintegrations/support-ai/format"
)

func (s *Server) BasicResponse(c *gin.Context){

	query := "Checkout Telex.im and summarize it" // should be dynamic to receive payload from telex
	
	ai_response,_ := s.AIService.GetAIResponse(query)

	var parsedResponse map[string]interface{}
	parseErr := json.Unmarshal([]byte(ai_response), &parsedResponse)
	if parseErr != nil {
		// If unmarshalling fails, wrap it in a response struct
		parsedResponse = map[string]interface{}{"response": ai_response}
	}

	formattedAnswer, err := format.FormatResponse(parsedResponse)
	if err != nil {
		fmt.Println("Failed to format response: ", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to format response", "detail": err.Error()})
		return
	}

	//TODO POST DATA BACK TO TELEX SHOULD HAPPEN HERE ON THE FORMATTED RESPONSE


	c.Data(http.StatusOK, "application/json", formattedAnswer)
}