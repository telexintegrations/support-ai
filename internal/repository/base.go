package repository

import (
	dbinterface "github.com/telexintegrations/support-ai/internal/repository/dbInterface"
	"go.mongodb.org/mongo-driver/mongo"
)

type DatabaseManager struct {
	Mongo  dbinterface.MongoManager
	Chroma dbinterface.ChromaManager
}

var DB = DatabaseManager{}

func ConnectDB(client *mongo.Client) *DatabaseManager {
	return &DatabaseManager{}
}
