package chromadb

import (
	"fmt"

	chromago "github.com/amikos-tech/chroma-go"
	"github.com/telexintegrations/support-ai/internal/repository"
)

func ConnectionToChroma(uri string) (*chromago.Client, error) {
	client, err := chromago.NewClient(uri)

	if err != nil {
		fmt.Println("failed to connect to DB", err)
		return nil, err
	}
	repository.DB.ChromaDB = NewChromeDB(client)
	return client, nil
}
