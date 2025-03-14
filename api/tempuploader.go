package api

import (
	"encoding/json"
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/microcosm-cc/bluemonday"
	"github.com/telexintegrations/support-ai/format"
	"github.com/telexintegrations/support-ai/telexcom"
)

func (s *Server) UploadFiles(ctx *gin.Context) {

	// var return_response string
	// txc := telexcom.NewTelexCom()
	var userQuery string
	var req telexcom.TelexChatPayload
	contentType := ctx.GetHeader("Content-Type")
	p := bluemonday.StrictPolicy()
	var task string

	if strings.HasPrefix(contentType, "multipart/form-data") {
		// Handle multipart/form-data
		log.Println("request content type is multitype")
		jsonDataString := ctx.PostForm("jsonData")

		if jsonDataString == "" {
			log.Println("jsonData is empty")
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "jsonData is missing"})
			return
		}
		jsonDataBytes, err := json.MarshalIndent(req, "", "  ")
		if err != nil {
			log.Printf("Error marshaling JSON for logging: %v", err)
			ctx.JSON(http.StatusBadRequest, gin.H{
				"error": "Failed to get JSON",
			})
		}
		log.Printf("Parsed JSON: %s", jsonDataBytes)
		if err := json.Unmarshal([]byte(jsonDataString), &req); err != nil {
			slog.Error("Invalid JSON in jsonData", "error", err)
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON in jsonData"})
			return
		}

		file, _ := ctx.FormFile("file")
		if file != nil {
			// Handle the uploaded file
			fmt.Printf("File uploaded filename: %s", file.Filename)
		}
		if file == nil {
			// Handle the uploaded file
			fmt.Printf("No File uploaded")
			ctx.JSON(http.StatusBadRequest, gin.H{
				"error": "Failed to read file",
			})
		}
		//TODO - handle the file and read content
		fileContent, err := format.ExtractText(file)
		if err != nil {
			fmt.Println("failed to read file: %w", err)
			ctx.JSON(http.StatusBadRequest, gin.H{
				"error": "Failed to read file",
			})
		}
		fmt.Printf("File content:%s", fileContent)
		combinedQuery := fmt.Sprintln(req.Message + fileContent)
		userQuery = p.Sanitize(combinedQuery)

	} else if strings.HasPrefix(contentType, "application/json") {
		// Handle application/json
		log.Println("request content type is json")
		if err := ctx.ShouldBindJSON(&req); err != nil {
			slog.Error("Invalid payload", "error", err)
			ctx.JSON(http.StatusBadRequest, gin.H{
				"error": "Invalid request payload",
			})
			return
		}
		userQuery = p.Sanitize(req.Message)
	}

	userQuery, task = processQuery(userQuery)
	fmt.Printf("Task is: %s", task)
	ctx.JSON(200, gin.H{
		"event_name": task,
		"message":    userQuery,
		"status":     "success",
		"thread_id":  "",
		"username":   "Support AI",
	})

}
