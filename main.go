package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MongoStore struct {
	MongoClient *mongo.Client
	DB          *mongo.Database
	StartTime   time.Time
	DBName      string
}

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
		MongoClient: client,
		DB:          db,
		DBName:      database,
		StartTime:   time.Now().UTC(),
	}, nil
}

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatal(err)
	}

	fmt.Println("Starting MongoDB Setup. Fuck SQL. SQL is for fags.")

	st, err := NewMongoStore(os.Getenv("MONGO_DATABASE"))
	if err != nil {
		log.Fatal("Error connecting to MongoDB", err)
	}
	constr := fmt.Sprintf("MongoDB Connection Established at [%s] using database: [%s]\n", st.StartTime, st.DBName)
	fmt.Println(constr)

	fmt.Println("Dropping Database")
	err = st.DB.Drop(context.Background())
	if err != nil {
		fmt.Println("Error dropping database:", err)
		log.Fatal(err)
	}

	fmt.Println("Creating Collections")
	st.GenCollections()
	fmt.Println("Collections Created")

	// dev accounts
	st.CreateAccount(NewAccount("supafiya", "devduncan89@gmail.com", "admin", "123"))
	//  fake accountss
	st.GenAccounts(100)

	// generate boards
	st.GenBoards()

	// generate articles
	st.GenArticles()

	// generate threades
	st.GenThreads(100, 200)
}

// primarily functions as an C++ assert
func ensureEnvVarVaild(val string) string {
	if val == "" {
		log.Fatal("Invalid Environment Variable")
	}
	return val
}

// uses uri if available, otherwise uses individual env vars to construct the uri
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
