package auth

import (
	"time"

	"github.com/golang-jwt/jwt/v5"

	"github.com/Bl00mGuy/Avito-Internship/avito-shop/internal/domain"
)

func GenerateJWT(user *domain.User, secret string) (string, error) {
	claims := jwt.MapClaims{
		"userId":   user.ID,
		"username": user.Username,
		"exp":      time.Now().Add(24 * time.Hour).Unix(), // токен действителен 24 часа
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(secret))
}
