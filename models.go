package main

import (
	"context"
	"fmt"
	"log"
	"time"
)

// collections to generate
var collections []string = []string{"accounts", "boards", "threads", "posts", "articles", "article_comments", "article_authors", "identities", "asset_sources", "assets", "sessions"}

// Drops and recreates all collections for a clean slate
func (s *MongoStore) SetupDB() {
	hrPrint("Setup - Reset & Regenerate")
	fmt.Printf(" - Connected to MongoDB using database: %s\n", s.DBName)

	fmt.Print(" - Dropping Collections")
	err := s.DB.Drop(context.Background())
	if err != nil {
		log.Fatal("Error dropping database:", err)
	}

	s.GenCollections()
	hrPrint("Setup Finished - Now Generating Data")
}

// Generate Collections
func (s *MongoStore) GenCollections() {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	for i, name := range collections {
		fmt.Print("\033[G\033[K")
		fmt.Printf(" - Recreating Collections: %v/%v", i+1, len(collections))
		err := s.DB.CreateCollection(ctx, name)
		if err != nil {
			fmt.Println("Error creating collection:", name, err)
			continue
		}
	}
}

// Generic document persistance
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

// Persists all generated data to the database via their respective Persist fns
func (s *MongoStore) PersistAll() {
	hrPrint("Finished Generating Data - Persisting to Database")

	storeFns := []func() error{
		s.PersistAccounts,
		s.PersistSessions,
		s.PersistArticles,
		s.PersistArticleAuthors,
		s.PersistArticleComments,
		s.PersistBoards,
		s.PersistThreads,
		s.PersistPosts,
		s.PersistIdentities,
		s.PersistAssetSrc,
		s.PersistAssets,
	}

	for _, fn := range storeFns {
		if err := fn(); err != nil {
			log.Fatal(err)
		}
	}
}
