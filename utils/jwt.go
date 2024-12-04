package utils

import (
	"errors"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
)

var secretKey = []byte("your_secret_key")

// Claims represents the payload of the JWT token
type Claims struct {
	UserID string `json:"user_id"`
	jwt.StandardClaims
}

// GenerateJWT generates a JWT with a user ID
func GenerateJWT(userID string) (string, error) {
	expirationTime := time.Now().Add(24 * time.Hour)
	claims := Claims{
		UserID: userID,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
			Issuer:    "trip-planner",
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(secretKey)
}

// ValidateJWT validates the token and extracts claims
func ValidateJWT(tokenString string) (*Claims, error) {
	if !strings.HasPrefix(tokenString, "Bearer ") {
		return nil, errors.New("invalid token format")
	}
	tokenString = strings.TrimPrefix(tokenString, "Bearer ")

	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(t *jwt.Token) (interface{}, error) {
		return secretKey, nil
	})
	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(*Claims)
	if !ok || !token.Valid {
		return nil, errors.New("invalid token")
	}

	return claims, nil
}
