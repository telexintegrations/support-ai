package mongo

import (
	"context"
	"fmt"
	"time"

	"github.com/telexintegrations/support-ai/api"
	"github.com/telexintegrations/support-ai/internal/repository"
	dbinterface "github.com/telexintegrations/support-ai/internal/repository/dbInterface"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

func ConnectToMongo(uri api.EnvConfig) (*mongo.Client, error) {
	// Connect to MongoDB
	//

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)

	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(uri.MONGODB_URI))

	defer func() {
		if err = client.Disconnect(ctx); err != nil {
			fmt.Println(err)
			panic(err)
		}
	}()

	err = client.Ping(ctx, readpref.Primary())
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	repository.DB.Mongo = NewMongoManager(client)
	// repository.DB.
	// instantiate db
	fmt.Println("Database Connection successful")
	return client, nil
}

func NewMongoManager(client *mongo.Client) dbinterface.MongoManager {
	return &MongoDB{
		MongoClient: client,
	}
}
