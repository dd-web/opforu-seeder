package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MongoStore struct {
	MongoClient *mongo.Client
	DB          *mongo.Database
	StartTime   time.Time
	DBName      string

	// tracks the current post number for each board
	PostRefs map[string]int

	// tracking of created documents for faster reference
	cBoards    []*Board
	cThreads   []*Thread
	cPosts     []*Post
	cIdentites []*Identity
	cAccounts  []*Account
	cArticles  []*Article

	// reference to admin account id's
	cAdmins []*primitive.ObjectID

	// map of each users identity for each thread
	// map is in the form of map[threadID]map[userID]*Identity - userID is the account id
	cUserThreadIdentitys map[primitive.ObjectID]map[primitive.ObjectID]*Identity
}

// MongoDB Store - seeder core engine
func NewMongoStore(database string) (*MongoStore, error) {
	uri := ensureEnvVarVaild(getMongoDBConnectionString())
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(uri))
	if err != nil {
		fmt.Println("Error connecting to MongoDB", err)
		if err = client.Disconnect(ctx); err != nil {
			log.Fatal(err)
		}
		return nil, err
	}

	db := client.Database(database)

	return &MongoStore{
		MongoClient:          client,
		DB:                   db,
		DBName:               database,
		PostRefs:             make(map[string]int),
		cUserThreadIdentitys: make(map[primitive.ObjectID]map[primitive.ObjectID]*Identity),
		StartTime:            time.Now().UTC(),
	}, nil
}

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatal(err)
	}

	store, err := NewMongoStore(os.Getenv("MONGO_DATABASE"))
	if err != nil {
		log.Fatal("Error connecting to MongoDB", err)
	}

	fmt.Printf("MongoDB Connection Established at [%s] using database: [%s]\n", store.StartTime, store.DBName)

	fmt.Println("Resetting Database")
	err = store.DB.Drop(context.Background())
	if err != nil {
		log.Fatal("Error dropping database:", err)
	}

	fmt.Println("Creating Collections")
	store.GenCollections()

	fmt.Println("Creating Accounts")
	store.GenerateAccounts(150, 300)

	// create our specific accounts
	devAccount := NewAccount("supafiya", "devduncan89@gmail.com", "admin", "123")
	store.cAccounts = append(store.cAccounts, devAccount)

	fmt.Println("Generating Boards")
	store.GenerateBoards()

	fmt.Println("Generating Articles")
	store.GenerateArticles(20, 60)

	fmt.Println("Generating Threads")
	store.GenerateThreads(200, 500) // between 200 and 500 total threads (all boards)

	fmt.Println("Generating Posts")
	store.GeneratePosts(5, 60) // generates between 5 and 60 posts per thread

	fmt.Printf("\n-----------------------------------------------------\n")
	fmt.Println("Finished Generating Data - Persisting to Database")
	fmt.Printf("-----------------------------------------------------\n\n")

	if err := store.PersistAccounts(); err != nil {
		log.Fatal(err)
	}

	if err := store.PersistArticles(); err != nil {
		log.Fatal(err)
	}

	if err := store.PersistBoards(); err != nil {
		log.Fatal(err)
	}

	if err := store.PersistThreads(); err != nil {
		log.Fatal(err)
	}

	if err := store.PersistPosts(); err != nil {
		log.Fatal(err)
	}

	if err := store.PersistIdentities(); err != nil {
		log.Fatal(err)
	}

	fmt.Printf("\n\n *** Finsihed Seeding Database *** \n\n")
}

// primarily functions as an C++ assert
func ensureEnvVarVaild(val string) string {
	if val == "" {
		log.Fatal("Invalid Environment Variable")
	}
	return val
}

// uses uri if available, otherwise uses individual env vars to construct the uri for connection
func getMongoDBConnectionString() string {
	if os.Getenv("MONGO_URI") != "" {
		return os.Getenv("MONGO_URI")
	}
	envKeys := []string{
		"MONGO_HOST",
		"MONGO_PORT",
	}
	return fmt.Sprintf("mongodb://%s:%s/", os.Getenv(envKeys[0]), os.Getenv(envKeys[1]))
}
