package logger

import "fmt"

type stdoutLogger struct {
}

func NewStdoutLogger() Logger {
	return &stdoutLogger{}
}

func (l *stdoutLogger) Output(text string) {
	fmt.Println(text)
}

func (l *stdoutLogger) Recycle() {

}
