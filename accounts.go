package main

import (
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Account struct {
	ID       primitive.ObjectID `json:"_id" bson:"_id,omitempty"`
	Username string             `json:"username" bson:"username"`
	Email    string             `json:"email" bson:"email"`

	Role   AccountRole   `json:"role" bson:"role"`
	Status AccountStatus `json:"status" bson:"status"`

	Password string `json:"password_hash" bson:"password_hash"`

	CreatedAt *time.Time `bson:"created_at" json:"created_at"`
	UpdatedAt *time.Time `bson:"updated_at" json:"updated_at"`
	DeletedAt *time.Time `bson:"deleted_at,omitempty" json:"deleted_at,omitempty"`
}

// new empty account ptr
func NewEmptyAccount() *Account {
	ts := time.Now().UTC()
	return &Account{
		CreatedAt: &ts,
	}
}

// randomize accoutn values
func (a *Account) Randomize() {
	ts := time.Now().UTC()
	a.Username = GetUsername()
	a.ID = primitive.NewObjectID()
	a.Email = GetEmail()
	a.Role = GetWeightedRole()
	a.Status = GetWeightedAccountStatus()
	a.UpdatedAt = &ts
}

// create an account with sepcific values
// username - email - role - password
func NewAccount(u string, e string, r AccountRole, st AccountStatus) *Account {
	ts := time.Now().UTC()
	return &Account{
		ID:        primitive.NewObjectID(),
		Status:    st,
		Username:  u,
		Email:     e,
		Role:      r,
		CreatedAt: &ts,
		UpdatedAt: &ts,
	}
}

// Generate Accounts
func (s *MongoStore) GenerateAccounts(min, max int) {
	accountCount := RandomIntBetween(min, max)
	for i := 0; i < accountCount; i++ {
		account := NewEmptyAccount()
		account.Randomize()
		s.cAccounts = append(s.cAccounts, account)

		if account.Role == AccountRoleAdmin {
			s.cAdmins = append(s.cAdmins, &account.ID)
		}

	}
}

// Get Random Admin ID
func (s *MongoStore) GetRandomAdminID() primitive.ObjectID {
	return *s.cAdmins[RandomIntBetween(0, len(s.cAdmins))]
}

// Get Random Account ID
func (s *MongoStore) GetRandomAccountID() primitive.ObjectID {
	return s.cAccounts[RandomIntBetween(0, len(s.cAccounts))].ID
}

// Persist Accounts To DB
func (s *MongoStore) PersistAccounts() error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	docs := []interface{}{}

	for _, account := range s.cAccounts {
		docs = append(docs, account)
	}

	accountCol := s.DB.Collection("accounts")
	response, err := accountCol.InsertMany(ctx, docs)
	if err != nil {
		return err
	}

	fmt.Printf(" - Persisted %d %s documents to database\n", len(response.InsertedIDs), "account")

	return nil
}

// enums

type AccountStatus string

const (
	AccountStatusUnknown   AccountStatus = "unknown"
	AccountStatusActive    AccountStatus = "active"
	AccountStatusSuspended AccountStatus = "suspended"
	AccountStatusBanned    AccountStatus = "banned"
	AccountStatusDeleted   AccountStatus = "deleted"
)

func (a AccountStatus) String() string {
	switch a {
	case AccountStatusActive:
		return "active"
	case AccountStatusSuspended:
		return "suspended"
	case AccountStatusBanned:
		return "banned"
	case AccountStatusDeleted:
		return "deleted"
	default:
		return "unknown"
	}
}

type AccountRole string

const (
	AccountRoleUnknown AccountRole = "unknown"
	AccountRolePublic  AccountRole = "public"
	AccountRoleUser    AccountRole = "user"
	AccountRoleMod     AccountRole = "mod"
	AccountRoleAdmin   AccountRole = "admin"
)

func (a AccountRole) String() string {
	switch a {
	case AccountRolePublic:
		return "public"
	case AccountRoleUser:
		return "user"
	case AccountRoleMod:
		return "mod"
	case AccountRoleAdmin:
		return "admin"
	default:
		return "unknown"
	}
}
