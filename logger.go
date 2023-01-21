package mapper

import "fmt"

var logger Logger

type Logger interface {
	Log(string, ...any)
}

func SetLogger(l Logger) {
	logger = l
}

func init() {
	SetLogger(StdoutLogger{})
}

type StdoutLogger struct{}

var _ Logger = StdoutLogger{}

func (StdoutLogger) Log(template string, fields ...any) {
	fmt.Printf(template, fields...)
	fmt.Print("\n")
}
