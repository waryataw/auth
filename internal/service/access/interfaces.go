package access

import (
	"context"

	"github.com/waryataw/auth/internal/models"
)

// Repository Интерфейс Auth репозитория.
type Repository interface {
	GetUserClaims(accessToken string) (*models.UserClaims, error)
	GetAccessibleRoles(ctx context.Context, path string) (map[string]models.Role, error)
}
