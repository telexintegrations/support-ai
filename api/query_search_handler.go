package api

import (
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/telexintegrations/support-ai/format"
	"github.com/telexintegrations/support-ai/telexcom"
)

func (s *Server) MakeQuerySearch(ctx *gin.Context) {
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

	fmt.Printf("\n Request details: %v \n", req)
	fmt.Printf("Channel Id: %v \n", req.ChannelID)
	fmt.Printf("Org Id: %v \n", req.OrgId)
	fmt.Printf("Thread Id: %v \n", req.ThreadID)
	fmt.Printf("Message: %v \n", format.StripHTMLTags(req.Message))

	go func(txc *telexcom.TelexCom) {
		if len(req.Media) > 0 {
			go func() {
				err := txc.ProcessTelexDownloadAndExtraction(ctx, req)
				if err != nil {
					slog.Error("failed to handle telex request", "details", err)
				}
			}()
		}

		err := txc.ProcessTelexQuery(ctx, req)
		if err != nil {
			slog.Error("failed to handle telex request", "details", err)
		}
	}(txc)
	ctx.Status(202)
}
