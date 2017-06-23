# middleware
Middleware Framework For Go net/http.

Fully compatible with ```http.Handler```. To use full featured middleware, you only need to add a no parameter ```next()``` call.

## Getting Started
```bash
go get github.com/withwind8/middleware
```

```go
package main

import (
  "log"
  "time"
  "net/http"
    
  "github.com/withwind8/middleware"
)

func main(){
  app := middleware.New()
  
  //a simple log middleware using func
  app.UseFunc(func(w http.ResponseWriter, r *http.Request, next func()){
    start := time.Now()
    
    next()
    
    log.Printf("%s %v",r.URL.Path,time.Since(start))
  })
  
  mux := http.NewServeMux()
  mux.HandleFunc("/",func(w http.ResponseWriter, r *http.Request) {
    w.Write([]byte("Hello World!"))
  })
  
  //can using standard http.Handler as middleware
  app.UseHandler(mux)
  
  log.Fatal(app.Listen(":8080"))
}  
```

## API
### func ```New()```
Return a middleware container instance
### method ```Use(middleware Middleware)```
Using a Middleware interface as middleware
```go
type Middleware interface {
	ServeHTTP(w http.ResponseWriter, r *http.Request, next func())
}
```
### method ```UseFunc(middlewareFunc MiddlewareFunc)```
Using a MiddlewareFunc as middleware
```go
type MiddlewareFunc func(w http.ResponseWriter, r *http.Request, next func())
```
### method ```UseHandler(handler http.Handler)```
Using a http.Handler interface as middleware, next() will be executed automatically at the end
### method ```UseHandlerFunc(handlerFunc http.HandlerFunc)```
Using a http.HandlerFunc as middleware, next() will be executed automatically at the end
### method ```Listen(addr string) error```
Start the server in with given address.
```go
app.Listen(":8080")
```
is Simply sugar for the following:
```go
http.ListenAndServe(":8080", app)
```

## Usage
Suppose you have 2 middleware (1 interface 1 func) and 2 handler (1 interface 1 func), and use them in the following order:
```go
app.Use(middleware1)
app.UseHandler(handler1)
app.UseFunc(middleware2)
app.UseHandlerFunc(handler2)
```
Finally, the execution sequence is as follows：
1. middleware1_before_next
2. handler1
3. middleware2_before_next
4. handler2
5. middleware2_after_next
6. middleware1_after_next



