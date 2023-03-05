// Package logger provider logging functionality
package logger

import (
	"fmt"
	"os"
)

type Logger interface {
	Log(string, ...any)
}

var logger Logger

func Log(template string, fields ...any) {
	logger.Log(template, fields...)
}

func SetLogger(l Logger) {
	logger = l
}

func init() {
	SetLogger(StderrLogger{})
}

type StderrLogger struct{}

var _ Logger = StderrLogger{}

func (StderrLogger) Log(template string, fields ...any) {
	_, _ = fmt.Fprintf(os.Stderr, template, fields...)
	_, _ = fmt.Fprintf(os.Stderr, "\n")
}
