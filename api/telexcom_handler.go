package api

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/telexintegrations/support-ai/telexcom"
)

// sendIntegrationJson returns the integration.json required by telex
func (s *Server) sendIntegrationJson(ctx *gin.Context) {
	ctx.Header("Content-Type", "application/json")
	ctx.JSON(http.StatusOK, telexcom.IntegrationJson)
}
func (s *Server) sendNgrokJson(ctx *gin.Context) {
	ctx.Header("Content-Type", "application/json")
	ctx.JSON(http.StatusOK, telexcom.NgrokIntegrationJson)
}

func (s *Server) sendChromaJson(ctx *gin.Context) {
	ctx.Header("Content-Type", "application/json")
	ctx.JSON(http.StatusOK, telexcom.ChromaIntegrationJson)
}

func (s *Server) receiveChatQueries(ctx *gin.Context) {
	var req telexcom.TelexChatPayload
	if err := ctx.ShouldBindJSON(&req); err != nil {
		slog.Error("Invalid payload", "error", err)
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid request payload",
		})
		return
	}
	telexComClient := http.Client{
		Timeout: time.Second * 3,
	}
	txc := telexcom.NewTelexCom(s.AIService, s.DB, s.CDB, telexComClient)
	go func(txc *telexcom.TelexCom) {
		err := txc.ProcessTelexInputRequest(ctx, req)
		if err != nil {
			slog.Error("failed to handle telex request", "details", err)
		}
	}(txc)
	ctx.Status(202)

}

func (s *Server) FetchEmbeddings(ctx *gin.Context) {
	results, _ := s.DB.GetContentEmbeddings(ctx)
	jsonData, err := json.Marshal(results)
	if err != nil {
		// panic(err)
		fmt.Printf("error is:%v", err)
	}
	response := string(jsonData)
	fmt.Printf("response from the db is: %s", response)
	ctx.JSON(200, gin.H{
		"message": response,
	})
}
