package config

import (
	"github.com/joho/godotenv"
)

// GRPCConfig GRPC config
type GRPCConfig interface {
	Address() string
}

// HTTPConfig Http config
type HTTPConfig interface {
	Address() string
}

// PGConfig Postgres config
type PGConfig interface {
	DSN() string
}

// SwaggerConfig Swagger config
type SwaggerConfig interface {
	Address() string
}

// Load Configs.
func Load(path string) error {
	err := godotenv.Load(path)
	if err != nil {
		return err
	}

	return nil
}
