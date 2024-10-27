package storage

// storageOptions contains the options for a storage implementation.
type storageOptions struct {
	filePath      string
	cacheSize     int64
	encryptionKey string
}
