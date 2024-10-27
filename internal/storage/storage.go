package storage

import (
	"context"

	. "github.com/wirvii/gopherdb/internal/consts"
	"github.com/wirvii/gopherdb/internal/query"
	"github.com/wirvii/gopherdb/v1/metadata"
)

// Storage is the interface that wraps the basic methods of a storage.
type Storage interface {
	GetDatabase(dbName string) (*metadata.DatabaseInfo, error)
	CreateDatabase(dbName string) error
	DropDatabase(dbName string) error
	ListDatabases(ctx context.Context) ([]*metadata.DatabaseInfo, error)

	GetCollection(dbName, collectionName string) (*metadata.CollectionInfo, error)
	CreateCollection(dbName, collectionName string) error
	DropCollection(dbName, collectionName string) error
	ListCollections(dbName string) ([]*metadata.CollectionInfo, error)

	InsertDocument(dbName, collectionName, id string, document *metadata.DocumentInfo) error
	GetDocument(dbName, collectionName, id string) ([]*metadata.DocumentInfo, error)
	UpdateDocument(dbName, collectionName, id string, document []*metadata.DocumentInfo) error
	DeleteDocument(dbName, collectionName, id string) error

	FindDocuments(dbName, collectionName string, parsedQuery *query.ParsedQuery) ([]*metadata.DocumentInfo, error)
}

// New creates a new storage with the provided options.
func New(opts ...Option) (Storage, error) {
	storageOpts := &storageOptions{
		filePath:      "",
		cacheSize:     OneHundredsMB,
		encryptionKey: "",
	}

	for _, opt := range opts {
		opt.apply(storageOpts)
	}

	return newBadgerStorage(storageOpts)
}
