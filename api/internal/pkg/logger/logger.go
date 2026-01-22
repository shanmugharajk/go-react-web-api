package logger

import (
	"log/slog"
	"os"
)

var log *slog.Logger

// Init initializes the global logger.
func Init() {
	handler := slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	})
	log = slog.New(handler)
	slog.SetDefault(log)
}

// Get returns the global logger instance.
func Get() *slog.Logger {
	if log == nil {
		Init()
	}
	return log
}

// Info logs an info message.
func Info(msg string, args ...any) {
	Get().Info(msg, args...)
}

// Error logs an error message.
func Error(msg string, args ...any) {
	Get().Error(msg, args...)
}

// Debug logs a debug message.
func Debug(msg string, args ...any) {
	Get().Debug(msg, args...)
}

// Warn logs a warning message.
func Warn(msg string, args ...any) {
	Get().Warn(msg, args...)
}
