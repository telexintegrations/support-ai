package mongo

import (
	"context"
	"fmt"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"gopkg.in/mgo.v2/bson"
)

type Vector []float32

type ContentEmbeddings struct {
	Content   string   `bson:"content"`
	Embedding []Vector `bson:"embedding"`
}

// Expects a slice of interface containing both the content and embeddings data to be inserted into the collection
func (m *MongoDB) InsertIntoEmbeddingCollection(ctx context.Context, content []string, embeddings [][]Vector) error {
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

type Organization struct {
	ID string `bson:"org_id"`
}

func (m *MongoDB) CreateCompanyCollection(data Organization) error {

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	_, err := m.DB().Database("support-ai").Collection("organizations").InsertOne(ctx, data)

	if err != nil {
		fmt.Println(err)
		return err
	}
	return nil
}

func (m *MongoDB) SearchVectorFromContentEmbedding(ctx context.Context, queryVector Vector, limit uint32) ([]bson.M, error) {
	pipeline := mongo.Pipeline{
		{{Key: "$vectorSearch", Value: bson.D{
			{Name: "index", Value: "vector_search_index"},
			{Name: "path", Value: "embedding"},
			{Name: "queryVector", Value: queryVector},
			{Name: "numCandidates", Value: 100},
			{Name: "limit", Value: limit},
			{Name: "similarity", Value: "dotProduct"},
		}}},
	}
	cursor, err := m.DB().Database("support-ai").Aggregate(ctx, pipeline)
	if err != nil {
		log.Println("Vector search aggregation error:", err)
		return nil, err
	}
	defer cursor.Close(ctx)

	var results []bson.M
	if err := cursor.All(ctx, &results); err != nil {
		log.Println("Error decoding search results:", err)
		return nil, err
	}

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
