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
	ID    primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	Title string             `json:"title" bson:"title"`

	// short 2-4 letter board identifier (used in url)
	Short     string               `json:"short" bson:"short"`
	Desc      string               `json:"desc" bson:"desc"`
	Threads   []primitive.ObjectID `json:"threads" bson:"threads"`
	CreatedAt time.Time            `json:"created_at" bson:"created_at"`
	UpdatedAt time.Time            `json:"updated_at" bson:"updated_at"`

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
	return &Board{
		Title:     t,
		Short:     s,
		Desc:      d,
		Threads:   []primitive.ObjectID{},
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
		PostRef:   0,
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
		board, err := GetBoardIndex(i)
		if err != nil {
			log.Fatal(err)
		}
		board.ID = primitive.NewObjectID()
		s.cBoards = append(s.cBoards, board)
	}
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

// Persis Boards
func (s *MongoStore) PersistBoards() error {
	docs := []interface{}{}

	for _, board := range s.cBoards {
		docs = append(docs, board)
	}

	if err := s.PersistDocuments(docs, "boards"); err != nil {
		return err
	}

	return nil
}
