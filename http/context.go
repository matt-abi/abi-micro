package http

import (
	"net/http"

	"github.com/ability-sh/abi-micro/micro"
)

func NewContext(p micro.Payload, r *http.Request, w http.ResponseWriter) (micro.Context, error) {
	trace := r.Header.Get("Trace")
	if trace == "" {
		trace = r.Header.Get("trace")
	}
	if trace == "" {
		trace = micro.NewTrace()
		w.Header().Add("Trace", trace)
	}
	return p.NewContext(r.URL.Path, trace)
}
