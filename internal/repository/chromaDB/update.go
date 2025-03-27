package chromadb

import (
	"context"
	"errors"
	"fmt"
	"log"
	"strings"

	"github.com/amikos-tech/chroma-go/types"
)

var ErrNoOrgId = errors.New("no organisation ID provided")
var ErrNoDataInOrg = errors.New("Organisation has no data provided")

func (c *ChromaDB) ReplaceEmbeddingContext(ctx context.Context, newContent []string, newEmbeddings []*types.Embedding, orgId string) error {
	if orgId == "" {
		return ErrNoOrgId
	}

	collectionAsOrgId := fmt.Sprintf("%s-%s", collectionPrefix, orgId)
	col, err := c.ChromaDB().GetCollection(ctx, collectionAsOrgId, nil)
	if err != nil {
		fmt.Println(err)
		if strings.Contains(err.Error(), fmt.Sprintf("%s does not exist", collectionAsOrgId)) {
			return ErrNoDataInOrg
		}
		return err
	}

	fmt.Println(col)
	err = c.DeleteEntireOrganisationContext(ctx, orgId)
	if err != nil {
		log.Println("Error deleting old embeddings:", err)
		return err
	}

	err = c.InsertIntoChromaEmbeddingCollection(ctx, ChromaContentEmbeddings{
		OrgId:         orgId,
		ContentChunks: newContent,
	})
	if err != nil {
		log.Println("Error inserting new embeddings:", err)
		return err
	}

	log.Println("Successfully replaced embedding context for org:", orgId)
	return nil
}
