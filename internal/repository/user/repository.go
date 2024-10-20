package user

import (
	"context"
	"fmt"
	"time"

	sq "github.com/Masterminds/squirrel"
	"github.com/waryataw/auth/internal/client/db"
	"github.com/waryataw/auth/internal/model"
	"github.com/waryataw/auth/internal/repository"
)

const (
	tableName = "users"

	idColumn              = "id"
	nameColumn            = "name"
	emailColumn           = "email"
	passwordColumn        = "password"
	passwordConfirmColumn = "password_confirm"
	roleColumn            = "role"
	createdAtColumn       = "created_at"
	updatedAtColumn       = "updated_at"
)

type repo struct {
	db db.Client
}

// NewRepository Конструктор репозитория пользователя
func NewRepository(db db.Client) repository.UserRepository {
	return &repo{db: db}
}

func (r *repo) Create(ctx context.Context, user *model.User) (int64, error) {
	builder := sq.Insert(tableName).
		Columns(
			nameColumn,
			emailColumn,
			passwordColumn,
			passwordConfirmColumn,
			roleColumn,
			createdAtColumn,
		).
		Values(
			user.Name,
			user.Email,
			user.Password,
			user.PasswordConfirm,
			user.Role,
			time.Now().UTC(),
		).
		Suffix("RETURNING id")

	sql, args, err := builder.PlaceholderFormat(sq.Dollar).ToSql()
	if err != nil {
		return 0, fmt.Errorf("failed to build query: %w", err)
	}

	query := db.Query{
		Name:     "user_repository.Create",
		QueryRaw: sql,
	}

	var id int64
	err = r.db.DB().QueryRowContext(ctx, query, args...).Scan(&id)
	if err != nil {
		return 0, fmt.Errorf("failed to execute query: %w", err)
	}

	return id, nil
}

func (r *repo) Get(ctx context.Context, id int64) (*model.User, error) {
	builder := sq.Select(
		idColumn,
		nameColumn,
		emailColumn,
		passwordColumn,
		passwordConfirmColumn,
		roleColumn,
		createdAtColumn,
		updatedAtColumn,
	).
		From(tableName).
		Where(sq.Eq{idColumn: id})

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

func (r *repo) Update(ctx context.Context, user *model.User) error {
	builder := sq.Update(tableName).
		Set(nameColumn, user.Name).
		Set(emailColumn, user.Email).
		Set(roleColumn, user.Role).
		Set(updatedAtColumn, time.Now().UTC()).
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

func (r *repo) Delete(ctx context.Context, id int64) error {
	builder := sq.
		Delete(tableName).
		Where(sq.Eq{idColumn: id})

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
