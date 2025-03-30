package api

import (
	"fmt"
	"net/http"

	"time"
	"github.com/gin-gonic/gin"
	"github.com/telexintegrations/support-ai/format"
	"github.com/telexintegrations/support-ai/telexcom"
)

var GlobalReq telexcom.TelexChatPayload

func (s *Server) UploadFilesToDb(ctx *gin.Context) {
	form, err := ctx.MultipartForm()
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "error", "error": "Failed to parse form data"})
		return
	}

	files := form.File["files"]
	if len(files) == 0 {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "error", "error": "Bad request. No files uploaded."})
		return
	}

	var extractedTexts string

	for _, fileHeader := range files {
		file, err := fileHeader.Open()
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"status": "error", "error": "Failed to open file"})
			return
		}
		defer file.Close()

		text, err := format.ExtractText(fileHeader)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"status": "error", "error": fmt.Sprintf("Failed to extract text. %v", err)})
			return
		}
		if text == "" {
			ctx.JSON(http.StatusInternalServerError, gin.H{"status": "error", "error": "Failed to extract text. File may be broken."})
			return
		}

		extractedTexts += text
	}

	telexComClient := http.Client{
		Timeout: time.Second * 3,
	}
	txc := telexcom.NewTelexCom(s.AIService, s.DB, s.CDB, telexComClient)
	uploadErr := txc.ProcessTelexUpload(ctx, extractedTexts, GlobalReq.ChannelID)
	if uploadErr != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"status": "error", "error": fmt.Sprintf("Process failed. %v", uploadErr)})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"status": "success",
		"message": "File uploaded succesfully!"})
}
