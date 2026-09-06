package config

import (
	"log/slog"
	"os"
	"testing"

	"github.com/joho/godotenv"
)

const AvatarEndpoint string = "/avatar"

type Config struct {
	Domain    string
	Port      string
	AdminPort string
	LogLevel  string

	DBPort     string
	DBHost     string
	DBUser     string
	DBName     string
	DBPassword string

	RedisPort     string
	RedisHost     string
	RedisPassword string

	HashSecret string

	MailAPIKey  string
	MailAddress string
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
		Domain:    os.Getenv("DOMAIN"),
		Port:      os.Getenv("APP_PORT"),
		AdminPort: os.Getenv("ADMIN_PORT"),
		LogLevel:  os.Getenv("LOG_LEVEL"),

		DBPort:     os.Getenv("DB_PORT"),
		DBHost:     os.Getenv("DB_HOST"),
		DBUser:     os.Getenv("DB_USER"),
		DBName:     os.Getenv("DB_NAME"),
		DBPassword: os.Getenv("DB_PASSWORD"),

		RedisPort:     os.Getenv("REDIS_PORT"),
		RedisHost:     os.Getenv("REDIS_HOST"),
		RedisPassword: os.Getenv("REDIS_PASSWORD"),

		HashSecret: os.Getenv("HASH_SECRET"),

		MailAPIKey:  os.Getenv("MAIL_API_KEY"),
		MailAddress: os.Getenv("MAIL_ADDRESS"),
	}
	slog.Info("Loaded environment variables")
}
