package mongo

import (
	"context"
	"fmt"
	"log"

	"go.mongodb.org/mongo-driver/mongo"
	"gopkg.in/mgo.v2/bson"
)

type Vector []float32

type ContentEmbeddings struct {
	Content   string   `bson:"content"`
	Embedding []float32 `bson:"embedding"`
	OrgId     string   `bson:"org_id"`
}

func (m *MongoDB)GetContentEmbeddings(ctx context.Context) ([]bson.M, error){
	// Select the database inside the handler
	
	cursor, err := m.DB().Database("support-ai").Collection(ContentEmbeddingsCollection).Find(ctx, bson.M{})
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	var response []bson.M
	if err := cursor.All(ctx, &response); err != nil {
		return nil, err
	}

	return response, nil

}

// Expects a slice of interface containing both the content and embeddings data to be inserted into the collection
func (m *MongoDB) InsertIntoEmbeddingCollection(ctx context.Context, content []string, embeddings [][]float32, orgId string) error {
	if orgId == "" {
		orgId = "018f6b36-bcc2-7d5a-b3c1-afe15c6d2"
	}

	dataEmbeddings := make([]interface{}, len(embeddings))

	for i, data := range content {
		dataEmbeddings[i] = ContentEmbeddings{
			Content:   data,
			Embedding: embeddings[i],
			OrgId:     orgId,
		}
	}
	_, err := m.DB().Database("support-ai").Collection(ContentEmbeddingsCollection).InsertMany(ctx, dataEmbeddings)

	if err != nil {
		fmt.Println(err)
		return err
	}
	return nil
}

type Organization struct {
	ID string `bson:"org_id"`
}

func (m *MongoDB) CreateCompanyCollection(ctx context.Context, data Organization) error {
	// TODO are we creating collections for each org?
	_, err := m.DB().Database("support-ai").Collection("organizations").InsertOne(ctx, data)

	if err != nil {
		fmt.Println(err)
		return err
	}
	return nil
}

func (m *MongoDB) SearchVectorFromContentEmbedding(ctx context.Context, queryVector []float32, limit uint32) ([]ContentEmbeddings, error) {
	// pipeline := mongo.Pipeline{
	// 	{{Key: "$vectorSearch", Value: bson.D{
	// 		{Name: "index", Value: "3072support"},
	// 		{Name: "path", Value: "embedding"},
	// 		{Name: "queryVector", Value: queryVector},
	// 		{Name: "numCandidates", Value: 100},
	// 		{Name: "limit", Value: limit},
	// 	}}},
	// }
	fmt.Println("Vector search starting")

	pipeline2 := mongo.Pipeline{
		{
			{Key: "$vectorSearch", Value: bson.M{
				"index": "supportive_index",
				"queryVector": queryVector,
				"path": "embedding",
				"numCandidates": 100,
				"limit": 10,
			}},
		},
	}
	cursor, err := m.DB().Database("support-ai").Collection(ContentEmbeddingsCollection).Aggregate(ctx, pipeline2)
	fmt.Println("Vector search dataset returned")
	if err != nil {
		log.Println("Vector search aggregation error:", err)
		return nil, err
	}
	defer cursor.Close(ctx)

	var results []ContentEmbeddings
	fmt.Println("Decoding Vector search results")
	if err := cursor.All(ctx, &results); err != nil {
		log.Println("Error decoding search results:", err)
		return nil, err
	}
	fmt.Println("Decoding results completed")

	return results, nil
}

func (m *MongoDB) DeleteOrganization(ctx context.Context, orgID string) error {

	_, err := m.DB().Database("support-ai").Collection("organizations").DeleteOne(ctx, bson.M{"org_id": orgID})
	if err != nil {
		fmt.Println(err)
		return err
	}
	return nil
}
