package oss

import (
	"github.com/ability-sh/abi-micro/micro"
)

func init() {
	micro.Reg("oss", func(name string, config interface{}) (micro.Service, error) {
		return newOSSService(name, config), nil
	})
}
