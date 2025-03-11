package mongo

import (
	"context"
	"fmt"
	"time"
)

type ContentEmbeddings struct {
	Content   string
	Embedding []float32
}

// Expects a slice of interface containing both the content and embeddings data to be inserted into the collection
func (m *MongoDB) InsertIntoEmbeddingCollection(content []string, embeddings [][]float32) error {
	dataEmbeddings := make([]interface{}, len(embeddings))

	for i, data := range content {
		dataEmbeddings[i] = ContentEmbeddings{
			Content:   data,
			Embedding: embeddings[i],
		}
	}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	_, err := m.DB().Database("support-ai").Collection(ContentEmbeddingsCollection).InsertMany(ctx, dataEmbeddings)

	if err != nil {
		fmt.Println(err)
		return err
	}
	return nil
}

func (m *MongoDB) CreateCompanyCollection() {

}
