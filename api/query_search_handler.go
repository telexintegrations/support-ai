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

var GlobalReq telexcom.TelexChatPayload

func (s *Server) MakeQuerySearch(ctx *gin.Context) {
	if err := ctx.ShouldBindJSON(&GlobalReq); err != nil {
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
	fmt.Printf("Channel Id: %v \n", GlobalReq.ChannelID)
	fmt.Printf("Org Id: %v \n", GlobalReq.OrgId)
	fmt.Printf("Thread Id: %v \n", GlobalReq.ThreadID)
	fmt.Printf("Message: %v \n", format.StripHTMLTags(GlobalReq.Message))

	go func(txc *telexcom.TelexCom) {
		err := txc.ProcessTelexQuery(ctx, GlobalReq)
		if err != nil {
			slog.Error("failed to handle telex request", "details", err)
		}
	}(txc)
	ctx.Status(202)
}
