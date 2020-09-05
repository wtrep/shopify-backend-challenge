package backend

import (
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"time"
)

type User struct {
	Username string
	PwHash   []byte
}

type UserActiveSession struct {
	UUID       uuid.UUID `json:"token"`
	Username   string    `json:"username"`
	CreatedAt  time.Time `json:"created_at"`
	Expiration time.Time `json:"expiration"`
}

type UserSessionRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

const (
	defaultCost = 10
)

func NewUser(username string, password string) (*User, error) {
	pwHash, err := bcrypt.GenerateFromPassword([]byte(password), defaultCost)
	if err != nil {
		return nil, err
	}
	return &User{
		Username: username,
		PwHash:   pwHash,
	}, nil
}

func (u *User) goodPassword(password string) bool {
	err := bcrypt.CompareHashAndPassword(u.PwHash, []byte(password))
	if err != nil {
		return false
	}
	return true
}
