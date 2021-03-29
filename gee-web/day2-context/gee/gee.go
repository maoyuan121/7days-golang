package gee

import (
	"log"
	"net/http"
)

// HandlerFunc 定义了请求 handler
type HandlerFunc func(*Context)

// Engine 实现 ServeHTTP 接口
type Engine struct {
	router *router
}

// New 是 Engine 的构造函数
func New() *Engine {
	return &Engine{router: newRouter()}
}

// 底层调用 router 的 addRoute 方法
func (engine *Engine) addRoute(method string, pattern string, handler HandlerFunc) {
	log.Printf("Route %4s - %s", method, pattern)
	engine.router.addRoute(method, pattern, handler)
}

// 添加 get 路由
func (engine *Engine) GET(pattern string, handler HandlerFunc) {
	engine.addRoute("GET", pattern, handler)
}

// 添加 post 路由
func (engine *Engine) POST(pattern string, handler HandlerFunc) {
	engine.addRoute("POST", pattern, handler)
}

// 运行 http 服务
func (engine *Engine) Run(addr string) (err error) {
	return http.ListenAndServe(addr, engine)
}

// 实现 ServeHTTP 接口
func (engine *Engine) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	c := newContext(w, req)
	engine.router.handle(c)
}
