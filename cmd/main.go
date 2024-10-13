package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net"
	"time"

	desc "github.com/waryataw/auth/pkg/auth_v1"

	sq "github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/waryataw/auth/internal/config"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

const errNoRows = "no rows in result set"

var configPath string

func init() {
	flag.StringVar(&configPath, "config-path", ".env", "path to config file")
}

type server struct {
	desc.UnimplementedAuthServiceServer
	pool *pgxpool.Pool
}

func getUserQuery(ID int64, name string) sq.SelectBuilder {
	query := sq.Select("id", "name", "email", "role", "created_at", "updated_at").
		From("users").
		PlaceholderFormat(sq.Dollar).
		Limit(1)

	if ID > 0 {
		query = query.Where(sq.Eq{"id": ID})
	}

	if name != "" {
		query = query.Where(sq.Like{"name": name})
	}

	return query
}

// GetUser Получение существующего пользователя
func (s *server) GetUser(ctx context.Context, req *desc.GetUserRequest) (*desc.GetUserResponse, error) {
	sql, args, err := getUserQuery(req.GetId(), req.GetName()).ToSql()
	if err != nil {
		return nil, fmt.Errorf("failed to build query: %v", err)
	}

	var id int64
	var roleID int
	var name, email string
	var createdAt, updatedAt *time.Time

	err = s.pool.QueryRow(ctx, sql, args...).Scan(&id, &name, &email, &roleID, &createdAt, &updatedAt)
	if err != nil {
		return nil, fmt.Errorf("failed to select user: %v", err)
	}

	role, err := roleByRoleID(roleID)
	if err != nil {
		return nil, fmt.Errorf("failed to select user: %v", err)
	}

	var createdAtProto, updatedAtProto *timestamppb.Timestamp

	if createdAt != nil {
		createdAtProto = timestamppb.New(*createdAt)
	}

	if updatedAt != nil {
		updatedAtProto = timestamppb.New(*updatedAt)
	}

	return &desc.GetUserResponse{
		Id:        id,
		Name:      name,
		Email:     email,
		Role:      *role,
		CreatedAt: createdAtProto,
		UpdatedAt: updatedAtProto,
	}, nil
}

// roleByRoleID Получение роли по идентификатору роли, если идентификатор вне диапазона возвращаем ошибку
func roleByRoleID(id int) (*desc.Role, error) {

	roles := []desc.Role{
		desc.Role_UNKNOWN,
		desc.Role_USER,
		desc.Role_ADMIN,
	}

	if id < 0 || id >= len(roles) {
		return nil, fmt.Errorf("invalid role value: %d", id)
	}

	return &roles[id], nil
}

// CreateUser Добавление нового пользователя
func (s *server) CreateUser(ctx context.Context, req *desc.CreateUserRequest) (*desc.CreateUserResponse, error) {

	query := sq.Insert("users").
		PlaceholderFormat(sq.Dollar).
		Columns("name", "email", "password", "password_confirm", "role").
		Values(req.Name, req.Email, req.Password, req.PasswordConfirm, req.Role).
		Suffix("RETURNING id")

	sql, args, err := query.ToSql()
	if err != nil {
		return nil, fmt.Errorf("failed to build query: %v", err)
	}

	var userID int64
	err = s.pool.QueryRow(ctx, sql, args...).Scan(&userID)
	if err != nil {
		return nil, fmt.Errorf("failed to insert user: %v", err)
	}

	log.Printf("inserted user with id: %d", userID)

	return &desc.CreateUserResponse{
		Id: userID,
	}, nil
}

// UpdateUser Изменения существующего пользователя
func (s *server) UpdateUser(ctx context.Context, req *desc.UpdateUserRequest) (*emptypb.Empty, error) {

	querySelect := sq.Select("1").
		From("users").
		PlaceholderFormat(sq.Dollar).
		Where(sq.Eq{"id": req.GetId()}).
		Limit(1)

	sql, args, err := querySelect.ToSql()
	if err != nil {
		return nil, fmt.Errorf("failed to build query: %v", err)
	}

	var exist int

	err = s.pool.QueryRow(ctx, sql, args...).Scan(&exist)
	if err != nil {
		if err.Error() == errNoRows {
			return nil, fmt.Errorf("user: %d not founded", req.GetId())
		}
		return nil, fmt.Errorf("failed to execute query: %v", err)
	}

	queryUpdate := sq.Update("users").
		PlaceholderFormat(sq.Dollar).
		Set("name", req.GetName()).
		Set("email", req.GetEmail()).
		Set("role", req.GetRole()).
		Set("updated_at", time.Now().UTC()).
		Where(sq.Eq{"id": req.GetId()})

	sql, args, err = queryUpdate.ToSql()
	if err != nil {
		return nil, fmt.Errorf("failed to build query: %v", err)
	}

	_, err = s.pool.Exec(ctx, sql, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to update user: %v", err)
	}

	return &emptypb.Empty{}, nil
}

// DeleteUser Удаление существующего пользователя
func (s *server) DeleteUser(ctx context.Context, req *desc.DeleteUserRequest) (*emptypb.Empty, error) {

	querySelect := sq.Select("1").
		From("users").
		PlaceholderFormat(sq.Dollar).
		Where(sq.Eq{"id": req.GetId()}).
		Limit(1)

	sql, args, err := querySelect.ToSql()
	if err != nil {
		return nil, fmt.Errorf("failed to build query: %v", err)
	}

	var exist int

	err = s.pool.QueryRow(ctx, sql, args...).Scan(&exist)
	if err != nil {
		if err.Error() == errNoRows {
			return nil, fmt.Errorf("user: %d not founded", req.GetId())
		}
		return nil, fmt.Errorf("failed to execute query: %v", err)
	}

	queryDelete := sq.Delete("users").PlaceholderFormat(sq.Dollar).Where(sq.Eq{"id": req.GetId()})

	sql, args, err = queryDelete.ToSql()
	_, err = s.pool.Exec(ctx, sql, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to execute query: %v", err)
	}

	return &emptypb.Empty{}, nil
}

func main() {

	flag.Parse()
	ctx := context.Background()

	err := config.Load(configPath)
	if err != nil {
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

	pool, err := pgxpool.Connect(ctx, pgConfig.DSN())
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}
	defer pool.Close()

	s := grpc.NewServer()
	reflection.Register(s)
	desc.RegisterAuthServiceServer(s, &server{pool: pool})

	log.Printf("server listening at %v", lis.Addr())

	if err = s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
