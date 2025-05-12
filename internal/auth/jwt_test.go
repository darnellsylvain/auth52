package auth

import (
	"os"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestMakeJWT_ValidToken(t *testing.T) {
	assert := assert.New(t)

	userID := uuid.New()
	email := "alice@example.com"
	// secret := "test-secret"
	// ttl := time.Minute

	// Generate a token
	tokenString, err := MakeJWT(userID, email)
	assert.NoError(err, "MakeJWT should not return an error")
	assert.NotEmpty(tokenString, "MakeJWT should return a non-empty token")

	// Parse it to inspect claims
	claims := &Auth52Claims{}
	parsed, err := jwt.ParseWithClaims(tokenString, claims, func(tok *jwt.Token) (interface{}, error) {
		return []byte(os.Getenv("AUTH52_JWT_SECRET")), nil
	})
	assert.NoError(err, "ParseWithClaims should not return an error")
	assert.True(parsed.Valid, "token should be valid")

	// Verify the payload
	assert.Equal(userID.String(), claims.UserId, "UserId claim mismatch")
	assert.Equal(email, claims.Email, "Email claim mismatch")
	assert.Equal("Auth52", claims.Issuer, "Issuer claim mismatch")

	now := time.Now()
	assert.WithinDuration(now, claims.IssuedAt.Time, time.Second, "IssuedAt should be ~now")
	assert.WithinDuration(now.Add(time.Hour), claims.ExpiresAt.Time, time.Second, "ExpiresAt should be ~now+ttl")
}

func TestValidateJWT_Success(t *testing.T) {
	assert := assert.New(t)

	userID := uuid.New()
	email := "bob@example.com"

	// Create a valid token
	tokenString, err := MakeJWT(userID, email)
	assert.NoError(err)
	assert.NotEmpty(tokenString)

	// Validate it
	returnedID, err := ValidateJWT(tokenString)
	assert.NoError(err, "ValidateJWT should not return an error for a valid token")
	assert.Equal(userID, returnedID, "ValidateJWT should return the original user ID")
}

// func TestValidateJWT_Expired(t *testing.T) {
// 	assert := assert.New(t)

// 	userID := uuid.New()
// 	email := "eve@example.com"

// 	// Create an already-expired token
// 	tokenString, err := MakeJWT(userID, email)
// 	assert.NoError(err)

// 	// Validate should fail
// 	_, err = ValidateJWT(tokenString)
// 	assert.Error(err, "ValidateJWT should return an error for expired token")
// 	assert.Contains(err.Error(), "expired", "Error should mention expiration")
// }

// func TestValidateJWT_WrongSecret(t *testing.T) {
// 	assert := assert.New(t)

// 	userID := uuid.New()
// 	email := "carol@example.com"
// 	secret := "correct-secret"

// 	// Create a token with one secret...
// 	tokenString, err := MakeJWT(userID, email, secret, time.Minute)
// 	assert.NoError(err)

// 	// ...but validate with another
// 	_, err = ValidateJWT(tokenString, "wrong-secret")
// 	assert.Error(err, "ValidateJWT should return an error when the secret is wrong")
// 	assert.Contains(err.Error(), "invalid token", "Error should mention invalid token")
// }