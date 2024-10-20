package model

import (
	"time"
)

// Role Роль пользователя
type Role int32

const (
	// RoleUnknown Роль не определена
	RoleUnknown Role = iota
	// RoleUser Роль Пользователь
	RoleUser
	// RoleAdmin Роль админ
	RoleAdmin
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
