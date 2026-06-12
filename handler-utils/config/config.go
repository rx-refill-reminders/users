package config

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/caarlos0/env/v11"
)

type Config struct {
	Environment

	AWSConfig aws.Config
}

type Environment struct {
	UsersTable string `env:"USERS_TABLE"`
}

func Load(ctx context.Context) (*Config, error) {
	cfg := Config{}

	// Read environment variables
	err := env.Parse(&cfg.Environment)
	if err != nil {
		return nil, fmt.Errorf("error parsing environment variables: %w", err)
	}

	// Load default AWS config
	cfg.AWSConfig, err = config.LoadDefaultConfig(ctx)
	if err != nil {
		return nil, fmt.Errorf("error loading AWS config: %w", err)
	}

	return &cfg, nil
}
