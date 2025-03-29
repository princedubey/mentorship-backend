package utils

import (
	"fmt"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type TokenType string

const (
	AccessToken  TokenType = "access"
	RefreshToken TokenType = "refresh"
)

type TokenClaims struct {
	UserID string    `json:"user_id"`
	Type   TokenType `json:"type"`
	jwt.RegisteredClaims
}

// GenerateTokenPair generates both access and refresh tokens for a user
func GenerateTokenPair(userID string) (string, string, error) {
	accessToken, err := generateToken(userID, AccessToken, 15*time.Minute)
	if err != nil {
		return "", "", fmt.Errorf("error generating access token: %v", err)
	}

	refreshToken, err := generateToken(userID, RefreshToken, 7*24*time.Hour)
	if err != nil {
		return "", "", fmt.Errorf("error generating refresh token: %v", err)
	}

	return accessToken, refreshToken, nil
}

// generateToken creates a new JWT token
func generateToken(userID string, tokenType TokenType, expiration time.Duration) (string, error) {
	key := []byte(os.Getenv("JWT_SECRET_KEY"))
	claims := TokenClaims{
		UserID: userID,
		Type:   tokenType,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(expiration)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
			ID:        uuid.New().String(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(key)
	if err != nil {
		return "", fmt.Errorf("error signing token: %v", err)
	}

	return tokenString, nil
}

// ValidateToken validates and parses a JWT token
func ValidateToken(tokenString string) (*TokenClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &TokenClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(os.Getenv("JWT_SECRET_KEY")), nil
	})

	if err != nil {
		return nil, fmt.Errorf("error parsing token: %v", err)
	}

	if claims, ok := token.Claims.(*TokenClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, fmt.Errorf("invalid token claims")
}

// RefreshTokenPair generates new access and refresh tokens using a valid refresh token
func RefreshTokenPair(refreshTokenString string) (string, string, error) {
	claims, err := ValidateToken(refreshTokenString)
	if err != nil {
		return "", "", fmt.Errorf("invalid refresh token: %v", err)
	}

	if claims.Type != RefreshToken {
		return "", "", fmt.Errorf("token is not a refresh token")
	}

	return GenerateTokenPair(claims.UserID)
}
