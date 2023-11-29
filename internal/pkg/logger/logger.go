package logger

import (
	"fmt"
	"log/slog"
)

type Logger struct {
}

func NewLogger() *Logger {
	return &Logger{}
}

func (*Logger) Errorf(format string, args ...interface{}) {
	slog.Error(fmt.Sprintf("error: "+format, args...))
}
