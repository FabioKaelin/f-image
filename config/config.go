package config

import (
	"github.com/joho/godotenv"
)

type Config struct {
	// DatabaseHost     string `mapstructure:"DATABASE_HOST"`
}

var StartConfig Config

func LoadConfig(path string) (config Config, err error) {
	godotenv.Load("app.env")

	config = Config{
		// DatabaseHost:       os.Getenv("DATABASE_HOST"),
	}

	return
}
