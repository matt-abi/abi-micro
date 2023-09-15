package logger

type emptyLogger struct {
}

func NewEmptyLogger() Logger {
	return &emptyLogger{}
}

func (l *emptyLogger) Output(text string) {

}

func (l *emptyLogger) Recycle() {

}
