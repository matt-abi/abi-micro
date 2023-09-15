package logger

import (
	"fmt"

	"github.com/matt-abi/abi-lib/dynamic"
	"github.com/matt-abi/abi-micro/micro"
)

type loggerService struct {
	config interface{}
	name   string
	logger Logger
}

func newLoggerService(name string, config interface{}) LoggerService {
	return &loggerService{name: name, config: config}
}

/**
* 服务名称
**/
func (s *loggerService) Name() string {
	return s.name
}

/**
* 服务配置
**/
func (s *loggerService) Config() interface{} {
	return s.config
}

/**
* 初始化服务
**/
func (s *loggerService) OnInit(ctx micro.Context) error {

	var err error = nil
	driver := dynamic.StringValue(dynamic.Get(s.config, "driver"), "empty")

	if driver == "empty" {
		s.logger = NewEmptyLogger()
	} else if driver == "stdout" {
		s.logger = NewStdoutLogger()
	} else if driver == "syslog" {
		s.logger, err = NewSyslogLogger(s.config)
	} else if driver == "fs" {
		s.logger, err = NewFSLogger(s.config)
	} else {
		return fmt.Errorf("logger driver %s not supported", driver)
	}

	return err
}

/**
* 校验服务是否可用
**/
func (s *loggerService) OnValid(ctx micro.Context) error {
	return nil
}

func (s *loggerService) Logger() Logger {
	return s.logger
}

func (s *loggerService) Recycle() {
	if s.logger != nil {
		s.logger.Recycle()
		s.logger = nil
	}
}
