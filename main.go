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

	// post count reference for each board
	PostRefs map[string]int

	// collections
	cBoards    []*Board
	cThreads   []*Thread
	cPosts     []*Post
	cIdentites []*Identity
	cAccounts  []*Account
	cArticles  []*Article
	cMedia     []*Media

	// reference for admin accounts
	cAdmins []*primitive.ObjectID

	// hash map of user thread identitys
	cUserThreadIdentitys map[primitive.ObjectID]map[primitive.ObjectID]*Identity

	// hash map of media sources
	cMediaSourceMap map[int]*MediaSource
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
		cMediaSourceMap:      make(map[int]*MediaSource),
		StartTime:            time.Now().UTC(),
	}, nil
}

func main() {
	var hrSplit string = "\n-----------------------------------------------------\n"

	if err := godotenv.Load(); err != nil {
		log.Fatal(err)
	}

	store, err := NewMongoStore(os.Getenv("MONGO_DATABASE"))
	if err != nil {
		log.Fatal("Error connecting to MongoDB", err)
	}

	fmt.Printf("MongoDB Connection Established at [%s] using database: [%s]\n\n", store.StartTime, store.DBName)

	fmt.Println("Resetting Database")
	err = store.DB.Drop(context.Background())
	if err != nil {
		log.Fatal("Error dropping database:", err)
	}

	fmt.Println("Creating Collections")
	store.GenCollections()

	/**
	 * In full seeder we'll generate in parallel with as few mutex locks as possible (strict reference locks)
	 * as well as a buffer limit to prevent memory overflow on large data sets
	 */

	fmt.Printf("\n" + hrSplit)
	fmt.Printf("Collections Generated - Generating Data")
	fmt.Printf(hrSplit + "\n")

	fmt.Println(" - Generating Accounts")
	store.GenerateAccounts(150, 300)

	devAccount := NewAccount("supafiya", "devduncan89@gmail.com", AccountRoleAdmin, AccountStatusActive)
	store.cAccounts = append(store.cAccounts, devAccount)

	fmt.Println(" - Generating Boards")
	store.GenerateBoards()

	fmt.Println(" - Generating Articles")
	store.GenerateArticles(20, 60)

	fmt.Println(" - Generating Media Sources")
	store.GenerateMediaSources(400, 800) // between 400 and 800 media sources

	fmt.Println(" - Generating Threads")
	store.GenerateThreads(200, 500) // between 200 and 500 total threads (all boards)
	store.GeneratePosts(5, 60)      // generates between 5 and 60 posts per thread

	fmt.Printf("\n" + hrSplit)
	fmt.Printf("Finished Generating Data - Persisting to Database")
	fmt.Printf(hrSplit + "\n")

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

	if err := store.PersistMedia(); err != nil {
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
	return fmt.Sprintf("mongodb://%s:%s/", ensureEnvVarVaild(os.Getenv(envKeys[0])), ensureEnvVarVaild(os.Getenv(envKeys[1])))
}
