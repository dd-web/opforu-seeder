package main

import (
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Post struct {
	ID primitive.ObjectID `json:"_id" bson:"_id,omitempty"`

	// each board has it's own post count reference for each post
	PostNumber int    `json:"post_number" bson:"post_number"`
	Body       string `json:"body" bson:"body"`

	// identity of the user who created the post
	Creator   primitive.ObjectID `json:"creator" bson:"creator"`
	Board     primitive.ObjectID `json:"board" bson:"board"`
	Thread    primitive.ObjectID `json:"thread" bson:"thread"`
	CreatedAt time.Time          `json:"created_at" bson:"created_at"`
	UpdatedAt time.Time          `json:"updated_at" bson:"updated_at"`
}

// new empty post ptr
func NewEmptyPost() *Post {
	return &Post{}
}

// randomize post values
func (p *Post) Randomize(boardId, threadId, creatorId primitive.ObjectID) {
	p.ID = primitive.NewObjectID()
	p.Board = boardId
	p.Thread = threadId
	p.Creator = creatorId // identity
	p.Body = GetParagraphsBetween(1, 8)
	p.CreatedAt = time.Now().UTC()
	p.UpdatedAt = time.Now().UTC()
}

// Generate Posts for each thread
func (s *MongoStore) GeneratePosts(min, max int) {
	for index, thread := range s.cThreads {
		postCount := RandomIntBetween(min, max)
		fmt.Printf("Generating %v posts for thread %v of %v\n", postCount, index, len(s.cThreads))

		for i := 0; i < postCount; i++ {
			postBoard := s.GetBoardByID(thread.Board)
			s.PostRefs[postBoard.Short]++
			postBoard.PostRef = s.PostRefs[postBoard.Short]

			postCreatorAccount := s.GetRandomAccountID()
			postCreatorIdentity := s.GetUserThreadIdentity(postCreatorAccount, thread.ID)

			post := NewEmptyPost()
			post.Randomize(thread.Board, thread.ID, postCreatorIdentity.ID)
			post.PostNumber = s.PostRefs[postBoard.Short]

			if postCreatorIdentity.Role == "mod" && !thread.HasMod(postCreatorIdentity.ID) {
				thread.Mods = append(thread.Mods, postCreatorIdentity.ID)
			}

			thread.Posts = append(thread.Posts, post.ID)
			s.cPosts = append(s.cPosts, post)
		}
	}

}

// Persist Posts
func (s *MongoStore) PersistPosts() error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	docs := []interface{}{}

	for _, post := range s.cPosts {
		docs = append(docs, post)
	}

	postCol := s.DB.Collection("posts")
	response, err := postCol.InsertMany(ctx, docs)
	if err != nil {
		return err
	}

	fmt.Printf(" - Persisted %d %s documents to database\n", len(response.InsertedIDs), "post")
	return nil
}
