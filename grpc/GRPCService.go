package grpc

import (
	"fmt"

	"github.com/ability-sh/abi-micro/micro"
	G "google.golang.org/grpc"
)

type GRPCService interface {
	micro.Service
	Conn() G.ClientConnInterface
}

func GetConn(ctx micro.Context, name string) (G.ClientConnInterface, error) {
	s, err := ctx.GetService(name)
	if err != nil {
		return nil, err
	}
	ss, ok := s.(GRPCService)
	if ok {
		ctx.AddCount("grpc", 1)
		return ss.Conn(), nil
	}
	return nil, fmt.Errorf("service %s not instanceof GRPCService", name)
}
