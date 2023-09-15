package grpc

import (
	"github.com/ability-sh/abi-micro/micro"
)

func init() {
	micro.Reg("grpc", func(name string, config interface{}) (micro.Service, error) {
		return newGRPCService(name, config), nil
	})
}
