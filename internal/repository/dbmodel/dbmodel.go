package dbmodel

type ContentEmbeddings struct {
	Content   string    `bson:"content"`
	Embedding []float32 `bson:"embedding"`
	OrgId     string    `bson:"org_id"`
}

type OrgMetaData struct {
	ID string `bson:"org_id"`
}
