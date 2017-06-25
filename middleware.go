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
	// execute the first middleware
	m.head.ServeHTTP(w, r)
}

// create middleware container
func New() *middlewares {
	return &middlewares{}
}

// use middleware
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

// MiddlewareFunc implement Middleware interface
func (f MiddlewareFunc) ServeHTTP(w http.ResponseWriter, r *http.Request, next func()) {
	f(w, r, next)
}

// use MiddlewareFunc as middleware
func (m *middlewares) UseFunc(middlewareFunc MiddlewareFunc) {
	m.Use(middlewareFunc)
}

type middlewareForHandler struct {
	handler http.Handler
}

func (m *middlewareForHandler) ServeHTTP(w http.ResponseWriter, r *http.Request, next func()) {
	m.handler.ServeHTTP(w, r)
	next()
}

// use http.Handler as middleware
func (m *middlewares) UseHandler(handler http.Handler) {
	m.Use(&middlewareForHandler{handler})
}

// use http.HandlerFunc as middleware
func (m *middlewares) UseHandlerFunc(handlerFunc http.HandlerFunc) {
	m.Use(&middlewareForHandler{handlerFunc})
}

// just a sugar method
func (m *middlewares) Listen(addr string) error {
	return http.ListenAndServe(addr, m)
}
