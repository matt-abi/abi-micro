package micro

type Payload interface {
	SetConfig(config interface{}) error
	NewContext(name string, trace string) (Context, error)
	Exit()

	GetValue(key string) interface{}
	SetValue(key string, value interface{})
}
