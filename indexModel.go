package gopherdb

// IndexOptions represents the options for an index.
type IndexOptions struct {
	Name   string `json:"name"`
	Unique bool   `json:"unique"`
}

// IndexField represents a field in an index.
type IndexField struct {
	Name  string `json:"name"`
	Order int    `json:"order"`
}

// IndexModel represents a model for an index.
type IndexModel struct {
	Fields  []IndexField `json:"fields"`
	Options IndexOptions `json:"options"`
}
