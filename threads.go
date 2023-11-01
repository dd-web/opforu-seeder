package main

import (
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Thread struct {
	ID     primitive.ObjectID `json:"_id" bson:"_id,omitempty"`
	Status ThreadStatus       `json:"status" bson:"status"`

	Title string `json:"title" bson:"title"`
	Body  string `json:"body" bson:"body"`
	Slug  string `json:"slug" bson:"slug"`

	Board   primitive.ObjectID `json:"board" bson:"board"`
	Creator primitive.ObjectID `json:"creator" bson:"creator"`

	Posts []primitive.ObjectID `json:"posts" bson:"posts"`
	Mods  []primitive.ObjectID `json:"mods" bson:"mods"`

	Assets []primitive.ObjectID `json:"assets" bson:"assets"`

	Tags  []string     `json:"tags" bson:"tags"`
	Flags []ThreadFlag `json:"flags" bson:"flags"`

	CreatedAt *time.Time `json:"created_at" bson:"created_at"`
	UpdatedAt *time.Time `json:"updated_at" bson:"updated_at"`
	DeletedAt *time.Time `bson:"deleted_at,omitempty" json:"deleted_at,omitempty"`
}

// new empty thread ptr
func NewEmptyThread() *Thread {
	ts := time.Now().UTC()
	return &Thread{
		ID:        primitive.NewObjectID(),
		Title:     GetSentence(),
		Body:      GetParagraphsBetween(1, 4),
		Slug:      GetSlug(12, 16),
		Assets:    []primitive.ObjectID{},
		Posts:     []primitive.ObjectID{},
		Status:    GetWeightedThreadStatus(),
		Tags:      GetRandomTags(),
		CreatedAt: &ts,
	}
}

// Randomize thread values
func (t *Thread) Randomize(boardId, creatorId primitive.ObjectID) {
	ts := time.Now().UTC()
	t.Board = boardId
	t.Creator = creatorId
	t.Mods = []primitive.ObjectID{creatorId}
	t.UpdatedAt = &ts
}

// Generate Threads
func (s *MongoStore) GenerateThreads(min, max int) {
	threadCount := RandomIntBetween(min, max)

	for i := 0; i < threadCount; i++ {
		fmt.Print("\033[G\033[K")
		fmt.Printf(" - Generating Threads: %v/%v", i+1, threadCount)
		mediaCt := RandomIntBetween(0, 9)
		threadBoard := s.GetRandomBoard()

		threadCreatorAccount := s.GetRandomAccountID()
		threadCreatorIdentity := s.GenerateThreadIdentity(threadCreatorAccount, ThreadRoleCreator)

		s.cIdentites = append(s.cIdentites, threadCreatorIdentity)

		thread := NewEmptyThread()
		thread.Randomize(threadBoard.ID, threadCreatorIdentity.ID)

		threadCreatorIdentity.Thread = thread.ID
		threadBoard.Threads = append(threadBoard.Threads, thread.ID)

		s.cUserThreadIdentitys[thread.ID] = make(map[primitive.ObjectID]*Identity)

		pmedIds, err := s.GenerateAssetCount(mediaCt, threadCreatorAccount)
		if err != nil {
			fmt.Println(err)
			continue
		}

		thread.Assets = pmedIds

		s.cUserThreadIdentitys[thread.ID][threadCreatorAccount] = threadCreatorIdentity
		s.cThreads = append(s.cThreads, thread)
	}

	fmt.Print("\n")
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

// Enums

type ThreadStatus string

const (
	ThreadStatusUnknown  ThreadStatus = "unknown" // all enums should have a default unknown value
	ThreadStatusOpen     ThreadStatus = "open"
	ThreadStatusClosed   ThreadStatus = "closed"
	ThreadStatusArchived ThreadStatus = "archived"
	ThreadStatusDeleted  ThreadStatus = "deleted"
)

func (s ThreadStatus) String() string {
	switch s {
	case ThreadStatusOpen:
		return "open"
	case ThreadStatusClosed:
		return "closed"
	case ThreadStatusArchived:
		return "archived"
	case ThreadStatusDeleted:
		return "deleted"
	default:
		return "unknown"
	}
}

type ThreadRole string

const (
	ThreadRoleUnknown ThreadRole = "unknown"
	ThreadRoleUser    ThreadRole = "user"
	ThreadRoleMod     ThreadRole = "mod"
	ThreadRoleCreator ThreadRole = "creator"
)

func (s ThreadRole) String() string {
	switch s {
	case ThreadRoleUser:
		return "user"
	case ThreadRoleMod:
		return "mod"
	case ThreadRoleCreator:
		return "creator"
	default:
		return "unknown"
	}
}

// Bitfield flags for threads
type ThreadFlag uint

const (
	ThreadFlagNone   ThreadFlag = iota      // thread has no flags
	ThreadFlagSticky ThreadFlag = 1 << iota // thread is sticky
	ThreadFlagLocked ThreadFlag = 1 << iota // thread is locked
	ThreadFlagHidden ThreadFlag = 1 << iota // thread is hidden
	ThreadFlagNSFW   ThreadFlag = 1 << iota // thread is NSFW
)
