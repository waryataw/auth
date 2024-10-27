package models

import (
	"time"
)

// Role Роль пользователя
type Role int32

const (
	// RoleUnknown Роль не определена
	RoleUnknown Role = 0
	// RoleUser Роль Пользователь
	RoleUser = 1
	// RoleAdmin Роль админ
	RoleAdmin = 2
)

// IsValid Валидация Роли
func (r Role) IsValid() bool {
	switch r {
	case RoleUnknown, RoleUser, RoleAdmin:
		return true
	default:
		return false
	}
}

// User Пользователь.
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
