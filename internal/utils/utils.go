package utils

import (
	"fmt"

	"github.com/dgrijalva/jwt-go"
	"github.com/pkg/errors"
	"github.com/waryataw/auth/internal/models"
	"golang.org/x/crypto/bcrypt"
)

// VerifyPassword Проверка пароля.
func VerifyPassword(hashedPassword string, candidatePassword string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(candidatePassword))
	return err == nil
}

// HashPassword хеширует пароль с использованием bcrypt
func HashPassword(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", fmt.Errorf("error hashing password: %v", err)
	}
	return string(hash), nil
}

// VerifyToken Верификация токена.
func VerifyToken(tokenStr string, secretKey []byte) (*models.UserClaims, error) {
	token, err := jwt.ParseWithClaims(
		tokenStr,
		&models.UserClaims{},
		func(token *jwt.Token) (interface{}, error) {
			_, ok := token.Method.(*jwt.SigningMethodHMAC)
			if !ok {
				return nil, errors.Errorf("unexpected token signing method")
			}

			return secretKey, nil
		},
	)
	if err != nil {
		return nil, errors.Errorf("invalid token: %s", err.Error())
	}

	claims, ok := token.Claims.(*models.UserClaims)
	if !ok {
		return nil, errors.Errorf("invalid token claims")
	}

	return claims, nil
}
