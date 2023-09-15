package smtp

import (
	"github.com/ability-sh/abi-micro/micro"
)

func init() {
	micro.Reg("smtp", func(name string, config interface{}) (micro.Service, error) {
		return newSmtpService(name, config), nil
	})
}
