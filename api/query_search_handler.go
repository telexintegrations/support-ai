package api

import(
	"log/slog"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
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
	txc := telexcom.NewTelexCom(s.AIService, s.DB, telexComClient)
	go func(txc *telexcom.TelexCom) {
		err := txc.ProcessTelexQuery(ctx, req)
		if err != nil {
			slog.Error("failed to handle telex request", "details", err)
		}
	}(txc)
	ctx.Status(202)

}