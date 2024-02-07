package main

import (
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Article struct {
	ID primitive.ObjectID `json:"_id" bson:"_id,omitempty"`

	AuthorID  primitive.ObjectID   `json:"author" bson:"author"`         // ArticleAuthor id
	CoAuthors []primitive.ObjectID `json:"co_authors" bson:"co_authors"` // ArticleAuthor id's

	Status     ArticleStatus `bson:"status" json:"status"`
	CommentRef int           `json:"comment_ref" bson:"comment_ref"`

	Comments []primitive.ObjectID `json:"comments" bson:"comments"`
	Assets   []primitive.ObjectID `json:"assets" bson:"assets"`

	Title string   `json:"title" bson:"title"`
	Body  string   `json:"body" bson:"body"`
	Slug  string   `json:"slug" bson:"slug"`
	Tags  []string `json:"tags" bson:"tags"`

	CreatedAt *time.Time `bson:"created_at" json:"created_at"`
	UpdatedAt *time.Time `bson:"updated_at" json:"updated_at"`
	DeletedAt *time.Time `bson:"deleted_at,omitempty" json:"deleted_at,omitempty"`
}

type ArticleComment struct {
	ID         primitive.ObjectID `json:"_id" bson:"_id,omitempty"`
	AuthorID   primitive.ObjectID `json:"author" bson:"author"`                     // account id
	AuthorAnon bool               `json:"author_anonymous" bson:"author_anonymous"` // only admins/mods have the option

	CommentNumber int                  `json:"comment_number" bson:"comment_number"`
	Body          string               `json:"body" bson:"body"`
	Assets        []primitive.ObjectID `json:"assets" bson:"assets"`

	CreatedAt *time.Time `bson:"created_at" json:"created_at"`
	UpdatedAt *time.Time `bson:"updated_at" json:"updated_at"`
	DeletedAt *time.Time `bson:"deleted_at,omitempty" json:"deleted_at,omitempty"`
}

type ArticleAuthor struct {
	ID        primitive.ObjectID `json:"_id" bson:"_id,omitempty"`
	AuthorID  primitive.ObjectID `json:"author" bson:"author"`       // account id
	Anonymize bool               `json:"anonymize" bson:"anonymize"` // make author/coauthor with above id anonymous
}

// new article author
func NewArticleAuthor(author primitive.ObjectID, anonimize bool) *ArticleAuthor {
	return &ArticleAuthor{
		ID:        primitive.NewObjectID(),
		AuthorID:  author,
		Anonymize: anonimize,
	}
}

// new article comment
func NewArticleComment() *ArticleComment {
	ts := time.Now().UTC()
	return &ArticleComment{
		ID:            primitive.NewObjectID(),
		AuthorID:      primitive.NilObjectID,
		AuthorAnon:    false,
		CommentNumber: 0,
		Body:          GetParagraphsBetween(1, 5),
		Assets:        []primitive.ObjectID{},
		CreatedAt:     &ts,
		UpdatedAt:     &ts,
	}
}

// new article
func NewArticle() *Article {
	ts := time.Now().UTC()
	return &Article{
		ID:         primitive.NewObjectID(),
		AuthorID:   primitive.NilObjectID,
		CoAuthors:  []primitive.ObjectID{},
		Status:     GetWeightedArticleStatus(),
		CommentRef: 0,
		Comments:   []primitive.ObjectID{},
		Assets:     []primitive.ObjectID{},
		Title:      GetSentence(),
		Body:       GetParagraphsBetween(3, 10),
		Slug:       GetSlug(8, 16),
		Tags:       GetRandomTags(),
		CreatedAt:  &ts,
		UpdatedAt:  &ts,
	}
}

// Generate Articles
func (s *MongoStore) GenerateArticles(min, max int) {
	articleCount := RandomIntBetween(min, max)
	for i := 0; i < articleCount; i++ {
		fmt.Print("\033[G\033[K")
		fmt.Printf(" - Generating Articles: %v/%v", i+1, articleCount)

		article := NewArticle()
		articleAuthors := s.GetRandomModAdminIDList()

		for k, v := range articleAuthors {
			aa := NewArticleAuthor(v, RandomIntBetween(0, 100) > 90)
			if k == 0 {
				article.AuthorID = aa.ID
			} else {
				article.CoAuthors = append(article.CoAuthors, aa.ID)
			}
			s.cArticleAuthors = append(s.cArticleAuthors, aa)
		}

		commentCount := RandomIntBetween(0, 100)

		for j := 0; j < commentCount; j++ {
			comment := NewArticleComment()
			commentAuthor := s.GetRandomAccount()
			comment.AuthorID = commentAuthor.ID
			comment.CommentNumber = article.CommentRef + 1
			article.CommentRef++

			if IsStaffRole(commentAuthor.Role) && RandomIntBetween(0, 100) > 70 {
				comment.AuthorAnon = true
			}

			if RandomIntBetween(0, 100) > 80 {
				mediaCount := RandomIntBetween(0, 9)
				mediaIds, err := s.GenerateAssetCount(mediaCount, commentAuthor.ID)
				if err != nil {
					fmt.Println(err)
					continue
				}
				comment.Assets = mediaIds
			}
			article.Comments = append(article.Comments, comment.ID)
			s.cArticleComments = append(s.cArticleComments, comment)
		}

		if RandomIntBetween(0, 100) > 60 {
			mediaCount := RandomIntBetween(0, 9)
			mediaIds, err := s.GenerateAssetCount(mediaCount, article.AuthorID)
			if err != nil {
				fmt.Println(err)
				continue
			}
			article.Assets = mediaIds
		}

		s.cArticles = append(s.cArticles, article)
	}

	fmt.Print("\n")
}

// Persist Articles
func (s *MongoStore) PersistArticles() error {
	docs := []interface{}{}
	for _, article := range s.cArticles {
		docs = append(docs, article)
	}
	return s.PersistDocuments(docs, "articles")
}

// Persist Article Comments
func (s *MongoStore) PersistArticleComments() error {
	docs := []interface{}{}
	for _, comment := range s.cArticleComments {
		docs = append(docs, comment)
	}
	return s.PersistDocuments(docs, "article_comments")
}

// Persist Article Authors
func (s *MongoStore) PersistArticleAuthors() error {
	docs := []interface{}{}
	for _, author := range s.cArticleAuthors {
		docs = append(docs, author)
	}
	return s.PersistDocuments(docs, "article_authors")
}

type ArticleStatus string

const (
	ArticleStatusDraft     ArticleStatus = "draft"
	ArticleStatusPublished ArticleStatus = "published"
	ArticleStatusArchived  ArticleStatus = "archived"
	ArticleStatusDeleted   ArticleStatus = "deleted"
)

func GetWeightedArticleStatus() ArticleStatus {
	weight := RandomIntBetween(0, 100)
	if weight < 85 {
		return ArticleStatusPublished
	} else if weight < 90 {
		return ArticleStatusDraft
	} else if weight < 95 {
		return ArticleStatusArchived
	} else {
		return ArticleStatusDeleted
	}
}
