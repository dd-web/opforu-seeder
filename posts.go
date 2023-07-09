package main

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Post struct {
	ID         primitive.ObjectID `json:"_id" bson:"_id,omitempty"`
	PostNumber int                `json:"post_number" bson:"post_number"`
	Body       string             `json:"body" bson:"body"`
	Creator    primitive.ObjectID `json:"creator" bson:"creator"`
	Board      primitive.ObjectID `json:"board" bson:"board"`
	Thread     primitive.ObjectID `json:"thread" bson:"thread"`
	CreatedAt  time.Time          `json:"created_at" bson:"created_at"`
	UpdatedAt  time.Time          `json:"updated_at" bson:"updated_at"`
}

// new empty post ptr
func NewEmptyPost() *Post {
	return &Post{}
}

// randomize post values
func (p *Post) Randomize(boardId, threadId, creatorId primitive.ObjectID) {
	p.Board = boardId
	p.Thread = threadId
	p.Creator = creatorId // identity
	p.Body = GetParagraphsBetween(1, 8)
	p.CreatedAt = time.Now().UTC()
	p.UpdatedAt = time.Now().UTC()
}
