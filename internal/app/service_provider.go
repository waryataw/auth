package app

import (
	"context"
	"log"

	"github.com/waryataw/auth/internal/api/auth"

	"github.com/waryataw/auth/internal/client/db"
	"github.com/waryataw/auth/internal/client/db/pg"
	"github.com/waryataw/auth/internal/client/db/transaction"
	"github.com/waryataw/auth/internal/closer"
	"github.com/waryataw/auth/internal/config"
	"github.com/waryataw/auth/internal/repository"
	userRepository "github.com/waryataw/auth/internal/repository/user"
	"github.com/waryataw/auth/internal/service"
	userService "github.com/waryataw/auth/internal/service/user"
)

type serviceProvider struct {
	pgConfig   config.PGConfig
	grpcConfig config.GRPCConfig

	dbClient       db.Client
	txManager      db.TxManager
	userRepository repository.UserRepository

	userService service.UserService

	userImpl *auth.Implementation
}

func newServiceProvider() *serviceProvider {
	return &serviceProvider{}
}

func (s *serviceProvider) PGConfig() config.PGConfig {
	if s.pgConfig == nil {
		cfg, err := config.NewPGConfig()
		if err != nil {
			log.Fatalf("failed to get pg config: %s", err.Error())
		}

		s.pgConfig = cfg
	}

	return s.pgConfig
}

func (s *serviceProvider) GRPCConfig() config.GRPCConfig {
	if s.grpcConfig == nil {
		cfg, err := config.NewGRPCConfig()
		if err != nil {
			log.Fatalf("failed to get grpc config: %s", err.Error())
		}

		s.grpcConfig = cfg
	}

	return s.grpcConfig
}

func (s *serviceProvider) DBClient(ctx context.Context) db.Client {
	if s.dbClient == nil {
		cl, err := pg.New(ctx, s.PGConfig().DSN())
		if err != nil {
			log.Fatalf("failed to create db client: %v", err)
		}

		err = cl.DB().Ping(ctx)
		if err != nil {
			log.Fatalf("ping error: %s", err.Error())
		}
		closer.Add(cl.Close)

		s.dbClient = cl
	}

	return s.dbClient
}

func (s *serviceProvider) TxManager(ctx context.Context) db.TxManager {
	if s.txManager == nil {
		s.txManager = transaction.NewTransactionManager(s.DBClient(ctx).DB())
	}

	return s.txManager
}

func (s *serviceProvider) UserRepository(ctx context.Context) repository.UserRepository {
	if s.userRepository == nil {
		s.userRepository = userRepository.NewRepository(s.DBClient(ctx))
	}

	return s.userRepository
}

func (s *serviceProvider) UserService(ctx context.Context) service.UserService {
	if s.userService == nil {
		s.userService = userService.NewService(
			s.UserRepository(ctx),
		)
	}

	return s.userService
}

func (s *serviceProvider) AuthImpl(ctx context.Context) *auth.Implementation {
	if s.userImpl == nil {
		s.userImpl = auth.NewImplementation(s.UserService(ctx))
	}

	return s.userImpl
}
