package main

import (
	"context"
	"fmt"
	"net/http"

	"github.com/robfig/cron"
	"github.com/telexintegrations/support-ai/api"
	chromadb "github.com/telexintegrations/support-ai/internal/repository/chromaDB"
	"github.com/telexintegrations/support-ai/internal/repository/mongo"
)

func main() {

	//cron job to keep render alive
	c := cron.New()
	c.AddFunc("*/550 * * * *", func() {
		fmt.Println("Cronning")
		http.Get("https://support-ai-hsd0.onrender.com/")

	})
	c.Start()

	config, err := api.LoadEnvConfig()
	if err != nil {
		fmt.Println(err)
		panic(err)
	}
	
	dbClient, err := mongo.ConnectToMongo(config.MONGODB_DEV_URI, config.MONGODATABASE_NAME)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer func() {
		dbClient.Disconnect(context.Background())
	}()

	cdb, err := chromadb.ConnectionToChroma(config.CHROMADB_DEV_URI)

	if err != nil {
		fmt.Println(err)
		return
	}

	// defer func() {
	// 	cdb.Close()
	// }()

	server := api.NewServer(&config, dbClient, cdb)
	server.StartServer(":8080")
}
