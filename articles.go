package main

import (
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Article struct {
	ID primitive.ObjectID `json:"_id" bson:"_id,omitempty"`

	AuthorID  primitive.ObjectID   `json:"author" bson:"author"`
	CoAuthors []primitive.ObjectID `json:"co_authors" bson:"co_authors"`

	Title string   `json:"title" bson:"title"`
	Body  string   `json:"body" bson:"body"`
	Slug  string   `json:"slug" bson:"slug"`
	Tags  []string `json:"tags" bson:"tags"`

	CreatedAt *time.Time `bson:"created_at" json:"created_at"`
	UpdatedAt *time.Time `bson:"updated_at" json:"updated_at"`
	DeletedAt *time.Time `bson:"deleted_at,omitempty" json:"deleted_at,omitempty"`
}

// new empty article ptr
func NewEmptyArticle() *Article {
	ts := time.Now().UTC()
	return &Article{
		ID:        primitive.NewObjectID(),
		CoAuthors: []primitive.ObjectID{},
		Tags:      []string{},
		CreatedAt: &ts,
	}
}

// randomize article values
func (a *Article) Randomize() {
	ts := time.Now().UTC()
	a.ID = primitive.NewObjectID()
	a.Title = GetSentence()
	a.Tags = GetRandomTags(0, 5)
	a.Body = GetParagraphsBetween(3, 10)
	a.Slug = GetSlug(8, 16)
	a.UpdatedAt = &ts
}

// Generate Articles
func (s *MongoStore) GenerateArticles(min, max int) {
	articleCount := RandomIntBetween(min, max)
	for i := 0; i < articleCount; i++ {
		article := NewEmptyArticle()
		article.Randomize()
		article.AuthorID = s.GetRandomAdminID()
		s.cArticles = append(s.cArticles, article)
	}
}

// Persist Articles
func (s *MongoStore) PersistArticles() error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	docs := []interface{}{}

	for _, article := range s.cArticles {
		docs = append(docs, article)
	}

	articleCol := s.DB.Collection("articles")
	response, err := articleCol.InsertMany(ctx, docs)
	if err != nil {
		return err
	}

	fmt.Printf(" - Persisted %d %s documents to database\n", len(response.InsertedIDs), "article")

	return nil

}
