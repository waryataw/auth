package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net"
	"time"

	authv1 "github.com/waryataw/auth/pkg/authv1"

	sq "github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/waryataw/auth/internal/config"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

var configPath string

func init() {
	flag.StringVar(&configPath, "config-path", ".env", "path to config file")
}

type server struct {
	authv1.UnimplementedAuthServiceServer
	pool *pgxpool.Pool
}

func getUserQuery(ID int64, name string) sq.SelectBuilder {
	query := sq.Select(
		"id",
		"name",
		"email",
		"role",
		"created_at",
		"updated_at",
	).
		From("users")

	if ID > 0 {
		query = query.Where(sq.Eq{"id": ID})
	}

	if name != "" {
		query = query.Where(sq.Eq{"name": name})
	}

	return query
}

// GetUser Получение существующего пользователя
func (s *server) GetUser(ctx context.Context, req *authv1.GetUserRequest) (*authv1.GetUserResponse, error) {
	sql, args, err := getUserQuery(
		req.GetId(),
		req.GetName(),
	).
		PlaceholderFormat(sq.Dollar).
		ToSql()
	if err != nil {
		return nil, fmt.Errorf("failed to build query: %w", err)
	}

	var (
		id                   int64
		roleID               int
		name, email          string
		createdAt, updatedAt *time.Time
	)

	err = s.pool.QueryRow(ctx, sql, args...).
		Scan(
			&id,
			&name,
			&email,
			&roleID,
			&createdAt,
			&updatedAt,
		)
	if err != nil {
		return nil, fmt.Errorf("failed to select user: %w", err)
	}

	var createdAtProto, updatedAtProto *timestamppb.Timestamp

	if createdAt != nil {
		createdAtProto = timestamppb.New(*createdAt)
	}

	if updatedAt != nil {
		updatedAtProto = timestamppb.New(*updatedAt)
	}

	return &authv1.GetUserResponse{
		Id:        id,
		Name:      name,
		Email:     email,
		Role:      *getRole(roleID),
		CreatedAt: createdAtProto,
		UpdatedAt: updatedAtProto,
	}, nil
}

// getRole Получение роли по идентификатору
func getRole(id int) *authv1.Role {
	roles := []authv1.Role{
		authv1.Role_UNKNOWN,
		authv1.Role_USER,
		authv1.Role_ADMIN,
	}

	return &roles[id]
}

// CreateUser Добавление нового пользователя
func (s *server) CreateUser(ctx context.Context, req *authv1.CreateUserRequest) (*authv1.CreateUserResponse, error) {
	query := sq.Insert("users").
		Columns(
			"name",
			"email",
			"password",
			"password_confirm",
			"role",
		).
		Values(
			req.Name,
			req.Email,
			req.Password,
			req.PasswordConfirm,
			req.Role,
		).
		Suffix("RETURNING id")

	sql, args, err := query.PlaceholderFormat(sq.Dollar).ToSql()
	if err != nil {
		return nil, fmt.Errorf("failed to build query: %w", err)
	}

	var userID int64
	if err := s.pool.QueryRow(ctx, sql, args...).Scan(&userID); err != nil {
		return nil, fmt.Errorf("failed to insert user: %w", err)
	}

	log.Printf("inserted user with id: %d", userID)

	return &authv1.CreateUserResponse{
		Id: userID,
	}, nil
}

// UpdateUser Изменение существующего пользователя
func (s *server) UpdateUser(ctx context.Context, req *authv1.UpdateUserRequest) (*emptypb.Empty, error) {
	queryUpdate := sq.Update("users").
		Set("name", req.GetName()).
		Set("email", req.GetEmail()).
		Set("role", req.GetRole()).
		Set("updated_at", time.Now().UTC()).
		Where(sq.Eq{"id": req.GetId()})

	sql, args, err := queryUpdate.PlaceholderFormat(sq.Dollar).ToSql()
	if err != nil {
		return nil, fmt.Errorf("failed to build query: %w", err)
	}

	tag, err := s.pool.Exec(ctx, sql, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to update user: %w", err)
	}
	if tag.RowsAffected() == 0 {
		return nil, fmt.Errorf("failed to update user: %d not found", req.GetId())
	}

	return &emptypb.Empty{}, nil
}

// DeleteUser Удаление существующего пользователя
func (s *server) DeleteUser(ctx context.Context, req *authv1.DeleteUserRequest) (*emptypb.Empty, error) {
	queryDelete := sq.Delete("users").Where(sq.Eq{"id": req.GetId()})

	sql, args, err := queryDelete.PlaceholderFormat(sq.Dollar).ToSql()
	if err != nil {
		return nil, fmt.Errorf("failed to build query: %w", err)
	}

	tag, err := s.pool.Exec(ctx, sql, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to execute query: %w", err)
	}
	if tag.RowsAffected() == 0 {
		return nil, fmt.Errorf("failed to delete user: %d not found", req.GetId())
	}

	return &emptypb.Empty{}, nil
}

func main() {
	flag.Parse()
	ctx := context.Background()

	if err := config.Load(configPath); err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	grpcConfig, err := config.NewGRPCConfig()
	if err != nil {
		log.Fatalf("failed to get grpc config: %v", err)
	}

	pgConfig, err := config.NewPGConfig()
	if err != nil {
		log.Fatalf("failed to get pg config: %v", err)
	}

	lis, err := net.Listen("tcp", grpcConfig.Address())
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	pool, err := pgxpool.New(ctx, pgConfig.DSN())
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}
	defer pool.Close()

	s := grpc.NewServer()
	reflection.Register(s)
	authv1.RegisterAuthServiceServer(s, &server{pool: pool})

	log.Printf("server listening at %v", lis.Addr())

	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
