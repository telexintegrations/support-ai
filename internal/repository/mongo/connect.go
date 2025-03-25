package mongo

import (
	"context"
	"fmt"
	"slices"
	"time"

	"github.com/telexintegrations/support-ai/internal/repository"
	dbinterface "github.com/telexintegrations/support-ai/internal/repository/dbInterface"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"gopkg.in/mgo.v2/bson"
)

const (
	ContentEmbeddingsCollection = "content-embeddings"
	NumDimensions               = 3062
)

func ConnectToMongo(uri, db_name string) (*mongo.Client, error) {
	fmt.Println("Connecting to MongoDB...")
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

	dbName := db_name
	if !checkDBExists(client, dbName) {
		collection := client.Database(dbName).Collection(ContentEmbeddingsCollection)
		_, err := collection.InsertOne(context.Background(), bson.M{"name": "init"})
		if err != nil {
			fmt.Println("failed to connect to mongoDB")
			return nil, err
		}
	}

	repository.DB.Mongo = NewMongoManager(client)
	fmt.Println("Database Connection successful")
	return client, nil
}

func checkDBExists(client *mongo.Client, dbName string) bool {
	databases, err := client.ListDatabaseNames(context.Background(), bson.M{})
	if err != nil {
		return false
	}
	return slices.Contains(databases, dbName)
}

func NewMongoManager(client *mongo.Client) dbinterface.MongoManager {
	return &MongoDB{
		MongoClient: client,
	}
}
