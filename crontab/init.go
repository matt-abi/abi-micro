package crontab

import (
	"github.com/matt-abi/abi-micro/micro"
)

func init() {
	micro.Reg("crontab", func(name string, config interface{}) (micro.Service, error) {
		return newCrontabService(name, config), nil
	})
}
