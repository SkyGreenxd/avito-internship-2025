package logger

import (
	"log/slog"
	"os"
	"time"
)

type SlogLogger struct {
	logger *slog.Logger
}

func NewSlogLogger() *SlogLogger {
	handler := slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelDebug,
	})
	return &SlogLogger{
		logger: slog.New(handler),
	}
}

func (l *SlogLogger) Info(msg string, args ...any) {
	l.logger.Info(msg, args...)
}

func (l *SlogLogger) Warn(msg string, args ...any) {
	l.logger.Warn(msg, args...)
}

func (l *SlogLogger) Error(err error, msg string, args ...any) {
	args = append(args, "error", err)
	l.logger.Error(msg, args...)
}

func (l *SlogLogger) Debug(msg string, args ...any) {
	l.logger.Debug(msg, args...)
}

// Track выполняет переданную функцию, которая может вернуть ошибку,
// замеряет время ее выполнения и логирует результат.
func (l *SlogLogger) Track(operationName string, operation func() error) error {
	start := time.Now()

	// Выполняем переданную операцию и получаем ошибку, если она есть
	err := operation()

	// Логируем результат в любом случае
	l.Info(
		"Operation finished",
		"operation", operationName,
		"duration", time.Since(start).String(),
		"error", err, // slog элегантно обработает nil, если ошибки не было
	)

	// Возвращаем ошибку, если она была, чтобы вызывающий код мог ее обработать
	return err
}
