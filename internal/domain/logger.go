package domain

type Logger interface {
	Panic(v any)
	Panicf(format string, v ...any)
	Fatal(v ...any)
	Fatalf(format string, v ...any)
	Error(msg string, args ...any)
	Warn(msg string, args ...any)
	Info(msg string, args ...any)
	Infof(msg string, args ...any)
	Debug(msg string, args ...any)
}
