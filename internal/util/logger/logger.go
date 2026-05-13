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

func Debug(diagnostics string) {
	slog.Debug(diagnostics)
}

func Info(diagnostics string) {
	slog.Info(diagnostics)
}

func Warn(diagnostics string) {
	slog.Warn(diagnostics)
}

func Error(diagnostics string) {
	slog.Error(diagnostics)
	os.Exit(1)
}
