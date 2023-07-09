package main

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Thread struct {
	ID      primitive.ObjectID   `json:"_id" bson:"_id,omitempty"`
	Title   string               `json:"title" bson:"title"`
	Body    string               `json:"body" bson:"body"`
	Slug    string               `json:"slug" bson:"slug"`
	Board   primitive.ObjectID   `json:"board" bson:"board"`
	Creator primitive.ObjectID   `json:"creator" bson:"creator"`
	Posts   []primitive.ObjectID `json:"posts" bson:"posts"`
	Mods    []primitive.ObjectID `json:"mods" bson:"mods"`

	// open - closed - archived - deleted
	Status    string    `json:"status" bson:"status"`
	Tags      []string  `json:"tags" bson:"tags"`
	CreatedAt time.Time `json:"created_at" bson:"created_at"`
	UpdatedAt time.Time `json:"updated_at" bson:"updated_at"`
}

// new empty thread ptr
func NewEmptyThread() *Thread {
	return &Thread{}
}

// randomize thread values
func (t *Thread) Randomize(boardId, creatorId primitive.ObjectID) {
	t.Title = GetSentence()
	t.Body = GetParagraphsBetween(3, 10)
	t.Slug = GetSlug(8, 16)
	t.Board = boardId
	t.Creator = creatorId // identity
	t.Posts = []primitive.ObjectID{}
	t.Mods = []primitive.ObjectID{creatorId}
	t.Status = GetWeightedThreadStatus()
	t.Tags = GetRandomTags(0, 5)
	t.CreatedAt = time.Now().UTC()
	t.UpdatedAt = time.Now().UTC()
}
