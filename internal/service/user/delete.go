package user

import (
	"context"
	"fmt"
)

// Delete Метод удаления пользователя.
func (s service) Delete(ctx context.Context, id int64) error {
	if err := s.repository.Delete(ctx, id); err != nil {
		return fmt.Errorf("failed to delete user: %w", err)
	}

	return nil
}
