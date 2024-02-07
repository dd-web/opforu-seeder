package main

import (
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Post struct {
	ID primitive.ObjectID `json:"_id" bson:"_id,omitempty"`

	PostNumber int                `json:"post_number" bson:"post_number"`
	Creator    primitive.ObjectID `json:"creator" bson:"creator"`

	Body   string               `json:"body" bson:"body"`
	Assets []primitive.ObjectID `json:"assets" bson:"assets"`

	Board  primitive.ObjectID `json:"board" bson:"board"`
	Thread primitive.ObjectID `json:"thread" bson:"thread"`

	CreatedAt *time.Time `bson:"created_at" json:"created_at"`
	UpdatedAt *time.Time `bson:"updated_at" json:"updated_at"`
	DeletedAt *time.Time `bson:"deleted_at,omitempty" json:"deleted_at,omitempty"`
}

// New post
func NewPost() *Post {
	ts := time.Now().UTC()
	return &Post{
		ID:         primitive.NewObjectID(),
		PostNumber: 0,
		Creator:    primitive.NilObjectID,
		Body:       GetParagraphsBetween(1, 5),
		Assets:     []primitive.ObjectID{},
		Board:      primitive.NilObjectID,
		Thread:     primitive.NilObjectID,
		CreatedAt:  &ts,
		UpdatedAt:  &ts,
	}
}

// Generate Posts for each thread
func (s *MongoStore) GeneratePosts(min, max int) {
	for index, thread := range s.cThreads {
		postCount := RandomIntBetween(min, max)
		progress := int(float64(index) / float64(len(s.cThreads)) * float64(postCount*len(s.cThreads)-index))

		fmt.Print("\033[G\033[K")
		fmt.Printf(" - Generating Posts: %v/%v", progress, postCount*len(s.cThreads))

		for i := 0; i < postCount; i++ {
			mediaCt := RandomIntBetween(0, 9)
			postBoard := s.GetBoardByID(thread.Board)

			s.PostRefs[postBoard.Short]++
			postBoard.PostRef = s.PostRefs[postBoard.Short]

			postCreatorAccount := s.GetRandomAccount()
			postCreatorIdentity := s.GetUserThreadIdentity(postCreatorAccount.ID, thread.ID)

			post := NewPost()
			post.Board = thread.Board
			post.Thread = thread.ID
			post.Creator = postCreatorIdentity.ID
			post.PostNumber = s.PostRefs[postBoard.Short]

			pmedIds, err := s.GenerateAssetCount(mediaCt, postCreatorAccount.ID)
			if err != nil {
				fmt.Println(err)
				continue
			}

			post.Assets = pmedIds

			if postCreatorIdentity.Role == "mod" && !thread.HasMod(postCreatorIdentity.ID) {
				thread.Mods = append(thread.Mods, postCreatorIdentity.ID)
			}

			thread.Posts = append(thread.Posts, post.ID)
			s.cPosts = append(s.cPosts, post)
		}
	}

	fmt.Print("\033[G\033[K")
	fmt.Printf(" - Generating Posts: %v/%v", len(s.cPosts), len(s.cPosts))
}

// Persist Posts
func (s *MongoStore) PersistPosts() error {
	docs := []interface{}{}

	for _, post := range s.cPosts {
		docs = append(docs, post)
	}

	return s.PersistDocuments(docs, "posts")
}
