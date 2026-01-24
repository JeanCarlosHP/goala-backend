package logger

import (
	"context"
	"fmt"
	"log"
	"log/slog"
	"os"

	"github.com/jeancarloshp/calorieai/internal/domain"
	"github.com/jeancarloshp/calorieai/internal/domain/enum"
)

type loggerImpl struct {
	ctx    context.Context
	Logger *slog.Logger
}

func New(config *domain.Config) domain.Logger {
	var logger *slog.Logger

	slogHandlerOptions := &slog.HandlerOptions{
		ReplaceAttr: func(groups []string, a slog.Attr) slog.Attr {
			if a.Key == slog.MessageKey {
				a.Key = "message"
			}

			return a
		},
	}

	setLevel(slogHandlerOptions, config.LoggingLevel)

	if config.LoggingJSONFormat {
		logger = slog.New(slog.NewJSONHandler(os.Stdout, slogHandlerOptions))
	} else {
		logger = slog.New(slog.NewTextHandler(os.Stdout, slogHandlerOptions))
	}

	return &loggerImpl{
		Logger: logger,
	}
}

func setLevel(handler *slog.HandlerOptions, level string) {
	switch level {
	case enum.LoggingLevelDebug.String():
		handler.Level = slog.LevelDebug
	case enum.LoggingLevelInfo.String():
		handler.Level = slog.LevelInfo
	case enum.LoggingLevelWarning.String():
		handler.Level = slog.LevelWarn
	case enum.LoggingLevelError.String():
		handler.Level = slog.LevelError
	}
}

func (l *loggerImpl) Panic(v any) {
	panic(v)
}

func (l *loggerImpl) Panicf(format string, v ...any) {
	panic(fmt.Sprintf(format, v...))
}

func (l *loggerImpl) Fatal(v ...any) {
	log.Fatal(v...)
}

func (l *loggerImpl) Fatalf(format string, v ...any) {
	log.Fatalf(format, v...)
}

func (l *loggerImpl) Error(msg string, args ...any) {
	l.Logger.Error(msg, args...)
}

func (l *loggerImpl) Warn(msg string, args ...any) {
	l.Logger.Warn(msg, args...)
}

func (l *loggerImpl) Info(msg string, args ...any) {
	l.Logger.Log(l.ctx, slog.LevelInfo, msg, args...)
}

func (l *loggerImpl) Infof(msg string, args ...any) {
	l.Logger.Info(fmt.Sprintf(msg, args...))
}

func (l *loggerImpl) Debug(msg string, args ...any) {
	l.Logger.Debug(msg, args...)
}
