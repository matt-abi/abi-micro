package crontab

import (
	"log"

	"github.com/ability-sh/abi-lib/dynamic"
	"github.com/ability-sh/abi-lib/http"
	"github.com/ability-sh/abi-micro/micro"
)

type HttpJob struct {
	p micro.Payload `json:"-"`

	Method string            `json:"method"`
	Query  map[string]string `json:"query"`
	Body   map[string]string `json:"body"`
	Url    string            `json:"url"`
}

func NewHttpJob(p micro.Payload, config interface{}) *HttpJob {
	r := &HttpJob{p: p}
	dynamic.SetValue(r, config)
	return r
}

func (job *HttpJob) Run() {

	ctx, err := job.p.NewContext("__job__", micro.NewTrace())

	if err != nil {
		log.Println(err)
		return
	}

	defer ctx.Recycle()

	ctx.Println(job.Method, job.Url, job.Query)

	req := http.NewHTTPRequest(job.Method).SetURL(job.Url, job.Query).SetHeaders(map[string]string{"Connection": "keepalive", "Trace": ctx.Trace()})

	if job.Method == "POST" {
		req.SetUrlencodeBody(job.Body)
	}

	res, err := req.Send()

	if err != nil {
		ctx.Println("[err:1]", err)
	} else {
		body, _ := res.PraseBody()
		ctx.Printf("%d %s", res.Code(), body)
	}

}
