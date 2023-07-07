package main

import (
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

var collNames []string = []string{"accounts", "boards", "threads", "posts", "articles", "identities"}

// generate database "columns" (collections)
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

// return all account ids with admin role
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

// create and return a new identity
func (s *MongoStore) GenIdentityFor(userId primitive.ObjectID, role string) (*Identity, error) {
	identity := NewEmptyIdentity()
	identity.Randomize(userId, role)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	col := s.DB.Collection("identities")
	result, err := col.InsertOne(ctx, identity)
	if err != nil {
		return nil, err
	}

	identity.ID = result.InsertedID.(primitive.ObjectID)

	return identity, nil
}

// generate a random number of threads between min and max per board
func (s *MongoStore) GenThreads(min, max int) error {
	fmt.Println("Generating threads")
	boardList, err := s.GetAllBoards()
	if err != nil {
		return fmt.Errorf("error getting all boards: %v", err)
	}

	for _, board := range boardList {
		fmt.Println("Generating threads for board:", board.Title)

		threadDocs := []interface{}{}
		threadCount := RandomIntBetween(min, max)

		// generate each thread data (save outside of loop to save time)
		for i := 0; i < threadCount; i++ {

			author, err := s.GetRandomUser()
			if err != nil {
				fmt.Println("Error getting random user:", err)
				return fmt.Errorf("error getting random user: %v", err)
			}

			authorIdentity, err := s.GenIdentityFor(author.ID, "creator")
			if err != nil {
				fmt.Println("Error generating identity for author:", err)
				return fmt.Errorf("error generating identity for author: %v", err)
			}

			thread := NewEmptyThread()
			thread.Randomize(board.ID, authorIdentity.ID)

			threadDocs = append(threadDocs, thread)
		}

		// save threads to db
		threadCol := s.DB.Collection("threads")
		threadCtx, tCancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer tCancel()

		threadResponse, err := threadCol.InsertMany(threadCtx, threadDocs)
		if err != nil {
			return fmt.Errorf("error inserting thread documents: %v", err)
		}

		var threadIdList []primitive.ObjectID

		threadIDs := threadResponse.InsertedIDs
		for _, id := range threadIDs {
			// fmt.Println("Thread ID:", id)
			threadIdList = append(threadIdList, id.(primitive.ObjectID))
		}

		board.Threads = threadIdList
		fltr := bson.M{"_id": board.ID}
		updt := bson.M{"$set": bson.M{"threads": threadIdList}}

		boardCtx, bCancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer bCancel()

		boardCol := s.DB.Collection("boards")
		bResult, err := boardCol.UpdateOne(boardCtx, fltr, updt)
		if err != nil {
			fmt.Println("Error updating board with thread ids:", err)
			return fmt.Errorf("error updating board with thread ids: %v", err)
		}

		fmt.Println("Generated Threads in Board:", board.Title, "Thread Count:", len(threadIdList))
		fmt.Println("Updated board with thread ids:", bResult.ModifiedCount)

		// update board with thread ids

	}

	return nil

}

// return a random user
func (s *MongoStore) GetRandomUser() (*Account, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	col := s.DB.Collection("accounts")
	pipeline := []bson.D{bson.D{{"$sample", bson.D{{"size", 1}}}}}

	cur, err := col.Aggregate(ctx, pipeline)
	if err != nil {
		return nil, err
	}

	account := Account{}

	for cur.Next(ctx) {
		if err := cur.Decode(&account); err != nil {
			return nil, err
		}
	}

	return &account, nil
}

// return a random board
func (s *MongoStore) GetRandomBoard() (*Board, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	col := s.DB.Collection("boards")
	pipeline := []bson.D{bson.D{{"$sample", bson.D{{"size", 1}}}}}

	cur, err := col.Aggregate(ctx, pipeline)
	if err != nil {
		return nil, err
	}

	board := Board{}

	for cur.Next(ctx) {
		if err := cur.Decode(&board); err != nil {
			return nil, err
		}
	}
	return &board, nil
}

// return all boards
func (s *MongoStore) GetAllBoards() ([]*Board, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	col := s.DB.Collection("boards")
	cur, err := col.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}

	var boards []*Board

	for cur.Next(ctx) {
		var board Board
		if err := cur.Decode(&board); err != nil {
			return nil, err
		}
		boards = append(boards, &board)
	}

	return boards, nil
}
