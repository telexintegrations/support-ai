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
	queryToLower := strings.ToLower(query)
	if strings.HasPrefix(queryToLower, "/upload") {
		remainingContent = strings.TrimPrefix(query, "/upload")
		task = "/upload"
		fmt.Println("Processing upload with:", strings.TrimSpace(remainingContent))
	} else if strings.HasPrefix(queryToLower, "/help") {
		remainingContent = strings.TrimPrefix(query, "/help")
		task = "/help"
		fmt.Println("Processing help with:", strings.TrimSpace(remainingContent))
	} else if strings.HasPrefix(queryToLower, "/change-contex") {
		remainingContent = strings.TrimPrefix(query, "/change-context")
		task = "/change-context"
		fmt.Println("Processing /change-context:", strings.TrimSpace(remainingContent))
	} else if strings.HasPrefix(queryToLower, "use") {
		task = ""
	} else {
		remainingContent = query
		task = "use"
	}
	return remainingContent, task
}

var lastResponseToTelex string

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

	var req telexcom.TelexChatPayload
	if err := ctx.ShouldBindJSON(&req); err != nil {
		slog.Error("Invalid payload", "error", err)
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid request payload",
		})
		return
	}
	if lastResponseToTelex == req.Message {
		ctx.Status(http.StatusAccepted)
		return
	}
	txc := telexcom.NewTelexCom()
	p := bluemonday.StrictPolicy()
	userQuery := p.Sanitize(req.Message)

	var task string
	userQuery, task = processQuery(userQuery)

	switch task {
	case "/upload":
		chunks := format.ChunkTextByParagraph(userQuery, 20)
		var embedded_chunks [][]float32
		for _, chunk := range chunks {
			chunk_embedding, err := s.AIService.GetGeminiEmbedding(chunk)
			if err != nil {
				fmt.Println(err)
				lastResponseToTelex = "sorry, couldn't process your upload"
				go txc.GenerateResponseToQuery(ctx, lastResponseToTelex, req.ChannelID)
			}
			embedded_chunks = append(embedded_chunks, chunk_embedding)
		}

		err := s.DB.InsertIntoEmbeddingCollection(ctx, chunks, embedded_chunks, req.ChannelID)
		if err != nil {
			go txc.GenerateResponseToQuery(ctx, "Error uploading", req.ChannelID)
		} else {
			lastResponseToTelex = "Content Uploaded, you can use /help to send queries"
			go txc.GenerateResponseToQuery(ctx, lastResponseToTelex, req.ChannelID)
		}
	case "/help":
		query_embedding, err := s.AIService.GetGeminiEmbedding(userQuery)
		if err != nil {
			fmt.Println(err)
			lastResponseToTelex = "sorry couldn't process your query"
			go txc.GenerateResponseToQuery(ctx, lastResponseToTelex, req.ChannelID)
			return
		}

		db_raw, err := s.DB.SearchVectorFromContentEmbedding(ctx, query_embedding, req.ChannelID, 10)
		if err != nil {
			fmt.Println(err)
			lastResponseToTelex = "sorry couldn't process your query"
			go txc.GenerateResponseToQuery(ctx, lastResponseToTelex, req.ChannelID)
			return
		}
		var db_response string
		for _, res := range db_raw {
			db_response += res.Content
		}
		fmt.Println("Vector search completed")

		fmt.Printf("db_response is: %s", db_response)
		//fine tune response
		ai_response, err := s.AIService.FineTunedResponse(userQuery, db_response)
		if err != nil {
			fmt.Println("failed to generate AI response")
			lastResponseToTelex = "sorry couldn't process your query"
			go txc.GenerateResponseToQuery(ctx, lastResponseToTelex, req.ChannelID)
		}
		//post to telex
		lastResponseToTelex = ai_response
		go txc.GenerateResponseToQuery(ctx, lastResponseToTelex, req.ChannelID)
	case "/change-context":
		chunks := format.ChunkTextByParagraph(userQuery, 20)
		var embedded_chunks [][]float32
		for _, chunk := range chunks {
			chunk_embedding, err := s.AIService.GetGeminiEmbedding(chunk)
			fmt.Println(chunk_embedding)
			if err != nil {
				fmt.Println(err)
				lastResponseToTelex = "sorry, couldn't process your upload"
				go txc.GenerateResponseToQuery(ctx, lastResponseToTelex, req.ChannelID)
			}
			embedded_chunks = append(embedded_chunks, chunk_embedding)
		}

		err := s.DB.ReplaceEmbeddingContextTxn(ctx, chunks, embedded_chunks, req.ChannelID)
		if err != nil {
			lastResponseToTelex = "Error uploading"
			go txc.GenerateResponseToQuery(ctx, lastResponseToTelex, req.ChannelID)
		} else {
			lastResponseToTelex = "Content Uploaded you can use /help to send queries"
			go txc.GenerateResponseToQuery(ctx, lastResponseToTelex, req.ChannelID)
		}
	case "use":
		go txc.GenerateResponseToQuery(ctx,
			`Use:\n/upload followed by your text to upload context
/help followed by a query to interact with AI
/change-context followed by text/document to change the AI's Knowledge base`,
			req.ChannelID)
	default:
		ctx.Status(202)
		return

	}
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
