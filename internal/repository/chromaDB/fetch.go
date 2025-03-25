package chromadb

import (
	"context"
	"fmt"
)

type SearchResult struct {
	Text       []string                 `json:"text"`
	Similarity []float32                `json:"similarity"`
	Metadata   []map[string]interface{} `json:"metadata"`
}

func (c *ChromaDB) SearchVectorFromContentEmbedding(ctx context.Context, query string, topK int32, orgId string) ([]SearchResult, error) {
	ef, err := functionEmbeddings()
	if err != nil {
		return nil, err
	}
	col, err := c.ChromaDB().GetCollection(ctx, ContentEmbeddingsCollection, ef)

	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	queryTexts := []string{query}

	where := map[string]interface{}{
		"OrgId": orgId,
	}
	results, err := col.Query(ctx, queryTexts, topK, where, nil, nil)

	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	var searchResults []SearchResult
	for i, document := range results.Documents {
		searchResults = append(searchResults, SearchResult{
			Text:       document,
			Similarity: results.Distances[i],
			Metadata:   results.Metadatas[i],
		})
	}

	fmt.Println(searchResults)
	return searchResults, nil
}
