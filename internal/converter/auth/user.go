package auth

import (
	"github.com/waryataw/auth/internal/models"
	"github.com/waryataw/auth/pkg/userv1"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// ToGetUserResponse Метод конвертации пользователя в ответ метода получения пользователя
func ToGetUserResponse(user *models.User) *userv1.GetUserResponse {
	var updatedAt *timestamppb.Timestamp
	if user.UpdatedAt != nil {
		updatedAt = timestamppb.New(*user.UpdatedAt)
	}

	return &userv1.GetUserResponse{
		Id:        user.ID,
		Name:      user.Name,
		Email:     user.Email,
		Role:      *getRole(user.Role),
		CreatedAt: timestamppb.New(*user.CreatedAt),
		UpdatedAt: updatedAt,
	}
}

// getRole Получение роли по идентификатору.
func getRole(id models.Role) *userv1.Role {
	roles := []userv1.Role{
		userv1.Role_UNKNOWN,
		userv1.Role_USER,
		userv1.Role_ADMIN,
	}

	return &roles[id]
}

// ToUser Метод конвертации CreateUserRequest в пользователя.
func ToUser(req *userv1.CreateUserRequest) *models.User {
	return &models.User{
		Name:            req.Name,
		Email:           req.Email,
		Password:        req.Password,
		PasswordConfirm: req.PasswordConfirm,
		Role:            models.Role(req.Role),
	}
}

// ToUserForUpdate Метод конвертации UpdateUserRequest в пользователя.
func ToUserForUpdate(req *userv1.UpdateUserRequest) *models.User {
	return &models.User{
		ID:    req.GetId(),
		Name:  req.Name,
		Email: req.Email,
		Role:  models.Role(req.Role),
	}
}
