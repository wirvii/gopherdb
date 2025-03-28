package gopherdb

import "github.com/wirvii/gopherdb/internal/queryengine"

// QueryPlan is the plan for a query.
type QueryPlan struct {
	IndexUsed   *IndexModel
	IndexFilter map[string]any
	IsExact     bool
}

// QueryPlanner is the planner for a query.
type QueryPlanner struct {
	indexes []IndexModel
}

// NewQueryPlanner creates a new QueryPlanner.
func NewQueryPlanner(indexes []IndexModel) *QueryPlanner {
	return &QueryPlanner{
		indexes: indexes,
	}
}

// Plan plans a query.
func (qp *QueryPlanner) Plan(filter map[string]any) *QueryPlan {
	flatFilter := qp.flattenFilter(filter)

	indexUsed := new(IndexModel)
	indexFilter := make(map[string]any)

	for _, index := range qp.indexes {
		indexFilterLocal := make(map[string]any)

		for _, field := range index.Fields {
			if v, ok := flatFilter[field.Name]; ok {
				if eq, ok := v.(map[string]any); ok {
					if val, ok := eq[queryengine.OperatorEqual.String()]; ok {
						indexFilterLocal[field.Name] = val
					}
				}
			}
		}

		if len(indexFilterLocal) > len(indexFilter) {
			indexUsed = &index
			indexFilter = indexFilterLocal
		}
	}

	return &QueryPlan{
		IndexUsed:   indexUsed,
		IndexFilter: indexFilter,
		IsExact:     (indexUsed != nil && len(indexFilter) == len(indexUsed.Fields)),
	}
}

// flattenFilter flattens a filter.
func (qp *QueryPlanner) flattenFilter(filter map[string]any) map[string]any {
	out := map[string]any{}

	for k, v := range filter {
		switch val := v.(type) {
		case map[string]any:
			out[k] = val
		default:
			out[k] = map[string]any{queryengine.OperatorEqual.String(): val}
		}
	}

	return out
}
