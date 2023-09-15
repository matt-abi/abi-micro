package mq

import (
	"context"
	"errors"
)

type Driver interface {
	Send(topic string, name string, data interface{}) error
	On(queue string, fn func(name string, data interface{}) bool) context.CancelFunc
	Recycle()
}

type DriverCreator func(driver string, config interface{}) (Driver, error)

var driverSet = map[string]DriverCreator{}

func Reg(driver string, driverCreator DriverCreator) {
	driverSet[driver] = driverCreator
}

func NewDriver(dirver string, config interface{}) (Driver, error) {
	fn, ok := driverSet[dirver]
	if !ok {
		return nil, errors.New("not found mq driver")
	}
	return fn(dirver, config)
}
