package access

import (
	"context"

	"github.com/waryataw/auth/internal/models"
)

var accessibleRoles map[string]models.Role

func (r repo) AccessibleRoles(_ context.Context) (map[string]models.Role, error) {
	if accessibleRoles == nil {
		accessibleRoles = make(map[string]models.Role)
		accessibleRoles[models.CreateChatPath] = models.RoleAdmin
	}

	return accessibleRoles, nil
}
