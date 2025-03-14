package mongo

import (
	"context"
	"crypto/tls"
	"fmt"
	"log"
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
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)

	defer cancel()
	clientOptions := options.Client().ApplyURI(uri)
	clientOptions.SetTLSConfig(&tls.Config{InsecureSkipVerify: true})
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

	// checks if a database exists
	dbName := db_name
	if !checkDBExists(client, dbName) {
		collection := client.Database(dbName).Collection(ContentEmbeddingsCollection) // create a database and a dummy collection
		// create vector index in the collection
		_, err := collection.InsertOne(context.Background(), bson.M{"name": "init"})
		if err != nil {
			fmt.Println("failed to connect to mongoDB")
			return nil, err
		}
		err = CreateVectorEmbeddingIndexes(collection, ctx)
		if err != nil {
			fmt.Println(err)
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

func CreateVectorEmbeddingIndexes(coll *mongo.Collection, ctx context.Context) error {
	// TODO create vector index dynamically, depending on collection
	type vectorDefinitionField struct {
		Type          string `bson:"type"`
		Path          string `bson:"path"`
		NumDimensions int    `bson:"numDimensions"`
		Similarity    string `bson:"similarity"`
		Quantization  string `bson:"quantization"`
	}

	type vectorDefinition struct {
		Fields []vectorDefinitionField `bson:"fields"`
	}

	indexName := "3072support"
	if ok, err := checkIfIndexExists(ctx, coll, indexName); err == nil && !ok {
		opts := options.SearchIndexes().SetName(indexName).SetType("vectorSearch")
		indexModel := mongo.SearchIndexModel{
			Definition: vectorDefinition{
				Fields: []vectorDefinitionField{{
					Type:          "vector",
					Path:          "embeddings",
					NumDimensions: NumDimensions,
					Similarity:    "dotProduct"}},
			},
			Options: opts,
		}
		searchIndexName, err := coll.SearchIndexes().CreateOne(ctx, indexModel)
		if err != nil {
			log.Printf("failed to create the search index: %v", err)
			return err
		}
		log.Println("New search index named " + searchIndexName + " is building.")
	} else {
		fmt.Println(err)
	}
	return nil
}

// Function to check if a search index exists
func checkIfIndexExists(ctx context.Context, coll *mongo.Collection, indexName string) (bool, error) {
	cursor, err := coll.SearchIndexes().List(ctx, nil)
	if err != nil {
		return false, err
	}
	defer cursor.Close(ctx)

	for cursor.Next(ctx) {
		var index bson.M
		if err := cursor.Decode(&index); err != nil {
			return false, err
		}

		if index["name"] == indexName {
			return true, nil
		}
	}
	return false, nil
}
