package redis

import (
	"time"
)

type Redis interface {
	Get(key string) (string, error)
	Set(key string, value string, expiration time.Duration) error
	TTL(key string) (time.Duration, error)
	Expire(key string, expiration time.Duration) error
	Del(key string) error
	HSet(key string, itemKey string, value string) error
	HGet(key string, itemKey string) (string, error)
	HDel(key string, itemKey string) error
}

func IsNil(err error) bool {
	return err.Error() == "redis: nil"
}
