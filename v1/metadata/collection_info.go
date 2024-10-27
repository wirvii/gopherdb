package metadata

import "go.mongodb.org/mongo-driver/bson"

type CollectionInfo struct {
	DbName    string `json:"db_name" bson:"db_name"`
	Name      string `json:"name" bson:"name"`
	Documents int64
}

// Marshal is a function that marshals the collection info.
func (c *CollectionInfo) Marshal() ([]byte, error) {
	return bson.Marshal(c)
}

// Unmarshal is a function that unmarshals the collection info.
func (c *CollectionInfo) Unmarshal(data []byte) error {
	return bson.Unmarshal(data, c)
}
