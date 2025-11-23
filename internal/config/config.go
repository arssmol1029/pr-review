package config

import (
	"log"
	"log/slog"
	"os"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
	"github.com/joho/godotenv"
)

const (
	envLocal = "local"
	envDev   = "dev"
	envProd  = "prod"
)

type Config struct {
	Env        string `env:"ENV" env-default:"local"`
	LogLevel   string `env:"LOG_LEVEL" env-default:"info"`
	HTTPServer HTTPServerConfig
	Database   DatabaseConfig
}

type HTTPServerConfig struct {
	Address     string        `env:"HTTP_ADDRESS" env-default:":8080"`
	Timeout     time.Duration `env:"HTTP_TIMEOUT" env-default:"10s"`
	IdleTimeout time.Duration `env:"HTTP_IDLE_TIMEOUT" env-default:"60s"`
}

type DatabaseConfig struct {
	Path        string        `env:"DB_PATH" env-default:"storage.db"`
	InitTimeout time.Duration `env:"DB_INIT_TIMEOUT" env-default:"10s"`
	PingTimeout time.Duration `env:"DB_PING_TIMEOUT" env-default:"5s"`
}

func MustLoad() *Config {
	if _, err := os.Stat(".env-default"); err == nil {
		if err := godotenv.Load(".env-default"); err != nil {
			log.Fatalf("cannot load .env-default: %s", err)
		}
	}

	if _, err := os.Stat(".env"); err == nil {
		if err := godotenv.Load(".env"); err != nil {
			log.Fatalf("cannot load .env: %s", err)
		}
	}

	var cfg Config

	if err := cleanenv.ReadEnv(&cfg); err != nil {
		log.Fatalf("cannot read config from environment: %s", err)
	}

	return &cfg
}

func (c *Config) GetSlogLevel() slog.Level {
	switch c.Env {
	case envLocal:
		return slog.LevelDebug
	case envDev:
		return slog.LevelDebug
	case envProd:
		return slog.LevelInfo
	default:
		return slog.LevelInfo
	}
}
