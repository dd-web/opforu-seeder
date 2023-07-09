package main

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Article struct {
	ID        primitive.ObjectID `json:"_id" bson:"_id,omitempty"`
	AuthorID  primitive.ObjectID `json:"author" bson:"author"`
	Title     string             `json:"title" bson:"title"`
	Body      string             `json:"body" bson:"body"`
	Slug      string             `json:"slug" bson:"slug"`
	Tags      []string           `json:"tags" bson:"tags"`
	CreatedAt time.Time          `json:"created_at" bson:"created_at"`
	UpdatedAt time.Time          `json:"updated_at" bson:"updated_at"`
}

// new empty article ptr
func NewEmptyArticle() *Article {
	return &Article{}
}

// randomize article values
func (a *Article) Randomize() {
	a.Title = GetSentence()
	a.Tags = GetRandomTags(0, 5)
	a.Body = GetParagraphsBetween(3, 10)
	a.Slug = GetSlug(8, 16)
	a.CreatedAt = time.Now().UTC()
	a.UpdatedAt = time.Now().UTC()
}

// create an article with sepcific values
// title - body - author_id
func NewArticle(t, b string, author primitive.ObjectID) *Article {
	return &Article{
		Title:     t,
		Body:      b,
		AuthorID:  author,
		Slug:      GetSlug(8, 16),
		Tags:      GetRandomTags(0, 5),
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
	}
}
