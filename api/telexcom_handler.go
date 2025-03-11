package api

import (
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/microcosm-cc/bluemonday"
	"github.com/telexintegrations/support-ai/telexcom"
)

// sendIntegrationJson returns the integration.json required by telex
func (s *Server) sendIntegrationJson(ctx *gin.Context) {
	ctx.Header("Content-Type", "application/json")
	ctx.JSON(http.StatusOK, telexcom.IntegrationJson)
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

	p := bluemonday.StrictPolicy()
	userQuery := p.Sanitize(req.Message)
	slog.Info("Sanitized user query", "message", userQuery)

	// respPayload := telexcom.TelexResponsePayload{
	// 	Message:   "Cannot provide support without setting your channel ID",
	// 	EventName: "info",
	// 	Username:  "TelexAI",
	// 	Status:    "info",
	// }
	ctx.Status(202)

}
