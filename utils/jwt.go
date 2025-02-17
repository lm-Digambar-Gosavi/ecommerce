package utils

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var secretKey = []byte("secret-key")

type JWTVerifier struct{} // struct that provides a method to verify tokens

func (j JWTVerifier) VerifyToken(tokenString string) (string, error) {
	token, err := jwt.Parse(tokenString, func(t *jwt.Token) (interface{}, error) { // decoding and verifying a JWT token
		return secretKey, nil
	})
	if err != nil {
		return "", err
	}
	// extracts claims (payload data) from a JWT token
	// token.Claims holds the decoded claim
	// .(jwt.MapClaims)) checks if the claims are of type jwt.MapClaims
	// If token.Claims is successfully converted to jwt.MapClaims, ok = true.
	// Otherwise, ok = false
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		username, _ := claims["username"].(string) // .(string)) ensures it's a string.
		return username, nil
	}

	return "", errors.New("invalid token")
}

func CreateToken(username string) (string, error) {
	claims := jwt.MapClaims{
		"username": username,
		"exp":      time.Now().Add(time.Hour * 2).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims) // HMAC SHA-256 (HS256) as the signing algorithm.
	tokenString, err := token.SignedString(secretKey)
	if err != nil {
		return " ", err
	}
	return tokenString, err
}
