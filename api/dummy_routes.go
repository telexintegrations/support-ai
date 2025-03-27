package api

import (
	"context"
	"fmt"

	"github.com/gin-gonic/gin"
	chromadb "github.com/telexintegrations/support-ai/internal/repository/chromaDB"
)

type Receiver struct {
	Text string
}

func (s *Server) DummyRoute(ctx *gin.Context) {
	var r Receiver

	fmt.Println("For Cupiddddd.........")
	err := ctx.BindJSON(&r)

	if err != nil {
		fmt.Println(err)
		ctx.JSON(400, gin.H{"error": "Invalid request payload"})
		return
	}

	va := chromadb.ChromaContentEmbeddings{OrgId: "546", ContentChunks: []string{r.Text}}
	fmt.Printf("%+v", va)
	cotx := context.TODO()

	err = s.CDB.InsertIntoChromaEmbeddingCollection(cotx, va)

	if err != nil {
		fmt.Println(err)
		ctx.JSON(400, gin.H{"error": "Failed to insert into db"})
		return
	}
	ctx.JSON(200, gin.H{
		"message": "Hello World",
	})
}

func (s *Server) SearchDummyRoutes(ctx *gin.Context) {
	var Search struct {
		Search string `json:"search"`
	}
	err := ctx.BindJSON(&Search)

	if err != nil {
		fmt.Println(err)
		ctx.JSON(400, gin.H{"error": "Invalid request payload"})
		return
	}

	cotx := context.TODO()

	sr, err := s.CDB.SearchVectorFromContentEmbedding(cotx, Search.Search, 1, "546")

	if err != nil {
		fmt.Println(err)
		ctx.JSON(400, gin.H{"error": "Failed to insert into db"})
		return
	}
	ctx.JSON(200, gin.H{
		"message": sr,
	})
}
