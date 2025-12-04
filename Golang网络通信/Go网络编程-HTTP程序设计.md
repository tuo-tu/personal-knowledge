# HTTP程序设计

Go编写HTTP服务器，用 Go实现一个 `http server`非常容易，Go 语言标准库 `net/http`自带了一系列结构和方法来帮助开发者简化 HTTP 服务开发的相关流程。因此，我们不需要依赖任何第三方组件就能构建并启动一个高并发的 HTTP 服务器。

## 简单的HTTP服务器

函数：

```go
// http.ListenAndServe
func ListenAndServe(addr string, handler Handler) error
```

用于启动HTTP服务器，监听addr，并使用handler来处理请求。返回启动错误。其中：

- addr，TCP address，形式为 IP:port，IP省略表示监听全部网络接口
- **handler，经常的被设置为nil，表示使用DefaultServeMux（默认服务复用器）来处理请求。**
- DefaultServeMux要使用以下两个函数来添加请求处理器
  - `func Handle(pattern string, handler Handler)`  ，第二个参数**Handler是一个接口。**
  - `func HandleFunc(pattern string, handler func(ResponseWriter, *Request))`，第二个参数是一个函数

示例代码：

httpServerSimple.go

```go
func HttpServerSimple() {
    // 一：设置不同路由（path）对应不同的处理器
    // /ping <-> pong
    http.HandleFunc("/ping", handlePing)

    // 三：使用http.Handle设置处理器对象
    http.Handle("/info", infoHandler)
    infoHandler := InfoHandler{
        info: "Welcome to Mashibing Go classroom.",
    }
    // 二：启动监听并提供服务
    addr := ":8088"
    log.Println("http server is listening on ", addr)
    err := http.ListenAndServe(addr, nil)
    log.Fatalln(err)
}

// http.ResponseWriter 响应Writer
// *http.Request 请求对象，包含了请求信息
func handlePing(w http.ResponseWriter, r *http.Request) {
    fmt.Fprintf(w, "pong")
}

// InfoHandler 实现Handler接口的类型
type InfoHandler struct {
    info string
}

func (h InfoHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
    fmt.Fprintf(w, h.info)
}
```

其中：Handler 接口的定义为：

```go
type Handler interface {
    ServeHTTP(ResponseWriter, *Request)
}
```

**我们的 InfoHandler实现了Handler接口，可以作为 http.Handle()的第二个参数来使用。**

测试，通过main.main() 启动服务器：

httpServerSimple.go

```go
func main() {
    // 简单的HTTP服务器
    HttpServerSimple()
}
```

运行

```shell
go run httpServerSimple.go
2023/03/02 21:00:29 http server is listening on  :8088
```

请求测试：

```
curl http://localhost:8088/ping
pong

curl http://localhost:8088/info
Welcome to Mashibing Go classroom.
```

## 复杂的HTTP服务器

定制性的HTTP服务器，通过 Server 类型进行设置。其定义如下：

```go
// net/http
type Server struct {
    // TCP Address
    Addr string
    Handler Handler // handler to invoke, http.DefaultServeMux if nil
    // LSConfig optionally provides a TLS configuration for use
	// by ServeTLS and ListenAndServeTLS
    TLSConfig *tls.Config
    // 读请求超时时间
    ReadTimeout time.Duration
    // 读请求头超时时间
    ReadHeaderTimeout time.Duration
    // 写响应超时时间
    WriteTimeout time.Duration
    // 空闲超时时间
    IdleTimeout time.Duration
    // Header最大长度
    MaxHeaderBytes int
    // 其他字段略
}
```

该类型的 `func (srv *Server) ListenAndServe() error` 函数用于监听和服务。

示例代码：

```go
// @file: HttpServerCustom.go
func HttpServerCustom() {
    // 1.定义http处理器，myHandler实现了ServeHTTP方法
    myHandler := CustomHandler{message: "http.Server"}
    // 2.HTTP 服务器配置（关键步骤），嵌入myHandler作为处理器
    s := &http.Server{
        Addr:           ":8080",
        Handler:        myHandler,
        ReadTimeout:    10 * time.Second,
        WriteTimeout:   10 * time.Second,
        MaxHeaderBytes: 1 << 20,
    }
    // 3.启动http服务器
    log.Fatal(s.ListenAndServe())
}

type CustomHandler struct {
    message string
}

func (h CustomHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
    time.Sleep(10 * time.Second)
    fmt.Fprintf(w, h.message)
}
```
