package router

type RouterCtx interface {
	Param(string) string
	Bind(interface{}) error
	JSON(int, interface{}) error
}
