package auth

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetBearerToken_ValidToken(t *testing.T) {
	headers := http.Header{}
	headers.Set("Authorization", "Bearer abcDEF123")

	token, err := GetBearerToken(headers)
	assert.NoError(t, err)
	assert.Equal(t, "abcDEF123", token)
}

func TestGetBearerToken_NoHeader(t *testing.T) {
	headers := http.Header{}
	token, err := GetBearerToken(headers)
	assert.Error(t, err)
	assert.Equal(t, ErrNoAuthHeader, err)
	assert.Empty(t, token)
}

func TestGetBearerToken_InvalidBearer(t *testing.T) {
	headers := http.Header{}
	headers.Set("Authorization", "Beerer abcDEF123")

	token, err := GetBearerToken(headers)
	assert.Error(t, err)
	assert.Equal(t, ErrBadAuthHeader, err)
	assert.Empty(t, token)
}

func TestGetBearerToken_EmptyToken(t *testing.T) {
	headers := http.Header{}
	headers.Set("Authorization", "Bearer ")

	token, err := GetBearerToken(headers)
	assert.Error(t, err)
	assert.Equal(t, ErrBadAuthHeader, err)
	assert.Empty(t, token)
}
