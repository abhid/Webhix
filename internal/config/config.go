package config

import (
	"os"

	"github.com/ilyakaznacheev/cleanenv"
	"github.com/joho/godotenv"
)

type Config struct {
	Addr    string `env:"WEBHIX_ADDR"     env-default:":8080"`
	BaseURL string `env:"WEBHIX_BASE_URL" env-default:"http://localhost:8080"`
	DBPath  string `env:"WEBHIX_DB_PATH"  env-default:"./data"`

	Password  string `env:"WEBHIX_PASSWORD"`
	SecretKey string `env:"WEBHIX_SECRET_KEY"`

	TrustedProxies []string `env:"WEBHIX_TRUSTED_PROXIES"`
}

func LoadConfig() (*Config, error) {
	var cfg Config

	if err := godotenv.Load("config/.env"); err != nil && !os.IsNotExist(err) {
		return nil, err
	}

	if err := cleanenv.ReadEnv(&cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}
