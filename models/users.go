package models

import (
	"time"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID        			uuid.UUID `json:"-"`
	CreatedAt 			time.Time `json:"created_at"`
	Name				string	  `json:"name"`
	Email     			string    `json:"email"`
	EncryptedPassword  	[]byte 	  `json:"-"`
	Activated 			bool      `json:"-"`
	Provider  			string    `json:"-"`
}

func NewUser(email, password string) (*User, error) {
	id := uuid.New()

	pw, err := hashPassword(password)
	if err != nil {
		return nil, err
	}

	user := &User{
		ID: 	id,
		Email: 	email,
		CreatedAt: time.Now(),
		EncryptedPassword: pw,
		Provider: "email",
	}

	return user, nil
}

func (user *User) isActivated() bool {
	return user.Activated
}

func hashPassword(plaintextPassword string) ([]byte, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(plaintextPassword), 12)
	if err != nil {
		return nil, err
	}

	return hash, nil
}
