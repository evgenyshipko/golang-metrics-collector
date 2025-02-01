package logger

import (
	"log/slog"
)

var Logger = newLogger()

type SlogLogger struct {
	logger *slog.Logger
}

func newLogger() *SlogLogger {
	logger := SetupPrettySlog()

	return &SlogLogger{
		logger: logger,
	}
}

func Error(msg string, args ...any) {
	Logger.logger.Error(msg, args...)
}

func Info(msg string, args ...any) {
	Logger.logger.Info(msg, args...)
}

func Debug(msg string, args ...any) {
	Logger.logger.Debug(msg, args...)
}

func Warn(msg string, args ...any) {
	Logger.logger.Warn(msg, args...)
}
