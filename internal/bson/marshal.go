package bson

import "go.mongodb.org/mongo-driver/bson"

// Marshal marshals a value to BSON.
func Marshal(v any) ([]byte, error) {
	return bson.Marshal(v)
}

// Unmarshal unmarshals a BSON value to a value.
func Unmarshal(data []byte, v any) error {
	return bson.Unmarshal(data, v)
}
