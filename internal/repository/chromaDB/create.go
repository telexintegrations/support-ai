package chromadb

import (
	"context"
	"fmt"

	g "github.com/amikos-tech/chroma-go/pkg/embeddings/gemini"
	"github.com/amikos-tech/chroma-go/types"
	"github.com/google/uuid"
)

type ChromaContentEmbeddings struct {
	OrgId         string
	ContentChunks []string
	documentID    string
}

const ContentEmbeddingsCollection = "content_embeddings"

func generateEmbeddings(content []string) ([]*types.Embedding, error) {
	ef, err := g.NewGeminiEmbeddingFunction(g.WithEnvAPIKey(), g.WithDefaultModel("text-embedding-004"))
	if err != nil {
		fmt.Printf("Error creating Gemini embedding function: %s \n", err)
		return nil, err
	}
	resp, err := ef.EmbedDocuments(context.Background(), content)
	if err != nil {
		fmt.Printf("Error embedding documents: %s \n", err)
		return nil, err
	}
	return resp, nil
}

func functionEmbeddings() (types.EmbeddingFunction, error) {
	embedFunc, err := g.NewGeminiEmbeddingFunction(g.WithEnvAPIKey(), g.WithDefaultModel("text-embedding-004"))
	if err != nil {
		fmt.Println(err, "---------->")
		return nil, err
	}
	return embedFunc, nil
}

var collectionPrefix string = "org"

// collection should be orgId
func (c *ChromaDB) InsertIntoChromaEmbeddingCollection(ctx context.Context, chromaColl ChromaContentEmbeddings) error {

	if chromaColl.OrgId == "" {
		return ErrNoOrgId
	}
	metadata := map[string]interface{}{
		"OrgId": chromaColl.OrgId,
	}

	collectionAsOrgId := fmt.Sprintf("%s-%s", collectionPrefix, chromaColl.OrgId)
	col, err := c.ChromaDB().CreateCollection(ctx, collectionAsOrgId, metadata, true, nil, types.L2)

	if err != nil {
		fmt.Println("add data err, -------->", err)
		fmt.Println(err)
		return err
	}

	chromaColl.documentID = uuid.New().String()
	ids := make([]string, len(chromaColl.ContentChunks))
	metadatas := make([]map[string]interface{}, len(chromaColl.ContentChunks))

	for i, chunk := range chromaColl.ContentChunks {
		ids[i] = uuid.New().String()
		metadatas[i] = map[string]interface{}{
			"OrgId":      chromaColl.OrgId,
			"DocumentID": chromaColl.documentID,
			"ChunkIndex": i,
			"Text":       chunk,
		}
	}
	_, err = col.Add(ctx, nil, metadatas, chromaColl.ContentChunks, ids)

	if err != nil {
		fmt.Println("add dataa", err)
		return err
	}

	fmt.Println(col.Count(ctx))
	return nil
}
