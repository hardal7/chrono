package config

import (
	"log/slog"
	"os"
	"testing"

	"github.com/joho/godotenv"
)

type Config struct {
	Port        string
	AdminPort   string
	LogLevel    string
	DB_PORT     string
	DB_HOST     string
	DB_USER     string
	DB_NAME     string
	DB_PASSWORD string
	JWT_SECRET  string
}

var App Config

func Load() {
	var file string
	if !testing.Testing() {
		file = ".env"
	} else {
		file = "/srv/.env.test"
		slog.Warn("Running in test environment")
	}
	err := godotenv.Load(file)
	if err != nil {
		slog.Error("Failed to load environment variables")
		slog.Debug(err.Error())
		os.Exit(1)
	}

	App = Config{
		Port:        os.Getenv("APP_PORT"),
		AdminPort:   os.Getenv("ADMIN_PORT"),
		LogLevel:    os.Getenv("LOG_LEVEL"),
		DB_PORT:     os.Getenv("DB_PORT"),
		DB_HOST:     os.Getenv("DB_HOST"),
		DB_USER:     os.Getenv("DB_USER"),
		DB_NAME:     os.Getenv("DB_NAME"),
		DB_PASSWORD: os.Getenv("DB_PASSWORD"),
		JWT_SECRET:  os.Getenv("JWT_SECRET"),
	}
	slog.Info("Loaded environment variables")
}
