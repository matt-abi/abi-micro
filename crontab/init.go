package crontab

import (
	"github.com/ability-sh/abi-micro/micro"
)

func init() {
	micro.Reg("crontab", func(name string, config interface{}) (micro.Service, error) {
		return newCrontabService(name, config), nil
	})
}
