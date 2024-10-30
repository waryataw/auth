package user

import (
	"context"
	"fmt"

	sq "github.com/Masterminds/squirrel"
	"github.com/waryataw/auth/internal/models"
	"github.com/waryataw/platform_common/pkg/db"
)

// Get Метод получения пользователя.
func (r repo) Get(ctx context.Context, id int64, name string) (*models.User, error) {
	builder := sq.Select(
		"id",
		"name",
		"email",
		"password",
		"password_confirm",
		"role",
		"created_at",
		"updated_at",
	).
		From("users")

	if id > 0 {
		builder = builder.Where(sq.Eq{"id": id})
	}

	if name != "" {
		builder = builder.Where(sq.Eq{"name": name})
	}

	sql, args, err := builder.PlaceholderFormat(sq.Dollar).ToSql()
	if err != nil {
		return nil, fmt.Errorf("failed to build query: %w", err)
	}

	query := db.Query{
		Name:     "user_repository.Get",
		QueryRaw: sql,
	}

	var user models.User
	err = r.db.DB().QueryRowContext(ctx, query, args...).
		Scan(
			&user.ID,
			&user.Name,
			&user.Email,
			&user.Password,
			&user.PasswordConfirm,
			&user.Role,
			&user.CreatedAt,
			&user.UpdatedAt,
		)
	if err != nil {
		return nil, fmt.Errorf("failed to execute query: %w", err)
	}

	return &user, nil
}
