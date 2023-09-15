package micro

/**
 * 运行时
 **/
type Runtime interface {
	/**
	 * 配置信息
	 **/
	Config() interface{}
	/**
	 * 名称
	 **/
	Name() string
	/**
	 * 节点
	 **/
	Node() string
	/**
	 * 创建唯一ID
	 **/
	NewID() int64
	/**
	 * 区域ID
	 **/
	Aid() int64
	/**
	 * 节点ID
	 **/
	Nid() int64
	/**
	 * 新建上下文
	 **/
	NewContext(path string, trace string) Context
	/**
	 * 获取服务
	 **/
	GetService(name string) (Service, error)
	/**
	 * 获取服务运行器
	 **/
	GetExecutor(name string) (Executor, error)
	/**
	 * 退出
	 **/
	Exit()
	/**
	 * 退出, 退出后 C <- 1
	 **/
	ExitWait(C chan int8)

	GetValue(key string) interface{}
	SetValue(key string, value interface{})
}
