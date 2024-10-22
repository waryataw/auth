package user

import (
	"context"
	"fmt"

	sq "github.com/Masterminds/squirrel"
	"github.com/waryataw/auth/internal/client/db"
	"github.com/waryataw/auth/internal/model"
)

func (r *repo) Get(ctx context.Context, id int64, name string) (*model.User, error) {
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

	var user model.User
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
