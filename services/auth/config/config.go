package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	DebugMode bool
	Port      string
	Host      string
	Database  struct {
		Host     string
		Port     string
		User     string
		Password string
		Name     string
	}
	JWTSecret     string
	JWTExpiration string
	Redis         struct {
		Host     string
		Port     string
		Password string
		DB       string
	}
}

func LoadConfig() *Config {

	err := godotenv.Load("services/auth/.env")

	if err != nil {
		log.Fatalf("Error loading .env file")
	}

	return &Config{
		DebugMode: os.Getenv("DEBUG_MODE") == "true",
		Port:      os.Getenv("PORT"),
		Host:      os.Getenv("HOST"),
		Database: struct {
			Host     string
			Port     string
			User     string
			Password string
			Name     string
		}{
			Host:     os.Getenv("DB_HOST"),
			Port:     os.Getenv("DB_PORT"),
			User:     os.Getenv("DB_USER"),
			Password: os.Getenv("DB_PASSWORD"),
			Name:     os.Getenv("DB_NAME"),
		},
		JWTSecret:     os.Getenv("JWT_SECRET"),
		JWTExpiration: os.Getenv("JWT_EXPIRATION"),
		Redis: struct {
			Host     string
			Port     string
			Password string
			DB       string
		}{
			Host:     os.Getenv("REDIS_HOST"),
			Port:     os.Getenv("REDIS_PORT"),
			Password: os.Getenv("REDIS_PASSWORD"),
			DB:       os.Getenv("REDIS_DB"),
		},
	}
}
