package lrucache

import (
	"github.com/ability-sh/abi-micro/micro"
)

func init() {
	micro.Reg("lrucache", func(name string, config interface{}) (micro.Service, error) {
		return newLRUCacheService(name, config), nil
	})
}
