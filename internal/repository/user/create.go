package user

import (
	"context"
	"fmt"
	"time"

	sq "github.com/Masterminds/squirrel"
	"github.com/waryataw/auth/internal/models"
	"github.com/waryataw/auth/pkg/client/db"
)

// Create Метод создания пользователя.
func (r repo) Create(ctx context.Context, user *models.User) (int64, error) {
	builder := sq.Insert("users").
		Columns(
			"name",
			"email",
			"password",
			"password_confirm",
			"role",
			"created_at",
		).
		Values(
			user.Name,
			user.Email,
			user.Password,
			user.PasswordConfirm,
			user.Role,
			time.Now().UTC(),
		).
		Suffix("RETURNING id")

	sql, args, err := builder.PlaceholderFormat(sq.Dollar).ToSql()
	if err != nil {
		return 0, fmt.Errorf("failed to build query: %w", err)
	}

	query := db.Query{
		Name:     "user_repository.Create",
		QueryRaw: sql,
	}

	var id int64
	err = r.db.DB().QueryRowContext(ctx, query, args...).Scan(&id)
	if err != nil {
		return 0, fmt.Errorf("failed to execute query: %w", err)
	}

	return id, nil
}
