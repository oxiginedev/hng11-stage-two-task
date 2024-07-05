package config

import (
	"os"

	"github.com/joho/godotenv"
)

type Configuration struct {
	Port uint16
	DB   DatabaseConfiguration
	JWT  JWTConfiguration
}

type DatabaseConfiguration struct {
	Driver   string
	DSN      string
	Host     string
	Port     uint16
	Username string
	Password string
	Database string
}

type JWTConfiguration struct {
	Secret string
	Expiry int64
}

func LoadConfig() *Configuration {
	godotenv.Load()

	return &Configuration{
		Port: env("PORT", 9000).(uint16),
		DB: DatabaseConfiguration{
			Driver:   env("DB_DRIVER", "postgres").(string),
			DSN:      env("DB_DSN", "").(string),
			Host:     env("DB_HOST", "localhost").(string),
			Port:     env("DB_PORT", 5432).(uint16),
			Username: env("DB_USERNAME", "root").(string),
			Password: env("DB_PASSWORD", "").(string),
			Database: env("DB_DATABASE", "").(string),
		},
		JWT: JWTConfiguration{
			Secret: env("JWT_SECRET", "").(string),
			Expiry: env("JWT_EXPIRY", 3600).(int64),
		},
	}
}

func env(key string, fallback any) any {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}

	return fallback
}
