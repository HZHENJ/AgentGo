package utils

import (
	"errors"
	// "os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// TODO: goble secret key for JWT signing and verification
// var jwtSecret = []byte(os.Getenv("JWT_SECRET"))
var jwtSecret = []byte("JWT_SECRET")

type Claims struct {
	UserId    uint   `json:"user_id"`
	Email     string `json:"email"`
	Username  string `json:"username"`
	// Authority int    `json:"authority"`
	jwt.RegisteredClaims
}

// GenerateToken generates a JWT token for the given user information, including user ID, email, and username. The token is valid for 1 hour.
func GenerateToken(userId uint, email string, username string) (string, error) {
	now := time.Now()
	exp := now.Add(time.Hour * 1) // Token valid for 1 hours

	claims := Claims {
		UserId:   userId,
		Email:    email,
		Username: username,
		RegisteredClaims: jwt.RegisteredClaims {
			IssuedAt: jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(exp),
			Issuer: "agentgo",
		},
	}

	// HS256 signing method
	tokenClaims := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	token, err := tokenClaims.SignedString(jwtSecret)
	return token, err
}

// ParseToken parses and validates a JWT token string, returning the claims if the token is valid.
func ParseToken(tokenString string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return jwtSecret, nil
	})
	
	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims, nil
	}
	
	return nil, errors.New("invalid token")
}