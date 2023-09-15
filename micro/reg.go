package micro

import (
	"fmt"
)

type ServiceFactory func(name string, config interface{}) (Service, error)

var serviceSet = map[string]ServiceFactory{}

/**
 * 注册服务
 **/
func Reg(stype string, factory ServiceFactory) {
	serviceSet[stype] = factory
}

/**
 * 创建服务
 **/
func NewService(stype string, name string, config interface{}) (Service, error) {
	ss := serviceSet[stype]
	if ss == nil {
		return nil, fmt.Errorf("micro service %s not found", stype)
	}
	return ss(name, config)
}
