package mongo

import (
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

func NewDBService(client *mongo.Client) *MongoDB {
	return &MongoDB{
		MongoClient: client,
	}
}
