package main

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Account struct {
	ID        primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	Username  string             `json:"username" bson:"username"`
	Email     string             `json:"email" bson:"email"`
	Role      string             `json:"role" bson:"role"`
	Password  string             `json:"password_hash" bson:"password_hash"`
	CreatedAt time.Time          `json:"created_at" bson:"created_at"`
	UpdatedAt time.Time          `json:"updated_at" bson:"updated_at"`
}

// new empty account ptr
func NewEmptyAccount() *Account {
	return &Account{}
}

// randomize accoutn values
func (a *Account) Randomize(password string) {
	a.Username = GetUsername()
	a.Email = GetEmail()
	a.Role = GetWeightedRole()
	a.Password = password
	a.CreatedAt = time.Now().UTC()
	a.UpdatedAt = time.Now().UTC()
}

// create an account with sepcific values
// username - email - role - password
func NewAccount(u, e, r, p string) *Account {
	return &Account{
		Username:  u,
		Email:     e,
		Role:      r,
		Password:  p,
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
	}
}
