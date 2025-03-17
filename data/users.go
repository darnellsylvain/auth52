package data

import (
	"encoding/json"
	"io"
	"time"

	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID        int64     `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	Password  password  `json:"-"`
	Activated bool      `json:"activated"`
	Version   int       `json:"-"`
}

type password struct {
	plaintext *string
	hash      []byte
}

type Users []*User

func (u *User) FromJSON(r io.Reader) error {
	decoder := json.NewDecoder(r)
	return decoder.Decode(u)
}
func (u *User) SetPassword(plaintextPassword string) error {
	hash, err := bcrypt.GenerateFromPassword([]byte(plaintextPassword), 12)
	if err != nil {
		return err
	}
	u.Password.plaintext = &plaintextPassword
	u.Password.hash = hash
	return nil
}

func (u *Users) ToJSON(w io.Writer) error {
	encoder := json.NewEncoder(w)
	return encoder.Encode(u)
}

func GetUsers() Users {
	return userList
}

var userList = []*User{
	{
		ID:        1,
		CreatedAt: time.Now(),
		Name:      "Alice Doe",
		Email:     "alice@example.com",
		Password:  password{
			plaintext: nil,
			hash:      []byte("password"),
		},
		Activated: true,
		Version:   1,
	},
	{
		ID:        2,
		CreatedAt: time.Now(),
		Name:      "Bob Smith",
		Email:     "bob@example.com",
		Password:  password{
			plaintext: nil,
			hash:      []byte("password"),
		},
		Activated: false,
		Version:   1,
	},
}


