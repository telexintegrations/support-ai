package chromadb

import (
	"context"
	"fmt"
	"strings"
)

type SearchResult struct {
	Text       []string                 `json:"text"`
	Similarity []float32                `json:"similarity"`
	Metadata   []map[string]interface{} `json:"metadata"`
}

func (c *ChromaDB) SearchVectorFromContentEmbedding(ctx context.Context, query string, topK int32, orgId string) ([]SearchResult, error) {

	if orgId == "" {
		return nil, ErrNoOrgId
	}

	collectionAsOrgId := fmt.Sprintf("%s-%s", collectionPrefix, orgId)
	col, err := c.ChromaDB().GetCollection(ctx, collectionAsOrgId, nil)

	if err != nil {
		fmt.Println(err)
		if strings.Contains(err.Error(), fmt.Sprintf("%s does not exist", collectionAsOrgId)) {
			return nil, ErrNoDataInOrg
		}
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
