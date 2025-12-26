package xlog

import "context"

type Config struct {
	Level   int    `json:"level"`
	Name    string `json:"name"`
	Suffix  string `json:"suffix"`
	Path    string `json:"path"`
	MaxSize int    `json:"max_size"`
}

type ILogger interface {
	SetLevel(ctx context.Context, level int)
	SetSkip(ctx context.Context, skip int)
	Debug(ctx context.Context, format string, v ...any)
	Info(ctx context.Context, format string, v ...any)
	Warn(ctx context.Context, format string, v ...any)
	Error(ctx context.Context, format string, v ...any)
	Fatal(ctx context.Context, format string, v ...any)
	Stack(ctx context.Context)
}

const (
	LevelDebug = iota
	LevelInfo
	LevelWarn
	LevelError
	LevelFatal
	LevelStack
)

var levelDescMap = map[int]string{
	LevelDebug: "DEBUG",
	LevelInfo:  "INFO",
	LevelWarn:  "WARN",
	LevelError: "ERROR",
	LevelFatal: "FATAL",
	LevelStack: "STACK",
}
