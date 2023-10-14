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

	PostRefs map[string]int

	cBoards    []*Board
	cThreads   []*Thread
	cPosts     []*Post
	cIdentites []*Identity
	cAccounts  []*Account
	cArticles  []*Article
	cMedia     []*Media
	cSessions  []*Session

	cAdmins []*primitive.ObjectID
	cMods   []*primitive.ObjectID

	cUserThreadIdentitys map[primitive.ObjectID]map[primitive.ObjectID]*Identity
	cMediaSourceMap      map[int]*MediaSource
}

func main() {
	store := NewMongoStore()

	store.SetupDB()

	store.GenerateAccounts(150, 300)     // between 150 and 300 accounts
	store.GenerateBoards()               // creates all default boards
	store.GenerateArticles(20, 60)       // between 20 and 60 articles
	store.GenerateMediaSources(400, 800) // between 400 and 800 media sources``
	store.GenerateThreads(200, 500)      // between 200 and 500 total threads (all boards)
	store.GeneratePosts(5, 60)           // generates between 5 and 60 posts per thread

	store.PersistAll()

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

// prints a string with a horizontal split above and below
func hrPrint(str string) {
	fmt.Printf("\n" + HrSplit)
	fmt.Print("  ", str)
	fmt.Printf(HrSplit + "\n")
}

// MongoDB Store - seeder core engine
func NewMongoStore() *MongoStore {
	if err := godotenv.Load(); err != nil {
		log.Fatal(err)
	}

	database := ensureEnvVarVaild(os.Getenv("MONGO_DATABASE"))

	uri := ensureEnvVarVaild(getMongoDBConnectionString())
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(uri))
	if err != nil {
		fmt.Println("Error connecting to MongoDB", err)
		if err = client.Disconnect(ctx); err != nil {
			log.Fatal(err)
		}
		log.Fatal(err)
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
	}
}
