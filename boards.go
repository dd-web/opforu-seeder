package main

import (
	"fmt"
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
	ID        primitive.ObjectID   `json:"id,omitempty" bson:"_id,omitempty"`
	Title     string               `json:"title" bson:"title"`
	Short     string               `json:"short" bson:"short"`
	Desc      string               `json:"desc" bson:"desc"`
	Threads   []primitive.ObjectID `json:"threads" bson:"threads"`
	CreatedAt time.Time            `json:"createdAt" bson:"createdAt"`
	UpdatedAt time.Time            `json:"updatedAt" bson:"updatedAt"`
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
	}
}

// returns the index of the configured board
func GetBoardIndex(index int) (*Board, error) {
	if index < 0 || index > len(defaultBoards) {
		return nil, fmt.Errorf("invalid board index %d", index)
	}
	return NewBoard(defaultBoards[index][0], defaultBoards[index][1], defaultBoards[index][2]), nil
}
