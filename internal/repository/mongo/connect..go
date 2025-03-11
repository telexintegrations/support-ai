package mongo

import (
	"context"
	"fmt"
	"log"
	"slices"
	"time"

	"github.com/telexintegrations/support-ai/api"
	"github.com/telexintegrations/support-ai/internal/repository"
	dbinterface "github.com/telexintegrations/support-ai/internal/repository/dbInterface"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"gopkg.in/mgo.v2/bson"
)

const (
	ContentEmbeddingsCollection = "content-embeddings"
)

func ConnectToMongo(uri api.EnvConfig) (*mongo.Client, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)

	defer cancel()
	clientOptions := options.Client().ApplyURI(uri.MONGODB_DEV_URI)

	clientOptions.SetAuth(options.Credential{
		Username: uri.MONGO_USERNAME,
		Password: uri.MONGO_PASSWORD,
	})
	client, err := mongo.Connect(ctx, clientOptions)

	defer func() {
		client.Disconnect(ctx)
	}()

	err = client.Ping(ctx, readpref.Primary())
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	// checks if a database exists
	dbName := uri.MONGODATABASE_NAME
	if !checkDBExists(client, dbName) {
		fmt.Println("Database does not exist")
		collection := client.Database(dbName).Collection(ContentEmbeddingsCollection) // create a database and a dummy collection
		// create vector index in the collection

		err = CreateVectorEmbeddingIndexes(collection, ctx)
		if err != nil {
			fmt.Println(err)
			return nil, err
		}
		_, err := collection.InsertOne(context.Background(), bson.M{"name": "init"})
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
	databases, err := client.ListDatabaseNames(context.Background(), nil)
	if err != nil {
		fmt.Println(err)
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

	indexName := "vector_search_index"
	opts := options.SearchIndexes().SetName(indexName).SetType("vectorSearch")
	indexModel := mongo.SearchIndexModel{
		Definition: vectorDefinition{
			Fields: []vectorDefinitionField{{
				Type:          "vector",
				Path:          "plot_embedding",
				NumDimensions: 1536,
				Similarity:    "dotProduct",
				Quantization:  "scalar"}},
		},
		Options: opts,
	}

	searchIndexName, err := coll.SearchIndexes().CreateOne(ctx, indexModel)
	if err != nil {
		log.Printf("failed to create the search index: %v", err)
		return err
	}
	log.Println("New search index named " + searchIndexName + " is building.")
	// Await the creation of the index.
	return nil
}
