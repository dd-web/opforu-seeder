package main

import (
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

var devAccounts [][]string = [][]string{
	{"supafiya", "devduncan89@gmail.com", "123"},
	{"nyro", "nyronic@gmail.com", "123"},
}

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

// New empty account ptr
func NewEmptyAccount() *Account {
	ts := time.Now().UTC()
	return &Account{
		CreatedAt: &ts,
	}
}

// Randomize account fields
func (a *Account) Randomize() {
	pw, _ := HashPassword("123")
	ts := time.Now().UTC()
	a.Username = GetUsername()
	a.ID = primitive.NewObjectID()
	a.Email = GetEmail()
	a.Role = GetWeightedRole()
	a.Status = GetWeightedAccountStatus()
	a.UpdatedAt = &ts
	a.Password = pw
}

// Create a new account with provided Username, Email, Role and Status
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
	s.GenerateDevAccounts()
	accountCount := RandomIntBetween(min, max)
	for i := 0; i < accountCount; i++ {
		fmt.Print("\033[G\033[K")
		fmt.Printf(" - Generating Accounts & Sessions (expensive hashes are slow): %v/%v", i+1, accountCount)

		account := NewEmptyAccount()
		account.Randomize()
		session := NewSessionFromAccount(account.ID)

		s.cSessions = append(s.cSessions, session)
		s.cAccounts = append(s.cAccounts, account)

		if account.Role == AccountRoleAdmin {
			s.cAdmins = append(s.cAdmins, &account.ID)
		} else if account.Role == AccountRoleMod {
			s.cMods = append(s.cMods, &account.ID)
		}
	}

	fmt.Print("\n")
}

// Generate Dev Accounts
func (s *MongoStore) GenerateDevAccounts() {
	for i := 0; i < len(devAccounts); i++ {
		fmt.Print("\033[G\033[K")
		fmt.Printf(" - Generating Dev Accounts: %v/%v", i+1, len(devAccounts))
		account := NewAccount(devAccounts[i][0], devAccounts[i][1], AccountRoleAdmin, AccountStatusActive)
		pwh, _ := HashPassword(devAccounts[i][2])
		account.Password = pwh
		s.cAccounts = append(s.cAccounts, account)
	}
	fmt.Print("\n")
}

// Get Random Admin ID
func (s *MongoStore) GetRandomAdminID() primitive.ObjectID {
	return *s.cAdmins[RandomIntBetween(0, len(s.cAdmins))]
}

// Get Random Account ID
func (s *MongoStore) GetRandomAccountID() primitive.ObjectID {
	return s.cAccounts[RandomIntBetween(0, len(s.cAccounts))].ID
}

// Gets a map of at least one admin account and up to 2 mod accounts (for article co_authors)
func (s *MongoStore) GetRandomModAdminIDList() map[int]primitive.ObjectID {
	idMap := make(map[int]primitive.ObjectID)

	adminIndex := RandomIntBetween(0, len(s.cAdmins))
	aditionalAdminCt := RandomIntBetween(0, 3)

	idMap[0] = *s.cAdmins[adminIndex]

	for i := 1; i < aditionalAdminCt; i++ {
		var adminId primitive.ObjectID

		if adminIndex+i >= len(s.cAdmins) {
			adminId = *s.cAdmins[adminIndex-i]
		} else {
			adminId = *s.cAdmins[adminIndex+i]
		}

		idMap[len(idMap)] = adminId
	}

	if RandomIntBetween(0, 100) > 70 {
		modIndex := RandomIntBetween(0, len(s.cMods))
		idMap[len(idMap)] = *s.cMods[modIndex]

		if RandomIntBetween(0, 100) > 80 {
			var modId primitive.ObjectID

			if modIndex+1 >= len(s.cMods) {
				modId = *s.cMods[modIndex-1]
			} else {
				modId = *s.cMods[modIndex+1]
			}

			idMap[len(idMap)] = modId
		}
	}

	return idMap
}

// Persist Accounts To DB
func (s *MongoStore) PersistAccounts() error {
	docs := []interface{}{}
	for _, account := range s.cAccounts {
		docs = append(docs, account)
	}
	return s.PersistDocuments(docs, "accounts")
}

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
