package utils

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var secreteKey = []byte("secrete-key")

func CreateToken(username string) (string, error) {
	claims := jwt.MapClaims{
		"username": username,
		"exp":      time.Now().Add(time.Hour * 2).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(secreteKey)
	if err != nil {
		return " ", err
	}
	return tokenString, err
}

func VerifyToken(tokenString string) (string, error) {
	token, err := jwt.Parse(tokenString, func(t *jwt.Token) (interface{}, error) {
		return secreteKey, nil
	})
	if err != nil {
		return "", err
	}
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		username, _ := claims["username"].(string)
		return username, nil
	}

	return "", errors.New("invalid token")
}
