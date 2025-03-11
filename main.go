package main

import (
	"github.com/telexintegrations/support-ai/api"
)

func main() {
	config, _ := api.LoadEnvConfig()

	server := api.NewServer(&config)
	server.StartServer(":8080")
}
