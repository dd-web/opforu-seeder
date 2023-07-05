package main

import (
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

var collNames []string = []string{"accounts", "boards", "threads", "posts", "articles", "identities"}

func (s *MongoStore) GenCollections() {
	for _, name := range collNames {

		err := s.DB.CreateCollection(context.TODO(), name)
		if err != nil {
			fmt.Println("Error creating collection:", name, err)
			continue
		}
		fmt.Println("Creating Collection:", name)
	}
}

// generate fake accounts
func (s *MongoStore) GenAccounts(count int) {
	fmt.Println("Generating", count, "fake accounts")
	docs := []interface{}{}

	for i := 0; i < count; i++ {
		account := NewEmptyAccount()
		account.Randomize(GetDefaultPassword())
		docs = append(docs, account)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	col := s.DB.Collection("accounts")
	res, err := col.InsertMany(ctx, docs)
	if err != nil {
		fmt.Println("Error inserting documents into accounts collection:", err)
		return
	}
	fmt.Println("Inserted", len(res.InsertedIDs), "documents into accounts collection")
}

// generate boards
func (s *MongoStore) GenBoards() error {
	fmt.Println("Generating boards")
	docs := []interface{}{}

	for i := 0; i < 7; i++ {
		board, err := GetBoardIndex(i)
		if err != nil {
			return fmt.Errorf("error getting board index %d: %v", i, err)
		}
		docs = append(docs, board)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	col := s.DB.Collection("boards")
	res, err := col.InsertMany(ctx, docs)
	if err != nil {
		return fmt.Errorf("error inserting documents into boards collection: %v", err)
	}

	fmt.Println("Inserted", len(res.InsertedIDs), "documents into boards collection")
	return nil
}

// generate articles
func (s *MongoStore) GenArticles() error {
	fmt.Println("Generating articles")
	adminIDS, err := s.GetAdminAccountIDs()
	if err != nil {
		return fmt.Errorf("error getting admin account ids: %v", err)
	}

	docs := []interface{}{}

	for i := 0; i < 20; i++ {
		art := NewEmptyArticle()
		art.Randomize()
		art.AuthorID = adminIDS[RandomIntBetween(0, len(adminIDS)-1)]

		docs = append(docs, art)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	col := s.DB.Collection("articles")
	res, err := col.InsertMany(ctx, docs)
	if err != nil {
		return fmt.Errorf("error inserting documents into articles collection: %v", err)
	}

	fmt.Println("Inserted", len(res.InsertedIDs), "documents into articles collection")
	return nil

}

// create a specific account
func (s *MongoStore) CreateAccount(account *Account) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	col := s.DB.Collection("accounts")
	_, err := col.InsertOne(ctx, account)
	if err != nil {
		return
	}
}

func (s *MongoStore) GetAdminAccountIDs() ([]primitive.ObjectID, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	col := s.DB.Collection("accounts")
	cur, err := col.Find(ctx, bson.M{"role": "admin"})
	if err != nil {
		return nil, err
	}

	adminIDs := []primitive.ObjectID{}
	for cur.Next(ctx) {
		var account Account
		if err := cur.Decode(&account); err != nil {
			return nil, err
		}
		adminIDs = append(adminIDs, account.ID)
	}

	return adminIDs, nil
}

// make array ptr for object id's
// func NewObjectIDArray() *[]primitive.ObjectID {
// 	return &[]primitive.ObjectID{}
// }
