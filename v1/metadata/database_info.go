package metadata

import "go.mongodb.org/mongo-driver/bson"

// DatabaseInfo represents the metadata of a database.
type DatabaseInfo struct {
	ID          string   `json:"id" bson:"_id,omitempty"`
	Name        string   `json:"name" bson:"name"`
	Collections []string `json:"collections" bson:"collections"`
}

// Marshal returns the BSON encoding of the database info.
func (d *DatabaseInfo) Marshal() ([]byte, error) {
	return bson.Marshal(d)
}

// Unmarshal parses the BSON-encoded data and stores the result in the value pointed to by d.
func (d *DatabaseInfo) Unmarshal(data []byte) error {
	return bson.Unmarshal(data, d)
}
