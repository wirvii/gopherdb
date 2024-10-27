package storage

import (
	"context"
	"fmt"
	"strings"

	"github.com/dgraph-io/badger/v4"
	. "github.com/wirvii/gopherdb/internal/consts"
	"github.com/wirvii/gopherdb/internal/hash"
	"github.com/wirvii/gopherdb/internal/query"
	"github.com/wirvii/gopherdb/v1/errors"
	"github.com/wirvii/gopherdb/v1/metadata"
)

// badgerStorage is a storage implementation using Badger.
type badgerStorage struct {
	db *badger.DB
}

// newBadgerStorage creates a new Badger storage with the provided options.
func newBadgerStorage(sOpts *storageOptions) (*badgerStorage, error) {
	opts := badger.DefaultOptions(sOpts.filePath).
		WithIndexCacheSize(sOpts.cacheSize)

	sOpts.encryptionKey = strings.TrimSpace(sOpts.encryptionKey)
	if sOpts.encryptionKey != "" {
		opts = opts.WithEncryptionKey([]byte(sOpts.encryptionKey))
	}

	db, err := badger.Open(opts)
	if err != nil {
		return nil, err
	}

	return &badgerStorage{
		db: db,
	}, nil
}

// GetDatabase is a function that returns a database.
func (b *badgerStorage) GetDatabase(dbName string) (*metadata.DatabaseInfo, error) {
	dbName = strings.TrimSpace(dbName)
	if dbName == "" {
		return nil, errors.ErrInvalidDatabaseName
	}

	var dbInfo metadata.DatabaseInfo

	err := b.db.View(func(txn *badger.Txn) error {
		item, err := txn.Get(databaseKey(dbName))
		if err != nil {
			return err
		}

		return item.Value(func(val []byte) error {
			return dbInfo.Unmarshal(val)
		})
	})

	if err != nil {
		if err == badger.ErrKeyNotFound {
			return nil, errors.ErrDatabaseNotFound
		}

		return nil, err
	}

	return &dbInfo, nil
}

// CreateDatabase is a function that creates a new database.
func (b *badgerStorage) CreateDatabase(dbName string) error {
	dbName = strings.TrimSpace(dbName)
	if dbName == "" {
		return errors.ErrInvalidDatabaseName
	}

	db, err := b.GetDatabase(dbName)
	if err != nil && !errors.IsDatabaseNotFound(err) {
		return err
	}

	if err == nil && db != nil {
		return errors.ErrDatabaseAlreadyExists
	}

	return b.db.Update(func(txn *badger.Txn) error {
		dbInfo := &metadata.DatabaseInfo{
			ID:          hash.SHA256(dbName),
			Name:        dbName,
			Collections: make([]string, 0),
		}

		bs, err := dbInfo.Marshal()
		if err != nil {
			return err
		}

		return txn.Set(databaseKey(dbName), bs)
	})
}

func (b *badgerStorage) DropDatabase(dbName string) error {
	panic("unimplemented")
}

// ListDatabases is a function that returns a list of databases.
func (b *badgerStorage) ListDatabases(ctx context.Context) ([]*metadata.DatabaseInfo, error) {
	databases := make([]*metadata.DatabaseInfo, 0)

	err := b.db.View(func(txn *badger.Txn) error {
		opts := badger.DefaultIteratorOptions
		opts.PrefetchSize = 10

		it := txn.NewIterator(opts)
		defer it.Close()

		prefix := []byte(fmt.Sprintf("%s/", PrefixDatabase))
		for it.Seek(prefix); it.ValidForPrefix(prefix); it.Next() {
			item := it.Item()

			var dbInfo metadata.DatabaseInfo

			err := item.Value(func(val []byte) error {
				return dbInfo.Unmarshal(val)
			})

			if err != nil {
				return err
			}

			databases = append(databases, &dbInfo)
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return databases, nil
}

// GetCollection is a function that returns a collection.
func (b *badgerStorage) GetCollection(dbName, collectionName string) (*metadata.CollectionInfo, error) {
	dbName = strings.TrimSpace(dbName)
	if dbName == "" {
		return nil, errors.ErrInvalidDatabaseName
	}

	collectionName = strings.TrimSpace(collectionName)
	if collectionName == "" {
		return nil, errors.ErrInvalidCollectionName
	}

	var collectionInfo metadata.CollectionInfo

	err := b.db.View(func(txn *badger.Txn) error {
		item, err := txn.Get(collectionKey(dbName, collectionName))
		if err != nil {
			return err
		}

		return item.Value(func(val []byte) error {
			return collectionInfo.Unmarshal(val)
		})
	})

	if err != nil {
		if err == badger.ErrKeyNotFound {
			return nil, errors.ErrCollectionNotFound
		}

		return nil, err
	}

	return &collectionInfo, nil
}

// CreateCollection is a function that creates a new collection.
func (b *badgerStorage) CreateCollection(dbName, collectionName string) error {
	dbName = strings.TrimSpace(dbName)
	if dbName == "" {
		return errors.ErrInvalidDatabaseName
	}

	collectionName = strings.TrimSpace(collectionName)
	if collectionName == "" {
		return errors.ErrInvalidCollectionName
	}

	collection, err := b.GetCollection(dbName, collectionName)
	if err != nil && !errors.IsCollectionNotFound(err) {
		return err
	}

	if err == nil && collection != nil {
		return errors.ErrCollectionAlreadyExists
	}

	return b.db.Update(func(txn *badger.Txn) error {
		collectionInfo := &metadata.CollectionInfo{
			DbName:    dbName,
			Name:      dbName,
			Documents: 0,
		}

		bs, err := collectionInfo.Marshal()
		if err != nil {
			return err
		}

		return txn.Set(collectionKey(dbName, collectionName), bs)
	})
}

func (b *badgerStorage) DropCollection(dbName, collectionName string) error {
	panic("unimplemented")
}

// ListCollections is a function that returns a list of collections.
func (b *badgerStorage) ListCollections(dbName string) ([]*metadata.CollectionInfo, error) {
	dbName = strings.TrimSpace(dbName)
	if dbName == "" {
		return nil, errors.ErrInvalidDatabaseName
	}

	collections := make([]*metadata.CollectionInfo, 0)

	err := b.db.View(func(txn *badger.Txn) error {
		opts := badger.DefaultIteratorOptions
		opts.PrefetchSize = 10

		it := txn.NewIterator(opts)
		defer it.Close()

		prefix := []byte(fmt.Sprintf("%s/%s/", PrefixCollection, dbName))
		for it.Seek(prefix); it.ValidForPrefix(prefix); it.Next() {
			item := it.Item()

			var collectionInfo metadata.CollectionInfo

			err := item.Value(func(val []byte) error {
				return collectionInfo.Unmarshal(val)
			})

			if err != nil {
				return err
			}

			collections = append(collections, &collectionInfo)
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return collections, nil
}

func (b *badgerStorage) InsertDocument(dbName, collectionName, id string, document *metadata.DocumentInfo) error {
	panic("unimplemented")
}

func (b *badgerStorage) GetDocument(dbName, collectionName, id string) ([]*metadata.DocumentInfo, error) {
	panic("unimplemented")
}

func (b *badgerStorage) UpdateDocument(dbName, collectionName, id string, document []*metadata.DocumentInfo) error {
	panic("unimplemented")
}

func (b *badgerStorage) DeleteDocument(dbName, collectionName, id string) error {
	panic("unimplemented")
}

func (b *badgerStorage) FindDocuments(
	dbName,
	collectionName string,
	parsedQuery *query.ParsedQuery,
) ([]*metadata.DocumentInfo, error) {
	panic("unimplemented")
}
