package runtime

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"log"
	"sort"
	"strconv"
	"time"

	"github.com/ability-sh/abi-lib/dynamic"
	"github.com/ability-sh/abi-lib/errors"
	"github.com/ability-sh/abi-lib/http"
	"github.com/ability-sh/abi-micro/micro"
)

type acContainer struct {
	Info interface{} `json:"info"`
	Ver  int         `json:"ver"`
}

type acPayload struct {
	c chan int8
	p micro.Payload
}

func getSignAcData(data map[string]string, secret string) string {

	m := md5.New()

	keys := []string{}

	for key, _ := range data {
		keys = append(keys, key)
	}

	sort.Strings(keys)

	op_and := []byte("&")
	op_eq := []byte("=")

	for i, key := range keys {
		if i != 0 {
			m.Write(op_and)
		}
		m.Write([]byte(key))
		m.Write(op_eq)
		m.Write([]byte(data[key]))
	}

	m.Write(op_and)
	m.Write([]byte(secret))

	return hex.EncodeToString(m.Sum(nil))
}

func getAcContainerInfo(baseURL string, containerId string, secret string, ver int) (*acContainer, error) {

	data := map[string]string{"id": containerId, "timestamp": strconv.FormatInt(time.Now().UnixNano()/int64(time.Millisecond), 10), "ver": strconv.Itoa(ver)}
	data["sign"] = getSignAcData(data, secret)

	res, err := http.NewHTTPRequest("GET").
		SetURL(fmt.Sprintf("%s/store/container/info/get.json", baseURL), data).
		Send()

	// log.Println(fmt.Sprintf("%s/store/container/info/get.json", baseURL), data)

	if err != nil {
		return nil, err
	}

	rs, err := res.PraseBody()

	if err != nil {
		return nil, err
	}

	// log.Println(rs)

	errno := dynamic.IntValue(dynamic.Get(rs, "errno"), 0)

	if errno != 200 {
		return nil, errors.Errorf(int32(errno), dynamic.StringValue(dynamic.Get(rs, "errmsg"), "Internal service error"))
	}

	r := &acContainer{}

	dynamic.SetValue(r, dynamic.Get(rs, "data"))

	return r, nil
}

func NewAcPayload(baseURL string, containerId string, secret string, p micro.Payload) (micro.Payload, error) {

	acContainer, err := getAcContainerInfo(baseURL, containerId, secret, 0)

	if err != nil {
		return nil, err
	}

	err = p.SetConfig(acContainer.Info)

	if err != nil {
		return nil, err
	}

	c := make(chan int8)

	rs := &acPayload{p: p, c: c}

	go func() {

		T := time.NewTicker(time.Second * 12)

		running := true

		for running {
			select {
			case <-T.C:
				{
					r, err := getAcContainerInfo(baseURL, containerId, secret, acContainer.Ver)
					if err != nil {
						log.Println("wait 12 seconds and try again", err)
					} else if r.Ver != acContainer.Ver {
						acContainer.Info = r.Info
						acContainer.Ver = r.Ver
						err = p.SetConfig(acContainer.Info)
						if err != nil {
							log.Println("wait 12 seconds and try again", err)
						}
					}
				}
			case <-c:
				running = false
			}
		}

		close(c)

	}()

	return rs, nil
}

func (p *acPayload) SetConfig(config interface{}) error {
	return p.p.SetConfig(config)
}

func (p *acPayload) NewContext(name string, trace string) (micro.Context, error) {
	return p.p.NewContext(name, trace)
}

func (p *acPayload) Exit() {
	p.c <- 1
	p.p.Exit()
}

func (p *acPayload) GetValue(key string) interface{} {
	return p.p.GetValue(key)
}

func (p *acPayload) SetValue(key string, value interface{}) {
	p.p.SetValue(key, value)
}
