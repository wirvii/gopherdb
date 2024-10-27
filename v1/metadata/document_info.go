package metadata

type DocumentInfo struct {
	ID   string `json:"id" bson:"_id,omitempty"`
	Data []byte `json:"data" bson:"data"`
}
