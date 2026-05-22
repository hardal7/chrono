package logger

import (
	"log/slog"
	"os"
	"time"

	"github.com/lmittmann/tint"
)

func Init() {
	w := os.Stderr

	slog.SetDefault(slog.New(
		tint.NewHandler(w, &tint.Options{
			Level:      slog.LevelDebug,
			TimeFormat: time.Kitchen,
		}),
	))
}

func Debug(msg string) {
	slog.Debug(msg)
}

func Info(msg string) {
	slog.Info(msg)
}

func Warn(msg string) {
	slog.Warn(msg)
}

func Error(msg string) {
	slog.Error(msg)
}

func Fatal(diagnostic string, err error) {
	slog.Error(diagnostic)
	slog.Debug(err.Error())
	os.Exit(1)
}
