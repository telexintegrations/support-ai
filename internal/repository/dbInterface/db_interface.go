package dbinterface

import (
	"go.mongodb.org/mongo-driver/mongo"
)

type MongoManager interface {
	DB() *mongo.Client
}
