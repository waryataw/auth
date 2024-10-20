package auth

import (
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/waryataw/auth/internal/model"
	"github.com/waryataw/auth/pkg/authv1"
)

// ToGetUserResponseFromUser Метод конвертации пользователя в ответ метода получения пользователя
func ToGetUserResponseFromUser(user *model.User) *authv1.GetUserResponse {
	var updatedAt *timestamppb.Timestamp
	if user.UpdatedAt != nil {
		updatedAt = timestamppb.New(*user.UpdatedAt)
	}

	return &authv1.GetUserResponse{
		Id:        user.ID,
		Name:      user.Name,
		Email:     user.Email,
		Role:      *getRole(user.Role),
		CreatedAt: timestamppb.New(*user.CreatedAt),
		UpdatedAt: updatedAt,
	}
}

// getRole Получение роли по идентификатору
func getRole(id model.Role) *authv1.Role {
	roles := []authv1.Role{
		authv1.Role_UNKNOWN,
		authv1.Role_USER,
		authv1.Role_ADMIN,
	}

	return &roles[id]
}

// ToUserFromCreateUserRequest Метод конвертации CreateUserRequest в пользователя
func ToUserFromCreateUserRequest(req *authv1.CreateUserRequest) *model.User {
	return &model.User{
		Name:            req.Name,
		Email:           req.Email,
		Password:        req.Password,
		PasswordConfirm: req.PasswordConfirm,
		Role:            model.Role(req.Role),
	}
}

// ToUserFromUpdateUserRequest Метод конвертации UpdateUserRequest в пользователя
func ToUserFromUpdateUserRequest(req *authv1.UpdateUserRequest) *model.User {
	return &model.User{
		ID:    req.GetId(),
		Name:  req.Name,
		Email: req.Email,
		Role:  model.Role(req.Role),
	}
}
