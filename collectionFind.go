package gopherdb

import (
	"fmt"

	"github.com/wirvii/gopherdb/internal/bson"
	"github.com/wirvii/gopherdb/internal/queryengine"
	"github.com/wirvii/gopherdb/internal/storage"
	"github.com/wirvii/gopherdb/options"
)

// FindByID finds a document by its ID.
func (c *Collection) FindByID(id any) FindOneResult {
	key := c.buildDocumentKey(fmt.Sprintf("%v", id))
	data, err := c.storage.Get(key)

	if err != nil {
		return FindOneResult{
			Err: fmt.Errorf("document not found: %w", err),
		}
	}

	return FindOneResult{
		raw: storage.KV{
			Key:   key,
			Value: data,
		},
	}
}

// FindOne finds a single document by a filter.
func (c *Collection) FindOne(
	filter map[string]any,
) FindOneResult {
	opts := options.Find().SetLimit(1)

	result := c.Find(filter, opts)
	if result.Err != nil {
		return FindOneResult{
			Err: result.Err,
		}
	}

	if len(result.raw) == 0 {
		return FindOneResult{
			Err: ErrDocumentNotFound,
		}
	}

	return FindOneResult{
		raw:       result.raw[0],
		IndexUsed: result.IndexUsed,
	}
}

// Find finds documents by a filter.
func (c *Collection) Find(
	filter map[string]any,
	opts ...*options.FindOptions,
) FindResult {
	c.IndexManager.loadMetadata()

	opt := options.Find()
	if len(opts) > 0 {
		opt = opt.Merge(opts...)
	}

	planner := NewQueryPlanner(c.IndexManager.metadata.Indexes)
	plan := planner.Plan(filter, opt.Sort)

	expr, err := queryengine.ParseFilter(filter)
	if err != nil {
		return FindResult{
			Err: fmt.Errorf("invalid filter: %w", err),
		}
	}

	rawDocs := make([]map[string]any, 0)
	totalCount := int64(0)

	if plan.IndexUsed != nil {
		docKeys := make([]string, 0)

		if plan.IndexFilter != nil && len(plan.IndexFilter) > 0 {
			docKeys, err = c.IndexManager.getDocumentIndexKeysByIndexAndFilter(*plan.IndexUsed, plan.IndexFilter)

			if err != nil {
				return FindResult{
					Err: fmt.Errorf("get document keys by index failed: %w", err),
				}
			}
		} else {
			docKeys, err = c.IndexManager.getDocumentIndexKeysByIndex(*plan.IndexUsed)
			if err != nil {
				return FindResult{
					Err: fmt.Errorf("get all document keys by index failed: %w", err),
				}
			}
		}

		if plan.UsedForSort {
			// ✅ Aplica paginación directamente
			start := 0
			end := len(docKeys)

			if opt.Skip != nil {
				start = int(*opt.Skip)
			}

			if opt.Limit != nil && start+int(*opt.Limit) < end {
				end = start + int(*opt.Limit)
			}

			if start < len(docKeys) {
				docKeys = docKeys[start:end]
			} else {
				docKeys = nil
			}
		}

		for _, docKey := range docKeys {
			k, err := c.IndexManager.getDocumentIdFromIndexKey(docKey)
			if err != nil {
				return FindResult{
					Err: fmt.Errorf("get document id from index key failed: %w", err),
				}
			}

			result := c.FindByID(k)
			if result.Err != nil && result.Err != ErrDocumentNotFound {
				return FindResult{
					Err: fmt.Errorf("get document by key failed: %w", result.Err),
				}
			}

			if expr.Evaluate(result.Document()) {
				rawDocs = append(rawDocs, result.Document())
				totalCount++
			}
		}
	} else {
		documentsKey := c.IndexManager.buildDocumentsKey()

		docs, err := c.storage.Scan(documentsKey)
		if err != nil {
			return FindResult{
				Err: fmt.Errorf("scan keys failed: %w", err),
			}
		}

		for _, kv := range docs {
			var doc map[string]any
			if err := bson.Unmarshal(kv.Value, &doc); err != nil {
				return FindResult{
					Err: fmt.Errorf("json unmarshal failed: %w", err),
				}
			}

			if expr.Evaluate(doc) {
				rawDocs = append(rawDocs, doc)
				totalCount++
			}
		}
	}

	result := FindResult{
		documents:  rawDocs,
		IndexUsed:  plan.IndexUsed,
		TotalCount: totalCount,
	}

	if opt.Sort != nil && !plan.UsedForSort {
		c.sortDocuments(result.Documents(), opt)
	}

	if opt.Skip != nil {
		if int(*opt.Skip) < len(result.Documents()) {
			result.documents = result.Documents()[*opt.Skip:]
		} else {
			result.documents = nil
		}
	}

	if opt.Limit != nil {
		if int(*opt.Limit) < len(result.Documents()) {
			result.documents = result.Documents()[:*opt.Limit]
		}
	}

	return result
}
