package config

import (
	"errors"
	"os"
)

type Config struct {
	DatabaseURL string
	HTTPAddr    string
	APIKey      string
}

func Load() (Config, error) {
	c := Config{DatabaseURL: os.Getenv("DATABASE_URL"), HTTPAddr: os.Getenv("HTTP_ADDR"), APIKey: os.Getenv("API_KEY")}
	if c.DatabaseURL == "" {
		return Config{}, errors.New("DATABASE_URL is required")
	}
	if c.HTTPAddr == "" {
		c.HTTPAddr = ":8080"
	}
	return c, nil
}
