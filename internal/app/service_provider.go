package app

import (
	"context"
	"log"

	"github.com/IBM/sarama"
	userApi "github.com/waryataw/auth/internal/api/user"
	"github.com/waryataw/auth/internal/config"
	"github.com/waryataw/auth/internal/config/env"
	userRepository "github.com/waryataw/auth/internal/repository/user"
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

	dbClient       db.Client
	txManager      db.TxManager
	userRepository userService.Repository

	userService userApi.MainService

	userSaverConsumer serviceConsumer.Service

	consumer             kafka.Consumer
	consumerGroup        sarama.ConsumerGroup
	consumerGroupHandler *consumer.GroupHandler

	controller *userApi.Controller
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

func (s *serviceProvider) UserSaverConsumer(ctx context.Context) serviceConsumer.Service {
	if s.userSaverConsumer == nil {
		s.userSaverConsumer = user_saver.NewService(
			s.UserRepository(ctx),
			s.Consumer(),
		)
	}

	return s.userSaverConsumer
}

func (s *serviceProvider) UserService(ctx context.Context) userApi.MainService {
	if s.userService == nil {
		s.userService = userService.NewService(
			s.UserRepository(ctx),
		)
	}

	return s.userService
}

func (s *serviceProvider) AuthController(ctx context.Context) *userApi.Controller {
	if s.controller == nil {
		s.controller = userApi.NewController(s.UserService(ctx))
	}

	return s.controller
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
