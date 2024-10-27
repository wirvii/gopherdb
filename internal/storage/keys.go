package storage

import (
	"fmt"

	. "github.com/wirvii/gopherdb/internal/consts"
)

// databaseKey returns the key for a database.
func databaseKey(dbName string) []byte {
	return []byte(fmt.Sprintf("%s/%s", PrefixDatabase, dbName))
}

// collectionKey returns the key for a collection.
func collectionKey(dbName, collectionName string) []byte {
	return []byte(fmt.Sprintf("%s/%s/%s", PrefixCollection, dbName, collectionName))
}

// documentKey returns the key for a document.
func documentKey(dbName, collectionName, docID string) []byte {
	return []byte(fmt.Sprintf("%s/%s/%s/%s", PrefixDocument, dbName, collectionName, docID))
}

// indexKey returns the key for an index.
func indexKey(dbName, collectionName, indexName, indexedValue, docID string) []byte {
	return []byte(fmt.Sprintf("%s/%s/%s/%s/%s/%s", PrefixIndex, dbName, collectionName, indexName, indexedValue, docID))
}
