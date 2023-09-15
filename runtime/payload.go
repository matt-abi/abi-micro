package runtime

import (
	"fmt"
	"sync"

	"github.com/matt-abi/abi-micro/micro"
)

type payload struct {
	lock   sync.RWMutex
	curr   micro.Runtime
	values map[string]interface{}
}

func NewPayload() micro.Payload {
	return &payload{values: map[string]interface{}{}}
}

func (c *payload) SetConfig(config interface{}) error {
	r, err := NewRuntime(config, c)
	if err != nil {
		return err
	}
	c.lock.Lock()
	defer c.lock.Unlock()
	if c.curr != nil {
		c.curr.Exit()
	}
	c.curr = r
	return nil
}

func (c *payload) NewContext(name string, trace string) (micro.Context, error) {
	c.lock.RLock()
	defer c.lock.RUnlock()
	if c.curr == nil {
		return nil, fmt.Errorf("not found runtime")
	}
	return c.curr.NewContext(name, trace), nil
}

func (c *payload) Exit() {

	var C chan int8

	c.lock.Lock()
	for _, v := range c.values {
		r := v.(micro.Recycle)
		if r != nil {
			r.Recycle()
		}
	}
	if c.curr != nil {
		C := make(chan int8)
		c.curr.ExitWait(C)
		c.curr = nil
	}
	c.lock.Unlock()

	if C != nil {
		<-C
		close(C)
	}
}

func (c *payload) GetValue(key string) interface{} {
	c.lock.RLock()
	defer c.lock.RUnlock()
	return c.values[key]
}

func (c *payload) SetValue(key string, value interface{}) {
	c.lock.Lock()
	defer c.lock.Unlock()
	c.values[key] = value
}
