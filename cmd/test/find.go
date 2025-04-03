package main

import (
	"log"
	"time"

	"github.com/wirvii/gopherdb"
)

// runFindAll runs the find all benchmark.
func runFindAll(coll *gopherdb.Collection) {
	log.Println()
	log.Println("Running find all...")

	benchmarkResult := FindResultBenchmark{
		Title: "Find all",
	}

	filter := map[string]any{}

	start := time.Now()
	results := coll.Find(filter)

	if results.Err != nil {
		log.Fatal(results.Err)
	}

	benchmarkResult.QueryTime = time.Since(start)

	if results.IndexUsed != nil {
		benchmarkResult.IndexUsed = results.IndexUsed.Options.Name
	}

	people := make([]Person, 0)

	start = time.Now()
	err := results.Unmarshal(&people)
	benchmarkResult.TotalCount = int64(len(results.Documents()))

	if err != nil {
		log.Fatal(err)
	}

	benchmarkResult.UnmarshalTime = time.Since(start)

	printFindAllResult(benchmarkResult)
}

// printFindAllResult prints the result of the find all benchmark.
func printFindAllResult(result FindResultBenchmark) {
	log.Println()
	log.Println()
	log.Println("--------------------------------")
	log.Println(result.Title)
	log.Println("Query time:", result.QueryTime)
	log.Println("Unmarshal time:", result.UnmarshalTime)
	log.Println("Index used:", result.IndexUsed)
	log.Println("Total count:", result.TotalCount)
	log.Println("--------------------------------")
}
