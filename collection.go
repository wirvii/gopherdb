package gopherdb

import (
	"fmt"
	"reflect"

	"github.com/google/uuid"
	"github.com/wirvii/gopherdb/internal/bson"
	"github.com/wirvii/gopherdb/internal/consts"
	"github.com/wirvii/gopherdb/internal/queryengine"
	"github.com/wirvii/gopherdb/internal/storage"
)

// Collection is a collection of documents.
type Collection struct {
	dbname       string
	collname     string
	storage      storage.Storage
	initialized  bool
	IndexManager *IndexManager
}

// newCollection creates a new collection.
func newCollection(storage storage.Storage, dbname, collname string) (*Collection, error) {
	idxMgr := newIndexManager(storage, dbname, collname)
	if err := idxMgr.loadMetadata(); err != nil {
		return nil, err
	}

	return &Collection{
		dbname:       dbname,
		collname:     collname,
		storage:      storage,
		initialized:  (idxMgr.metadata.DocumentCount > 0),
		IndexManager: idxMgr,
	}, nil
}

// ensureInitialized ensures that the collection is initialized.
func (c *Collection) ensureInitialized() error {
	if c.initialized {
		return nil
	}

	if c.IndexManager.metadata.DocumentCount > 0 {
		if err := c.IndexManager.saveMetadata(); err != nil {
			return err
		}

		c.initialized = true
	}

	return nil
}

// buildDocumentKey builds the key for a document.
func (c *Collection) buildDocumentKey(docID string) string {
	return fmt.Sprintf(consts.DocumentKeyStringFormat, c.dbname, c.collname, docID)
}

// ensureDocumentID ensures that the document has an ID.
func (c *Collection) ensureDocumentID(doc map[string]any) string {
	val, ok := doc[consts.DocumentFieldID]
	if !ok {
		doc[consts.DocumentFieldID] = uuid.NewString()
	} else {
		tval := reflect.ValueOf(val)
		if tval.IsZero() && tval.Kind() == reflect.String {
			doc[consts.DocumentFieldID] = uuid.NewString()
		}
	}

	return fmt.Sprintf("%v", doc[consts.DocumentFieldID])
}

// InsertOne inserts a single document into the collection.
func (c *Collection) InsertOne(doc any) (*InsertOneResult, error) {
	c.IndexManager.loadMetadata()

	// 1. Convertimos a BSON (map[string]interface{})
	parsed, err := bson.ConvertToMap(doc)
	if err != nil {
		return nil, fmt.Errorf("bson conversion failed: %w", err)
	}

	// 2. Generamos ID único
	docID := c.ensureDocumentID(parsed)

	// 3. Verificamos unicidad en índices
	if err := c.IndexManager.checkUniqueness(parsed); err != nil {
		return nil, err
	}

	// 4. Serializamos a JSON
	data, err := bson.Marshal(parsed)
	if err != nil {
		return nil, fmt.Errorf("json marshal failed: %w", err)
	}

	// 5. Guardamos el documento
	key := c.buildDocumentKey(docID)
	if err := c.storage.Put(key, data); err != nil {
		return nil, fmt.Errorf("storage put failed: %w", err)
	}

	c.IndexManager.metadata.DocumentCount++

	// 6. Registramos índices secundarios
	if err := c.IndexManager.insertIndexes(parsed); err != nil {
		return nil, fmt.Errorf("index registration failed: %w", err)
	}

	// 7. Persistimos metadata si es la primera vez
	if err := c.ensureInitialized(); err != nil {
		return nil, fmt.Errorf("metadata initialization failed: %w", err)
	}

	return &InsertOneResult{InsertedID: docID}, nil
}

// FindByID finds a document by its ID.
func (c *Collection) FindByID(id any) (map[string]any, error) {
	key := c.buildDocumentKey(fmt.Sprintf("%v", id))
	data, err := c.storage.Get(key)

	if err != nil {
		return nil, fmt.Errorf("document not found: %w", err)
	}

	var doc map[string]any
	if err := bson.Unmarshal(data, &doc); err != nil {
		return nil, fmt.Errorf("json unmarshal failed: %w", err)
	}

	return doc, nil
}

// Find finds documents by a filter.
func (c *Collection) Find(filter map[string]any) ([]map[string]any, error) {
	c.IndexManager.loadMetadata()
	planner := NewQueryPlanner(c.IndexManager.metadata.Indexes)
	plan := planner.Plan(filter)

	expr, err := queryengine.ParseFilter(filter)
	if err != nil {
		return nil, fmt.Errorf("invalid filter: %w", err)
	}

	documents := make([]map[string]any, 0)

	if plan.IndexUsed != nil {
		docIds, err := c.IndexManager.getDocumentIdsByIndex(*plan.IndexUsed, plan.IndexFilter)
		if err != nil {
			return nil, fmt.Errorf("get document ids by index failed: %w", err)
		}

		for _, docId := range docIds {
			doc, err := c.FindByID(docId)
			if err != nil {
				return nil, fmt.Errorf("find by id failed: %w", err)
			}

			if expr.Evaluate(doc) {
				documents = append(documents, doc)
			}
		}
	} else {
		documentsKey := c.IndexManager.buildDocumentsKey()

		docs, err := c.storage.Scan(documentsKey)
		if err != nil {
			return nil, fmt.Errorf("scan keys failed: %w", err)
		}

		for _, bdoc := range docs {
			var doc map[string]any
			if err := bson.Unmarshal(bdoc, &doc); err != nil {
				return nil, fmt.Errorf("json unmarshal failed: %w", err)
			}

			if expr.Evaluate(doc) {
				documents = append(documents, doc)
			}
		}
	}

	return documents, nil
}
