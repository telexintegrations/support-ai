package chromadb

import (
	"context"
	"errors"
	"fmt"
	"log"

	"github.com/amikos-tech/chroma-go/types"
)

var ErrNoOrgId = errors.New("no organisation ID provided")

func (c *ChromaDB) ReplaceEmbeddingContext(ctx context.Context, newContent []string, newEmbeddings []*types.Embedding, orgId string) error {
	if orgId == "" {
		return ErrNoOrgId
	}

	ef, err := functionEmbeddings()
	if err != nil {
		return nil
	}
	col, err := c.ChromaDB().GetCollection(ctx, ContentEmbeddingsCollection, ef)
	if err != nil {
		log.Println("Failed to get ChromaDB collection:", err)
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
