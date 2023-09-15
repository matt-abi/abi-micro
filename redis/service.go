package redis

import (
	"crypto/tls"
	"strings"
	"time"

	R "github.com/go-redis/redis/v8"
	"github.com/matt-abi/abi-lib/dynamic"
	"github.com/matt-abi/abi-micro/micro"
)

type redisConfig struct {
	Addr         string `json:"addr"`
	UserName     string `json:"username"`
	Password     string `json:"password"`
	DB           int    `json:"db"`
	PoolSize     int    `json:"pool-size"`
	MinIdleConns int    `json:"min-idle-conns"`
	IdleTimeout  int    `json:"idle-timeout"`
	Tls          bool   `json:"tls"`
	Cluster      bool   `json:"cluster"`
}

type redisService struct {
	config        interface{}
	name          string
	client        *R.Client
	clusterClient *R.ClusterClient
}

func newRedisService(name string, config interface{}) RedisService {
	return &redisService{name: name, config: config}
}

/**
* 服务名称
**/
func (s *redisService) Name() string {
	return s.name
}

/**
* 服务配置
**/
func (s *redisService) Config() interface{} {
	return s.config
}

/**
* 初始化服务
**/
func (s *redisService) OnInit(ctx micro.Context) error {

	var cfg = &redisConfig{}

	dynamic.SetValue(cfg, s.config)

	if cfg.Cluster {

		opt := &R.ClusterOptions{
			Addrs:        strings.Split(cfg.Addr, ","),
			Password:     cfg.Password, // no password set
			PoolSize:     cfg.PoolSize,
			Username:     cfg.UserName,
			MinIdleConns: cfg.MinIdleConns,
			IdleTimeout:  time.Duration(cfg.IdleTimeout) * time.Second,
		}

		if cfg.Tls {
			opt.TLSConfig = &tls.Config{}
		}

		s.clusterClient = R.NewClusterClient(opt)

		return nil
	}

	opt := &R.Options{
		Addr:         cfg.Addr,
		Password:     cfg.Password, // no password set
		DB:           cfg.DB,       // use default DB,
		PoolSize:     cfg.PoolSize,
		Username:     cfg.UserName,
		MinIdleConns: cfg.MinIdleConns,
		IdleTimeout:  time.Duration(cfg.IdleTimeout) * time.Second,
	}

	if cfg.Tls {
		opt.TLSConfig = &tls.Config{}
	}

	s.client = R.NewClient(opt)

	return nil
}

/**
* 校验服务是否可用
**/
func (s *redisService) OnValid(ctx micro.Context) error {
	return nil
}

func (s *redisService) Client() *R.Client {
	return s.client
}

func (s *redisService) ClusterClient() *R.ClusterClient {
	return s.clusterClient
}

func (s *redisService) Recycle() {
	if s.client != nil {
		s.client.Close()
		s.client = nil
	}
	if s.clusterClient != nil {
		s.clusterClient.Close()
		s.clusterClient = nil
	}
}
