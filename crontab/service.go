package crontab

import (
	"fmt"

	"github.com/matt-abi/abi-lib/dynamic"
	"github.com/matt-abi/abi-micro/micro"
	"github.com/robfig/cron/v3"
)

type crontabService struct {
	config interface{}
	name   string
	c      *cron.Cron
}

func newCrontabService(name string, config interface{}) *crontabService {
	return &crontabService{name: name, config: config}
}

/**
* 服务名称
**/
func (s *crontabService) Name() string {
	return s.name
}

/**
* 服务配置
**/
func (s *crontabService) Config() interface{} {
	return s.config
}

/**
* 初始化服务
**/
func (s *crontabService) OnInit(ctx micro.Context) error {

	var err error = nil

	p := ctx.Payload()

	s.c = cron.New(cron.WithSeconds())

	dynamic.Each(dynamic.Get(s.config, "jobs"), func(_ interface{}, job interface{}) bool {

		spec := dynamic.StringValue(dynamic.Get(job, "spec"), "")
		t := dynamic.StringValue(dynamic.Get(job, "type"), "http")

		if t == "http" {
			job := NewHttpJob(p, job)
			_, err = s.c.AddJob(spec, job)
		} else {
			err = fmt.Errorf("not support job type %s", t)
		}

		return err == nil
	})

	if err != nil {
		return err
	}

	s.c.Start()

	return nil
}

/**
* 校验服务是否可用
**/
func (s *crontabService) OnValid(ctx micro.Context) error {
	return nil
}

func (s *crontabService) Recycle() {
	if s.c != nil {
		s.c.Stop()
		s.c = nil
	}
}
