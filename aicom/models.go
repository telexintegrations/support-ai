package aicom

// Vector represents the structure of our stored vectors
type Vector struct {
	ID     string    `bson:"_id"`
	Text   string    `bson:"text"`
	Vector []float64 `bson:"vector"`
}

// Request structure
type AIRequest struct {
	Query         string   `json:"query" binding:"required"`
	RetrievedDocs []string `json:"retrieved_docs" binding:"required"`
}

// Response structure
type AIResponse struct {
	Response string `json:"response"`
}
