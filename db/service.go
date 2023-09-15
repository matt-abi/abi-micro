package db

import (
	"database/sql"
	"time"

	"github.com/matt-abi/abi-lib/dynamic"
	"github.com/matt-abi/abi-micro/micro"
)

type dbConfig struct {
	Driver       string `json:"driver"`
	Url          string `json:"url"`
	MaxIdleConns int    `json:"maxIdleConns"`
	MaxOpenConns int    `json:"maxOpenConns"`
	MaxLifeTime  int    `json:"maxLifeTime"`
}

type dbService struct {
	config interface{}
	name   string
	db     *sql.DB
}

func newDBService(name string, config interface{}) DBService {
	return &dbService{name: name, config: config}
}

/**
* 服务名称
**/
func (s *dbService) Name() string {
	return s.name
}

/**
* 服务配置
**/
func (s *dbService) Config() interface{} {
	return s.config
}

/**
* 初始化服务
**/
func (s *dbService) OnInit(ctx micro.Context) error {

	var err error = nil
	cfg := dbConfig{}

	dynamic.SetValue(&cfg, s.config)

	db, err := sql.Open(cfg.Driver, cfg.Url)

	if err != nil {
		return err
	}

	db.SetMaxIdleConns(cfg.MaxIdleConns)
	db.SetMaxOpenConns(cfg.MaxOpenConns)
	db.SetConnMaxLifetime(time.Duration(cfg.MaxLifeTime) * time.Second)

	s.db = db

	return err
}

/**
* 校验服务是否可用
**/
func (s *dbService) OnValid(ctx micro.Context) error {
	return nil
}

func (s *dbService) GetDB() *sql.DB {
	return s.db
}

func (s *dbService) Recycle() {
	if s.db != nil {
		s.db.Close()
		s.db = nil
	}
}
