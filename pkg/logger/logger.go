package logger

type Logger interface {
	Info(msg string, args ...any)
	Warn(msg string, args ...any)
	Error(err error, msg string, args ...any)
	Debug(msg string, args ...any)
	Track(operationName string, operation func() error) error
}
