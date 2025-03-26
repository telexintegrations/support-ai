package api

import (
	"fmt"
	"net/http"

	"time"

	"github.com/gin-gonic/gin"
	"github.com/telexintegrations/support-ai/telexcom"
)

func (s *Server) UploadTextToDb(ctx *gin.Context) {
	var requestData telexcom.UploadTextRequestData

	if err := ctx.ShouldBindJSON(&requestData); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "error", "error": "Invalid JSON input"})
		return
	}

	telexComClient := http.Client{
		Timeout: time.Second * 3,
	}
	txc := telexcom.NewTelexCom(s.AIService, s.DB, s.CDB, telexComClient)
	uploadErr := txc.ProcessTelexUpload(ctx, requestData.FileText, "0195d2fb-997a-7665-a413-ea5a653bb240")
	if uploadErr != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"status": "error", "error": fmt.Sprintf("Process failed. %v", uploadErr)})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"status": "success",
		"message": "File uploaded succesfully!"})
}
