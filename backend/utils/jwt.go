// backend/utils/jwt.go
package utils

import (
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var (
	accessSecret  = []byte(os.Getenv("JWT_SECRET"))
	refreshSecret = []byte(os.Getenv("JWT_REFRESH_SECRET"))
)

func CreateAccessToken(userID string) (string, error) {
	claims := jwt.MapClaims{
		"sub": userID,
		"exp": time.Now().Add(15 * time.Minute).Unix(),
	}
	tok := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return tok.SignedString(accessSecret)
}

func CreateRefreshToken(userID string) (string, error) {
	claims := jwt.MapClaims{
		"sub": userID,
		"exp": time.Now().Add(7 * 24 * time.Hour).Unix(),
	}
	tok := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return tok.SignedString(refreshSecret)
}

func ParseToken(tokenString string, isRefresh bool) (*jwt.Token, error) {
	secret := accessSecret
	if isRefresh {
		secret = refreshSecret
	}
	return jwt.Parse(tokenString, func(t *jwt.Token) (interface{}, error) {
		return secret, nil
	})
}
