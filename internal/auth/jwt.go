package auth

import (
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type Auth52Claims struct {
	jwt.RegisteredClaims
	UserId string `json:"user_id"`
	Email  string `json:"email"`
}

type TokenType string

const (
	TokenIssuer TokenType = "Auth52"
)

var (
	jwtSecret  = []byte(os.Getenv("AUTH52_JWT_SECRET"))
	jwtExpires = time.Hour
)

func CreateToken(userId uuid.UUID, email string) (string, *Auth52Claims, error) {
	claims := &Auth52Claims{
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   userId.String(),
			Issuer:    string(TokenIssuer),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(jwtExpires)),
		},
		UserId: userId.String(),
		Email:  email,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	signedString, err := token.SignedString(jwtSecret)
	if err != nil {
		return "", nil, err
	}

	return signedString, claims, nil
}

func ValidateToken(tokenString string) (*Auth52Claims, error) {
	claims := &Auth52Claims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		if token.Method != jwt.SigningMethodHS256 {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Method)
		}
		return jwtSecret, nil
	})
	if err != nil {
		if errors.Is(err, jwt.ErrTokenExpired) {
			return nil, fmt.Errorf("token expired: %w", err)
		}
		return nil, fmt.Errorf("invalid token: %w", err)
	}

	userIDString, err := token.Claims.GetSubject()
	if err != nil {
		return nil, err
	}

	issuer, err := token.Claims.GetIssuer()
	if err != nil {
		return nil, err
	}
	if issuer != string(TokenIssuer) {
		return nil, errors.New("invalid issuer")
	}

	if !token.Valid {
		return nil, errors.New("invalid token")
	}

	// parse the user ID out of the claim
	_, err = uuid.Parse(userIDString)
	if err != nil {
		return nil, fmt.Errorf("invalid user id in token: %w", err)
	}
	return claims, nil
}
