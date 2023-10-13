package main

import (
	"fmt"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Board constants
var defaultBoards [][]string = [][]string{
	{"general", "gen", "general discussion on general topics, generally."},
	{"mathematics", "math", "math is for cool kids"},
	{"programming", "pro", "i wrote a javascript C++ parser"},
	{"technology", "tech", "technology is cool"},
	{"science", "sci", "can we go to mars yet?"},
	{"politics", "pol", "politics is a mess"},
	{"history", "his", "history is cool"},
}

// returns number of default boards
func GetDefaultBoardCount() int {
	return len(defaultBoards)
}

type Board struct {
	ID primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`

	Title       string `json:"title" bson:"title"`
	Short       string `json:"short" bson:"short"`
	Description string `json:"description" bson:"description"`

	Threads []primitive.ObjectID `json:"threads" bson:"threads"`

	CreatedAt *time.Time `json:"created_at" bson:"created_at"`
	UpdatedAt *time.Time `json:"updated_at" bson:"updated_at"`
	DeletedAt *time.Time `bson:"deleted_at,omitempty" json:"deleted_at"`

	// Keeps track of the current post number for that board
	PostRef int `json:"post_ref" bson:"post_ref"`
}

// new empty board ptr
func NewEmptyBoard() *Board {
	return &Board{}
}

// create a board with specific values
// title - short - desc
func NewBoard(t, s, d string) *Board {
	ts := time.Now().UTC()
	return &Board{
		Title:       t,
		Short:       s,
		Description: d,
		Threads:     []primitive.ObjectID{},
		CreatedAt:   &ts,
		UpdatedAt:   &ts,
		PostRef:     0,
	}
}

// returns the index of the configured board
func GetBoardIndex(index int) (*Board, error) {
	if index < 0 || index > len(defaultBoards) {
		return nil, fmt.Errorf("invalid board index %d", index)
	}
	return NewBoard(defaultBoards[index][0], defaultBoards[index][1], defaultBoards[index][2]), nil
}

// Generate boards
func (s *MongoStore) GenerateBoards() {
	for i := 0; i < len(defaultBoards); i++ {
		fmt.Print("\033[G\033[K")
		fmt.Printf(" - Generating Boards: %v/%v", i+1, len(defaultBoards))
		board, err := GetBoardIndex(i)
		if err != nil {
			log.Fatal(err)
		}
		board.ID = primitive.NewObjectID()
		s.cBoards = append(s.cBoards, board)
	}
	fmt.Print("\n")
}

// Get Random Board ID
func (s *MongoStore) GetRandomBoardID() primitive.ObjectID {
	return s.cBoards[RandomIntBetween(0, len(s.cBoards))].ID
}

// Get Random Board
func (s *MongoStore) GetRandomBoard() *Board {
	return s.cBoards[RandomIntBetween(0, len(s.cBoards))]
}

// Get Board by ID
func (s *MongoStore) GetBoardByID(id primitive.ObjectID) *Board {
	var board *Board

	for _, b := range s.cBoards {
		if b.ID == id {
			board = b
			break
		}
	}

	return board
}

// Persist Boards
func (s *MongoStore) PersistBoards() error {
	docs := []interface{}{}

	for _, board := range s.cBoards {
		docs = append(docs, board)
	}

	return s.PersistDocuments(docs, "boards")
}
