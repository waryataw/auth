package app

import (
	"context"
	"log"

	"github.com/IBM/sarama"
	accessApi "github.com/waryataw/auth/internal/api/access"
	authApi "github.com/waryataw/auth/internal/api/auth"
	userApi "github.com/waryataw/auth/internal/api/user"
	"github.com/waryataw/auth/internal/config"
	"github.com/waryataw/auth/internal/config/env"
	authRepository "github.com/waryataw/auth/internal/repository/auth"
	userRepository "github.com/waryataw/auth/internal/repository/user"
	accessService "github.com/waryataw/auth/internal/service/access"
	authService "github.com/waryataw/auth/internal/service/auth"
	serviceConsumer "github.com/waryataw/auth/internal/service/consumer"
	"github.com/waryataw/auth/internal/service/consumer/user_saver"
	userService "github.com/waryataw/auth/internal/service/user"
	"github.com/waryataw/platform_common/pkg/closer"
	"github.com/waryataw/platform_common/pkg/db"
	"github.com/waryataw/platform_common/pkg/db/pg"
	"github.com/waryataw/platform_common/pkg/db/transaction"
	"github.com/waryataw/platform_common/pkg/kafka"
	"github.com/waryataw/platform_common/pkg/kafka/consumer"
)

type serviceProvider struct {
	pgConfig            config.PGConfig
	grpcConfig          config.GRPCConfig
	httpConfig          config.HTTPConfig
	kafkaConsumerConfig config.KafkaConsumerConfig
	swaggerConfig       config.SwaggerConfig
	refreshTokenConfig  config.RefreshTokenConfig

	dbClient       db.Client
	txManager      db.TxManager
	userRepository userService.Repository
	authRepository authService.Repository

	userService   userApi.Service
	authService   authApi.Service
	accessService accessApi.Service

	userSaverConsumer serviceConsumer.Service

	consumer             kafka.Consumer
	consumerGroup        sarama.ConsumerGroup
	consumerGroupHandler *consumer.GroupHandler

	userController   *userApi.Controller
	authController   *authApi.Controller
	accessController *accessApi.Controller
}

func newServiceProvider() *serviceProvider {
	return &serviceProvider{}
}

func (s *serviceProvider) PGConfig() config.PGConfig {
	if s.pgConfig == nil {
		cfg, err := env.NewPGConfig()
		if err != nil {
			log.Fatalf("failed to get pg config: %s", err.Error())
		}

		s.pgConfig = cfg
	}

	return s.pgConfig
}

func (s *serviceProvider) GRPCConfig() config.GRPCConfig {
	if s.grpcConfig == nil {
		cfg, err := env.NewGRPCConfig()
		if err != nil {
			log.Fatalf("failed to get grpc config: %s", err.Error())
		}

		s.grpcConfig = cfg
	}

	return s.grpcConfig
}

func (s *serviceProvider) HTTPConfig() config.HTTPConfig {
	if s.httpConfig == nil {
		cfg, err := env.NewHTTPConfig()
		if err != nil {
			log.Fatalf("failed to get http config: %s", err.Error())
		}

		s.httpConfig = cfg
	}

	return s.httpConfig
}

func (s *serviceProvider) KafkaConsumerConfig() config.KafkaConsumerConfig {
	if s.kafkaConsumerConfig == nil {
		cfg, err := env.NewKafkaConsumerConfig()
		if err != nil {
			log.Fatalf("failed to get kafka consumer config: %s", err.Error())
		}

		s.kafkaConsumerConfig = cfg
	}

	return s.kafkaConsumerConfig
}

func (s *serviceProvider) SwaggerConfig() config.SwaggerConfig {
	if s.swaggerConfig == nil {
		cfg, err := env.NewSwaggerConfig()
		if err != nil {
			log.Fatalf("failed to get swagger config: %s", err.Error())
		}

		s.swaggerConfig = cfg
	}

	return s.swaggerConfig
}

func (s *serviceProvider) RefreshTokenConfig() config.RefreshTokenConfig {
	if s.refreshTokenConfig == nil {
		cfg, err := env.NewRefreshTokenConfig()
		if err != nil {
			log.Fatalf("failed to get refresh token config: %s", err.Error())
		}

		s.refreshTokenConfig = cfg
	}

	return s.refreshTokenConfig
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

func (s *serviceProvider) UserRepository(ctx context.Context) userService.Repository {
	if s.userRepository == nil {
		s.userRepository = userRepository.NewRepository(s.DBClient(ctx))
	}

	return s.userRepository
}

func (s *serviceProvider) AuthRepository(_ context.Context) authService.Repository {
	if s.authRepository == nil {
		s.authRepository = authRepository.NewRepository(s.RefreshTokenConfig())
	}

	return s.authRepository
}

func (s *serviceProvider) UserSaverConsumer(ctx context.Context) serviceConsumer.Service {
	if s.userSaverConsumer == nil {
		s.userSaverConsumer = user_saver.NewService(
			s.UserRepository(ctx),
			s.Consumer(),
		)
	}

	return s.userSaverConsumer
}

func (s *serviceProvider) UserService(ctx context.Context) userApi.Service {
	if s.userService == nil {
		s.userService = userService.NewService(
			s.UserRepository(ctx),
		)
	}

	return s.userService
}

func (s *serviceProvider) AuthService(ctx context.Context) authApi.Service {
	if s.authService == nil {
		s.authService = authService.NewService(s.AuthRepository(ctx), s.UserRepository(ctx))
	}

	return s.authService
}

func (s *serviceProvider) AccessService(_ context.Context) accessApi.Service {
	if s.accessService == nil {
		s.accessService = accessService.NewService()
	}

	return s.accessService
}

func (s *serviceProvider) UserController(ctx context.Context) *userApi.Controller {
	if s.userController == nil {
		s.userController = userApi.NewController(s.UserService(ctx))
	}

	return s.userController
}

func (s *serviceProvider) AuthController(ctx context.Context) *authApi.Controller {
	if s.authController == nil {
		s.authController = authApi.NewController(s.AuthService(ctx))
	}

	return s.authController
}

func (s *serviceProvider) AccessController(ctx context.Context) *accessApi.Controller {
	if s.accessController == nil {
		s.accessController = accessApi.NewController(s.AccessService(ctx))
	}

	return s.accessController
}

func (s *serviceProvider) Consumer() kafka.Consumer {
	if s.consumer == nil {
		s.consumer = consumer.NewConsumer(
			s.ConsumerGroup(),
			s.ConsumerGroupHandler(),
		)
		closer.Add(s.consumer.Close)
	}

	return s.consumer
}

func (s *serviceProvider) ConsumerGroup() sarama.ConsumerGroup {
	if s.consumerGroup == nil {
		consumerGroup, err := sarama.NewConsumerGroup(
			s.KafkaConsumerConfig().Brokers(),
			s.KafkaConsumerConfig().GroupID(),
			s.KafkaConsumerConfig().Config(),
		)
		if err != nil {
			log.Fatalf("failed to create consumer group: %v", err)
		}

		s.consumerGroup = consumerGroup
	}

	return s.consumerGroup
}

func (s *serviceProvider) ConsumerGroupHandler() *consumer.GroupHandler {
	if s.consumerGroupHandler == nil {
		s.consumerGroupHandler = consumer.NewGroupHandler()
	}

	return s.consumerGroupHandler
}
