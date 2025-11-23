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
	Host     string `env:"DB_HOST" env-default:"localhost"`
	Port     string `env:"DB_PORT" env-default:"5432"`
	User     string `env:"DB_USER" env-default:"pr_review_user"`
	Password string `env:"DB_PASSWORD" env-default:"pr_review_password"`
	Database string `env:"DB_NAME" env-default:"pr_review_db"`
	SSLMode  string `env:"DB_SSL_MODE" env-default:"disable"`

	MaxOpenConns    int           `env:"DB_MAX_OPEN_CONNS" env-default:"25"`
	MaxIdleConns    int           `env:"DB_MAX_IDLE_CONNS" env-default:"5"`
	ConnMaxLifetime time.Duration `env:"DB_CONN_MAX_LIFETIME" env-default:"1h"`

	InitTimeout time.Duration `env:"DB_INIT_TIMEOUT" env-default:"10s"`
	PingTimeout time.Duration `env:"DB_PING_TIMEOUT" env-default:"5s"`

	Path string `env:"DB_PATH" env-default:""`
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
