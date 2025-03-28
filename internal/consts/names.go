package consts

import "github.com/wirvii/gopherdb/internal/pathmatcher"

var (
	DatabaseKeyPathmatcher            = pathmatcher.NewPath("dbs/{db}")
	DatabaseKeyStringFormat           = "dbs/%s"
	CollectionKeyPathmatcher          = pathmatcher.NewPath("dbs/{db}/colls/{collection}")
	CollectionKeyStringFormat         = "dbs/%s/colls/%s"
	DocumentKeyPathmatcher            = pathmatcher.NewPath("dbs/{db}/colls/{collection}/docs/{docId}")
	DocumentKeyStringFormat           = "dbs/%s/colls/%s/docs/%s"
	IndexKeyPathmatcher               = pathmatcher.NewPath("dbs/{db}/colls/{collection}/idxs/{indexName}/{fieldsKeyHash}/{valuesKeyHash}/{docId}")
	IndexKeyStringFormat              = "dbs/%s/colls/%s/idxs/%s/%s/%s/%s"
	MetadataDatabaseKeyPathmatcher    = pathmatcher.NewPath("meta/dbs/{db}")
	MetadataDatabaseKeyStringFormat   = "meta/dbs/%s"
	MetadataCollectionKeyPathmatcher  = pathmatcher.NewPath("meta/dbs/{db}/colls/{collection}")
	MetadataCollectionKeyStringFormat = "meta/dbs/%s/colls/%s"
)
