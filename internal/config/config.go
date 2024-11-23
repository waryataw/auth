package config

import (
	"github.com/IBM/sarama"
	"github.com/joho/godotenv"
)

// GRPCConfig GRPC config.
type GRPCConfig interface {
	Address() string
}

// HTTPConfig Http config.
type HTTPConfig interface {
	Address() string
}

// PGConfig Postgres config.
type PGConfig interface {
	DSN() string
}

// SwaggerConfig Swagger config.
type SwaggerConfig interface {
	Address() string
}

// KafkaConsumerConfig Kafka consumer config.
type KafkaConsumerConfig interface {
	Brokers() []string
	GroupID() string
	Config() *sarama.Config
}

// AuthConfig Refresh token config.
type AuthConfig interface {
	AuthPrefix() string
	RefreshTokenSecretKey() string
	RefreshTokenExpirationMinutes() int64
	AccessTokenSecretKey() string
	AccessTokenExpirationMinutes() int64
}

// Load Configs.
func Load(path string) error {
	err := godotenv.Load(path)
	if err != nil {
		return err
	}

	return nil
}
