package models

import (
	"time"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID        			uuid.UUID `json:"id"`
	CreatedAt 			time.Time `json:"created_at"`
	Name				string	
	Email     			string    `json:"email"`
	EncryptedPassword  	[]byte  	`json:"-"`
	Activated 			bool      `json:"activated"`
	Provider  			string
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



var userList = []*User{
	{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		Email:     "alice@example.com",
		EncryptedPassword:  []byte("password"),
		Activated: true,
	},
	{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		Email:     "bob@example.com",
		EncryptedPassword:  []byte("password"),
		Activated: false,
	},
}
