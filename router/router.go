package router

type RouterCtx interface {
	Bind(interface{}) error
	JSON(int, interface{}) error
}
