package envconfig

import "github.com/caarlos0/env/v11"

type Config struct {
	AWSRegion  string `env:"AWS_REGION"`
	UsersTable string `env:"USERS_TABLE"`
}

func Load() (*Config, error) {
	cfg := Config{}
	err := env.Parse(&cfg)
	if err != nil {
		return nil, err
	}

	return &cfg, nil
}
