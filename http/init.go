package http

import (
	"github.com/matt-abi/abi-micro/micro"
)

func init() {
	micro.Reg("http", func(name string, config interface{}) (micro.Service, error) {
		return newHTTPService(name, config), nil
	})
}
