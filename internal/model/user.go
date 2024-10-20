package model

import (
	"time"
)

// User Пользователь
type User struct {
	ID              int64
	Name            string
	Email           string
	Password        string
	PasswordConfirm string
	Role            int32
	CreatedAt       *time.Time
	UpdatedAt       *time.Time
}
