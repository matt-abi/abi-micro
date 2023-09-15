package mongodb

import (
	"github.com/ability-sh/abi-micro/micro"
)

func init() {
	micro.Reg("mongodb", func(name string, config interface{}) (micro.Service, error) {
		return newMongoDBService(name, config), nil
	})
}
