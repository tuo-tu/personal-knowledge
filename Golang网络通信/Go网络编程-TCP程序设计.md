# TCP程序设计

## TCP概述

## 建立连接

![image.png](https://fynotefile.oss-cn-zhangjiakou.aliyuncs.com/fynote/fyfile/13080/1682659809092/db91c90e503c49d6b005b47bf4a8a705.png)

客户端和服务器端在建立连接时：

- **服务端是典型的监听+接受连接的模式**，就是Listen+Accept
- **客户端是主动建立连接的模式，就是Dial**

**Go语言中使用 `net`包实现网络的相关操作**，包括我们TCP的操作。

用于建立连接的典型方法如下：

```go
// 监听某一种网络的某一个地址
func Listen(network, address string) (Listener, error)
// 接受监听到的连接。
func (l *TCPListener) Accept() (Conn, error)

// 连接网络
func Dial(network, address string) (Conn, error)
// 带有超时的连接网络
func DialTimeout(network, address string, timeout time.Duration) (Conn, error)


func ListenTCP(network string, laddr *TCPAddr) (*TCPListener, error)
func (l *TCPListener) AcceptTCP() (*TCPConn, error)
```

### 服务端程序

示例代码：

```go
// 服务端
func TcpServer() {
	// A. 基于某个地址建立监听
	// 服务端地址（就是监听地址）
	address := "127.0.0.1:5678"
	listener, err := net.Listen(tcp, address)
	if err != nil {
		log.Fatalln(err)
	}
	// 关闭监听
	defer listener.Close()
	log.Printf("%s server is listening on %s\n", tcp, listener.Addr())

	// B. 接受连接请求
	// 循环接受
	for {
		// 阻塞接受
		conn, err := listener.Accept()
		if err != nil {
			log.Println(err)
		}

		// 处理连接，读写
		// 日志连接的远程地址（client addr）
		log.Printf("accept from %s\n", conn.RemoteAddr())
	}
}
```

其中：

- address的表述方式
  - IP:Port 明确的IP和端口。
  - IP: 明确了IP端口任意。
  - :port 明确了端口IP全部
- listener.Addr() 监听的地址
- conn.RemoteAddr() 连接的远程地址

### 客户端程序

示例代码：

```go
// 客户端
func TcpClient() {
	// tcp服务端地址
	address := "127.0.0.1:5678"
	// 模拟多客户端
	// 并发的客户端请求
	num := 10
	wg := sync.WaitGroup{}
	wg.Add(num)
	for i := 0; i < num; i++ {
		// 并发请求
		go func(wg *sync.WaitGroup) {
			defer wg.Done()
			// A. 建立连接
			conn, err := net.Dial(tcp, address)
			if err != nil {
				log.Println(err)
				return
			}
			// 保证关闭
			defer conn.Close()
			log.Printf("connection is establish, client addr is %s\n", conn.LocalAddr())
		}(&wg)
	}

	wg.Wait()
}
```

其中：

- conn.Close()，关闭连接，连接资源使用完毕要记得关闭
- conn.LocalAddr()， 用于获得客户端本地地址，会与服务端的RemoteAddr对应

### 测试

先开启服务端程序，再开启客户端程序：

```shell
func TestTcpServer(t *testing.T) {
	TcpServer()
}

func TestTcpClient(t *testing.T) {
	TcpClient()
}
```

Server:

```shell
netProgram> go test -run TcpServer
2023/04/28 14:24:12 tcp server is listening on 127.0.0.1:5678
2023/04/28 14:24:17 accept from 127.0.0.1:50690
2023/04/28 14:24:17 accept from 127.0.0.1:50689
2023/04/28 14:24:17 accept from 127.0.0.1:50694
2023/04/28 14:24:17 accept from 127.0.0.1:50695
2023/04/28 14:24:17 accept from 127.0.0.1:50692
2023/04/28 14:24:17 accept from 127.0.0.1:50687
2023/04/28 14:24:17 accept from 127.0.0.1:50688
2023/04/28 14:24:17 accept from 127.0.0.1:50696
2023/04/28 14:24:17 accept from 127.0.0.1:50691
2023/04/28 14:24:17 accept from 127.0.0.1:50693
```

Client:

```shell
netProgram> go test -run TcpClient
2023/04/28 14:24:17 connection is establish, client addr is 127.0.0.1:50695
2023/04/28 14:24:17 connection is establish, client addr is 127.0.0.1:50694
2023/04/28 14:24:17 connection is establish, client addr is 127.0.0.1:50689
2023/04/28 14:24:17 connection is establish, client addr is 127.0.0.1:50691
2023/04/28 14:24:17 connection is establish, client addr is 127.0.0.1:50696
2023/04/28 14:24:17 connection is establish, client addr is 127.0.0.1:50692
2023/04/28 14:24:17 connection is establish, client addr is 127.0.0.1:50693
2023/04/28 14:24:17 connection is establish, client addr is 127.0.0.1:50688
2023/04/28 14:24:17 connection is establish, client addr is 127.0.0.1:50687
2023/04/28 14:24:17 connection is establish, client addr is 127.0.0.1:50690
```

注：并由于发编程的调度，不能保证服务端的日志顺序与客户端一致。因为建立连接和输出日志不是在一个原子操作中进行的。

### tcp网络支持

函数：

```go
func Listen(network, address string) (Listener, error)
func Dial(network, address string) (Conn, error)
```

参数 network 表示网络类型, 支持的TCP类型字符串:

- tcp, 使用IPv4或IPv6
- tcp4, 仅使用IPv4
- tcp6, 仅使用IPv6
- **省略IP部分, 绑定可用的全部IP, 包括IPv4和IPv6**

客户端在建立连接时使用的网络类型，要与服务器监听的网络类型能够匹配。

示例代码：

```go
// tcp协议类型
// address := "127.0.0.1:5678" // IPv4
// address := "[::1]:5678" // IPv6
address := ":5678" // Any IP or version

```

### 连接失败

当客户端net.Dial()建立连接时, 还有可能会失败, 典型的失败原因:

- 服务器端未启动, 或网络连接失败
- 网络原因超时
- 并发连接的客户端太多, 服务端处理不完

示例错误: 服务器端未启动, 或网络连接失败, 连接超时等:

```shell
# 无连接目标可用
No connection could be made because the target machine actively refused it.

# 网络不可达
A socket operation was attempted to an unreachable network.

# 超时
dial tcp 127.0.0.1:56789: i/o timeout
```

#### net.DialTimetout

设置超时时间.

```go
// 带有超时的连接网络
func DialTimeout(network, address string, timeout time.Duration) (Conn, error)
```

示例:

```go
func TcpTimeoutClient() {
	// tcp服务端地址
	serverAddress := "192.168.110.123:5678" // IPv6 4

	// 模拟多客户端
	// 并发的客户端请求
	num := 10
	wg := sync.WaitGroup{}
	wg.Add(num)
	for i := 0; i < num; i++ {
		// 并发请求
		go func(wg *sync.WaitGroup) {
			defer wg.Done()
			// A. 建立连接
			conn, err := net.DialTimeout(tcp, serverAddress, time.Second)
			//conn, err := net.Dial(tcp, serverAddress)
			if err != nil {
				log.Println(err)
				return
			}
			// 保证关闭
			defer conn.Close()
			log.Printf("connection is establish, client addr is %s\n", conn.LocalAddr())
		}(&wg)
	}

	wg.Wait()
}
```

#### 未能及时Accept

客户端发出的连接,若服务器端未能及时Accept, 会被**缓存到队列中**。当队列存满时，就不会在接受客户端连接了.

**这个队列大小的配置,就叫Backlog**

示例：

```go
// 服务端
func TcpBacklogServer() {
	// A. 基于某个地址建立监听
	// 服务端地址
	address := ":5678" // Any IP or version
	listener, err := net.Listen(tcp, address)
	if err != nil {
		log.Fatalln(err)
	}
	// 关闭监听
	defer listener.Close()
	log.Printf("%s server is listening on %s\n", tcp, listener.Addr())

	// B. 接受连接请求
	// 循环接受
	for {
		// 阻塞接受
		conn, err := listener.Accept()
		if err != nil {
			log.Println(err)
		}

		// 处理连接，读写
		func(conn net.Conn) {
			// 日志连接的远程地址（client addr）
			log.Printf("accept from %s\n", conn.RemoteAddr())
			time.Sleep(time.Second)
		}(conn)
	}
}

func TcpBacklogClient() {
	// tcp服务端地址
	serverAddress := "127.0.0.1:5678" // IPv6 4

	// 模拟多客户端
	// 并发的客户端请求
	num := 256
	wg := sync.WaitGroup{}
	wg.Add(num)
	for i := 0; i < num; i++ {
		// 并发请求
		go func(wg *sync.WaitGroup, no int) {
			defer wg.Done()
			// A. 建立连接
			conn, err := net.DialTimeout(tcp, serverAddress, time.Second)
			//conn, err := net.Dial(tcp, serverAddress)
			if err != nil {
				log.Println(err)
				return
			}
			// 保证关闭
			defer conn.Close()
			log.Printf("%d: connection is establish, client addr is %s\n", no, conn.LocalAddr())
		}(&wg, i)

		time.Sleep(30 * time.Millisecond)
	}

	wg.Wait()
}
```

在授课的测试电脑中, Backlog的值为200, Linux系统通常为128.

到达上限,需要等待服务端Accept某个连接后,才会有新的客户端进入.

go中的典型解决方案为并发处理每个连接. 示例代码:

```go
		// 处理连接，读写
		func(conn net.Conn) {
			// 日志连接的远程地址（client addr）
			log.Printf("accept from %s\n", conn.RemoteAddr())
			time.Sleep(time.Second)
		}(conn)
```

## 读写操作

### 基本示例

当建立了客户端与服务端的连接后，就需要相互发送数据了。**TCP协议是全双工通信**，就是连接两端允许同时进行双向数据传输（读写）。

Go程序设计时，服务端通常使用独立的Goroutine处理每个客户端的连接及使用该连接的读写操作。

conn，提供了读写方法：

```go
// 从conn读内容至b， 返回读取长度和错误
Read(b []byte) (n int, err error)
// 向conn写入数据b，返回写入长度和错误
Write(b []byte) (n int, err error)

```

示例：

```go
// server
// 处理每个连接
func HandleConn(conn net.Conn) {
	// 日志连接的远程地址（client addr）
	log.Printf("accept from %s\n", conn.RemoteAddr())

	// A.保证连接关闭
	defer conn.Close()

	// B.向客户端发送数据，Write
	wn, err := conn.Write([]byte("send some data from server" + "\n"))
	if err != nil {
		log.Println(err)
	}
	log.Printf("server write len is %d\n", wn)

	// C.从客户端接收数据，Read
	buf := make([]byte, 1024)
	rn, err := conn.Read(buf)
	if err != nil {
		log.Println(err)
	}
	log.Println("received from client data is:", string(buf[:rn]))
}


// client
func TcpClientRW() {
	// tcp服务端地址
	serverAddress := "127.0.0.1:5678" // IPv6 4

	// A. 建立连接
	conn, err := net.DialTimeout(tcp, serverAddress, time.Second)
	//conn, err := net.Dial(tcp, serverAddress)
	if err != nil {
		log.Println(err)
		return
	}
	// 保证关闭
	defer conn.Close()
	log.Printf("connection is establish, client addr is %s\n", conn.LocalAddr())

	// B.从服务端接收数据，Read
	buf := make([]byte, 1024)
	rn, err := conn.Read(buf)
	if err != nil {
		log.Println(err)
	}
	log.Println("received from server data is:", string(buf[:rn]))

	// C.向服务器端发送数据，Write
	wn, err := conn.Write([]byte("send some data from client" + "\n"))
	if err != nil {
		log.Println(err)
	}
	log.Printf("client write len is %d\n", wn)
}

```

测试结果：

```shell
# server
> go test -run TcpServerRW
2023/05/03 13:15:36 tcp server is listening on [::]:5678
2023/05/03 13:15:41 accept from 127.0.0.1:50932
2023/05/03 13:15:41 server write len is 27
2023/05/03 13:15:41 received from client data is: send some data from client

# client
> go test -run TcpClientRW
2023/05/03 13:15:41 connection is establish, client addr is 127.0.0.1:50932
2023/05/03 13:15:41 received from server data is: send some data from server

2023/05/03 13:15:41 client write len is 27
```

我们在Server和Client端，都可以使用Read和Write方法，基于conn完成读写操作。

### Write和Read的注意事项

**Write特点**

- 写成功, **err ==nil && wn == len(data)** 表示写入成功
- 写阻塞，当无法继续写时，Write会进入阻塞状态。无法继续写，**通常意味着TCP的窗口已满.**
- 已关闭的连接不能继续写入
- 可以使用如下方法控制Write的超时时长

  - `SetDeadline(t time.Time) error`
  - `SetWriteDeadline(t time.Time) error`

**Read特点**

- 当conn中无数据时，Read处于阻塞状态
- **当conn中有足够数据时，Read读满buf，并返回读取长度，需要循环读取，才可以读取全部内容**
- 当conn中有部分数据时，Read读部分数据，并返回读取长度
- 当conn已经关闭时，通常会返回EOF error
- 可以使用如下方法控制Read的超时时长
  - `SetDeadline(t time.Time) error`
  - `SetReadDeadline(t time.Time) error`

示例代码

```go
// 1. 严谨的判断是否写入成功
data := []byte("send some data from server" + "\n")
wn, err := conn.Write(data)
if err != nil {
	log.Println(err)
}
// 若要严谨的判断是否写入成功,需要:
if err == nil && wn == len(data) {
	log.Println("write success")
}

// 2. 写操作会被阻塞
for i := 0; i < 300000; i++ {
	data := []byte("send some data from server" + "\n")
	wn, err := conn.Write(data)
	if err != nil {
		log.Fatalln(err)
	}
	log.Printf("%d, server write len is %d\n", i, wn)
}
// 客户端,仅连接,未读操作


// 3. 循环读
for {

    buf := make([]byte, 10)
    rn, err := conn.Read(buf)
    if err != nil {
        log.Println(err)
        break
    }
    log.Println("received from server data is:", string(buf[:rn]))
}
```

### 并发读写

并发读写，指的是两方面：

- 读操作和写操作是并发执行的
- 可能出现多个Goroutine同时写或读

因此在Go中，要使用Goroutine完成。同一个连接的并发读或写操作是**Goroutine并发安全**的。指的是同时存在多个Goroutine并发的读写，之间是不会相互影响的，这个在实操中，主要针对Write操作。conn.Write()是通过锁来实现的。

示例：

```go
// 并发的读和写操作，全双工
func TcpServerRWConcurrency() {
	// A. 基于某个地址建立监听
	// 服务端地址
	address := ":5678" // Any IP or version
	listener, err := net.Listen(tcp, address)
	if err != nil {
		log.Fatalln(err)
	}
	// 关闭监听
	defer listener.Close()
	log.Printf("%s server is listening on %s\n", tcp, listener.Addr())

	// B. 接受连接请求
	// 循环接受
	for {
		// 阻塞接受
		conn, err := listener.Accept()
		if err != nil {
			log.Println(err)
		}

		// 处理连接，读写
		go HandleConnConcurrency(conn)
	}
}

// 处理每个连接
func HandleConnConcurrency(conn net.Conn) {
	// 日志连接的远程地址（client addr）
	log.Printf("accept from %s\n", conn.RemoteAddr())
	// A.保证连接关闭
	defer conn.Close()

	wg := sync.WaitGroup{}
	// 并发的写
	wg.Add(1)
	go SerWrite(conn, &wg, "abcd")
	wg.Add(1)
	go SerWrite(conn, &wg, "efgh")
	wg.Add(1)
	go SerWrite(conn, &wg, "ijkl")

	// 并发的读
	wg.Add(1)
	go SerRead(conn, &wg)

	wg.Wait()
}

func SerWrite(conn net.Conn, wg *sync.WaitGroup, data string) {
	defer wg.Done()
	// B.向客户端发送数据，SerWrite
	for {
		wn, err := conn.Write([]byte(data + "\n"))
		if err != nil {
			log.Println(err)
		}
		log.Printf("server write len is %d\n", wn)
		time.Sleep(1 * time.Second)
	}
}

func SerRead(conn net.Conn, wg *sync.WaitGroup) {
	defer wg.Done()
	for {
		// C.从客户端接收数据，SerRead
		buf := make([]byte, 1024)
		rn, err := conn.Read(buf)
		if err != nil {
			log.Println(err)
		}
		log.Println("received from client data is:", string(buf[:rn]))
	}
}
```

注意，一次Write操作，表示一个原子的业务单元，不能再分。否则在Goroutine调度时不能保证连续性。

锁示例代码：

GOROOT/src/internal/poll/fd_windows.go

```go
// Write implements io.Writer.
func (fd *FD) Write(buf []byte) (int, error) {
	if err := fd.writeLock(); err != nil {
		return 0, err
	}
	defer fd.writeUnlock()
	if fd.isFile {
		fd.l.Lock()
		defer fd.l.Unlock()
	}
```

## 格式化消息

在发送或接收消息时，需要对消息进行格式化处理，才能在应用程序中保证消息具有逻辑含义**。前面的例子，我们采用的是字符串传递消息，**也是一种格式，但能够包含的数据字段有限。

典型编程时，我们会将两端处理好的数据，使用特定格式进行发送。典型的有两类：

- 文本编码，例如JSON，YAML，CSV等
- 二进制编码，例如GOB（Go Binary），Protocol Buffer等

格式化消息的典型流程，如图：

![image.png](https://fynotefile.oss-cn-zhangjiakou.aliyuncs.com/fynote/fyfile/13080/1682659809092/79abaa093b0c4e879727363b9946dc91.png)

**注意：发送端需要创建编码器，接收端需要创建解码器**

示例：

服务端代码

```go
// 格式化传输
func TcpServerFormat() {
	// A. 基于某个地址建立监听
	// 服务端地址
	address := ":5678" // Any IP or version
	listener, err := net.Listen(tcp, address)
	if err != nil {
		log.Fatalln(err)
	}
	// 关闭监听
	defer listener.Close()
	log.Printf("%s server is listening on %s\n", tcp, listener.Addr())

	// B. 接受连接请求
	// 循环接受
	for {
		// 阻塞接受
		conn, err := listener.Accept()
		if err != nil {
			log.Println(err)
		}

		// 处理连接，读写
		go HandleConnFormat(conn)
	}
}

func HandleConnFormat(conn net.Conn) {
	// 日志连接的远程地址（client addr）
	log.Printf("accept from %s\n", conn.RemoteAddr())
	// A.保证连接关闭
	defer conn.Close()

	wg := sync.WaitGroup{}
	wg.Add(1)
	// 发送端，
	go SerWriteFormat(conn, &wg)
	wg.Wait()
}

func SerWriteFormat(conn net.Conn, wg *sync.WaitGroup) {
	defer wg.Done()
	for {
		// 向客户端发送数据
		// 数据编码后发送

		// 创建需要传递的数据
		// 自定义的消息结构类型
		type Message struct {
			ID      uint   `json:"id,omitempty"`
			Code    string `json:"code,omitempty"`
			Content string `json:"content,omitempty"`
		}
		message := Message{
			ID:      uint(rand.Int()),
			Code:    "SERVER-STANDARD",
			Content: "message from server",
		}
		// 编码后数据的展示
		var buf bytes.Buffer
		encoderData := json.NewEncoder(&buf)
		//encoderData := gob.NewEncoder(&buf)
		if err := encoderData.Encode(message); err != nil {
			log.Println(err)
			continue
		}
		log.Println(buf.String())

		// 1, JSON, 文本编码
		//// 创建编码器
		//encoder := json.NewEncoder(conn)
		//// 利用编码器进行编码
		//// encode 成功后，会写入到conn，已经完成了conn.Write()
		//if err := encoder.Encode(message); err != nil {
		//	log.Println(err)
		//	continue
		//}
		//log.Println("message was send")

		// 2, GOB, 二进制编码
		// 创建编码器
		encoder := gob.NewEncoder(conn)
		// 利用编码器进行编码
		// encode 成功后，会写入到conn，已经完成了conn.Write()
		if err := encoder.Encode(message); err != nil {
			log.Println(err)
			continue
		}
		log.Println("message was send")

		time.Sleep(1 * time.Second)
	}
}

```

客户端（接收端）代码：

```go
func TcpClientFormat() {
	// tcp服务端地址
	serverAddress := "127.0.0.1:5678" // IPv6 4

	// A. 建立连接
	conn, err := net.DialTimeout(tcp, serverAddress, time.Second)
	//conn, err := net.Dial(tcp, serverAddress)
	if err != nil {
		log.Println(err)
		return
	}
	// 保证关闭
	defer conn.Close()
	log.Printf("connection is establish, client addr is %s\n", conn.LocalAddr())

	wg := sync.WaitGroup{}

	// 接收端
	wg.Add(1)
	go CliReadFormat(conn, &wg)

	wg.Wait()
}

func CliReadFormat(conn net.Conn, wg *sync.WaitGroup) {
	defer wg.Done()
	for {
		// 从客户端接收数据
		// 接收到数据后，先解码

		// 传递的消息类型
		type Message struct {
			ID      uint   `json:"id,omitempty"`
			Code    string `json:"code,omitempty"`
			Content string `json:"content,omitempty"`
		}
		message := Message{}

		// 1, JSON解码
		//// 创建解码器
		//decoder := json.NewDecoder(conn)
		//// 利用解码器进行解码
		//// 解码操作，从conn中读取内容，成功会将解码后的结果，赋值到message变量
		//if err := decoder.Decode(&message); err != nil {
		//	log.Println(err)
		//	continue
		//}
		//log.Println(message)

		// 2, GOB解码
		// 创建解码器
		decoder := gob.NewDecoder(conn)
		// 利用解码器进行解码
		// 解码操作，从conn中读取内容，成功会将解码后的结果，赋值到message变量
		if err := decoder.Decode(&message); err != nil {
			log.Println(err)
			continue
		}
		log.Println(message)

	}
}
```

测试：

```go
// 格式化消息的测试
func TestTcpServerFormat(t *testing.T) {
	TcpServerFormat()
}
func TestTcpClientFormat(t *testing.T) {
	TcpClientFormat()
}
```

结果：

客户端，解码成功，得到原始数据：

```
> go test -run TcpServerFormat
2023/05/04 12:51:19 tcp server is listening on [::]:5678
```

```
> go test -run TcpClientFormat
2023/05/04 12:52:09 connection is establish, client addr is 127.0.0.1:53253
2023/05/04 12:52:09 {3841400281839720065 SERVER-STANDARD message from server}
2023/05/04 12:52:10 {2803185154894110739 SERVER-STANDARD message from server}
2023/05/04 12:52:11 {4755708381034294201 SERVER-STANDARD message from server}
```

## 短连接和长连接

程序开发时，连接管理通常分为长短连接：

- 短链接，连接建立**传输数据后立即关闭**，下次需要传输数据再次建立连接。
- 长连接，**连接建立完毕后，利用连接多次发送数据**，在发送数据的过程中，保持连接不被关闭。最后才关闭连接

短连接的操作步骤：

1. 建立连接
2. 传输数据
3. 关闭连接

如图：

![image.png](https://fynotefile.oss-cn-zhangjiakou.aliyuncs.com/fynote/fyfile/13080/1682659809092/205fa40d019a4f63b32ae2cd9df50eaa.png)

长连接的操作步骤：

1. 建立连接
2. 传输数据（重复）
3. 心跳检测（重复）
4. 关闭连接

如图：

![image.png](https://fynotefile.oss-cn-zhangjiakou.aliyuncs.com/fynote/fyfile/13080/1682659809092/0f7c0ff39d924cac8e2efed8428fbce1.png)

### 短连接示例

短连接的编码很直观，只要**在发送完数据后，主动断开连接**。那么对应的接收端当读取不到内容时，表示接收完毕，随之断开连接即可。

**接收端（读），当读取到错误io.EOF时，我们认为连接已经结束关闭。**

示例代码：

服务端：

```go
// 短连接编程示例
func TcpServerSort() {
	// A. 基于某个地址建立监听
	// 服务端地址
	address := ":5678" // Any IP or version
	listener, err := net.Listen(tcp, address)
	if err != nil {
		log.Fatalln(err)
	}
	// 关闭监听
	defer listener.Close()
	log.Printf("%s server is listening on %s\n", tcp, listener.Addr())

	// B. 接受连接请求
	// 循环接受
	for {
		// 阻塞接受
		conn, err := listener.Accept()
		if err != nil {
			log.Println(err)
		}

		// 处理连接，读写
		go HandleConnSort(conn)
	}
}

func HandleConnSort(conn net.Conn) {
	// 日志连接的远程地址（client addr）
	log.Printf("accept from %s\n", conn.RemoteAddr())
	// A.保证连接关闭
	defer conn.Close()

	wg := sync.WaitGroup{}
	wg.Add(1)
	// 发送端，
	go SerWriteSort(conn, &wg)
	wg.Wait()
}

func SerWriteSort(conn net.Conn, wg *sync.WaitGroup) {
	defer wg.Done()

	// 创建需要传递的数据
	// 自定义的消息结构类型
	type Message struct {
		ID      uint   `json:"id,omitempty"`
		Code    string `json:"code,omitempty"`
		Content string `json:"content,omitempty"`
	}
	message := Message{
		ID:      uint(rand.Int()),
		Code:    "SERVER-STANDARD",
		Content: "message from server",
	}

	// GOB, 二进制编码
	// 创建编码器
	encoder := gob.NewEncoder(conn)
	// 利用编码器进行编码
	// encode 成功后，会写入到conn，已经完成了conn.Write()
	if err := encoder.Encode(message); err != nil {
		log.Println(err)
		return
	}
	log.Println("message was send")
	log.Println("link will be close")
	return
}
```

客户端，注意判断Read后的EOF错误：

```go
// 短连接示例
func TcpClientSort() {
	// tcp服务端地址
	serverAddress := "127.0.0.1:5678" // IPv6 4
	// A. 建立连接
	conn, err := net.DialTimeout(tcp, serverAddress, time.Second)
	if err != nil {
		log.Println(err)
		return
	}
	// 保证关闭
	defer conn.Close()
	log.Printf("connection is establish, client addr is %s\n", conn.LocalAddr())

	wg := sync.WaitGroup{}

	// 接收端
	wg.Add(1)
	go CliReadSort(conn, &wg)

	wg.Wait()
}

func CliReadSort(conn net.Conn, wg *sync.WaitGroup) {
	defer wg.Done()
	// 传递的消息类型
	type Message struct {
		ID      uint   `json:"id,omitempty"`
		Code    string `json:"code,omitempty"`
		Content string `json:"content,omitempty"`
	}
	message := Message{}
	for {
		// 从客户端接收数据
		// 接收到数据后，先解码
		// GOB解码
		// 创建解码器
		decoder := gob.NewDecoder(conn)
		// 利用解码器进行解码
		// 解码操作，从conn中读取内容，成功会将解码后的结果，赋值到message变量
		err := decoder.Decode(&message)
		// 错误 io.EOF 时，表示连接被给关闭
		if err != nil && errors.Is(err, io.EOF) {
			log.Println(err)
			log.Println("link was closed")
			break
		}
		log.Println(message)
	}
}
```

测试：

```shell
# 服务端（发送端）
> go test -run TcpServerSort
2023/05/04 18:34:46 tcp server is listening on [::]:5678
2023/05/04 18:34:51 accept from 127.0.0.1:62893
2023/05/04 18:34:51 message was send
2023/05/04 18:34:51 link will be close


# 客户端（接收端）
> go test -run TcpClientSort
2023/05/04 18:34:51 connection is establish, client addr is 127.0.0.1:62893
2023/05/04 18:34:51 {5307031956865372045 SERVER-STANDARD message from server}
2023/05/04 18:34:51 EOF
2023/05/04 18:34:51 link was closed
```

### 长连接的心跳检测

在使用长连接时，通常需要使用规律性的发送数据包，以维持在线状态，称为心跳检测。一旦心跳检测不能正确响应，那么就意味着对方（或者己方）不在线，关闭连接。心跳检测用来解决半连接问题。

测试：将连接建立后，关闭客户端或服务器，查看另一端的状态。

发送心跳检测的发送端：

- 可以是客户端
- 也可以是服务端
- 甚至是两端都发

典型的有**两种发送策略**：

1. **建立连接后，就使用固定的频率发送**
2. **一段时间没有接收到数据后，发送检测包**。（TCP 层的**KeepAlive就是该策略**）

心跳检测包的数据内容：

- 可以无数据
- 可以携带数据，例如做时钟同步，业务状态同步
- 典型的 ping pong 结构

心跳检测包是否需要响应？

- 可以不响应，发送成功即可
- 可以响应，通常用于同步数据

总而言之，都是业务来决定。

示例， pingpong模式，在连接建立后持续心跳：

* 定时心跳
* 判断是否接收到正确心跳响应
* 当N次心跳检测失败后，断开连接
* Server端，发送ping包
* Client端，接收到ping后，响应pong
* Server端，要检测是否收到了正确的响应pong，进而判断是否要主动断开连接

消息类型：

```go

// 传递的消息类型
type MessageHB struct {
	ID      uint      `json:"id,omitempty"`
	Code    string    `json:"code,omitempty"`
	Content string    `json:"content,omitempty"`
	Time    time.Time `json:"time,omitempty"`
}
```

服务端：

```go
func HandleConnHB(conn net.Conn) {
	// 日志连接的远程地址（client addr）
	log.Printf("accept from %s\n", conn.RemoteAddr())
	// A.保证连接关闭
	defer func() {
		conn.Close()
		log.Println("connection be closed")
	}()

	wg := sync.WaitGroup{}

	// 独立的goroutine，在连接建立后，周期发送ping
	wg.Add(1)
	// 发送ping
	go SerPing(conn, &wg)
	wg.Wait()
}

func SerPing(conn net.Conn, wg *sync.WaitGroup) {
	defer wg.Done()

	// 启动接收pong
	ctx, cancel := context.WithCancel(context.Background())
	go SerReadPong(conn, ctx)

	// ping失败的次数
	const maxPingNum = 3
	pingErrCounter := 0

	//周期性的发送
	//利用 time.Ticker
	ticker := time.NewTicker(2 * time.Second)
	for t := range ticker.C {
		pingMsg := MessageHB{
			ID:   uint(rand.Int()),
			Code: "PING-SERVER",
			Time: t,
		}

		// GOB, 二进制编码
		encoder := gob.NewEncoder(conn)
		// encode 成功后，会写入到conn，已经完成了conn.Write()
		if err := encoder.Encode(pingMsg); err != nil {
			log.Println(err)
			// 连接有问题的情况
			// 累加错误计数器
			pingErrCounter++
			// 判断是否到达上限
			if pingErrCounter == maxPingNum {
				// 心跳失败
				// 终止pong的处理
				cancel()
				return
			}
		}
		log.Printf("ping send to %s, ping id is %d\n", conn.RemoteAddr(), pingMsg.ID)
	}
}
func SerReadPong(conn net.Conn, ctx context.Context) {

	for {
		// 处理Ping结束
		select {
		case <-ctx.Done():
			return
		default:
			message := MessageHB{}
			// GOB解码
			decoder := gob.NewDecoder(conn)
			// 解码操作，从conn中读取内容，成功会将解码后的结果，赋值到message变量
			err := decoder.Decode(&message)
			// 错误 io.EOF 时，表示连接被给关闭
			if err != nil && errors.Is(err, io.EOF) {
				log.Println(err)
				break
			}
			// 判断是为为 pong 类型消息
			if message.Code == "PONG-CLIENT" {
				log.Printf("receive pong from %s, %s\n", conn.RemoteAddr(), message.Content)
			}
		}
	}
}
```

客户端：

```go
func CliReadPing(conn net.Conn, wg *sync.WaitGroup) {
	defer wg.Done()
	// 传递的消息类型
	message := MessageHB{}
	for {
		// GOB解码
		decoder := gob.NewDecoder(conn)
		// 解码操作，从conn中读取内容，成功会将解码后的结果，赋值到message变量
		err := decoder.Decode(&message)
		// 错误 io.EOF 时，表示连接被给关闭
		if err != nil && errors.Is(err, io.EOF) {
			log.Println(err)
			break
		}
		// 判断是为为 ping 类型消息
		if message.Code == "PING-SERVER" {
			log.Println("receive ping from", conn.RemoteAddr())
			CliWritePong(conn, message)
		}
	}
}

func CliWritePong(conn net.Conn, pingMsg MessageHB) {
	pongMsg := MessageHB{
		ID:      uint(rand.Int()),
		Code:    "PONG-CLIENT",
		Content: fmt.Sprintf("pingID:%d", pingMsg.ID),
		Time:    time.Now(),
	}

	// GOB, 二进制编码
	// 创建编码器
	encoder := gob.NewEncoder(conn)
	// 利用编码器进行编码
	// encode 成功后，会写入到conn，已经完成了conn.Write()
	if err := encoder.Encode(pongMsg); err != nil {
		log.Println(err)
		return
	}
	log.Println("pong was send to", conn.RemoteAddr())
	return
}
```

测试:

开启服务, 多开几个客户端, 关闭其中某些客户端.

服务器端检测时,会主动断开连接.

```shell
# server
2023/05/08 16:16:58 receive pong from 127.0.0.1:57726, pingID:1076147737332978911
2023/05/08 16:17:00 write tcp 127.0.0.1:5678->127.0.0.1:57726: wsasend: An existing connection was forcibly closed by the remote host.
2023/05/08 16:17:00 ping send to 127.0.0.1:57726, ping id is 7403348597707339775
2023/05/08 16:17:02 write tcp 127.0.0.1:5678->127.0.0.1:57726: wsasend: An existing connection was forcibly closed by the remote host.
2023/05/08 16:17:04 ping send to 127.0.0.1:57726, ping id is 8800556496508959890
2023/05/08 16:17:04 write tcp 127.0.0.1:5678->127.0.0.1:57726: wsasend: An existing connection was forcibly closed by the remote host.
2023/05/08 16:17:04 connection be closed
```

## 连接池（使用channel存储空闲连接）

### 核心结构

TCP连接的每次建立，都要进行三次握手，**为了避免频繁创建和销毁连接，因此复用连接，通常使用连接池技术**:

<img src="https://fynotefile.oss-cn-zhangjiakou.aliyuncs.com/fynote/fyfile/13080/1682659809092/e61cac61c498495c83760718a517474a.png" alt="image.png" style="zoom: 67%;" />

连接池基本操作

- 客户端(连接发起端), 通过连接池获取连接,**Get操作**
- 当暂时使用完毕后,将连接归还连接池, **Put操作**
- 其他客户端,需要连接同样去池中获取, Get操作, **只要连接没有被其他客户端占用,就可以重复使用**
- 少量长链接, 维护大量客户端的目的。否则，每个客户端，就需要1个连接。（？）

典型的,数据库连接池, 消息队列连接池等.

连接池的必要功能:

- New, 初始化连接池
- Get,获取连接
- Put, 放回连接

示例接口如下:

```go
type Pool interface {
	// 获取连接
	Get() (net.Conn, error)
	// 放回连接
	Put(conn net.Conn) error
	// 释放池(全部连接)
	Release() error
	// 有效连接个数
	Len() int
}
```

除此之外,**连接池还应该有能力创建新的连接**. 在Get操作时,若没有空闲可用的连接, 在数量允许的情况下,会创造新的连接. 该方法成为为连接工厂, 示例接口结构为:

```go
type ConnFactory interface {
	// 构造连接
	Factory(addr string) (net.Conn, error)
	// 关闭连接的方法
	Close(net.Conn) error
	// 检查连接是否有效的方法
	Ping(net.Conn) error
}
```

除了必要的操作, 连接池典型的配置有:

- **初始连接数**, 池初始化时的连接数
- **最大连接数**, 池中最多支持多少连接
- **最大空闲连接数**, 池中最多有多少可用的连接
- **空闲连接超时时间**, 多久后空闲连接会被释放

示例配置结构如下:

```go
type PoolConfig struct {
	//初始连接数, 池初始化时的连接数
	InitConnNum int
	//最大连接数, 池中最多支持多少连接
	MaxConnNum int
	//最大空闲连接数, 池中最多有多少可用的连接
	MaxIdleNum int
	//空闲连接超时时间, 多久后空闲连接会被释放
	IdleTimeout time.Duration

	// 连接地址
	addr string

	// 连接工厂
	Factory ConnFactory
}
```

由于**需要判断连接的空闲时间**（空闲连接在一段时间后会被释放），因此，需要记录连接被放入到连接池中的时间, 我们封装连接结构：

```go
// 空闲连接结构
type IdleConn struct {
	// 连接本身
	conn net.Conn
	// 放回时间
	putTime time.Time
}
```

有了基本操作和典型配置后, 连接池的结构设计如下:

- 要实现TcpPool接口
- 可以找到Factory
- 记录当前池信息，例如当前正在使用的连接数量，空闲的连接队列等

```go
type TcpPool struct {
	// 相关配置
	config *PoolConfig

	// 开放使用的连接数量
	openingConnNum int
	// 空闲的连接队列
	idleList chan *IdleConn

	// 并发安全锁
	mu sync.RWMutex
}
```

### 生产工厂的实现

工厂类型，实现ConnFactory接口，创建的对象用在配置中。

实现如下：

```go
// Tcp连接工厂类型
type TcpConnFactory struct{}

// 产生连接方法
func (*TcpConnFactory) Factory(addr string) (net.Conn, error) {
	// 校验参数的合理性
	if addr == "" {
		return nil, errors.New("addr is empty")
	}

	// 建立连接
	conn, err := net.DialTimeout("tcp", addr, 5*time.Second)
	if err != nil {
		return nil, err
	}

	// return
	return conn, nil
}

// 关闭连接
func (*TcpConnFactory) Close(conn net.Conn) error {
	return conn.Close()
}

// 检查连接是否有效
func (*TcpConnFactory) Ping(conn net.Conn) error {
	return nil
}
```

### 完善连接池基本结构

先依据Pool接口，将TcpPool的方法集实现完整。

```go
// TcpPool 实现 Pool 接口
func (*TcpPool) Get() (net.Conn, error) {
	return nil, nil
}

func (*TcpPool) Put(conn net.Conn) error {
	return nil
}
func (*TcpPool) Release() error {
	return nil
}
func (*TcpPool) Len() int {
	return 0
}
```

### 创建连接池函数

定义函数New，用于创建TcpPool，具体的功能有如下几步：

1. 校验参数
2. 初始化TcpPool
3. 初始化连接，关键步骤
4. 响应

示例代码：

```go
// 创建TcpPool对象
func NewTcpPool(addr string, poolConfig PoolConfig) (*TcpPool, error) {
	// 1校验参数
	if addr == "" {
		return nil, errors.New("addr is empty")
	}

	// 校验工厂的存在
	if poolConfig.Factory == nil {
		return nil, errors.New("factory is not exists")
	}

	// 最大连接数
	if poolConfig.MaxConnNum == 0 {
		//a,return错误
		//return nil, errors.New("max conn num is zero")
		//b,人为修改一个合理的
		poolConfig.MaxConnNum = defaultMaxConnNum
	}
	// 初始化连接数
	if poolConfig.InitConnNum == 0 {
		poolConfig.InitConnNum = defaultInitConnNum
	} else if poolConfig.InitConnNum > poolConfig.MaxConnNum {
		poolConfig.InitConnNum = poolConfig.MaxConnNum
	}
	// 合理化最大空闲连接数
	if poolConfig.MaxIdleNum == 0 {
		poolConfig.MaxIdleNum = poolConfig.InitConnNum
	} else if poolConfig.MaxIdleNum > poolConfig.MaxConnNum {
		poolConfig.MaxIdleNum = poolConfig.MaxConnNum
	}

	// 2初始化TcpPool对象
	pool := TcpPool{
		config:         poolConfig,
		openingConnNum: 0,
		idleList:       make(chan *IdleConn, poolConfig.MaxIdleNum),
		addr:           addr,
		mu:             sync.RWMutex{},
	}

	// 3初始化连接
	// 根据InitConnNum的配置来创建
	for i := 0; i < poolConfig.InitConnNum; i++ {
		conn, err := pool.config.Factory.Factory(addr)
		if err != nil {
			// 通常意味着，连接池初始化失败
			// 释放可能已经存在的连接
			pool.Release()
			return nil, err
		}
		// 连接创建成功
		// 加入到空闲连接队列中
		pool.idleList <- &IdleConn{
			conn:    conn,
			putTime: time.Now(),
		}
	}

	// 4返回
	return &pool, nil
}
```

**测试创建结果**

测试代码：

服务端：

```go
// 测试连接池服务端
func TcpServerPool() {
	// A. 基于某个地址建立监听
	// 服务端地址
	address := ":5678" // Any IP or version
	listener, err := net.Listen(tcp, address)
	if err != nil {
		log.Fatalln(err)
	}
	// 关闭监听
	defer listener.Close()
	log.Printf("%s server is listening on %s\n", tcp, listener.Addr())

	// B. 接受连接请求
	// 循环接受
	for {
		// 阻塞接受
		conn, err := listener.Accept()
		if err != nil {
			log.Println(err)
		}

		// 处理连接，读写
		go HandleConnPool(conn)
	}
}

func HandleConnPool(conn net.Conn) {
	// 日志连接的远程地址（client addr）
	log.Printf("accept from %s\n", conn.RemoteAddr())
	// A.保证连接关闭
	defer func() {
		conn.Close()
		log.Println("connection be closed")
	}()

	select {}
}
```

客户端：

```go
// 连接池使用
func TcpClientPool() {
	// tcp服务端地址
	serverAddress := "127.0.0.1:5678" // IPv6 4
	// A，建立连接池
	pool, err := NewTcpPool(serverAddress, PoolConfig{
		Factory:     &TcpConnFactory{},
		InitConnNum: 10,
		MaxIdleNum:  20,
	})
	if err != nil {
		log.Fatalln(err)
	}
	log.Println(pool, len(pool.idleList))
	// B, 复用连接池中的连接
}

```

测试结果，基于配置的初始化连接数量，会创建对应数量的连接：

服务端：

```shell
> go test -run TcpServerPool
2023/05/09 16:35:16 tcp server is listening on [::]:5678
2023/05/09 16:35:31 accept from 127.0.0.1:65458
2023/05/09 16:36:13 accept from 127.0.0.1:65464
2023/05/09 16:36:13 accept from 127.0.0.1:65465
2023/05/09 16:36:13 accept from 127.0.0.1:65466
2023/05/09 16:36:13 accept from 127.0.0.1:65467
2023/05/09 16:36:13 accept from 127.0.0.1:65468
2023/05/09 16:36:13 accept from 127.0.0.1:65469
2023/05/09 16:36:13 accept from 127.0.0.1:65470
2023/05/09 16:36:13 accept from 127.0.0.1:65471
2023/05/09 16:36:13 accept from 127.0.0.1:65472
```

客户端：

```shell
> go test -run TcpClientPool
2023/05/09 16:36:13 &{{10 100 20 0 0x120d0c0} 0 0xc000086120 {{0 0} 0 0 {{} 0} {{} 0}}} 10
```

### 从连接池中获取连接

编码实现 TcpPool.Get 方法， 其核心步骤为：

1. 并发安全锁
2. 利用for+select结构从chan *IdleConn中获取空闲连接
3. 判断连接的超时状态
4. 若不存在空闲连接，则利用工厂新建连接
5. 记录使用的连接数量
6. 最大连接数上限错误处理

示例代码：

```go
// TcpPool 实现 Pool 接口
func (pool *TcpPool) Get() (net.Conn, error) {
	// 1锁定
	pool.mu.Lock()
	defer pool.mu.Unlock()

	// 2获取空闲连接，若没有则创建连接
	for {
		select {
		// 获取空闲连接
		case idleConn, ok := <-pool.idleList:
			// 判断channel是否被关闭
			if !ok {
				return nil, errors.New("idle list closed")
			}

			// 判断连接是否超时
			//pool.config.IdleTimeout, idleConn.putTime
			if pool.config.IdleTimeout > 0 { // 设置了超时时间
				// putTime + timeout 是否在 now 之前
				if idleConn.putTime.Add(pool.config.IdleTimeout).Before(time.Now()) {
					// 关闭连接，继续查找下一个连接
					_ = pool.config.Factory.Close(idleConn.conn)
					continue
				}
			}

			// 判断连接是否可用
			if err := pool.config.Factory.Ping(idleConn.conn); err != nil {
				// ping 失败，连接不可用
				// 关闭连接，继续查找
				_ = pool.config.Factory.Close(idleConn.conn)
				continue
			}

			// 找到了可用的空闲连接
			log.Println("get conn from Idle")
			// 使用的连接计数
			pool.openingConnNum++
			// 返回连接
			return idleConn.conn, nil

		// 创建连接
		default:
			// a判断是否还可以继续创建
			// 基于开放的连接是否已经达到了连接池最大的连接数
			if pool.openingConnNum >= pool.config.MaxConnNum {
				return nil, errors.New("max opening connection")
				// 另一种方案，就是阻塞
				//continue
			}

			// b创建连接
			conn, err := pool.config.Factory.Factory(pool.addr)
			if err != nil {
				return nil, err
			}

			// c正确创建了可用的连接
			log.Println("get conn from Factory")
			// 使用的连接计数
			pool.openingConnNum++
			// 返回连接
			return conn, nil
		}
	}
}
```

**测试**

更新客户端测试代码：

```go
func TcpClientPool() {
	// tcp服务端地址
	serverAddress := "127.0.0.1:5678" // IPv6 4
	// A，建立连接池
	pool, err := NewTcpPool(serverAddress, PoolConfig{
		Factory:     &TcpConnFactory{},
		InitConnNum: 4,
	})
	if err != nil {
		log.Fatalln(err)
	}
	log.Println(pool, len(pool.idleList))

	wg := sync.WaitGroup{}
	clientNum := 18
	wg.Add(clientNum)
	// B, 复用连接池中的连接
	for i := 0; i < clientNum; i++ {
		// goroutine 模拟独立的客户端
		go func(wg *sync.WaitGroup) {
			defer wg.Done()
			log.Println(pool.Get())
		}(&wg)
	}
	wg.Wait()
}

```

`for select` 如果没有 `break` 或其他退出机制（如 `return`），`for select` 会永远循环下去。

测试设置：defaultMaxConnNum = 10

测试结果：

服务端：一次测试，最多建立10个连接

```
> go test -run TcpServerPool
2023/05/09 18:14:37 tcp server is listening on [::]:5678
2023/05/09 18:14:44 accept from 127.0.0.1:50405
2023/05/09 18:14:44 accept from 127.0.0.1:50406
2023/05/09 18:14:44 accept from 127.0.0.1:50407
2023/05/09 18:14:44 accept from 127.0.0.1:50408
2023/05/09 18:14:44 accept from 127.0.0.1:50409
2023/05/09 18:14:44 accept from 127.0.0.1:50410
2023/05/09 18:14:44 accept from 127.0.0.1:50411
2023/05/09 18:14:44 accept from 127.0.0.1:50412
2023/05/09 18:14:44 accept from 127.0.0.1:50413
2023/05/09 18:14:44 accept from 127.0.0.1:50414
```

客户端：4个连接来自Idle， 6个连接来自工厂创建。剩下的获取连接失败：

```
> go test -run TcpClientPool
2023/05/09 18:14:44 &{{4 10 4 0 0xa7e0c0} 0 0xc00005c0c0 127.0.0.1:5678 {{0 0} 0 0 {{} 0} {{} 0}}} 4
2023/05/09 18:14:44 get conn from Idle
2023/05/09 18:14:44 &{{0xc000108f00}} <nil>
2023/05/09 18:14:44 get conn from Idle
2023/05/09 18:14:44 &{{0xc000109180}} <nil>
2023/05/09 18:14:44 get conn from Idle
2023/05/09 18:14:44 &{{0xc000109400}} <nil>
2023/05/09 18:14:44 get conn from Idle
2023/05/09 18:14:44 &{{0xc000109680}} <nil>
2023/05/09 18:14:44 get conn from Factory
2023/05/09 18:14:44 &{{0xc000212000}} <nil>
2023/05/09 18:14:44 get conn from Factory
2023/05/09 18:14:44 &{{0xc00019e000}} <nil>
2023/05/09 18:14:44 get conn from Factory
2023/05/09 18:14:44 &{{0xc000109900}} <nil>
2023/05/09 18:14:44 get conn from Factory
2023/05/09 18:14:44 &{{0xc00019e280}} <nil>
2023/05/09 18:14:44 get conn from Factory
2023/05/09 18:14:44 &{{0xc000109b80}} <nil>
2023/05/09 18:14:44 get conn from Factory
2023/05/09 18:14:44 &{{0xc00019e500}} <nil>
2023/05/09 18:14:44 <nil> max opening connection
2023/05/09 18:14:44 <nil> max opening connection
2023/05/09 18:14:44 <nil> max opening connection
2023/05/09 18:14:44 <nil> max opening connection
2023/05/09 18:14:44 <nil> max opening connection
2023/05/09 18:14:44 <nil> max opening connection
2023/05/09 18:14:44 <nil> max opening connection
2023/05/09 18:14:44 <nil> max opening connection
```

### 将连接放回池中

编码实现 TcpPool.Put 方法， 其核心步骤为：

1. 并发安全锁
2. 利用select结构向chan *IdleConn中发送连接
3. 若空闲连接列表已满，则关闭连接
4. 更新开放的连接数量
5. 做一些校验
   1. channel是否可用
   2. 连接是否可用等

示例代码：

```go
func (pool *TcpPool) Put(conn net.Conn) error {
	// 1锁
	pool.mu.Lock()
	defer pool.mu.Unlock()

	// 2做一些校验
	if conn == nil {
		return errors.New("connection is not exists")
	}
	// 判断空闲连接列表是否存在
	if pool.idleList == nil {
		// 关闭连接
		_ = pool.config.Factory.Close(conn)
		return errors.New("idle list is not exists")
	}

	// 3放回连接
	select {
	// 放回连接
	case pool.idleList <- &IdleConn{
		conn:    conn,
		putTime: time.Now(),
	}:
		// 只要可以发送成功，任务完成
		// 更新开放的连接数量
		pool.openingConnNum--
		return nil
	// 关闭连接
	default:
		_ = pool.config.Factory.Close(conn)
		return nil
	}
}
```

**测试**

在客户端，Get后Put，测试是否支持连接复用：`clientNum := 50`

```go
// 连接池使用
func TcpClientPool() {
	// tcp服务端地址
	serverAddress := "127.0.0.1:5678" // IPv6 4
	// A，建立连接池
	pool, err := NewTcpPool(serverAddress, PoolConfig{
		Factory:     &TcpConnFactory{},
		InitConnNum: 4,
	})
	if err != nil {
		log.Fatalln(err)
	}
	log.Println(pool, len(pool.idleList))

	wg := sync.WaitGroup{}
	clientNum := 50
	wg.Add(clientNum)
	// B, 复用连接池中的连接
	for i := 0; i < clientNum; i++ {
		// goroutine 模拟独立的客户端
		go func(wg *sync.WaitGroup) { // 有几个客户端就创建几个goroutine
			defer wg.Done()
			// 获取连接
			conn, err := pool.Get()
			if err != nil {
				log.Println(err)
				return
			}
			//log.Println(conn)
			// 回收连接
			pool.Put(conn)
		}(&wg)
	}
	wg.Wait()
}
```

以上例子，可以看到大部分连接来自于Idle。

结果：

```shell
> go test -run TcpClientPool
2023/05/09 18:56:53 &{{4 10 4 0 0x59e0c0} 0 0xc00005c0c0 127.0.0.1:5678 {{0 0} 0 0 {{} 0} {{} 0}}} 4
2023/05/09 18:56:53 get conn from Idle
2023/05/09 18:56:53 get conn from Idle
2023/05/09 18:56:53 get conn from Idle
2023/05/09 18:56:53 get conn from Idle
2023/05/09 18:56:53 get conn from Idle
2023/05/09 18:56:53 get conn from Idle
2023/05/09 18:56:53 get conn from Idle
2023/05/09 18:56:53 get conn from Idle
2023/05/09 18:56:53 get conn from Idle
2023/05/09 18:56:53 get conn from Idle
2023/05/09 18:56:53 get conn from Idle
2023/05/09 18:56:53 get conn from Idle
2023/05/09 18:56:53 get conn from Idle
2023/05/09 18:56:53 get conn from Idle
2023/05/09 18:56:53 get conn from Idle
2023/05/09 18:56:53 get conn from Idle
2023/05/09 18:56:53 get conn from Idle
2023/05/09 18:56:53 get conn from Idle
2023/05/09 18:56:53 get conn from Idle
2023/05/09 18:56:53 get conn from Idle
2023/05/09 18:56:53 get conn from Idle
2023/05/09 18:56:53 get conn from Idle
2023/05/09 18:56:53 get conn from Idle
2023/05/09 18:56:53 get conn from Idle
2023/05/09 18:56:53 get conn from Idle
2023/05/09 18:56:53 get conn from Idle
2023/05/09 18:56:53 get conn from Idle
2023/05/09 18:56:53 get conn from Idle
2023/05/09 18:56:53 get conn from Idle
2023/05/09 18:56:53 get conn from Idle
2023/05/09 18:56:53 get conn from Idle
2023/05/09 18:56:53 get conn from Idle
2023/05/09 18:56:53 get conn from Idle
2023/05/09 18:56:53 get conn from Idle
2023/05/09 18:56:53 get conn from Idle
2023/05/09 18:56:53 get conn from Idle
2023/05/09 18:56:53 get conn from Idle
2023/05/09 18:56:53 get conn from Idle
2023/05/09 18:56:53 get conn from Factory
2023/05/09 18:56:53 get conn from Factory
2023/05/09 18:56:53 get conn from Factory
2023/05/09 18:56:53 get conn from Factory
2023/05/09 18:56:53 get conn from Factory
2023/05/09 18:56:53 get conn from Factory
2023/05/09 18:56:53 max opening connection
2023/05/09 18:56:53 max opening connection
2023/05/09 18:56:53 max opening connection
2023/05/09 18:56:53 max opening connection
2023/05/09 18:56:53 max opening connection
2023/05/09 18:56:53 max opening connection
```

并发数量继续增加，还会继续看到获取连接成功的输出的。

### 释放连接池

当连接池不再使用时，需要将池中的全部的连接关闭，该操作称为释放连接池操作。

核心步骤：

* 关闭 Idle List
* 将 Idle List 中的连接全部关闭
* 配合Put操作,关闭全部连接

示例代码:

```go
// 释放连接池
func (pool *TcpPool) Release() error {
	// 1并发安全锁
	pool.mu.Lock()
	defer pool.mu.Unlock()

	// 2确定连接池是否被释放
	if pool.idleList == nil {
		return nil
	}

	// 3关闭IdleList
	close(pool.idleList)

	// 4释放全部空闲连接
	// 继续接收已关闭channel中的元素
	for idleConn := range pool.idleList {
		// 关闭连接
		_ = pool.config.Factory.Close(idleConn.conn)
	}

	return nil
}
```

**测试**

客户端连接池使用后，释放连接池：

```go
// 连接池使用
func TcpClientPool() {
	// 其他代码与之前的测试一致，略

	wg.Wait()

	// 释放连接池
	pool.Release()
	log.Println(pool, pool.idleList)
}
```

### 获取连接池长度

也就是获取 pool.idleList的长度。

示例代码：

```go
// 获取连接池长度
// 当前的可用连接数
func (pool *TcpPool) Len() int {
	return len(pool.idleList)
}
```

自行测试即可！

```
log.Println(pool, pool.Len())
```

### 总结

* 连接池作用：**复用连接**
* 设计池与生产隔离
  * 池的管理
  * 生产连接管理
  * 适用于任何资源的池
* 编码
  * channel
  * select
  * sync.Mutex, sync.RWMutex
  * Interface
* 使用连接池
  * 使用多goroutine并发模拟使用

扩展：将TCP连接池，扩展为支持任何类型的资源。

## TCP黏包

### 粘包现象

TCP粘包是指在基于TCP协议进行数据传输时，发送方的多条数据可能会粘在一起，**接收方在读取数据时可能会把多条消息当成一条消息读取的现象**。与之相关的还有“分包”现象，即一条完整的消息被拆分成了多次读取。

从接收缓冲区看，后一包数据的头紧接着前一包数据的尾。

如图：

<img src="https://fynotefile.oss-cn-zhangjiakou.aliyuncs.com/fynote/fyfile/13080/1682659809092/d4292178976a45338736f9a0e8fc81b7.png" alt="image.png" style="zoom: 50%;" />

其实TCP是面向字节流的协议，就是没有界限的一串数据，本没有“包”的概念， 包可以当作一个逻辑上的数据单元。**“粘包”和“拆包”是逻辑上的概念。**

粘包示例：

服务端，连续发送多个数据包：

```go
func HandleConnSticky(conn net.Conn) {
	// 日志连接的远程地址（client addr）
	log.Printf("accept from %s\n", conn.RemoteAddr())
	// A.保证连接关闭
	defer func() {
		conn.Close()
		log.Println("connection be closed")
	}()

	// 连续发送数据
	data := "package data."
	for i := 0; i < 50; i++ {
		_, err := conn.Write([]byte(data))
		if err != nil {
			log.Println(err)
		}
	}
}
```

客户端，接收数据：

```go
func HandleConnSticky(conn net.Conn) {
	// 日志连接的远程地址（client addr）
	log.Printf("accept from %s\n", conn.RemoteAddr())
	// A.保证连接关闭
	defer func() {
		conn.Close()
		log.Println("connection be closed")
	}()

	// 连续发送数据
	data := "package data."
	for i := 0; i < 50; i++ {
		_, err := conn.Write([]byte(data))
		if err != nil {
			log.Println(err)
		}
	}
}
```

输出结果：

```shell
# 客户端（接收端）
> go test -run TcpClientSticky
2023/05/10 17:17:05 connection is establish, client addr is 127.0.0.1:51627
2023/05/10 17:17:05 received data: package data.package data.package data.package data.package data.package data.package data.package data.package data.package data.package data.package data.package data.package data.package data.package data.package data.package data.package data.package data.package data.package data.package data.package data.package data.package data.package data.package data.package data.package data.package data.package data.package data.package data.
2023/05/10 17:17:05 received data: package data.package data.package data.package data.package data.package data.package data.package data.package data.package data.
2023/05/10 17:17:05 received data: package data.package data.package data.
2023/05/10 17:17:05 received data: package data.package data.
2023/05/10 17:17:05 received data: package data.
2023/05/10 17:17:05 EOF
```

从结果上看，读取到的数据连在一起了，称为粘包。

### 粘包/分包原因

1. **TCP特性**：TCP是**面向字节流**的协议，没有消息边界的概念，数据以流的形式发送。因此**发送方可能会将多次`send`的数据进行合并**，接收方可能一次性接收这些数据。

2. **Nagle算法（发送端）**：为了提高网络传输效率，TCP使用`Nagle`算法来减少传输的报文数量，**Nagle算法会将小数据包进行合并后再发送**，可能导致粘包。

3. **接收方读取逻辑**：接收方并不一定能完全按照发送方的消息边界来读取数据，可能会**一次性读取多条消息**，或者一条消息读取不完整。

当发送的多个数据包之间需要逻辑隔离，那么就需要处理粘包问题。反之若发送的数据本身就是一个连续的整体，那么不需要处理粘包问题。

### 解决方案

数据包连着发送这个是不能改变的。我们需要**在数据包层面标注包与包的分离方案**，来解决粘包现象带来的问题。

典型的方案有：

- 每个包都**封装成固定的长度**。读取到内容时，依据长度进行分割即可
- 数据包使用**特定分隔符**。读取到内容时，依据分隔符分割数据即可，例如HTTP,FTP协议的\r\n。
- **将消息封装为Header+Body形式**，Header通常时固定长度，Header中包含Body的长度信息。读取到期待长度时，才表示成功。

不论何种方案，在编码实现时，**通常采用定义编码器（接收端用解码器）的方案来实现**。就类似JSON和GOB编码。

示例编码，以方案三，Header+Body的模式为例：

Header的长度为4个字节，用于表示Body的长度。

```go
// 逻辑数据
13package data.

// 传输数据
[4]byte(int32(13))[]byte("package data.")
binary(int32(13))[]byte("package data.")

```

先定义编码解码器：

```go
// 定义编码器（发送端）
type Encoder struct {
	// 编码结束后，写入的目标
	w io.Writer
}

// 创建编码器函数
func NewEncoder(w io.Writer) *Encoder {
	return &Encoder{
		w: w,
	}
}

// 编码，将编码的结果，写入到w io.Writer
// binary(int32(13))[]byte("package data.")
func (enc *Encoder) Encode(message string) error {
	// 1.获取message的长度
	l := int32(len(message))

	// 2.构建一个数据包缓存
	buf := new(bytes.Buffer)

	// 3.在buf中写入长度，将数据以小端格式（Little Endian）写入缓冲区（buf）。
	if err := binary.Write(buf, binary.LittleEndian, l); err != nil {
		return err
	}

	// 4.将数据主体Body写入
	//if err := binary.Write(buf, binary.LittleEndian, []byte(message)); err != nil {
	//	return err
	//}
	if _, err := buf.Write([]byte(message)); err != nil {
		return err
	}

	// 4.利用io.Writer发送数据
	if n, err := enc.w.Write(buf.Bytes()); err != nil {
		log.Println(n, err)
		return err
	}

	return nil
}

// 定义解码器（接收端）
// 解码器
type Decoder struct {
	// Reader
	r io.Reader
}

// 创建Decoder
func NewDecoder(r io.Reader) *Decoder {
	return &Decoder{
		r: r,
	}
}

// 从Reader中读取内容，解码。
// binary(int32(13))[]byte("package data.")
func (dec *Decoder) Decode(message *string) error {
	// 1.读取header（读取前4个字节）
	header := make([]byte, 4)
	hn, err := dec.r.Read(header)
	if err != nil {
		return err
	}
	if hn != 4 {
		return errors.New("header is not enough")
	}

	// 2.将前4个字节转换为int32类型，确定了body的长度。
	var l int32
    // 将 header包装为一个 *bytes.Buffer 对象
	headerBuf := bytes.NewBuffer(header)
    // 将数据从 headerBuf 读取并存储到 l 中（也是使用小端格式读取）。
	if err := binary.Read(headerBuf, binary.LittleEndian, &l); err != nil {
		return err
	}

	// 3.读取body
	body := make([]byte, l)
	bn, err := dec.r.Read(body)
	if err != nil {
		return err
	}
	if bn != int(l) {
		return errors.New("body is not enough")
	}

	// 4.设置message
	*message = string(body)
    
	return nil
}

```

在利用编解码器，完成读写操作：

发送端，写，server：

```go
func HandleConnCoder(conn net.Conn) {
	// 日志连接的远程地址（client addr）
	log.Printf("accept from %s\n", conn.RemoteAddr())
	// A.保证连接关闭
	defer func() {
		conn.Close()
		log.Println("connection be closed")
	}()

	// 连续发送数据
	data := []string{
		"package data.",
		"package.",
		"package data data",
		"pack",
	}
    
	encoder := NewEncoder(conn)
	for i := 0; i < 50; i++ {
		// 创建编解码器
		// 利用编码器进行编码
		// encode 成功后，会写入到conn，已经完成了conn.Write()
		if err := encoder.Encode(data[rand.Intn(len(data))]); err != nil {
			log.Println(err)
		}
	}
}
```

接收端，读，client：

```go
// 粘包,编解码器解决
func TcpClientCoder() {
	// tcp服务端地址
	serverAddress := "127.0.0.1:5678" // IPv6 4

	// A. 建立连接
	conn, err := net.DialTimeout(tcp, serverAddress, time.Second)
	//conn, err := net.Dial(tcp, serverAddress)
	if err != nil {
		log.Println(err)
		return
	}
	// 保证关闭
	defer conn.Close()
	log.Printf("connection is establish, client addr is %s\n", conn.LocalAddr())

	// 从服务端接收数据，SerRead
	// 创建解码器
	decoder := NewDecoder(conn)
	data := ""
	i := 0
	for {
		// 将读取的数据存储到 data 里面，错误 io.EOF 时，表示连接被给关闭。
		if err := decoder.Decode(&data); err != nil {
			log.Println(err)
			break
		}

		log.Println(i, "received data:", data)
		i++
	}
}
```

测试：

```go
> go test -run TcpClientCoder
2023/05/10 20:06:03 connection is establish, client addr is 127.0.0.1:53269
2023/05/10 20:06:03 0 received data: pack
2023/05/10 20:06:03 1 received data: package.
2023/05/10 20:06:03 2 received data: package data data
2023/05/10 20:06:03 3 received data: pack
2023/05/10 20:06:03 4 received data: package data data
2023/05/10 20:06:03 5 received data: pack
2023/05/10 20:06:03 6 received data: pack
2023/05/10 20:06:03 7 received data: package.
```

## TCP专用方法

除了通用的Listen，Accept，Dial外，net包还提供了专门的TCP方法：

```go
// Listen
func Listen(network, address string) (Listener, error)
func ListenTCP(network string, laddr *TCPAddr) (*TCPListener, error)

// Accept
func (Listener) Accept() (Conn, error)
func (l *TCPListener) AcceptTCP() (*TCPConn, error)

// Dial
func Dial(network, address string) (Conn, error)
func DialTCP(network string, laddr, raddr *TCPAddr) (*TCPConn, error)
```

其中，TCP特定的类型：

```go
*TCPAddr
*TCPListene
*TCPConn
```

示例代码：

```go
// 服务端
// TCP特定方法
func TcpServerSpecial() {
	// 1建立监听
	// 获取本地地址（监听地址）
	laddr, err := net.ResolveTCPAddr("tcp", ":5678")
	if err != nil {
		log.Fatalln(err)
	}
	tcpListener, err := net.ListenTCP("tcp", laddr)
	if err != nil {
		log.Fatalln(err)
	}
	defer tcpListener.Close()
	log.Printf("%s server is listening on %s\n", tcp, tcpListener.Addr())

	// 2接收连接
	for {
		tcpConn, err := tcpListener.AcceptTCP()
		if err != nil {
			log.Println(err)
			continue
		}

		// 3处理每个连接
		go handleConnSpecial(tcpConn)
	}
}
func handleConnSpecial(tcpConn *net.TCPConn) {
	log.Printf("accept from %s\n", tcpConn.RemoteAddr())

	// 设置连接属性
	tcpConn.SetKeepAlive(true)

	// 写数据
	data := "tcp message."
	n, err := tcpConn.Write([]byte(data))
	if err != nil {
		log.Println(err)
		return
	}
	log.Println("Send len:", n)
}


// 客户端
func TcpClientSpecial() {
	// 1建立连接
	// raddr remote addr，服务端的地址
	raddr, err := net.ResolveTCPAddr("tcp", "127.0.0.1:5678")
	if err != nil {
		log.Fatalln(err)
	}
	// laddr, local addr, 客户端的地址，可以用于设置客户端的端口
	tcpConn, err := net.DialTCP("tcp", nil, raddr)
	if err != nil {
		log.Fatalln(err)
	}
	// 保证关闭
	defer tcpConn.Close()
	log.Printf("connection is establish, client addr is %s\n", tcpConn.LocalAddr())

	// 2读写数据
	buf := make([]byte, 1024)
	for {
		n, err := tcpConn.Read(buf)
		if err != nil {
			log.Println(err)
			return
		}
		log.Println("receive len:", n)
		log.Println("receive data:", string(buf[:n]))
	}
}
```

注意，几个当建立连接的相关方法即可。

建立连接后，传输数据的操作是通用的。

使用TCPConn的目的，是需要对TCP连接的特定属性进行配置，例如：

```go
// 设置连接属性
tcpConn.SetKeepAlive(true)


// SetKeepAlive sets whether the operating system should send
// keep-alive messages on the connection.
func (c *TCPConn) SetKeepAlive(keepalive bool) error
```

## TCP连接属性设置

`*net.TCPConn`提供如下几个设置连接熟悉的方法：

```go
// 设置读写操作的Deadline（截至时间）
func (c *conn) SetDeadline(t time.Time) error
func (c *conn) SetReadDeadline(t time.Time) error
func (c *conn) SetWriteDeadline(t time.Time) error

// 设置读缓冲尺寸
func (c *conn) SetReadBuffer(bytes int) error
// 设置写缓存尺寸
func (c *conn) SetWriteBuffer(bytes int) error

// 设置TCP连接关闭后的逗留时间
func (c *TCPConn) SetLinger(sec int) error
// 设置是否开启KeepAlive，在一定时间段内（7200s，取决于OS），连接上没有数据传输，会发送测试数据以用来探测对方的在线状态
func (c *TCPConn) SetKeepAlive(keepalive bool) error
// 设置KeepAlive的空闲时间
func (c *TCPConn) SetKeepAlivePeriod(d time.Duration) error
// 设置是否不延迟。默认false，表示有延迟，其实就是使用Nagle算法。true为无延迟。
func (c *TCPConn) SetNoDelay(noDelay bool) error
```

缓冲示例图：

![image.png](https://fynotefile.oss-cn-zhangjiakou.aliyuncs.com/fynote/fyfile/13080/1682659809092/60ae35ede268423db7739c095d6769dc.png)

延迟和不延迟：

![image.png](https://fynotefile.oss-cn-zhangjiakou.aliyuncs.com/fynote/fyfile/13080/1682659809092/78e05b638785455c89c8f960ddda917d.png)

**`conn.(*net.TCPConn)`**

可以将Conn接口断言为*net.TCPConn类型，使用其特定的方法。

```go
// 断言为TCPConn即可
	tcpConn, ok := conn.(*net.TCPConn)
	if !ok {
		log.Println("non tcp connection")
	}
	tcpConn.SetNoDelay(true)
```

## 文件传输案例

### 案例说明

- 客户端：发送文件
- 服务端：接收文件

技术支持：

- tcp网络编程
- 文件操作

用到的文件函数：

```go
package os
// 打开文件，用于读取
func Open(name string) (*File, error)
// 关闭文件
func (f *File) Close() error
// 创建文件
func Create(name string) (*File, error)
// 写入文件
func (f *File) Write(b []byte) (n int, err error)
```

### 编码实现

客户端：

```go

```

服务端：

```go

```

测试：
