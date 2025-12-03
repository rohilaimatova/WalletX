package utils

import (
	"WalletX/config"
	"fmt"
	"github.com/golang-jwt/jwt"
	"time"
)

type CustomClaims struct {
	UserID   int    ` json:"user_id"`
	Username string `json:"username"`
	jwt.StandardClaims
}

// GenerateToken — генерирует JWT токен
func GenerateToken(userID int, username string) (string, error) {
	auth := config.AppSettings.AuthParams

	claims := CustomClaims{
		UserID:   userID,
		Username: username,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Minute * time.Duration(auth.JwtTtlMinutes)).Unix(),
			IssuedAt:  time.Now().Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString([]byte(auth.JwtSecretKey))
}

// ParseToken — разбирает и валидирует токен
func ParseToken(tokenString string) (*CustomClaims, error) {
	auth := config.AppSettings.AuthParams

	token, err := jwt.ParseWithClaims(tokenString, &CustomClaims{}, func(t *jwt.Token) (interface{}, error) {
		// проверка метода подписи
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
		}
		return []byte(auth.JwtSecretKey), nil
	})

	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(*CustomClaims)
	if !ok || !token.Valid {
		return nil, fmt.Errorf("invalid token")
	}

	return claims, nil
}
