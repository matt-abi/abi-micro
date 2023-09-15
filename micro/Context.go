package micro

import "context"

type Step = func(format string, v ...interface{})

/**
 * 上下文
 **/
type Context interface {
	Recycle
	Id() string
	Path() string
	Trace() string
	GetValue(key string) string
	SetValue(key string, value string)
	Each(fn func(key string, value string) bool)
	/**
	 * 增加计数，用户统计
	 **/
	AddCount(key string, count int)
	AddTag(key string, value string)
	GetService(name string) (Service, error)
	Runtime() Runtime
	Payload() Payload
	/**
	 * 日志
	 **/
	Println(v ...interface{})
	Printf(format string, v ...interface{})
	Step(step string) Step
	BeginStep(step string)
	EndStep(format string, v ...interface{})
	Ctx() context.Context
	WithValue(key interface{}, value interface{}) context.Context
}

type Key struct{}

var ContextKey = Key{}

func WithContext(c context.Context, ctx Context) context.Context {
	return context.WithValue(c, ContextKey, ctx)
}

func GetContext(c context.Context) Context {
	v, ok := c.Value(ContextKey).(Context)
	if ok {
		return v
	}
	return nil
}
