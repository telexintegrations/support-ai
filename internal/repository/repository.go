package repository

import (
	"context"

	"github.com/telexintegrations/support-ai/internal/repository/dbmodel"
)

type VectorRepo interface {
	GetContentEmbeddings(ctx context.Context) ([]dbmodel.ContentEmbeddings, error)
	InsertIntoEmbeddingCollection(ctx context.Context, content []string, embeddings [][]float32, orgData dbmodel.OrgMetaData) error
	// Not yet in use in this application
	CreateCompanyCollection(ctx context.Context, orgData dbmodel.OrgMetaData) error
	SearchVectorFromContentEmbedding(ctx context.Context, queryVector []float32, orgData dbmodel.OrgMetaData, limit uint32) ([]dbmodel.ContentEmbeddings, error)
	ReplaceEmbeddingContextTxn(ctx context.Context, newContent []string, newEmbeddings [][]float32, orgData dbmodel.OrgMetaData) error
}
