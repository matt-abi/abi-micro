package oss

import (
	"fmt"

	"github.com/ability-sh/abi-lib/dynamic"
	"github.com/ability-sh/abi-micro/micro"
)

type ossService struct {
	config interface{}
	name   string
	oss    OSS
}

func newOSSService(name string, config interface{}) OSSService {
	return &ossService{name: name, config: config}
}

/**
* 服务名称
**/
func (s *ossService) Name() string {
	return s.name
}

/**
* 服务配置
**/
func (s *ossService) Config() interface{} {
	return s.config
}

/**
* 初始化服务
**/
func (s *ossService) OnInit(ctx micro.Context) error {

	var err error = nil
	driver := dynamic.StringValue(dynamic.Get(s.config, "driver"), "ali")

	if driver == "ali" {
		cfg := AliOSSConfig{}
		dynamic.SetValue(&cfg, s.config)
		s.oss, err = NewAliOSS(&cfg)
	} else if driver == "leveldb" {
		cfg := LevelDBConfig{}
		dynamic.SetValue(&cfg, s.config)
		s.oss, err = NewLevelDB(&cfg)
	} else if driver == "aws" {
		cfg := AwsOSSConfig{}
		dynamic.SetValue(&cfg, s.config)
		s.oss, err = NewAwsOSS(&cfg)
	} else {
		return fmt.Errorf("oss driver %s not supported", driver)
	}

	return err
}

/**
* 校验服务是否可用
**/
func (s *ossService) OnValid(ctx micro.Context) error {
	return nil
}

func (s *ossService) OSS() OSS {
	return s.oss
}

func (s *ossService) Recycle() {
	if s.oss != nil {
		s.oss.Recycle()
		s.oss = nil
	}
}
