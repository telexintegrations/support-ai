package chromadb

import (
	chromago "github.com/amikos-tech/chroma-go"
)

type ChromaDB struct {
	Chroma *chromago.Client
}

func NewChromeDB(client *chromago.Client) *ChromaDB {
	return &ChromaDB{Chroma: client}
}

func (c *ChromaDB) ChromaDB() *chromago.Client {
	return c.Chroma
}
