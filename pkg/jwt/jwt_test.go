package jwt

import (
	"testing"

	"github.com/oxiginedev/hng11-stage-two-task/pkg/models"
	"github.com/stretchr/testify/assert"
)

func TestGenerateAccessToken(t *testing.T) {
	// Set up test data
	secret := "secret"
	user := &models.User{
		ID:        "123",
		FirstName: "john_doe",
		Email:     "john@example.com",
	}
	ttl := int64(3600) // 1 hour

	// Generate access token
	accessToken, err := GenerateAccessToken(secret, user, ttl)

	assert.NoError(t, err, "Failed to generate access token")
	assert.NotEmpty(t, accessToken, "Access token is empty")
}

func TestValidateExpiredAccessToken(t *testing.T) {
	// Set up test data
	secret := "secret"
	user := &models.User{
		ID:        "123",
		FirstName: "john_doe",
		Email:     "john@example.com",
	}
	ttl := int64(-3600) // Expired token (-1 hour)

	// Generate access token
	accessToken, err := GenerateAccessToken(secret, user, ttl)
	assert.NoError(t, err, "Failed to generate access token")

	// Validate access token
	_, err = ValidateAccessToken(secret, accessToken)

	assert.Error(t, err, "Expected error for expired token")
	assert.EqualError(t, err, ErrTokenExpired.Error(), "Error message mismatch")
}
