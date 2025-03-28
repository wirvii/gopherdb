package storage

// Storage is a generic interface for a key-value store.
type Storage interface {
	// Get returns the value for a given key.
	Get(key string) ([]byte, error)
	// Put sets the value for a given key.
	Put(key string, value []byte) error
	// Delete deletes the value for a given key.
	Delete(key string) error
	// Scan scans the database for all keys that match the prefix.
	Scan(prefix string) (map[string][]byte, error)
	// ScanKeys scans the database for all keys that match the prefix.
	ScanKeys(prefix string) ([]string, error)
	// Close closes the storage engine.
	Close() error
}

// NewStorage creates a new storage engine.
func NewStorage(path string) (Storage, error) {
	return newBadgerEngine(path)
}
