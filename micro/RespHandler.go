package micro

import "net/http"

type RespFunc func(ctx Context, req *http.Request, w http.ResponseWriter, data interface{}) bool
type RespHandler interface {
	MatchRespHandler(ctx Context, name string) RespFunc
}
