package user

import (
	"context"
	"fmt"

	sq "github.com/Masterminds/squirrel"
	"github.com/waryataw/auth/pkg/client/db"
)

// Delete Метод добавления пользователя.
func (r repo) Delete(ctx context.Context, id int64) error {
	builder := sq.
		Delete("users").
		Where(sq.Eq{"id": id})

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
		return fmt.Errorf("failed to delete user: %d not found", id)
	}

	return nil
}
