package db

import (
	"github.com/ability-sh/abi-micro/micro"
)

func init() {
	micro.Reg("db", func(name string, config interface{}) (micro.Service, error) {
		return newDBService(name, config), nil
	})
}
