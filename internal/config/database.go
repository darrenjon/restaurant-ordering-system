package config

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

type DatabaseConfig struct {
	Host      string
	Port      string
	User      string
	Password  string
	DBName    string
	SSLMode   string
	JWTSecret string
}

func LoadDatabaseConfig() (*DatabaseConfig, error) {
	err := godotenv.Load()
	if err != nil {
		return nil, fmt.Errorf("error loading .env file: %w", err)
	}

	return &DatabaseConfig{
		Host:      os.Getenv("DB_HOST"),
		Port:      os.Getenv("DB_PORT"),
		User:      os.Getenv("DB_USER"),
		Password:  os.Getenv("DB_PASSWORD"),
		DBName:    os.Getenv("DB_NAME"),
		SSLMode:   os.Getenv("DB_SSLMODE"),
		JWTSecret: os.Getenv("JWT_SECRET"),
	}, nil
}
