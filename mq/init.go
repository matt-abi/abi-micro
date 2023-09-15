package mq

import (
	"github.com/ability-sh/abi-micro/micro"
)

func init() {
	micro.Reg("mq", func(name string, config interface{}) (micro.Service, error) {
		return newMQService(name, config), nil
	})
}
