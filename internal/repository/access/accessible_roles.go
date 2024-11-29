package access

import (
	"context"
	"fmt"

	sq "github.com/Masterminds/squirrel"
	"github.com/waryataw/auth/internal/models"
	"github.com/waryataw/platform_common/pkg/db"
)

func (r repo) GetAccessibleRoles(ctx context.Context, path string) (map[string]models.Role, error) {
	builder := sq.Select(
		"path",
		"role",
	).
		From("accessible_roles").
		Where(sq.Eq{"path": path})

	sql, args, err := builder.PlaceholderFormat(sq.Dollar).ToSql()
	if err != nil {
		return nil, fmt.Errorf("failed to build query: %w", err)
	}

	query := db.Query{
		Name:     "access_repository.GetAccessibleRoles",
		QueryRaw: sql,
	}

	accessibleRoles := make(map[string]models.Role)

	rows, err := r.db.DB().QueryContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to execute query: %w", err)
	}

	for rows.Next() {
		var path string
		var role int32

		if err := rows.Scan(&path, &role); err != nil {
			return nil, err
		}

		accessibleRoles[path] = models.Role(role)
	}

	return accessibleRoles, nil
}
