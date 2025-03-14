package mongo

import (
	"github.com/telexintegrations/support-ai/internal/repository"
	"go.mongodb.org/mongo-driver/mongo"
)

type MongoDB struct {
	MongoClient *mongo.Client
}

type MongoManager interface {
	DB() *mongo.Client
}

func (m *MongoDB) DB() *mongo.Client {
	return m.MongoClient
}

func NewDBService(client *mongo.Client) repository.VectorRepo {
	return &MongoDB{
		MongoClient: client,
	}
}
