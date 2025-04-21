package main

import (
	"context"
	"log"

	"github.com/wirvii/gopherdb"
)

func main() {
	db, err := gopherdb.NewDatabase("test", "./data")
	if err != nil {
		log.Fatal(err)
	}

	ctx := context.Background()
	coll, err := db.Collection("people")

	if err != nil {
		log.Fatal(err)
	}

	log.Println("Initializing benchmark...")

	runCreateIndexes(ctx, coll)
	//runInsert(ctx, coll)
	//runFindAll(coll)
	//runFindWithFilterWithSimpleIndex(coll)
	//runFindWithFilterWithCompoundIndex(coll)
	//runUpdate(coll)
	//runDelete(coll)
}
