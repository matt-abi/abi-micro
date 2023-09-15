package grpc

import (
	"context"
	"strings"

	"github.com/ability-sh/abi-micro/micro"
	G "google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

var traceKeys = []string{"trace", "Trace"}

func NewUnaryServerInterceptor(p micro.Payload) G.UnaryServerInterceptor {

	return func(ctx context.Context, req interface{}, info *G.UnaryServerInfo, handler G.UnaryHandler) (resp interface{}, err error) {

		md, ok := metadata.FromIncomingContext(ctx)

		if !ok {
			return nil, status.Errorf(codes.InvalidArgument, "Retrieving metadata is failed")
		}

		var trace string = ""

		for _, key := range traceKeys {

			vs := md[key]

			if len(vs) > 0 {
				trace = vs[0]
				break
			}
		}

		if trace == "" {
			trace = micro.NewTrace()
		}

		name := info.FullMethod[1:]
		i := strings.Index(name, "/")

		if i >= 0 {
			name = name[i+1:]
		}

		c, err := p.NewContext(name, trace)

		if err != nil {
			return nil, err
		}

		defer c.Recycle()

		for key, vs := range md {
			if len(vs) > 0 {
				c.SetValue(key, vs[0])
			}
		}

		rs, err := handler(micro.WithContext(ctx, c), req)

		if err != nil {
			c.Printf("[err:1] %s", err.Error())
		}

		return rs, err
	}
}

func NewServer(p micro.Payload) *G.Server {
	return G.NewServer(G.UnaryInterceptor(NewUnaryServerInterceptor(p)))
}
