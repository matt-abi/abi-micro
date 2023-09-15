package logger

import (
	"fmt"
	"os"

	"github.com/ability-sh/abi-lib/dynamic"
)

type fsLogger struct {
	Path string `json:"path"`
	fd   *os.File
}

func NewFSLogger(config interface{}) (Logger, error) {
	v := &fsLogger{}
	dynamic.SetValue(v, config)
	var err error
	v.fd, err = os.OpenFile(v.Path, os.O_APPEND, os.ModeAppend)
	if err != nil {
		v.fd, err = os.Create(v.Path)
	}
	if err != nil {
		return nil, err
	}
	return v, nil
}

func (l *fsLogger) Output(text string) {
	fmt.Fprintln(l.fd, text)
}

func (l *fsLogger) Recycle() {
	if l.fd != nil {
		l.fd.Close()
		l.fd = nil
	}
}
