package access

import (
	"context"

	"github.com/waryataw/auth/internal/models"
)

// Repository Интерфейс Auth репозитория.
type Repository interface {
	GetUserClaims(ctx context.Context, accessToken string) (*models.UserClaims, error)
	AccessibleRoles(ctx context.Context) (map[string]models.Role, error)
}
