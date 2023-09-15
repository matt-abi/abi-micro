package redis

import (
	"github.com/ability-sh/abi-micro/micro"
)

func init() {
	micro.Reg("redis", func(name string, config interface{}) (micro.Service, error) {
		return newRedisService(name, config), nil
	})
}
