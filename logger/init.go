package logger

import (
	"github.com/matt-abi/abi-micro/micro"
)

func init() {
	micro.Reg("logger", func(name string, config interface{}) (micro.Service, error) {
		return newLoggerService(name, config), nil
	})
}
