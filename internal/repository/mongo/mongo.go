package mongo

import (
	"context"
	"errors"
	"fmt"
	"log"

	"github.com/telexintegrations/support-ai/internal/repository"
	"github.com/telexintegrations/support-ai/internal/repository/dbmodel"
	"go.mongodb.org/mongo-driver/mongo"
	"gopkg.in/mgo.v2/bson"
)

var (
	ErrNoOrgId                        = errors.New("no log id passed")
	ErrFailedToReplaceExistingContext = errors.New("failed to replace existing context")
)

type Vector []float32

var _ repository.VectorRepo = (*MongoDB)(nil)

func (m *MongoDB) GetContentEmbeddings(ctx context.Context) ([]dbmodel.ContentEmbeddings, error) {
	// Select the database inside the handler

	cursor, err := m.DB().Database("support-ai").Collection(ContentEmbeddingsCollection).Find(ctx, bson.M{})
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	var results []dbmodel.ContentEmbeddings
	if err := cursor.All(ctx, &results); err != nil {
		return nil, err
	}

	return results, nil

}

// Expects a slice of interface containing both the content and embeddings data to be inserted into the collection
func (m *MongoDB) InsertIntoEmbeddingCollection(ctx context.Context, content []string, embeddings [][]float32, orgData dbmodel.OrgMetaData) error {
	if orgData.ID == "" {
		return ErrNoOrgId
	}

	dataEmbeddings := make([]interface{}, len(embeddings))

	for i, data := range content {
		dataEmbeddings[i] = dbmodel.ContentEmbeddings{
			Content:   data,
			Embedding: embeddings[i],
			OrgId:     orgData.ID,
		}
	}
	_, err := m.DB().Database("support-ai").Collection(ContentEmbeddingsCollection).InsertMany(ctx, dataEmbeddings)

	if err != nil {
		fmt.Println(err)
		return err
	}
	return nil
}

func (m *MongoDB) CreateCompanyCollection(ctx context.Context, data dbmodel.OrgMetaData) error {
	// TODO are we creating collections for each org?
	_, err := m.DB().Database("support-ai").Collection("organizations").InsertOne(ctx, data.ID)

	if err != nil {
		fmt.Println(err)
		return err
	}
	return nil
}

func (m *MongoDB) SearchVectorFromContentEmbedding(ctx context.Context, queryVector []float32, orgData dbmodel.OrgMetaData, limit uint32) ([]dbmodel.ContentEmbeddings, error) {
	fmt.Println("Vector search starting")

	pipeline2 := mongo.Pipeline{
		{{Key: "$vectorSearch", Value: bson.M{
			"index":         "supportive_index",
			"queryVector":   queryVector,
			"path":          "embedding",
			"numCandidates": 100,
			"limit":         limit,
		}}},
		{{Key: "$match", Value: bson.M{"org_id": orgData.ID}}},
	}

	cursor, err := m.DB().Database("support-ai").Collection(ContentEmbeddingsCollection).Aggregate(ctx, pipeline2)
	fmt.Println("Vector search dataset returned")
	if err != nil {
		log.Println("Vector search aggregation error:", err)
		return nil, err
	}
	defer cursor.Close(ctx)

	var results []dbmodel.ContentEmbeddings
	fmt.Println("Decoding Vector search results")
	if err := cursor.All(ctx, &results); err != nil {
		log.Println("Error decoding search results:", err)
		return nil, err
	}
	fmt.Println("Decoding results completed")

	return results, nil
}

func (m *MongoDB) deleteEntireOrganisationContext(ctx context.Context, orgID string) error {

	res, err := m.DB().Database("support-ai").Collection("organizations").DeleteMany(ctx, bson.M{"org_id": orgID})
	if err != nil {
		fmt.Println(err)
		return err
	}
	fmt.Printf("All context for organisation with id: %s, successfully delete. And %d, vector fields where deleted\n", orgID, res.DeletedCount)
	return nil
}

func (m *MongoDB) ReplaceEmbeddingContextTxn(ctx context.Context, newContent []string, newEmbeddings [][]float32, orgData dbmodel.OrgMetaData) error {
	if orgData.ID == "" {
		return ErrNoOrgId
	}

	orgId := orgData.ID
	client := m.DB()
	session, err := client.StartSession()
	if err != nil {
		log.Println("Error starting transaction session:", err)
		return err
	}
	defer session.EndSession(ctx)

	txnCallback := func(sessCtx mongo.SessionContext) (interface{}, error) {

		err := m.deleteEntireOrganisationContext(sessCtx, orgId)
		if err != nil {
			log.Println("Error deleting old embeddings:", err)
			return nil, err
		}

		err = m.InsertIntoEmbeddingCollection(sessCtx, newContent, newEmbeddings, orgData)
		if err != nil {
			log.Println("Error inserting new embeddings:", err)
			return nil, err
		}

		return nil, nil
	}

	_, err = session.WithTransaction(ctx, txnCallback)
	if err != nil {
		log.Println("Transaction failed:", err)
		return ErrFailedToReplaceExistingContext
	}

	log.Println("Successfully replaced embedding context for org:", orgId)
	return nil
}
