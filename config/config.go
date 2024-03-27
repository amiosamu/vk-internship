package config

import (
	"fmt"
	"path"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

type (
	Config struct {
		App  `yaml:"app"`
		HTTP `yaml:"http"`
		Log  `yaml:"log"`
		PG   `yaml:"postgres"`
		JWT  `yaml:"jwt"`
	}

	App struct {
		Name    string `env-required:"true" yaml:"name"    env:"APP_NAME"`
		Version string `env-required:"true" yaml:"version" env:"APP_VERSION"`
	}

	HTTP struct {
		Port string `env-required:"true" yaml:"port" env:"HTTP_PORT"`
	}

	Log struct {
		Level string `env-required:"true" yaml:"level" env:"LOG_LEVEL"`
	}

	PG struct {
		MaxPoolSize int    `env-required:"true" yaml:"max_pool_size" env:"PG_MAX_POOL_SIZE"`
		URL         string `env-required:"true"                      env:"PG_URL_LOCALHOST"`
	}

	JWT struct {
		SignKey  string        `env-required:"true"                  env:"JWT_SIGN_KEY"`
		TokenTTL time.Duration `env-required:"true" yaml:"token_ttl" env:"JWT_TOKEN_TTL"`
	}
)

func NewConfig(configPath string) (*Config, error) {
	cfg := &Config{}

	err := cleanenv.ReadConfig(path.Join("./", configPath), cfg)
	if err != nil {
		return nil, fmt.Errorf("error reading config file: %w", err)
	}

	err = cleanenv.UpdateEnv(cfg)
	if err != nil {
		return nil, fmt.Errorf("error updating env: %w", err)
	}

	return cfg, nil
}
