package models

import (
	"time"

	"github.com/darnellsylvain/auth52/internal/auth"
	"github.com/darnellsylvain/auth52/internal/database"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID        			uuid.UUID `json:"-"`
	CreatedAt 			time.Time `json:"created_at"`
	Name				*string	  `json:"name"`
	Email     			string    `json:"email"`
	EncryptedPassword  	[]byte 	  `json:"-"`
	Activated 			bool      `json:"-"`
	Provider  			string    `json:"-"`
}

func NewUser(email, password string) (*User, error) {
	id := uuid.New()

	pw, err := auth.HashPassword(password)
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

func (user *User) Authenticate(password string) bool {
	err := bcrypt.CompareHashAndPassword(user.EncryptedPassword, []byte(password))
	return err == nil
}

func FromDBUser(u database.FindUserByEmailRow) *User {
	return &User{
		ID:                	u.ID,
		CreatedAt:			u.CreatedAt.Time,
		Name:              	u.Name,
		Email:             	u.Email,
		EncryptedPassword: 	u.EncryptedPassword,
		Activated:         	u.Activated,
		Provider:         	u.Provider,
	}
}