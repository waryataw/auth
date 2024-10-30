package user

import (
	"context"
	"fmt"
	"time"

	sq "github.com/Masterminds/squirrel"
	"github.com/waryataw/auth/internal/models"
	"github.com/waryataw/platform_common/pkg/db"
)

// Update Метод изменения пользователя.
func (r repo) Update(ctx context.Context, user *models.User) error {
	builder := sq.Update("users").
		Set("name", user.Name).
		Set("email", user.Email).
		Set("role", user.Role).
		Set("updated_at", time.Now().UTC()).
		Where(sq.Eq{"id": user.ID})

	sql, args, err := builder.PlaceholderFormat(sq.Dollar).ToSql()
	if err != nil {
		return fmt.Errorf("failed to build query: %w", err)
	}

	query := db.Query{
		Name:     "user_repository.Update",
		QueryRaw: sql,
	}

	tag, err := r.db.DB().ExecContext(ctx, query, args...)
	if err != nil {
		return fmt.Errorf("failed to execute query: %w", err)
	}

	if tag.RowsAffected() == 0 {
		return fmt.Errorf("failed to update user: %d not found", user.ID)
	}

	return err
}
