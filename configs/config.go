package configs

import (
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	Db DbConfig
}

type DbConfig struct {
	// Host     string
	// Port     string
	// User     string
	// Password string
	// DbName   string
	DSN string
}

func LoadConfig() *Config {
	err := godotenv.Load()
	if err != nil {
		panic(err)
	}

	DSN := os.Getenv("DSN")
	return &Config{
		Db: DbConfig{
			DSN: os.Getenv("DSN"),
		}}
}