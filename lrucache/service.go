package lrucache

import (
	lru "github.com/hashicorp/golang-lru"
	"github.com/matt-abi/abi-lib/dynamic"
	"github.com/matt-abi/abi-micro/micro"
)

type lrucacheConfig struct {
	MaxSize int `json:"max-size"`
}

type lrucacheService struct {
	config interface{}
	name   string
	cache  *lru.Cache
}

func newLRUCacheService(name string, config interface{}) LRUCacheService {
	return &lrucacheService{name: name, config: config}
}

/**
* 服务名称
**/
func (s *lrucacheService) Name() string {
	return s.name
}

/**
* 服务配置
**/
func (s *lrucacheService) Config() interface{} {
	return s.config
}

/**
* 初始化服务
**/
func (s *lrucacheService) OnInit(ctx micro.Context) error {

	var err error = nil
	cfg := lrucacheConfig{}

	dynamic.SetValue(&cfg, s.config)

	s.cache, err = lru.New(cfg.MaxSize)

	return err
}

/**
* 校验服务是否可用
**/
func (s *lrucacheService) OnValid(ctx micro.Context) error {
	return nil
}

func (s *lrucacheService) Cache() LRUCache {
	return s.cache
}

func (s *lrucacheService) Recycle() {
	s.cache = nil
}
