package services

import (
	"time"

	"github.com/dgrijalva/jwt-go"
	"ssr-metaverse/internal/config"
)

func GenerateToken(userID int, roles []string) (string, error) {
	claims := jwt.MapClaims{
		"user_id": userID,
		"roles":   roles, // Include roles in the payload
		"exp":     time.Now().Add(24 * time.Hour).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(config.JwtSecret)
}