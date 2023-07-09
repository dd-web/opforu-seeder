package main

import (
	"context"
	"fmt"
	"time"
)

// collections to generate
var collections []string = []string{"accounts", "boards", "threads", "posts", "articles", "identities"}

// Generate Collections
func (s *MongoStore) GenCollections() {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	for _, name := range collections {
		err := s.DB.CreateCollection(ctx, name)
		if err != nil {
			fmt.Println("Error creating collection:", name, err)
			continue
		}
		fmt.Println("Creating Collection:", name)
	}
}

// Generic Document Persistance
func (s *MongoStore) PersistDocuments(docs []interface{}, colName string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	col := s.DB.Collection(colName)
	_, err := col.InsertMany(ctx, docs)
	if err != nil {
		return err
	}

	fmt.Printf(" - Persisted %d %s documents to database\n", len(docs), colName)

	return nil
}
