package micro

/**
 * 服务
 **/
type Service interface {
	Recycle
	/**
	 * 服务名称
	 **/
	Name() string
	/**
	 * 服务配置
	 **/
	Config() interface{}
	/**
	 * 初始化服务
	 **/
	OnInit(ctx Context) error
	/**
	 * 校验服务是否可用
	 **/
	OnValid(ctx Context) error
}
