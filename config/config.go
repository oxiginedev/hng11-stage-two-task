package config

import (
	"math/rand"
	"os"
	"strings"

	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
)

type Configuration struct {
	Port uint16
	DB   DatabaseConfiguration
	JWT  JWTConfiguration
}

type DatabaseConfiguration struct {
	Driver   string
	DSN      string
}

type JWTConfiguration struct {
	Secret string
	Expiry int64
}

func loadEnvironment(f string) error {
	var err error
	if len(strings.TrimSpace(f)) != 0 {
		err = godotenv.Overload(f)
	} else {
		err = godotenv.Load()

		if os.IsNotExist(err) {
			return nil
		}
	}

	return err
}

func Load(f string) (*Configuration, error) {
	if err := loadEnvironment(f); err != nil {
		return nil, err
	}

	config := new(Configuration)

	if err := envconfig.Process("", config); err != nil {
		return nil, err
	}

	if config.JWT.Secret == "" {
		config.JWT.Secret = randomString(32)
	}

	return config, nil
}

func randomString(len int) string {
	bytes := make([]byte, len)
	for i := 0; i < len; i++ {
		bytes[i] = byte(65 + rand.Intn(25)) //A=65 and Z = 65+25
	}

	return string(bytes)
}
