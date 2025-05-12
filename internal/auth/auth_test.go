package auth

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetBearerToken(t *testing.T) {
	headers := http.Header{}
	headers.Set("Authorization", "Bearer abcDEF123")

	token, err := GetBearerToken(headers)

	assert.NoError(t, err)
	assert.Equal(t, "abcDEF123", token)
}