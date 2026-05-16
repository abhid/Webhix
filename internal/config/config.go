package config

import (
	"os"

	"github.com/ilyakaznacheev/cleanenv"
	"github.com/joho/godotenv"
)

type Config struct {
	SQLiteDBPath string `env:"SQLITE_DB_PATH"`
	WebHixAddr   string `env:"WEBHIX_ADDR"`
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
