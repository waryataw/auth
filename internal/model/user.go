package model

import (
	"time"
)

// Role Роль пользователя
type Role int32

// User Пользователь
type User struct {
	ID              int64
	Name            string
	Email           string
	Password        string
	PasswordConfirm string
	Role            Role
	CreatedAt       *time.Time
	UpdatedAt       *time.Time
}
