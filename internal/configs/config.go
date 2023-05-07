package configs

import "github.com/caarlos0/env/v7"

type Config struct {
	AppPort     string `env:"APP_PORT" envDefault:"9000"`
	PostgresUrl string `env:"POSTGRES_URL"`
}

func NewConfig() (*Config, error) {
	cfg := &Config{}
	err := env.Parse(cfg)
	if err != nil {
		return nil, err
	}
	return cfg, nil
}
