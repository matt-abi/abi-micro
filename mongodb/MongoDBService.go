package mongodb

import (
	"fmt"

	"github.com/matt-abi/abi-micro/micro"
	"go.mongodb.org/mongo-driver/mongo"
)

type MongoDBService interface {
	micro.Service
	GetClient() *mongo.Client
	GetDB() *mongo.Database
}

func GetClient(ctx micro.Context, name string) (*mongo.Client, error) {
	s, err := ctx.GetService(name)
	if err != nil {
		return nil, err
	}
	ss, ok := s.(MongoDBService)
	if ok {
		ctx.AddCount("mongo", 1)
		return ss.GetClient(), nil
	}
	return nil, fmt.Errorf("service %s not instanceof MongoDBService", name)
}

func GetDB(ctx micro.Context, name string) (*mongo.Database, error) {
	s, err := ctx.GetService(name)
	if err != nil {
		return nil, err
	}
	ss, ok := s.(MongoDBService)
	if ok {
		ctx.AddCount("mongo", 1)
		return ss.GetDB(), nil
	}
	return nil, fmt.Errorf("service %s not instanceof MongoDBService", name)
}
