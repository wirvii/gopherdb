package storage

import (
	"errors"
	"os"

	"github.com/dgraph-io/badger/v4"
	. "github.com/wirvii/gopherdb/internal/consts"
)

// badgerEngine is an implementation of the StorageEngine interface that uses BadgerDB as the underlying storage engine.
type badgerEngine struct {
	db *badger.DB
}

// newBadgerEngine creates a new BadgerDB-based storage engine at the given path.
func newBadgerEngine(path string) (*badgerEngine, error) {
	if err := os.MkdirAll(path, P0755); err != nil {
		return nil, err
	}

	opts := badger.DefaultOptions(path)
	opts.Logger = nil

	db, err := badger.Open(opts)
	if err != nil {
		return nil, err
	}

	return &badgerEngine{db: db}, nil
}

// Put inserts a key-value pair into the storage engine.
func (e *badgerEngine) Put(key string, value []byte) error {
	return e.db.Update(func(txn *badger.Txn) error {
		return txn.Set([]byte(key), value)
	})
}

// Get retrieves a value from the storage engine for the given key.
func (e *badgerEngine) Get(key string) ([]byte, error) {
	var value []byte

	err := e.db.View(func(txn *badger.Txn) error {
		item, err := txn.Get([]byte(key))
		if err != nil {
			if errors.Is(err, badger.ErrKeyNotFound) {
				return ErrKeyNotFound
			}

			return err
		}

		value, err = item.ValueCopy(nil)
		if err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		return nil, err
	}

	return value, nil
}

// Delete removes a key-value pair from the storage engine.
func (e *badgerEngine) Delete(key string) error {
	return e.db.Update(func(txn *badger.Txn) error {
		return txn.Delete([]byte(key))
	})
}

// Scan scans the storage engine for all keys that match the given prefix.
func (e *badgerEngine) Scan(prefix string) (map[string][]byte, error) {
	results := make(map[string][]byte)

	err := e.db.View(func(txn *badger.Txn) error {
		it := txn.NewIterator(badger.DefaultIteratorOptions)
		defer it.Close()

		for it.Seek([]byte(prefix)); it.ValidForPrefix([]byte(prefix)); it.Next() {
			item := it.Item()
			k := item.Key()
			err := item.Value(func(v []byte) error {
				results[string(k)] = v

				return nil
			})

			if err != nil {
				return err
			}
		}

		return nil
	})
	if err != nil {
		return nil, err
	}

	return results, nil
}

// ScanKeys scans the storage engine for all keys that match the given prefix.
func (e *badgerEngine) ScanKeys(prefix string) ([]string, error) {
	results := make([]string, 0)

	err := e.db.View(func(txn *badger.Txn) error {
		opts := badger.DefaultIteratorOptions
		opts.PrefetchValues = false

		it := txn.NewIterator(opts)
		defer it.Close()

		for it.Seek([]byte(prefix)); it.ValidForPrefix([]byte(prefix)); it.Next() {
			item := it.Item()
			k := item.Key()
			results = append(results, string(k))
		}

		return nil
	})
	if err != nil {
		return nil, err
	}

	return results, nil
}

// Close closes the storage engine.
func (e *badgerEngine) Close() error {
	return e.db.Close()
}
