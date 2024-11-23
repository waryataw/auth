package models

import "github.com/dgrijalva/jwt-go"

// UserClaims Данные пользователя добавляемые в jwt токен.
type UserClaims struct {
	jwt.StandardClaims
	Username string `json:"username"`
	Role     Role   `json:"role"`
}
