package config

import (
	"errors"
	"os"
)

type Config struct {
	HTTP_ADDR           string
	MONGO_URI           string
	MONGO_DATABASE_NAME string
}

// Read reads config from environment.
func Read() (*Config, error) {
	config := Config{}
	HTTP_ADDR, exists := os.LookupEnv("HTTP_ADDR")
	if exists {
		config.HTTP_ADDR = HTTP_ADDR
	} else {
		return nil, errors.New("HTTP_ADDR is not set")
	}

	MONGO_URI, exists := os.LookupEnv("MONGO_URI")
	if exists {
		config.MONGO_URI = MONGO_URI
	} else {
		return nil, errors.New("MONGO_URI is not set")
	}

	MONGO_DATABASE_NAME, exists := os.LookupEnv("MONGO_DATABASE_NAME")
	if exists {
		config.MONGO_DATABASE_NAME = MONGO_DATABASE_NAME
	} else {
		return nil, errors.New("MONGO_DATABASE_NAME is not set")
	}

	return &config, nil
}
