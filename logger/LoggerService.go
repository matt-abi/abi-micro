package logger

import (
	"fmt"

	"github.com/ability-sh/abi-micro/micro"
)

type Logger interface {
	Output(text string)
	Recycle()
}

type LoggerService interface {
	micro.Service
	Logger() Logger
}

func GetLogger(ctx micro.Context, name string) (Logger, error) {
	s, err := ctx.GetService(name)
	if err != nil {
		return nil, err
	}
	ss, ok := s.(LoggerService)
	if ok {
		return ss.Logger(), nil
	}
	return nil, fmt.Errorf("service %s not instanceof LoggerService", name)
}
