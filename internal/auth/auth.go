package auth

import (
	"golang.org/x/crypto/bcrypt"
)

func HashPassword(plaintextPassword string) ([]byte, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(plaintextPassword), 12)
	if err != nil {
		return nil, err
	}
	return hash, nil
}

// func MakeJWT(userId uuid.UUID, tokenSecret string, expiresIn time.Duration) (string, error) {

// }

// func ValidateJWT(tokenString, tokenSecret string) (uuid.UUID, error) {

// }

// func GetBearerToken(headers http.Header) (string, error) {
	
// }