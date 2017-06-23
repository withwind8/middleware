package middleware

import "net/http"

type middlewares struct {
	head *middleware
	tail *middleware
}

type middleware struct {
	middleware Middleware
	next       *middleware
}

type MiddlewareFunc func(http.ResponseWriter, *http.Request, func())

type Middleware interface {
	ServeHTTP(http.ResponseWriter, *http.Request, func())
}

func (m *middleware) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	m.middleware.ServeHTTP(w, r, func() {
		if m.next != nil {
			m.next.ServeHTTP(w, r)
		}
	})
}

func (m *middlewares) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	//调用第一个
	m.head.ServeHTTP(w, r)
}

//返回中间件管理器
func New() *middlewares {
	return &middlewares{}
}

//使用中间件
func (m *middlewares) Use(mw Middleware) {
	n := &middleware{mw, nil}
	if m.head == nil && m.tail == nil {
		m.head = n
		m.tail = n
	} else {
		m.tail.next = n
		m.tail = n
	}
}

type fakeMiddleware struct {
	middlewareFunc MiddlewareFunc
}

func (m *fakeMiddleware) ServeHTTP(w http.ResponseWriter, r *http.Request, next func()) {
	m.middlewareFunc(w, r, next)
}

//通过中间件函数使用中间件
func (m *middlewares) UseFunc(middlewareFunc MiddlewareFunc) {
	m.Use(&fakeMiddleware{middlewareFunc})
}

type fakeMiddlewareForHandler struct {
	handler http.Handler
}

func (m *fakeMiddlewareForHandler) ServeHTTP(w http.ResponseWriter, r *http.Request, next func()) {
	m.handler.ServeHTTP(w, r)
	next()
}

//把http.Handler当中间件使用
func (m *middlewares) UseHandler(handler http.Handler) {
	m.Use(&fakeMiddlewareForHandler{handler})
}

type fakeMiddlewareForHandlerFunc struct {
	handlerFunc http.HandlerFunc
}

func (m *fakeMiddlewareForHandlerFunc) ServeHTTP(w http.ResponseWriter, r *http.Request, next func()) {
	m.handlerFunc(w, r)
	next()
}

//把http.HandlerFunc当中间件使用
func (m *middlewares) UseHandlerFunc(handlerFunc http.HandlerFunc) {
	m.Use(&fakeMiddlewareForHandlerFunc{handlerFunc})
}
