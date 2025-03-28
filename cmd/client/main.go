package main

import (
	"log"
	"time"

	"github.com/wirvii/gopherdb"
)

const (
	n25 = 25
	n10 = 10
	n5  = 5
	n2  = 2
	n1  = 1
)

func main() {
	db, err := gopherdb.NewDatabase("test", "./data")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	coll, err := db.Collection("users")
	if err != nil {
		log.Fatal(err)
	}

	coll.IndexManager.CreateMany(
		[]gopherdb.IndexModel{
			{
				Fields: []gopherdb.IndexField{
					{Name: "name", Order: 1},
				},
			},
		},
	)

	users, err := coll.Find(
		map[string]any{
			"name": "Juan Fernando",
		},
	)
	if err != nil {
		log.Fatal(err)
	}

	log.Println(users)

	/* 	users := []User{
	   		{
	   			Name: "Juan Fernando",
	   			Age:  n25,
	   		},
	   		{
	   			Name: "Genny Lorena",
	   			Age:  n25,
	   		},
	   		{
	   			Name: "Juanita",
	   			Age:  n10,
	   		},
	   		{
	   			Name: "Miguel",
	   			Age:  n10,
	   		},
	   	}

	   	for _, user := range users {
	   		result, err := coll.InsertOne(user)
	   		if err != nil {
	   			log.Fatal(err)
	   		}

	   		log.Println(result.InsertedID)
	   	} */

	db.Close()
}

type User struct {
	ID   string `bson:"_id"`
	Name string `bson:"name"`
	Age  int    `bson:"age"`
}

type Id struct {
	ID   string    `bson:"_id"`
	Time time.Time `bson:"time"`
}
