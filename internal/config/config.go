package config

import "github.com/caarlos0/env/v7"

type Config struct {
	AppPort     string `env:"APP_PORT" envDefault:"9000"`
	PostgresUrl string `env:"POSTGRES_URL"`
	DirAddress  string `env:"BIOCAD_DIR_ADDRESS"`
}

func NewConfig() (*Config, error) {
	cfg := &Config{}
	err := env.Parse(cfg)

	//for testing purposes
	cfg.PostgresUrl = "postgresql://biocad:password@localhost:5432/biocad"

	if err != nil {
		return nil, err
	}
	return cfg, nil
}
