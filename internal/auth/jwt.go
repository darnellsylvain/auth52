package auth

import (
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)


type Auth52Claims struct {
	jwt.RegisteredClaims
	UserId					string `json:"user_id"`
	Email 					string `json:"email"`
}

type TokenType string

const (
	TokenIssuer TokenType = "Auth52"
)


func MakeJWT(userId uuid.UUID, email, tokenSecret string, expiresIn time.Duration) (string, error) {
	claims := &Auth52Claims{
		RegisteredClaims: jwt.RegisteredClaims{
			Subject: userId.String(),
			Issuer: string(TokenIssuer),
			IssuedAt: jwt.NewNumericDate(time.Now()),
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(expiresIn)),
		},
		UserId: userId.String(),
		Email: email,
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
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		if token.Method != jwt.SigningMethodHS256 {
            return nil, fmt.Errorf("unexpected signing method: %v", token.Method)
		}
		return []byte(tokenSecret), nil
	})
    if err != nil {
        if errors.Is(err, jwt.ErrTokenExpired) {
            return uuid.Nil, fmt.Errorf("token expired: %w", err)
        }
        return uuid.Nil, fmt.Errorf("invalid token: %w", err)
    }

	userIDString, err := token.Claims.GetSubject()
	if err != nil {
		return uuid.Nil, err
	}

	issuer, err := token.Claims.GetIssuer()
	if err != nil {
		return uuid.Nil, err
	}
	if issuer != string(TokenIssuer) {
		return uuid.Nil, errors.New("invalid issuer")
	}

    if !token.Valid {
        return uuid.Nil, errors.New("invalid token")
    }

    // parse the user ID out of the claim
    userID, err := uuid.Parse(userIDString)
    if err != nil {
        return uuid.Nil, fmt.Errorf("invalid user id in token: %w", err)
    }
    return userID, nil
}
