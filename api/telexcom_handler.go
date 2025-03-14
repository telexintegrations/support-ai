package api

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/microcosm-cc/bluemonday"
	"github.com/telexintegrations/support-ai/format"
	"github.com/telexintegrations/support-ai/telexcom"
)

func processQuery(query string) (string, string) {
	var remainingContent string
	var task string
	if strings.HasPrefix(query, "/upload") {
		remainingContent = strings.TrimPrefix(query, "/upload")
		task = "/upload"
		fmt.Println("Processing upload with:", strings.TrimSpace(remainingContent))
	} else if strings.HasPrefix(query, "/help") {
		remainingContent = strings.TrimPrefix(query, "/help")
		task = "/help"
		fmt.Println("Processing help with:", strings.TrimSpace(remainingContent))
	} else {
		remainingContent = query
		task = ""
		// fmt.Println("Invalid command or unrecognized query.")
	}
	return remainingContent, task
}

// sendIntegrationJson returns the integration.json required by telex
func (s *Server) sendIntegrationJson(ctx *gin.Context) {
	ctx.Header("Content-Type", "application/json")
	ctx.JSON(http.StatusOK, telexcom.IntegrationJson)
}
func (s *Server) sendNgrokJson(ctx *gin.Context) {
	ctx.Header("Content-Type", "application/json")
	ctx.JSON(http.StatusOK, telexcom.NgrokIntegrationJson)
}

func (s *Server) receiveChatQueries(ctx *gin.Context) {
	var return_response string
	var req telexcom.TelexChatPayload
	txc := telexcom.NewTelexCom()
	if err := ctx.ShouldBindJSON(&req); err != nil {
		slog.Error("Invalid payload", "error", err)
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid request payload",
		})
		return
	}

	p := bluemonday.StrictPolicy()
	userQuery := p.Sanitize(req.Message)
	fmt.Printf("Request Message is: %s",userQuery)
	var task string
	userQuery, task = processQuery(userQuery)

	if task == "/upload"{
		chunks := format.ChunkTextByParagraph(userQuery, 20)
		var embedded_chunks [][]float32
		for _, chunk := range chunks{
			chunk_embedding, err := s.AIService.GetGeminiEmbedding(chunk)
			if err != nil{
				go txc.GenerateResponseToQuery(ctx, "Error uploading to AI, check Credits for AI or contact support", req.ChannelID)
			}
			embedded_chunks = append(embedded_chunks, chunk_embedding)
		}
	
		err := s.DB.InsertIntoEmbeddingCollection(ctx, chunks, embedded_chunks, "")
		if err != nil{
			go txc.GenerateResponseToQuery(ctx, "Error uploading", req.ChannelID)
		}else{
			go txc.GenerateResponseToQuery(ctx, "Content Uploaded, you can use /help to send queries", req.ChannelID)
		}
	}else if task == ""{
		fmt.Println("no slash command")
		return
	}else{
		query_embedding, err := s.AIService.GetGeminiEmbedding(userQuery)
		if err != nil{
			fmt.Println(err)
		}
		// fmt.Printf("query_embedding is: %v", query_embedding)
		//vector search
		db_raw, _ := s.DB.SearchVectorFromContentEmbedding(ctx, query_embedding, 3)
		var db_response string
		for _, res := range db_raw{
			db_response += res.Content
		}
		fmt.Println("Vector search completed")

  
	  	// db_response := string(jsonData)
		fmt.Printf("db_response is: %s",db_response)
		//fine tune response
		ai_response, err := s.AIService.FineTunedResponse(userQuery, db_response)
		return_response = ai_response
		//post to telex
		go txc.GenerateResponseToQuery(ctx, ai_response, req.ChannelID)
	}
	ctx.JSON(200, gin.H{
		"event_name": "Query Support",
		"message": return_response,
		"status": "success",
		// "reply_in_thread": true,
		"username": "Support AI",
	  })

}

func (s *Server) FetchEmbeddings(ctx *gin.Context){
	results, _ := s.DB.GetContentEmbeddings(ctx)
	jsonData, err := json.Marshal(results)
	if err != nil {
		// panic(err)
		fmt.Printf("error is:%v", err)
   }
   response := string(jsonData)
   fmt.Printf("response from the db is: %s",response)
	ctx.JSON(200, gin.H{
		"message": response,
	})
}

