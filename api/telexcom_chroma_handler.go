package api

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"log/slog"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/telexintegrations/support-ai/telexcom"
)

func (s *Server) ReceiveChatQueries2(ctx *gin.Context) {

	//read the raw body
	body, err := io.ReadAll(ctx.Request.Body)
	if err != nil {
		log.Printf("Error reading body: %v", err)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	fmt.Printf("Received payload: %s\n", body)             // Print to console
	ctx.Request.Body = io.NopCloser(bytes.NewBuffer(body)) //return the body back

	// get necessary details from payload
	var req telexcom.TelexChatPayload
	if err := ctx.ShouldBindJSON(&req); err != nil {
		slog.Error("Invalid payload", "error", err)
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid request payload",
		})
		return
	}
	//perform vector search
	query := req.Message
	fmt.Printf("query is %s", query)
	telexComClient := http.Client{
		Timeout: time.Second * 3,
	}
	txc := telexcom.NewTelexCom(s.AIService, s.DB, s.CDB, telexComClient)

	err = txc.ProcessTelexChromaInputRequest(ctx, req)
	if err != nil {
		slog.Error("failed to handle telex request", "details", err)
	}

	//fine tune response

	//send result to telex

	ctx.Status(202)
	// ctx.JSON(http.StatusOK, gin.H{"message": "Payload received"})

}
