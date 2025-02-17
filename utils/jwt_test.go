package utils

import (
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"
)

func TestCreateToken(t *testing.T) {
	username := "testuser"
	token, err := CreateToken(username)
	assert.NoError(t, err)
	assert.NotEmpty(t, token)
}

func TestVerifyToken(t *testing.T) {
	verifier := JWTVerifier{}
	t.Run("ValidToken", func(t *testing.T) {
		username := "testuser"
		token, err := CreateToken(username)
		assert.NoError(t, err)
		assert.NotEmpty(t, token)

		verifier := JWTVerifier{}
		verifiedUsername, err := verifier.VerifyToken(token)
		assert.NoError(t, err)
		assert.Equal(t, username, verifiedUsername)
	})
	t.Run("InvalidToken", func(t *testing.T) {
		_, err := verifier.VerifyToken("invalid.token.string")
		assert.Error(t, err)
	})
	t.Run("ExpiredToken", func(t *testing.T) {
		expiredClaims := jwt.MapClaims{
			"username": "testuser",
			"exp":      time.Now().Add(-time.Hour).Unix(),
		}
		token := jwt.NewWithClaims(jwt.SigningMethodHS256, expiredClaims)
		tokenString, err := token.SignedString(secretKey)
		assert.NoError(t, err)

		verifier := JWTVerifier{}
		_, err = verifier.VerifyToken(tokenString)
		assert.Error(t, err)
	})
}
