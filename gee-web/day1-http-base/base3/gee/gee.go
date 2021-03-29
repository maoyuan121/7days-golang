package gee

import (
	"fmt"
	"log"
	"net/http"
)

// HandlerFunc 定义了 request handler
type HandlerFunc func(http.ResponseWriter, *http.Request)

// Engine 实现 ServeHTTP 接口
type Engine struct {
	// 路由 键由请求方法和路径组成，值是对应的  HandlerFunc
	router map[string]HandlerFunc
}

// New 是 gee.Engine  的构造函数
func New() *Engine {
	return &Engine{router: make(map[string]HandlerFunc)}
}

// 添加路由
// @param: method。请求的 http method
// @param: pattern。请求的路由
// @param: handler。请求处理 handler
// @description 该方法为私有的，外部不能直接调用。对外暴露的是 GET POST 等方法，这些方法来调用 addRoute
func (engine *Engine) addRoute(method string, pattern string, handler HandlerFunc) {
	key := method + "-" + pattern
	log.Printf("Route %4s - %s", method, pattern)
	engine.router[key] = handler
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
	key := req.Method + "-" + req.URL.Path
	if handler, ok := engine.router[key]; ok {
		handler(w, req)
	} else {
		fmt.Fprintf(w, "404 NOT FOUND: %s\n", req.URL)
	}
}
