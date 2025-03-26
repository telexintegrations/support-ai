package dbinterface

import (
	chromago "github.com/amikos-tech/chroma-go"
	"go.mongodb.org/mongo-driver/mongo"
)

type MongoManager interface {
	DB() *mongo.Client
}

type ChromaManager interface {
	ChromaDB() *chromago.Client
}
