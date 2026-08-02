package logger

import (
	"context"
	"log/slog"
	"os"

	"github.com/hardal7/chrono/internal/util/config"
	"github.com/lmittmann/tint"
)

func Init() {
	w := os.Stderr
	logger := (slog.New(
		tint.NewHandler(w, &tint.Options{
			Level: getLevel(config.App.LogLevel),
			ReplaceAttr: func(groups []string, a slog.Attr) slog.Attr {
				if a.Key == slog.LevelKey && len(groups) == 0 {
					level, ok := a.Value.Any().(slog.Level)
					if ok && level <= LevelTrace {
						return tint.Attr(1, slog.String(a.Key, "TRC"))
					}
				}
				return a
			},
			TimeFormat: "12:00:00",
		}),
	))
	slog.SetDefault(logger)
}

const LevelTrace = slog.Level(-8)

func getLevel(level string) slog.Level {
	switch level {
	case "TRACE":
		return LevelTrace
	case "DEBUG":
		return slog.LevelDebug
	case "INFO":
		return slog.LevelInfo
	case "WARN":
		return slog.LevelWarn
	}
	return slog.LevelDebug
}

func Trace(msg string, args ...any) {
	slog.Log(context.Background(), LevelTrace, msg, args...)
}

func Debug(msg string, args ...any) {
	slog.Debug(msg, args...)
}

func Info(msg string, args ...any) {
	slog.Info(msg, args...)
}

func Warn(msg string, args ...any) {
	slog.Warn(msg, args...)
}

func Error(msg string, args ...any) {
	slog.Error(msg, args...)
}

func Fatal(msg string, err error, args ...any) {
	slog.Error(msg)
	slog.Debug(err.Error())
	os.Exit(1)
}
