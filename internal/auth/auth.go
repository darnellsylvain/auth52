package auth

import (
	"errors"
	"fmt"
	"time"

	"github.com/darnellsylvain/auth52/models"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)


type Auth52Claims struct {
	jwt.RegisteredClaims
	UserId					string `json:"user_id"`
	Email 					string `json:"email"`
}

func HashPassword(plaintextPassword string) ([]byte, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(plaintextPassword), 12)
	if err != nil {
		return nil, err
	}
	return hash, nil
}

func MakeJWT(user *models.User, tokenSecret string, expiresIn time.Duration) (string, error) {
	claims := &Auth52Claims{
		RegisteredClaims: jwt.RegisteredClaims{
			Subject: user.ID.String(),
			Issuer: "Auth52",
			IssuedAt: jwt.NewNumericDate(time.Now()),
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(expiresIn)),
		},
		Email: user.Email,
		UserId: user.ID.String(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	signedString, err := token.SignedString([]byte(tokenSecret))
	if err != nil {
		return "", err
	}

	return signedString, nil
}

func ValidateJWT(tokenString, tokenSecret string) (uuid.UUID, error) {
	claims := &Auth52Claims{}
	token, err := jwt.ParseWithClaims(tokenSecret, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(tokenSecret), nil
	})

	if err != nil {
		return uuid.Nil, fmt.Errorf("invalid token: %w", err)
	}

	if !token.Valid {
		return uuid.Nil, errors.New("invalid token")
	}

	id, err := uuid.Parse(claims.UserId)
	return id, nil
}

// func GetBearerToken(headers http.Header) (string, error) {
	
// }