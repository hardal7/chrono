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

	PostgresHost     string
	PostgresPort     string
	PostgresUser     string
	PostgresDB       string
	PostgresPassword string

	RedisPort     string
	RedisHost     string
	RedisPassword string

	OTelEndpoint string

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
		slog.Error("load environment variables")
		slog.Debug(err.Error())
		os.Exit(1)
	}

	App = Config{
		Domain:    os.Getenv("DOMAIN"),
		Port:      os.Getenv("APP_PORT"),
		AdminPort: os.Getenv("ADMIN_PORT"),
		LogLevel:  os.Getenv("LOG_LEVEL"),

		PostgresHost:     os.Getenv("POSTGRES_HOST"),
		PostgresPort:     os.Getenv("POSTGRES_PORT"),
		PostgresUser:     os.Getenv("POSTGRES_USER"),
		PostgresDB:       os.Getenv("POSTGRES_DB"),
		PostgresPassword: os.Getenv("POSTGRES_PASSWORD"),

		RedisPort:     os.Getenv("REDIS_PORT"),
		RedisHost:     os.Getenv("REDIS_HOST"),
		RedisPassword: os.Getenv("REDIS_PASSWORD"),

		OTelEndpoint: os.Getenv("OTEL_EXPORTER_OTLP_ENDPOINT"),

		HashSecret: os.Getenv("HASH_SECRET"),

		MailAPIKey:  os.Getenv("MAIL_API_KEY"),
		MailAddress: os.Getenv("MAIL_ADDRESS"),
	}
	slog.Info("Loaded environment variables")
}
