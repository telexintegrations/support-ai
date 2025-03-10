package mongo

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

func ConnectToMongo(host string, port string, username string, password string) (*MongoDB, error) {
	// Connect to MongoDB
	//

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)

	defer cancel()
	_, err := strconv.Atoi(port)
	if err != nil {
		fmt.Println(port, "port must be an integer value")
		panic(err)
	}

	client, err := mongo.NewClient(options.Client().ApplyURI("mongodb://" + username + ":" + password + "@" + host + ":" + port))

	defer func() {
		if err = client.Disconnect(ctx); err != nil {
			panic(err)
		}
	}()

	err = client.Ping(ctx, readpref.Primary())

	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	return nil, nil
}
