package runtime

import (
	"bytes"
	"fmt"
	"strings"
	"time"

	C "context"

	"github.com/ability-sh/abi-micro/logger"
	"github.com/ability-sh/abi-micro/micro"
)

const (
	SERVICE_LOGGER = "logger"
)

type contextTag struct {
	count int
	value string
}

type context struct {
	r       *runtime
	path    string
	trace   string
	values  map[string]string
	tags    map[string]*contextTag
	keys    []string
	b       *bytes.Buffer
	ctime   int64
	payload micro.Payload
	logger  logger.Logger
	id      string
	step    micro.Step
	ctx     C.Context
}

func newContext(r *runtime, path string, trace string, payload micro.Payload) micro.Context {
	c := &context{r: r,
		path:    path,
		trace:   trace,
		values:  map[string]string{"trace": trace},
		tags:    map[string]*contextTag{},
		b:       bytes.NewBuffer(nil),
		keys:    []string{},
		ctime:   time.Now().UnixNano(),
		payload: payload,
		id:      micro.NewTrace(),
		ctx:     C.Background(),
	}
	c.logger, _ = logger.GetLogger(c, SERVICE_LOGGER)
	r.ch <- 1
	c.Printf("[stat] [step:in]")
	return c
}

func (c *context) Path() string {
	return c.path
}

func (c *context) Trace() string {
	return c.trace
}

func (c *context) GetValue(key string) string {
	return c.values[key]
}

func (c *context) SetValue(key string, value string) {
	c.values[key] = value
}

func (c *context) Each(fn func(key string, value string) bool) {
	for key, value := range c.values {
		if !fn(key, value) {
			break
		}
	}
}

func (c *context) GetService(name string) (micro.Service, error) {
	s, err := c.r.GetService(name)
	if err != nil {
		return nil, err
	}
	err = s.OnValid(c)
	if err != nil {
		return nil, err
	}
	return s, nil
}

func (c *context) Runtime() micro.Runtime {
	return c.r
}

func (c *context) Payload() micro.Payload {
	return c.payload
}

func (c *context) beginLog() {

	c.b.Reset()

	c.b.WriteString(fmt.Sprintf("[UV] [%s] [%s] [%s] [%s] [%s] [node:%s]", time.Now().Format("2006-01-02 15:04:05.000"), c.r.Name(), c.trace, c.path, c.id, c.r.Node()))

	for _, key := range c.keys {
		t := c.tags[key]
		if t.value != "" {
			c.b.WriteString(fmt.Sprintf(" [%s:%s]", key, t.value))
		} else {
			c.b.WriteString(fmt.Sprintf(" [%s:%d]", key, t.count))
		}
	}

}

func (c *context) log(text string) {
	if c.logger != nil {
		c.logger.Output(text)
	} else {
		fmt.Println(text)
	}
}

func (c *context) Println(v ...interface{}) {

	c.beginLog()

	for _, i := range v {
		c.b.WriteString(strings.ReplaceAll(fmt.Sprintf(" %s", i), "\n", " "))
	}

	c.log(c.b.String())

}

func (c *context) Printf(format string, v ...interface{}) {

	c.beginLog()

	c.b.WriteString(" ")
	c.b.WriteString(strings.ReplaceAll(fmt.Sprintf(format, v...), "\n", " "))

	c.log(c.b.String())

}

func (c *context) Recycle() {
	c.Statf("done", "")
	c.r.ch <- -1
}

func (c *context) Statf(step string, format string, v ...interface{}) {
	tv := time.Now().UnixNano() - c.ctime
	c.Printf("[stat] [step:%s] [ms:%d] [us:%d] [ns:%d] %s", step, tv/int64(time.Millisecond), tv/int64(time.Microsecond), tv, fmt.Sprintf(format, v...))
}

func (c *context) Step(step string) micro.Step {
	ctime := time.Now().UnixNano()
	return func(format string, v ...interface{}) {
		tv := time.Now().UnixNano() - ctime
		c.Printf("[stat] [step:%s] [ms:%d] [us:%d] [ns:%d] %s", step, tv/int64(time.Millisecond), tv/int64(time.Microsecond), tv, fmt.Sprintf(format, v...))
	}
}

func (c *context) BeginStep(step string) {
	c.step = c.Step(step)
}

func (c *context) EndStep(format string, v ...interface{}) {
	if c.step != nil {
		c.step(format, v...)
		c.step = nil
	} else {
		c.Printf(format, v...)
	}
}

func (c *context) AddCount(key string, count int) {
	tag, ok := c.tags[key]
	if ok {
		tag.count += count
	} else {
		c.keys = append(c.keys, key)
		c.tags[key] = &contextTag{count: count}
	}
}

func (c *context) Id() string {
	return c.id
}

func (c *context) AddTag(key string, value string) {
	tag, ok := c.tags[key]
	if ok {
		tag.value = value
	} else {
		c.keys = append(c.keys, key)
		c.tags[key] = &contextTag{value: value}
	}
}

func (c *context) Ctx() C.Context {
	return c.ctx
}

func (c *context) WithValue(key interface{}, value interface{}) C.Context {
	c.ctx = C.WithValue(c.ctx, key, value)
	return c.ctx
}
