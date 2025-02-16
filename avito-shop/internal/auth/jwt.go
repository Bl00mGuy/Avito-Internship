package auth

import (
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v5"

	"github.com/Bl00mGuy/Avito-Internship/avito-shop/internal/domain"
)

func GenerateJWT(user *domain.User, secret string) (string, error) {
	claims := &jwt.RegisteredClaims{
		Subject:   strconv.FormatInt(user.ID, 10),
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	signedToken, err := token.SignedString([]byte(secret))
	if err != nil {
		return "", err
	}

	return signedToken, nil
}
