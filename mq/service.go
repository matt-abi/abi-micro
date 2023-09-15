package mq

import (
	"fmt"

	"github.com/matt-abi/abi-lib/dynamic"
	"github.com/matt-abi/abi-micro/micro"
)

type mqService struct {
	config interface{}
	name   string
	driver Driver
}

func newMQService(name string, config interface{}) *mqService {
	return &mqService{name: name, config: config}
}

/**
* 服务名称
**/
func (s *mqService) Name() string {
	return s.name
}

/**
* 服务配置
**/
func (s *mqService) Config() interface{} {
	return s.config
}

/**
* 初始化服务
**/
func (s *mqService) OnInit(ctx micro.Context) error {

	var err error = nil
	driver := dynamic.StringValue(dynamic.Get(s.config, "driver"), "kafka")

	if driver == "kafka" {
		v, err := NewDriver(driver, s.config)
		if err != nil {
			return err
		}
		s.driver = v
	} else {
		return fmt.Errorf("mq driver %s not supported", driver)
	}

	return err
}

/**
* 校验服务是否可用
**/
func (s *mqService) OnValid(ctx micro.Context) error {
	return nil
}

func (s *mqService) MQ() Driver {
	return s.driver
}

func (s *mqService) Recycle() {
	if s.driver != nil {
		s.driver.Recycle()
		s.driver = nil
	}
}
