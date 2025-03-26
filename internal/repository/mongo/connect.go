package mongo

import (
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

const (
	ContentEmbeddingsCollection = "content-embeddings"
	NumDimensions               = 3062
)

func ConnectToMongo(uri, db_name string) (*mongo.Client, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)

	defer cancel()
	clientOptions := options.Client().ApplyURI(uri)
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		fmt.Println("failed to connect to mongoDB")
		return nil, err
	}
	err = client.Ping(ctx, readpref.Primary())
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	fmt.Println("Database Connection successful")
	return client, nil
}
