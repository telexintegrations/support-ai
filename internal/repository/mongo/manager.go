package mongo

import "go.mongodb.org/mongo-driver/mongo"

type MongoDB struct {
	mongoClient *mongo.Client
}

type MongoManager interface {
	DB() *mongo.Client
}

func NewMongoManager(client *mongo.Client) MongoManager {
	return &MongoDB{
		mongoClient: client,
	}
}

func (m *MongoDB) DB() *mongo.Client {
	return m.mongoClient
}
