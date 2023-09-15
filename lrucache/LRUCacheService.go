package lrucache

import (
	"fmt"

	"github.com/matt-abi/abi-micro/micro"
)

type LRUCache interface {
	Contains(key interface{}) bool
	Get(key interface{}) (interface{}, bool)
	Add(key, value interface{}) bool
	Remove(key interface{}) bool
}

type LRUCacheService interface {
	micro.Service
	Cache() LRUCache
}

func GetCache(ctx micro.Context, name string) (LRUCache, error) {
	s, err := ctx.GetService(name)
	if err != nil {
		return nil, err
	}
	ss, ok := s.(LRUCacheService)
	if ok {
		return ss.Cache(), nil
	}
	return nil, fmt.Errorf("service %s not instanceof LRUCacheService", name)
}
