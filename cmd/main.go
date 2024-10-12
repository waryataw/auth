package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net"
	"time"

	desc "github.com/waryataw/auth/pkg/user_v1"

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
	desc.UnimplementedUserV1Server
	pool *pgxpool.Pool
}

type userQuery struct {
	ID   int64
	name string
}

func getUserQueryBuilderByQuery(uq userQuery) sq.SelectBuilder {

	// Строим запрос с использованием Squirrel
	bQ := sq.Select("id", "name", "email", "role", "created_at", "updated_at").
		From("users").
		PlaceholderFormat(sq.Dollar).
		Limit(1)

	if uq.ID > 0 {
		bQ = bQ.Where(sq.Eq{"id": uq.ID})
	}

	if uq.name != "" {
		bQ = bQ.Where(sq.Like{"name": uq.name})
	}

	return bQ
}

// Get Получение существующего пользователя
func (s *server) Get(ctx context.Context, req *desc.GetRequest) (*desc.GetResponse, error) {
	query, args, err := getUserQueryBuilderByQuery(userQuery{ID: req.GetId()}).ToSql()
	if err != nil {
		return nil, fmt.Errorf("failed to build query: %v", err)
	}

	// Объявляем переменные для сканирования данных из базы
	var id int64
	var roleID int
	var name, email string
	var createdAt, updatedAt *time.Time

	// Выполняем запрос и сканируем результат
	err = s.pool.QueryRow(ctx, query, args...).Scan(&id, &name, &email, &roleID, &createdAt, &updatedAt)
	if err != nil {
		return nil, fmt.Errorf("failed to select user: %v", err)
	}

	// Получаем роль по идентификатору Роли
	role, err := roleByRoleID(roleID)
	if err != nil {
		return nil, fmt.Errorf("failed to select user: %v", err)
	}

	// Возвращаем ответ в формате gRPC
	return &desc.GetResponse{
		User: &desc.User{
			Id:        id,
			Name:      name,
			Email:     email,
			Role:      *role,
			CreatedAt: convertTimeToProtoTime(createdAt),
			UpdatedAt: convertTimeToProtoTime(updatedAt),
		},
	}, nil
}

// GetByName Get Получение существующего пользователя по Имени
func (s *server) GetByName(ctx context.Context, req *desc.GetByNameRequest) (*desc.GetByNameResponse, error) {
	query, args, err := getUserQueryBuilderByQuery(userQuery{name: req.GetName()}).ToSql()
	if err != nil {
		return nil, fmt.Errorf("failed to build query: %v", err)
	}

	// Объявляем переменные для сканирования данных из базы
	var id int64
	var roleID int
	var name, email string
	var createdAt, updatedAt *time.Time

	// Выполняем запрос и сканируем результат
	err = s.pool.QueryRow(ctx, query, args...).Scan(&id, &name, &email, &roleID, &createdAt, &updatedAt)
	if err != nil {
		return nil, fmt.Errorf("failed to select user: %v", err)
	}

	// Получаем роль по идентификатору Роли
	role, err := roleByRoleID(roleID)
	if err != nil {
		return nil, fmt.Errorf("failed to select user: %v", err)
	}

	// Возвращаем ответ в формате gRPC
	return &desc.GetByNameResponse{
		User: &desc.User{
			Id:        id,
			Name:      name,
			Email:     email,
			Role:      *role,
			CreatedAt: convertTimeToProtoTime(createdAt),
			UpdatedAt: convertTimeToProtoTime(updatedAt),
		},
	}, nil
}

// roleByRoleID Получение роли по идентификатору роли, если идентификатор вне диапазона возвращаем ошибку
func roleByRoleID(id int) (*desc.Role, error) {

	// Определяем возможные роли
	roles := []desc.Role{
		desc.Role_UNKNOWN,
		desc.Role_USER,
		desc.Role_ADMIN,
	}

	// Проверяем границы значения role
	if id < 0 || id >= len(roles) {
		return nil, fmt.Errorf("invalid role value: %d", id)
	}

	return &roles[id], nil
}

// convertTimeToProtoTime Конвертируем Time для ответа в формате gRPC или null
func convertTimeToProtoTime(t *time.Time) *timestamppb.Timestamp {
	if t == nil {
		return nil
	}

	return timestamppb.New(*t)
}

// Create Добавление нового пользователя
func (s *server) Create(ctx context.Context, req *desc.CreateRequest) (*desc.CreateResponse, error) {
	// Делаем запрос на вставку записи в таблицу пользователя
	builderInsert := sq.Insert("users").
		PlaceholderFormat(sq.Dollar).
		Columns("name", "email", "password", "password_confirm", "role").
		Values(req.Name, req.Email, req.Password, req.PasswordConfirm, req.Role).
		Suffix("RETURNING id")

	query, args, err := builderInsert.ToSql()
	if err != nil {
		return nil, fmt.Errorf("failed to build query: %v", err)
	}

	var userID int64
	err = s.pool.QueryRow(ctx, query, args...).Scan(&userID)
	if err != nil {
		return nil, fmt.Errorf("failed to insert user: %v", err)
	}

	log.Printf("inserted user with id: %d", userID)

	return &desc.CreateResponse{
		Id: userID,
	}, nil
}

// Update Изменения существующего пользователя
func (s *server) Update(ctx context.Context, req *desc.UpdateRequest) (*emptypb.Empty, error) {

	bqs := sq.Select("1").
		From("users").
		PlaceholderFormat(sq.Dollar).
		Where(sq.Eq{"id": req.GetId()}).
		Limit(1)

	query, args, err := bqs.ToSql()
	if err != nil {
		return nil, fmt.Errorf("failed to build query: %v", err)
	}

	var exist int

	err = s.pool.QueryRow(ctx, query, args...).Scan(&exist)
	if err != nil {
		if err.Error() == errNoRows {
			return nil, fmt.Errorf("user: %d not founded", req.GetId())
		}
		return nil, fmt.Errorf("failed to execute query: %v", err)
	}

	bqu := sq.Update("users").
		PlaceholderFormat(sq.Dollar).
		Set("name", req.GetName()).
		Set("email", req.GetEmail()).
		Set("role", req.GetRole()).
		Set("updated_at", time.Now().UTC()).
		Where(sq.Eq{"id": req.GetId()})

	query, args, err = bqu.ToSql()
	if err != nil {
		return nil, fmt.Errorf("failed to build query: %v", err)
	}

	_, err = s.pool.Exec(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to update user: %v", err)
	}

	return &emptypb.Empty{}, nil
}

// Delete Удаление существующего пользователя
func (s *server) Delete(ctx context.Context, req *desc.DeleteRequest) (*emptypb.Empty, error) {

	bqs := sq.Select("1").
		From("users").
		PlaceholderFormat(sq.Dollar).
		Where(sq.Eq{"id": req.GetId()}).
		Limit(1)

	query, args, err := bqs.ToSql()
	if err != nil {
		return nil, fmt.Errorf("failed to build query: %v", err)
	}

	var exist int

	err = s.pool.QueryRow(ctx, query, args...).Scan(&exist)
	if err != nil {
		if err.Error() == errNoRows {
			return nil, fmt.Errorf("user: %d not founded", req.GetId())
		}
		return nil, fmt.Errorf("failed to execute query: %v", err)
	}

	bqd := sq.Delete("users").PlaceholderFormat(sq.Dollar).Where(sq.Eq{"id": req.GetId()})

	query, args, err = bqd.ToSql()
	_, err = s.pool.Exec(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to execute query: %v", err)
	}

	return &emptypb.Empty{}, nil
}

func main() {

	flag.Parse()
	ctx := context.Background()

	// Считываем переменные окружения
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

	// Создаем пул соединений с базой данных
	pool, err := pgxpool.Connect(ctx, pgConfig.DSN())
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}
	defer pool.Close()

	s := grpc.NewServer()
	reflection.Register(s)
	desc.RegisterUserV1Server(s, &server{pool: pool})

	log.Printf("server listening at %v", lis.Addr())

	if err = s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
