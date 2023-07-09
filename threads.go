package main

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Thread struct {
	ID    primitive.ObjectID `json:"_id" bson:"_id,omitempty"`
	Title string             `json:"title" bson:"title"`
	Body  string             `json:"body" bson:"body"`

	// randomized string for url path
	Slug  string             `json:"slug" bson:"slug"`
	Board primitive.ObjectID `json:"board" bson:"board"`

	// identity
	Creator primitive.ObjectID   `json:"creator" bson:"creator"`
	Posts   []primitive.ObjectID `json:"posts" bson:"posts"`

	// mods are identities - includes creator as default
	Mods []primitive.ObjectID `json:"mods" bson:"mods"`

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
	t.ID = primitive.NewObjectID()
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

// Generate Threads
func (s *MongoStore) GenerateThreads(min, max int) {
	threadCount := RandomIntBetween(min, max)

	for i := 0; i < threadCount; i++ {
		threadBoard := s.GetRandomBoard()

		threadCreatorAccount := s.GetRandomAccountID()
		threadCreatorIdentity := s.GenerateThreadIdentity(threadCreatorAccount, "creator")

		s.cIdentites = append(s.cIdentites, threadCreatorIdentity)

		thread := NewEmptyThread()
		thread.Randomize(threadBoard.ID, threadCreatorIdentity.ID)

		threadCreatorIdentity.Thread = thread.ID
		threadBoard.Threads = append(threadBoard.Threads, thread.ID)

		s.cUserThreadIdentitys[thread.ID] = make(map[primitive.ObjectID]*Identity)

		s.cUserThreadIdentitys[thread.ID][threadCreatorAccount] = threadCreatorIdentity
		s.cThreads = append(s.cThreads, thread)
	}
}

// Thread Mod list contains id?
func (t *Thread) HasMod(id primitive.ObjectID) bool {
	for _, mod := range t.Mods {
		if mod == id {
			return true
		}
	}
	return false
}

// Persist Threads
func (s *MongoStore) PersistThreads() error {
	docs := []interface{}{}

	for _, thread := range s.cThreads {
		docs = append(docs, thread)
	}

	if err := s.PersistDocuments(docs, "threads"); err != nil {
		return err
	}

	return nil
}
