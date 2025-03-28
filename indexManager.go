package gopherdb

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"slices"
	"strings"
	"sync"

	"github.com/wirvii/gopherdb/internal/consts"
	"github.com/wirvii/gopherdb/internal/storage"
)

// IndexManager is a manager for indexes.
type IndexManager struct {
	mu       sync.RWMutex
	dbname   string
	collname string
	storage  storage.Storage
	metadata CollectionMetadata
}

// newIndexManager creates a new IndexManager.
func newIndexManager(storage storage.Storage, dbname, collname string) *IndexManager {
	return &IndexManager{
		mu:       sync.RWMutex{},
		storage:  storage,
		dbname:   dbname,
		collname: collname,
	}
}

// buildMetadataKey builds the key for the collection metadata.
func (m *IndexManager) buildMetadataKey() string {
	return fmt.Sprintf(consts.MetadataCollectionKeyStringFormat, m.dbname, m.collname)
}

// loadMetadata loads the collection metadata from the storage.
func (m *IndexManager) loadMetadata() error {
	m.mu.RLock()
	defer m.mu.RUnlock()

	data, err := m.storage.Get(m.buildMetadataKey())
	if err != nil {
		if errors.Is(err, storage.ErrKeyNotFound) {
			m.metadata = CollectionMetadata{
				Name: fmt.Sprintf(
					consts.CollectionKeyStringFormat,
					m.dbname,
					m.collname,
				),
				Indexes: []IndexModel{
					{
						Fields: []IndexField{
							{
								Name:  consts.DocumentFieldID,
								Order: 1,
							},
						},
						Options: IndexOptions{
							Name:   "_id_",
							Unique: true,
						},
					},
				},
				DocumentCount: 0,
			}

			return nil
		}

		return err
	}

	return json.Unmarshal(data, &m.metadata)
}

// saveMetadata saves the collection metadata to the storage.
func (m *IndexManager) saveMetadata() error {
	m.mu.Lock()
	defer m.mu.Unlock()

	data, err := json.Marshal(m.metadata)
	if err != nil {
		return err
	}

	return m.storage.Put(m.buildMetadataKey(), data)
}

// getDocumentIdsByIndex gets the document ids for a given index.
func (m *IndexManager) getDocumentIdsByIndex(index IndexModel, indexFilter map[string]any) ([]string, error) {
	indexKeyPrefix, err := m.buildIndexKey(index, indexFilter, true)
	if err != nil {
		return nil, err
	}

	indexKeyPrefix = strings.TrimSuffix(
		indexKeyPrefix,
		fmt.Sprintf("/%v", consts.RemoverWildcard),
	)

	entries, err := m.storage.ScanKeys(indexKeyPrefix)
	if err != nil {
		return nil, err
	}

	docIds := make([]string, 0, len(entries))

	for _, entry := range entries {
		match, err := consts.IndexKeyPathmatcher.Match(entry)
		if err != nil {
			return nil, err
		}

		docIds = append(docIds, match["docId"])
	}

	return docIds, nil
}

// buildIndexFieldsKey builds the fields key for a given index.
func (m *IndexManager) buildIndexFieldsKey(index IndexModel) string {
	slices.SortFunc(index.Fields, func(a, b IndexField) int {
		return (a.Order - b.Order)
	})

	fields := make([]string, 0, len(index.Fields))

	for _, f := range index.Fields {
		fields = append(fields, f.Name)
	}

	joinedFields := sha256.Sum256([]byte(strings.Join(fields, "|")))

	fieldsKeyHash := hex.EncodeToString(joinedFields[:])

	key := fmt.Sprintf(
		consts.IndexKeyStringFormat,
		m.dbname,
		m.collname,
		index.Options.Name,
		fieldsKeyHash,
		consts.RemoverWildcard,
		consts.RemoverWildcard,
	)

	return strings.TrimSuffix(
		key,
		fmt.Sprintf(
			"%v/%v",
			consts.RemoverWildcard,
			consts.RemoverWildcard,
		),
	)
}

// buildIndexKey builds the index key for a document.
func (m *IndexManager) buildIndexKey(index IndexModel, doc map[string]any, isPartial bool) (string, error) {
	slices.SortFunc(index.Fields, func(a, b IndexField) int {
		return (a.Order - b.Order)
	})

	values := make([]string, 0, len(index.Fields))
	fields := make([]string, 0, len(index.Fields))

	for _, f := range index.Fields {
		val, ok := doc[f.Name]
		if !ok {
			if isPartial {
				break
			}

			return "", ErrMissingFieldForIndex
		}

		fields = append(fields, f.Name)
		values = append(values, fmt.Sprintf("%v", val))
	}

	joinedFields := sha256.Sum256([]byte(strings.Join(fields, "|")))
	joinedValues := sha256.Sum256([]byte(strings.Join(values, "|")))

	fieldsKeyHash := hex.EncodeToString(joinedFields[:])
	valuesKeyHash := hex.EncodeToString(joinedValues[:])

	docId := fmt.Sprintf("%v", doc[consts.DocumentFieldID])
	if isPartial {
		docId = consts.RemoverWildcard
	}

	return fmt.Sprintf(
		consts.IndexKeyStringFormat,
		m.dbname,
		m.collname,
		index.Options.Name,
		fieldsKeyHash,
		valuesKeyHash,
		docId,
	), nil
}

// checkUniqueness checks if the document violates the uniqueness constraint of the index.
func (m *IndexManager) checkUniqueness(doc map[string]any) error {
	for _, idx := range m.metadata.Indexes {
		if !idx.Options.Unique {
			continue
		}

		idxKey, err := m.buildIndexKey(idx, doc, false)
		if err != nil {
			return err
		}

		key := strings.TrimSuffix(idxKey, fmt.Sprintf("%v", doc[consts.DocumentFieldID]))
		entries, err := m.storage.ScanKeys(key)

		if err != nil {
			return err
		}

		if len(entries) > 0 {
			return fmt.Errorf("%w: fields %+v", ErrUniqueIndexViolation, idx.Fields)
		}
	}

	return nil
}

// buildIndexName builds the index name.
func (m *IndexManager) buildIndexName(idx IndexModel) string {
	name := ""

	for _, f := range idx.Fields {
		name += fmt.Sprintf("_%s_%d", f.Name, f.Order)
	}

	return strings.TrimPrefix(name, "_")
}

// CreateMany creates many indexes.
func (m *IndexManager) CreateMany(indexes []IndexModel) error {
	if indexes == nil || len(indexes) == 0 {
		return nil
	}

	m.loadMetadata()
	defer m.buildIndexes()

	for _, newidx := range indexes {
		newidx.Options.Name = strings.TrimSpace(newidx.Options.Name)

		if newidx.Options.Name == "" {
			newidx.Options.Name = m.buildIndexName(newidx)
		}

		for _, idx := range m.metadata.Indexes {
			if idx.Options.Name == newidx.Options.Name {
				return fmt.Errorf("%w: index name %s", ErrIndexAlreadyExists, newidx.Options.Name)
			}

			found := 0

			for _, newf := range newidx.Fields {
				for _, f := range idx.Fields {
					if f.Name == newf.Name {
						found++
					}
				}
			}

			if found == len(newidx.Fields) {
				return fmt.Errorf("%w: index fields %v", ErrIndexAlreadyExists, newidx.Fields)
			}
		}

		m.metadata.Indexes = append(m.metadata.Indexes, newidx)
	}

	return m.saveMetadata()
}

// insertIndexes inserts the indexes for a document.
func (m *IndexManager) insertIndexes(doc map[string]any) error {
	for _, idx := range m.metadata.Indexes {
		idxKey, err := m.buildIndexKey(idx, doc, false)
		if err != nil {
			return err
		}

		if err := m.storage.Put(idxKey, nil); err != nil {
			return err
		}
	}

	return nil
}

// deleteIndexes deletes the indexes for a document.
func (m *IndexManager) deleteIndexes(doc map[string]any) error {
	for _, idx := range m.metadata.Indexes {
		idxKey, err := m.buildIndexKey(idx, doc, false)
		if err != nil {
			return err
		}

		if err := m.storage.Delete(idxKey); err != nil {
			return err
		}
	}

	return nil
}

// buildCollectionKey builds the collection key.
func (m *IndexManager) buildCollectionKey() string {
	return fmt.Sprintf(consts.CollectionKeyStringFormat, m.dbname, m.collname)
}

// buildCollectionsKey builds the collections key.
func (m *IndexManager) buildCollectionsKey() string {
	key := fmt.Sprintf(consts.CollectionKeyStringFormat, m.dbname, consts.RemoverWildcard)

	return strings.TrimSuffix(
		key,
		consts.RemoverWildcard,
	)
}

// buildDocumentKey builds the document key.
func (m *IndexManager) buildDocumentKey(docID string) string {
	return fmt.Sprintf(consts.DocumentKeyStringFormat, m.dbname, m.collname, docID)
}

// buildDocumentsKey builds the documents key.
func (m *IndexManager) buildDocumentsKey() string {
	key := fmt.Sprintf(consts.DocumentKeyStringFormat, m.dbname, m.collname, consts.RemoverWildcard)

	return strings.TrimSuffix(
		key,
		consts.RemoverWildcard,
	)
}

// buildIndexes builds the indexes for a collection.
func (m *IndexManager) buildIndexes() {
	go func() {
		m.loadMetadata()

		docsPrefix := strings.TrimSuffix(
			m.buildDocumentKey(consts.RemoverWildcard),
			consts.RemoverWildcard,
		)

		docs, err := m.storage.Scan(docsPrefix)
		if err != nil {
			return
		}

		for _, doc := range docs {
			var docMap map[string]any
			if err := json.Unmarshal(doc, &docMap); err != nil {
				return
			}

			if err := m.insertIndexes(docMap); err != nil {
				return
			}
		}
	}()
}
