package grpc

import (
	"context"

	"github.com/ability-sh/abi-lib/dynamic"
	"github.com/ability-sh/abi-micro/micro"
	G "google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type grpcConn struct {
	conn *G.ClientConn
}

func (c *grpcConn) Invoke(cc context.Context, method string, args interface{}, reply interface{}, opts ...G.CallOption) error {
	ctx := GetContext(cc)
	if ctx == nil {
		return c.conn.Invoke(cc, method, args, reply, opts...)
	}
	st := ctx.Step("grpc.Invoke")
	err := c.conn.Invoke(cc, method, args, reply, opts...)
	if err != nil {
		st("grpc.Invoke", "[%s] [err:1] %s", method, err.Error())
	} else {
		st("grpc.Invoke", "[%s]", method)
	}
	return err
}

func (c *grpcConn) NewStream(cc context.Context, desc *G.StreamDesc, method string, opts ...G.CallOption) (G.ClientStream, error) {
	ctx := GetContext(cc)
	if ctx == nil {
		return c.conn.NewStream(cc, desc, method, opts...)
	}
	st := ctx.Step("grpc.NewStream")
	s, err := c.conn.NewStream(cc, desc, method, opts...)
	if err != nil {
		st("grpc.NewStream", "[%s] [err:1] %s", method, err.Error())
	} else {
		st("grpc.NewStream", "[%s]", method)
	}
	return s, err
}

type grpcConfig struct {
	Addr string `json:"addr"`
}

type grpcService struct {
	config interface{}
	name   string
	conn   *grpcConn
}

func newGRPCService(name string, config interface{}) GRPCService {
	return &grpcService{name: name, config: config}
}

/**
* 服务名称
**/
func (s *grpcService) Name() string {
	return s.name
}

/**
* 服务配置
**/
func (s *grpcService) Config() interface{} {
	return s.config
}

/**
* 初始化服务
**/
func (s *grpcService) OnInit(ctx micro.Context) error {

	var err error = nil
	cfg := grpcConfig{}

	dynamic.SetValue(&cfg, s.config)

	conn, err := G.Dial(cfg.Addr, G.WithTransportCredentials(insecure.NewCredentials()))

	if err != nil {
		return err
	}

	s.conn = &grpcConn{conn: conn}

	return nil
}

/**
* 校验服务是否可用
**/
func (s *grpcService) OnValid(ctx micro.Context) error {
	return nil
}

func (s *grpcService) Conn() G.ClientConnInterface {
	return s.conn
}

func (s *grpcService) Recycle() {
	if s.conn != nil {
		s.conn.conn.Close()
		s.conn = nil
	}
}
