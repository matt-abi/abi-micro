package runtime

import (
	"fmt"
	"log"
	"net"
	"os"
	"strconv"
	"strings"
	"sync"

	"github.com/ability-sh/abi-lib/dynamic"
	"github.com/ability-sh/abi-lib/eval"
	"github.com/ability-sh/abi-lib/iid"
	"github.com/ability-sh/abi-micro/micro"
)

type runtime struct {
	config  interface{}
	ss      map[string]micro.Service
	es      map[string]micro.Executor
	ch      chan int8
	name    string
	node    string
	C       chan int8
	lock    sync.RWMutex
	values  map[string]interface{}
	payload micro.Payload
	IID     *iid.IID
	aid     int64
	nid     int64
}

func GetLocalIP() (string, error) {

	addrs, err := net.InterfaceAddrs()

	if err != nil {
		return "", err
	}

	for _, addr := range addrs {

		ip, ok := addr.(*net.IPNet)

		if ok && !ip.IP.IsLoopback() && ip.IP.To4() != nil {
			ipv4 := ip.IP.String()
			if ipv4 == "127.0.0.1" {
				continue
			}
			return ipv4, nil
		}

	}

	return "", fmt.Errorf("not found local ip")

}

func NewRuntime(config interface{}, payload micro.Payload) (micro.Runtime, error) {

	var err error = nil

	ss := map[string]micro.Service{}
	es := map[string]micro.Executor{}

	dynamic.Each(dynamic.Get(config, "services"), func(key interface{}, item interface{}) bool {

		var s micro.Service = nil

		name := dynamic.StringValue(key, "")
		stype := dynamic.StringValue(dynamic.Get(item, "type"), "")

		s, err = micro.NewService(stype, name, item)

		if err != nil {
			return false
		}

		ss[name] = s

		log.Println(name, "=>", "Executor")

		es[name] = NewReflectExecutor(s)

		return true
	})

	if err != nil {
		return nil, err
	}

	name := eval.ParseEval(dynamic.StringValue(dynamic.Get(config, "name"), ""), func(key string) string {
		return os.Getenv(key)
	})

	node := eval.ParseEval(dynamic.StringValue(dynamic.Get(config, "node"), ""), func(key string) string {
		return os.Getenv(key)
	})

	if node == "" {
		node, err = GetLocalIP()
		if err != nil {
			node = fmt.Sprintf("%s-0", name)
		}
	}

	aid := dynamic.IntValue(dynamic.Get(config, "aid"), 0)
	nid := dynamic.IntValue(dynamic.Get(config, "nid"), 0)

	if nid == 0 {
		i := strings.LastIndex(node, "-")
		if i >= 0 {
			nid, _ = strconv.ParseInt(node[i+1:], 10, 64)
		} else {
			i := strings.LastIndex(node, ".")
			if i >= 0 {
				nid, _ = strconv.ParseInt(node[i+1:], 10, 64)
			}
		}
	}

	ch := make(chan int8, 32)

	r := &runtime{ss: ss,
		es:      es,
		config:  config,
		ch:      ch,
		name:    name,
		node:    node,
		aid:     aid,
		nid:     nid,
		IID:     iid.NewIID(aid, nid),
		values:  map[string]interface{}{},
		payload: payload}

	go func() {

		var s int8 = 0
		var count int64 = 0
		var loopbreak = false

		for count > 0 || !loopbreak {

			s = <-ch

			if s == 0 {
				loopbreak = true
			} else {
				count += int64(s)
			}

		}

		for _, i := range ss {
			i.Recycle()
		}

		for _, v := range r.values {
			i := v.(micro.Recycle)
			if i != nil {
				i.Recycle()
			}
		}

		close(ch)

		if r.C != nil {
			r.C <- 1
		}
	}()

	{

		trace := micro.NewTrace()

		ctx := r.NewContext("__init__", trace)

		defer ctx.Recycle()

		init_ss := map[string]bool{}

		var onInit func(key string, s micro.Service) error

		onInit = func(key string, s micro.Service) error {
			var err error = nil
			if !init_ss[key] {
				init_ss[key] = true
				dynamic.Each(dynamic.Get(s.Config(), "dependencies"), func(_ interface{}, name interface{}) bool {
					skey := dynamic.StringValue(name, "")
					if !init_ss[skey] {
						s := ss[skey]
						if s != nil {
							err = onInit(skey, s)
							if err != nil {
								return false
							}
						}
						return true
					}
					return true
				})
				if err != nil {
					return err
				}

				ctx.Printf("%s init ...", key)
				err = s.OnInit(ctx)
				if err != nil {
					ctx.Printf("%s init %s %s", key, s.Config(), err)
					return err
				}
				ctx.Printf("%s init done", key)
			}
			return err
		}

		for key, s := range ss {
			err := onInit(key, s)
			if err != nil {
				r.Exit()
				return nil, err
			}
		}

	}

	return r, nil
}

func (r *runtime) Config() interface{} {
	return r.config
}

func (r *runtime) Name() string {
	return r.name
}

func (r *runtime) NewContext(path string, trace string) micro.Context {
	return newContext(r, path, trace, r.payload)
}

func (r *runtime) GetService(name string) (micro.Service, error) {
	s := r.ss[name]
	if s == nil {
		return nil, fmt.Errorf("not found service %s", name)
	}
	return s, nil
}

func (r *runtime) GetExecutor(name string) (micro.Executor, error) {
	s := r.es[name]
	if s == nil {
		return nil, fmt.Errorf("not found service %s", name)
	}
	return s, nil
}

func (r *runtime) Exit() {
	r.ch <- 0
}

func (r *runtime) ExitWait(C chan int8) {
	r.C = C
	r.Exit()
}

func (r *runtime) GetValue(key string) interface{} {
	r.lock.RLock()
	defer r.lock.RUnlock()
	return r.values[key]
}

func (r *runtime) SetValue(key string, value interface{}) {
	r.lock.Lock()
	defer r.lock.Unlock()
	r.values[key] = value
}

func (r *runtime) Node() string {
	return r.node
}

func (r *runtime) NewID() int64 {
	return r.IID.NewID()
}

func (r *runtime) Aid() int64 {
	return r.aid
}

func (r *runtime) Nid() int64 {
	return r.nid
}
