package gopherdb

type DatabaseMetadata struct {
	Name string `json:"name"`
}

type CollectionMetadata struct {
	Name          string       `json:"name"`
	Indexes       []IndexModel `json:"indexes"`
	DocumentCount int64        `json:"document_count"`
}
