package grpc

import (
	"context"

	"github.com/matt-abi/abi-micro/micro"
	"google.golang.org/grpc/metadata"
)

func NewGRPCContext(ctx micro.Context) context.Context {
	md := metadata.MD{}
	ctx.Each(func(key string, value string) bool {
		md.Set(key, value)
		return true
	})
	return metadata.NewOutgoingContext(micro.WithContext(context.Background(), ctx), md)
}

func GetContext(c context.Context) micro.Context {
	return micro.GetContext(c)
}
