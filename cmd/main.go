package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net"
	"time"

	sq "github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/waryataw/auth/internal/config"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/timestamppb"

	desc "github.com/waryataw/auth/pkg/user_v1"
)

var configPath string

func init() {
	flag.StringVar(&configPath, "config-path", ".env", "path to config file")
}

type server struct {
	desc.UnimplementedUserV1Server
	pool *pgxpool.Pool
}

func (s *server) Get(ctx context.Context, req *desc.GetRequest) (*desc.GetResponse, error) {
	// Строим SQL-запрос с использованием Squirrel
	bQ := sq.Select("id", "name", "email", "role", "created_at", "updated_at").
		From("users").
		PlaceholderFormat(sq.Dollar).
		Where(sq.Eq{"id": req.GetId()}).
		Limit(1)

	// Генерируем SQL-запрос и аргументы
	query, args, err := bQ.ToSql()
	if err != nil {
		return nil, fmt.Errorf("failed to build query: %v", err)
	}

	// Объявляем переменные для сканирования данных из базы
	var id int64
	var role int
	var name, email string
	var createdAt, updatedAt *time.Time

	// Выполняем запрос и сканируем результат
	err = s.pool.QueryRow(ctx, query, args...).Scan(&id, &name, &email, &role, &createdAt, &updatedAt)
	if err != nil {
		return nil, fmt.Errorf("failed to select user: %v", err)
	}

	// Определяем возможные роли
	roles := []desc.Role{
		desc.Role_USER,
		desc.Role_ADMIN,
	}

	// Проверяем границы значения role
	if role < 0 || role >= len(roles) {
		return nil, fmt.Errorf("invalid role value: %d", role)
	}

	// Возвращаем ответ в формате gRPC
	return &desc.GetResponse{
		Id:        id,
		Name:      name,
		Email:     email,
		Role:      roles[role], // Преобразование значения role в enum
		CreatedAt: convertTimeToProtoTime(createdAt),
		UpdatedAt: convertTimeToProtoTime(updatedAt),
	}, nil
}

func convertTimeToProtoTime(t *time.Time) *timestamppb.Timestamp {
	if t == nil {
		return nil
	}

	return timestamppb.New(*t)
}

func (s *server) Create(ctx context.Context, req *desc.CreateRequest) (*desc.CreateResponse, error) {
	// Делаем запрос на вставку записи в таблицу пользователя
	builderInsert := sq.Insert("users").
		PlaceholderFormat(sq.Dollar).
		Columns("name", "email", "password", "password_confirm", "role").
		Values(req.Name, req.Email, req.Password, req.PasswordConfirm, req.Role).
		Suffix("RETURNING id")

	query, args, err := builderInsert.ToSql()
	if err != nil {
		log.Fatalf("failed to build query: %v", err)
	}

	var userID int64
	err = s.pool.QueryRow(ctx, query, args...).Scan(&userID)
	if err != nil {
		log.Fatalf("failed to insert user: %v", err)
	}

	log.Printf("inserted user with id: %d", userID)

	return &desc.CreateResponse{
		Id: userID,
	}, nil
}

func (s *server) Update(_ context.Context, _ *desc.UpdateRequest) (*emptypb.Empty, error) {
	return &emptypb.Empty{}, nil
}

func (s *server) Delete(_ context.Context, _ *desc.DeleteRequest) (*emptypb.Empty, error) {
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
