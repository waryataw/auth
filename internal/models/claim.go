package models

import "github.com/dgrijalva/jwt-go"

const (
	// CreateChatPath Метод создания Чата в сервисе Чат сервер.
	CreateChatPath = "/chat_server_v1.ChatServerService/CreateChat"
)

// UserClaims Данные пользователя добавляемые в jwt токен.
type UserClaims struct {
	jwt.StandardClaims
	Username string `json:"username"`
	Role     Role   `json:"role"`
}
