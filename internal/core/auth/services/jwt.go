package services

import (
	"time"

	"github.com/dgrijalva/jwt-go"
	"ssr-metaverse/internal/config"
	"ssr-metaverse/internal/core/error"
)

func GenerateToken(userID int, roles []string) (string, *error.APIError) {
	claims := jwt.MapClaims{
		"user_id": userID,
		"roles":   roles,
		"exp":     time.Now().Add(24 * time.Hour).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	signedToken, err := token.SignedString(config.JwtSecret)
	if err != nil {
		return "", &error.APIError{
			Code:    500,
			Message: "Failed to generate authentication token",
		}
	}

	return signedToken, nil
}
