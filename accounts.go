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

	// public - mod - admin
	Role string `json:"role" bson:"role"`

	// active - suspended - banned - deleted
	Status string `json:"status" bson:"status"`

	Password  string    `json:"password_hash" bson:"password_hash"`
	CreatedAt time.Time `json:"created_at" bson:"created_at"`
	UpdatedAt time.Time `json:"updated_at" bson:"updated_at"`
}

// new empty account ptr
func NewEmptyAccount() *Account {
	return &Account{}
}

// randomize accoutn values
func (a *Account) Randomize(password string) {
	a.Username = GetUsername()
	a.ID = primitive.NewObjectID()
	a.Email = GetEmail()
	a.Role = GetWeightedRole()
	a.Status = GetWeightedAccountStatus()
	a.Password = password
	a.CreatedAt = time.Now().UTC()
	a.UpdatedAt = time.Now().UTC()
}

// create an account with sepcific values
// username - email - role - password
func NewAccount(u, e, r, p string) *Account {
	return &Account{
		ID:        primitive.NewObjectID(),
		Status:    "active",
		Username:  u,
		Email:     e,
		Role:      r,
		Password:  p,
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
	}
}

// Generate Accounts
func (s *MongoStore) GenerateAccounts(min, max int) {
	accountCount := RandomIntBetween(min, max)
	for i := 0; i < accountCount; i++ {
		account := NewEmptyAccount()
		account.Randomize(GetDefaultPassword())
		s.cAccounts = append(s.cAccounts, account)

		if account.Role == "admin" {
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
	// fmt.Printf("Persisted %v accounts to DB \n", len(response.InsertedIDs))

	return nil
}
