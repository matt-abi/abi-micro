package oss

import (
	"fmt"
	"time"

	"github.com/matt-abi/abi-micro/micro"
)

type OSSService interface {
	micro.Service
	OSS() OSS
}

func GetOSS(ctx micro.Context, name string) (OSS, error) {
	s, err := ctx.GetService(name)
	if err != nil {
		return nil, err
	}
	ss, ok := s.(OSSService)
	if ok {
		ctx.AddCount("oss", 1)
		return &statOSS{oss: ss.OSS(), ctx: ctx}, nil
	}
	return nil, fmt.Errorf("service %s not instanceof OSSService", name)
}

type statOSS struct {
	ctx micro.Context
	oss OSS
}

func (s *statOSS) Get(key string) ([]byte, error) {
	st := s.ctx.Step("oss.Get")
	rs, err := s.oss.Get(key)
	st("")
	return rs, err
}

func (s *statOSS) GetURL(key string) string {
	st := s.ctx.Step("oss.GetURL")
	rs := s.oss.GetURL(key)
	st("oss.GetURL", "")
	return rs
}

func (s *statOSS) GetSignURL(key string, expires time.Duration) (string, error) {
	st := s.ctx.Step("oss.GetSignURL")
	rs, err := s.oss.GetSignURL(key, expires)
	st("")
	return rs, err
}

func (s *statOSS) Put(key string, data []byte, header map[string]string) error {
	st := s.ctx.Step("oss.Put")
	err := s.oss.Put(key, data, header)
	st("")
	return err
}

func (s *statOSS) PutSignURL(key string, expires time.Duration, header map[string]string) (string, error) {
	st := s.ctx.Step("oss.PutSignURL")
	rs, err := s.oss.PutSignURL(key, expires, header)
	st("")
	return rs, err
}

func (s *statOSS) PostSignURL(key string, expires time.Duration, maxSize int64, header map[string]string) (string, map[string]string, error) {
	st := s.ctx.Step("oss.PostSignURL")
	rs, data, err := s.oss.PostSignURL(key, expires, maxSize, header)
	st("")
	return rs, data, err
}

func (s *statOSS) Del(key string) error {
	st := s.ctx.Step("oss.Del")
	rs := s.oss.Del(key)
	st("")
	return rs
}

func (s *statOSS) Has(key string) (bool, error) {
	st := s.ctx.Step("oss.Has")
	rs, err := s.oss.Has(key)
	st("")
	return rs, err
}

func (s *statOSS) Recycle() {

}
