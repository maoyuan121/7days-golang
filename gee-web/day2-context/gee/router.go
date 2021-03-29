package gee

import (
	"net/http"
)

// 路由。不对外
type router struct {
	handlers map[string]HandlerFunc // key 为请求方法和路径， value 为对应的 handler
}

// 路由的构造函数
func newRouter() *router {
	return &router{handlers: make(map[string]HandlerFunc)}
}

// 添加路由
// @param： method。 请求方法
// @param: pattern。路径
// @handler: 处理程序
func (r *router) addRoute(method string, pattern string, handler HandlerFunc) {
	key := method + "-" + pattern
	r.handlers[key] = handler
}

// 处理请求
func (r *router) handle(c *Context) {
	key := c.Method + "-" + c.Path
	if handler, ok := r.handlers[key]; ok {
		handler(c)
	} else {
		c.String(http.StatusNotFound, "404 NOT FOUND: %s\n", c.Path)
	}
}
