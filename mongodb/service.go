package mongodb

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"io/ioutil"

	"github.com/matt-abi/abi-lib/dynamic"
	"github.com/matt-abi/abi-micro/micro"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type mongodbConfig struct {
	URI    string `json:"uri"`
	DB     string `json:"db"`
	CA     string `json:"ca"`
	CAFile string `json:"ca-file"`
}

type mongodbService struct {
	config interface{}
	name   string
	client *mongo.Client
	db     *mongo.Database
}

func newMongoDBService(name string, config interface{}) MongoDBService {
	return &mongodbService{name: name, config: config}
}

/**
* 服务名称
**/
func (s *mongodbService) Name() string {
	return s.name
}

/**
* 服务配置
**/
func (s *mongodbService) Config() interface{} {
	return s.config
}

/**
* 初始化服务
**/
func (s *mongodbService) OnInit(ctx micro.Context) error {

	var err error = nil
	cfg := mongodbConfig{}

	dynamic.SetValue(&cfg, s.config)

	opt := options.Client().ApplyURI(cfg.URI)

	if cfg.CA != "" {
		tlsConfig := &tls.Config{RootCAs: x509.NewCertPool()}
		ok := tlsConfig.RootCAs.AppendCertsFromPEM([]byte(cfg.CA))
		if !ok {
			return fmt.Errorf("Failed parsing pem %s", cfg.CA)
		}
		opt.SetTLSConfig(tlsConfig)
	}

	if cfg.CAFile != "" {

		tlsConfig := &tls.Config{}
		certs, err := ioutil.ReadFile(cfg.CAFile)

		if err != nil {
			return err
		}

		tlsConfig.RootCAs = x509.NewCertPool()
		ok := tlsConfig.RootCAs.AppendCertsFromPEM(certs)

		if !ok {
			return fmt.Errorf("Failed parsing pem file %s", cfg.CAFile)
		}
		opt.SetTLSConfig(tlsConfig)
	}

	s.client, err = mongo.NewClient(opt)

	if err != nil {
		return err
	}

	err = s.client.Connect(context.Background())

	if err != nil {
		return err
	}

	s.db = s.client.Database(cfg.DB)

	return err
}

/**
* 校验服务是否可用
**/
func (s *mongodbService) OnValid(ctx micro.Context) error {
	return nil
}

func (s *mongodbService) GetClient() *mongo.Client {
	return s.client
}

func (s *mongodbService) GetDB() *mongo.Database {
	return s.db
}

func (s *mongodbService) Recycle() {
	if s.client != nil {
		s.client.Disconnect(context.Background())
		s.client = nil
	}
}
