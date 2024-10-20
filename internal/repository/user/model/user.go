package model

import (
	"database/sql"
	"time"
)

// User Пользователь репо слоя
type User struct {
	ID              int64        `db:"id"`
	Name            string       `db:"name"`
	Email           string       `db:"email"`
	Password        string       `db:"password"`
	PasswordConfirm string       `db:"password_confirm"`
	Role            int32        `db:"role"`
	CreatedAt       time.Time    `db:"created_at"`
	UpdatedAt       sql.NullTime `db:"updated_at"`
}
