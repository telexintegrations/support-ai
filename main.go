package main

import (
	"fmt"

	"github.com/telexintegrations/support-ai/api"
	"github.com/telexintegrations/support-ai/internal/repository/mongo"
)

func main() {
	config, err := api.LoadEnvConfig()

	if err != nil {
		fmt.Println(err)
		panic(err)
	}
	_, err = mongo.ConnectToMongo(config)

	if err != nil {
		fmt.Println(err)
		return
	}

	server := api.NewServer(&config)
	server.StartServr(":8080")
}
