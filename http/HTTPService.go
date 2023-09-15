package http

import (
	"fmt"

	"github.com/ability-sh/abi-lib/http"
	"github.com/ability-sh/abi-micro/micro"
)

type HTTPService interface {
	micro.Service
	Request(ctx micro.Context, method string) http.HTTPRequest
}

func GetHTTPService(ctx micro.Context, name string) (HTTPService, error) {
	s, err := ctx.GetService(name)
	if err != nil {
		return nil, err
	}
	ss, ok := s.(HTTPService)
	if ok {
		ctx.AddCount("http", 1)
		return ss, nil
	}
	return nil, fmt.Errorf("service %s not instanceof HTTPService", name)
}
