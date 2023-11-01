package micro

import "net/http"

type ReqVerifyFunc func(ctx Context, req *http.Request, data interface{}) error
type Verifier interface {
	MatchReqVerify(ctx Context, name string) ReqVerifyFunc
}
