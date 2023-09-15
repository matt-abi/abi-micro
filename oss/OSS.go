package oss

import (
	"fmt"
	"time"

	"github.com/matt-abi/abi-micro/micro"
)

var ErrNoSuchKey = fmt.Errorf("no such key")

type OSS interface {
	micro.Recycle
	Get(key string) ([]byte, error)
	GetURL(key string) string
	GetSignURL(key string, expires time.Duration) (string, error)
	Put(key string, data []byte, header map[string]string) error
	PutSignURL(key string, expires time.Duration, header map[string]string) (string, error)
	PostSignURL(key string, expires time.Duration, maxSize int64, header map[string]string) (string, map[string]string, error)
	Del(key string) error
	Has(key string) (bool, error)
}
