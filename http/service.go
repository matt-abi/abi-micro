package http

import (
	xhttp "net/http"
	"net/url"

	"github.com/ability-sh/abi-lib/dynamic"
	"github.com/ability-sh/abi-lib/http"
	"github.com/ability-sh/abi-micro/micro"
)

var httpClient = xhttp.DefaultClient

type httpService struct {
	config interface{}
	name   string

	Proxy string `json:"proxy"`

	client *xhttp.Client
}

func newHTTPService(name string, config interface{}) HTTPService {
	return &httpService{name: name, config: config}
}

/**
* 服务名称
**/
func (s *httpService) Name() string {
	return s.name
}

/**
* 服务配置
**/
func (s *httpService) Config() interface{} {
	return s.config
}

/**
* 初始化服务
**/
func (s *httpService) OnInit(ctx micro.Context) error {

	dynamic.SetValue(s, s.config)

	if s.Proxy != "" {
		u, err := url.Parse(s.Proxy)
		if err != nil {
			return err
		}
		s.client = http.NewClientWithProxy(u)
	}

	return nil
}

/**
* 校验服务是否可用
**/
func (s *httpService) OnValid(ctx micro.Context) error {
	return nil
}

func (s *httpService) Request(ctx micro.Context, method string) http.HTTPRequest {
	req := http.NewHTTPRequest(method).SetHeaders(map[string]string{"Trace": ctx.Trace()})
	if s.client != nil {
		req.SetClient(s.client)
	}
	return req
}

func (s *httpService) Recycle() {

}
