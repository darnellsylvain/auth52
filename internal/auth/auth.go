package auth

import (
	"errors"
	"net/http"
	"strings"

	"golang.org/x/crypto/bcrypt"
)

var (
    ErrNoAuthHeader  = errors.New("authorization header missing or empty")
    ErrBadAuthHeader = errors.New("authorization header format must be 'Bearer <token>'")
)

func HashPassword(plaintextPassword string) ([]byte, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(plaintextPassword), 12)
	if err != nil {
		return nil, err
	}
	return hash, nil
}

func GetBearerToken(headers http.Header) (string, error) {
	auth := headers.Get("Authorization")
	if auth == "" {
		return "", ErrNoAuthHeader
	}

	parts := strings.Split(auth, " ")
	if len(parts) != 2 || parts[0] != "Bearer" {
		return "", ErrBadAuthHeader
	}

	return parts[1], nil
}
