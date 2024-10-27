package storage

import "strings"

// Option is a function that applies an option to a storageOptions.
type Option interface {
	apply(s *storageOptions)
}

// optionFunc is a type that implements the Option interface.
type optionFunc func(*storageOptions)

// apply calls the optionFunc with the provided storageOptions.
func (f optionFunc) apply(c *storageOptions) {
	f(c)
}

// WithFilePath sets the path of the storage.
func WithFilePath(path string) Option {
	return optionFunc(func(c *storageOptions) {
		c.filePath = strings.TrimSpace(path)
	})
}

// WithCacheSize sets the cache size of the storage.
func WithCacheSize(cacheSize int64) Option {
	return optionFunc(func(c *storageOptions) {
		c.cacheSize = cacheSize
	})
}

// WithEncryptionKey sets the encryption key of the storage.
func WithEncryptionKey(encryptionKey string) Option {
	return optionFunc(func(c *storageOptions) {
		c.encryptionKey = strings.TrimSpace(encryptionKey)
	})
}
