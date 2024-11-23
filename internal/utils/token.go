package utils

import (
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/waryataw/auth/internal/models"
)

// GenerateToken Генерация токена.
func GenerateToken(user models.User, secretKey []byte, duration time.Duration) (string, error) {
	claims := models.UserClaims{
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(duration).Unix(),
		},
		Username: user.Name,
		Role:     user.Role,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString(secretKey)
}
