# 网络轮询器 NetPoller

**网络轮询器是 Go 语言运行时用来处理 I/O 操作的关键组件**，它使用了操作系统提供的 **I/O 多路复用机制**增强程序的并发处理能力（通常用于处理大规模并发网络连接）。网络轮询器不仅用于监控网络 I/O，还能用于监控文件的 I/O，它利用了操作系统提供的 I/O 多路复用模型来提升 I/O 设备的利用率以及程序的性能。

<img src="https://fynotefile.oss-cn-zhangjiakou.aliyuncs.com/fynote/fyfile/13080/1684143948010/8a878fcb07654905bf92854694a36b97.png" alt="image.png" style="zoom:80%;" />

## I/O模型

操作系统中包含：

1. **阻塞I/O模型**，Blocking I/O Model

   * I/O 操作（如读写文件或网络数据）会让**调用进程**阻塞，直到操作完成。
   * **在阻塞期间，进程无法执行其他任务。**

2. **非阻塞I/O模型**，Non-Blocking I/O Model

   * I/O 操作立即返回，**不会阻塞进程**。

   - **如果数据未准备好，操作返回错误**，程序需要主动轮询再次尝试。
3. **信号驱动I/O模型**，Signal-Driven I/O Model，**是非阻塞的**

   - 当**文件描述符的状态发生变化**（例如可读、可写）时，**内核向应用程序发送信号，通知程序处理相应的 I/O 操作。**

   - **进程收到信号后执行对应的处理逻辑。**这说明文件描述符的状态变化是由内核监控的。
4. **异步I/O模型** ，Asynchronous I/O Model

   - **由内核完全负责 I/O 操作**

   - 应用程序发起 I/O 请求后立即返回（**剩余工作交给内核处理**）。

   - 内核完成 I/O 操作并通知应用程序（通过信号或回调），**程序可以直接处理数据或获取操作结果。**

5. **I/O多路复用模型**，Multiplexing I/O Model

   * 使用系统调用（如 `select`、`poll`、`epoll`）**监视多个文件描述符的状态**。（操作系统完成）

   - 当一个或多个文件描述符发生变化时，**内核返回这些就绪事件**，通知程序进行对应的 I/O 操作。
   - **内核负责： 监控文件描述符是否就绪**。这是 I/O 多路复用的核心功能，由内核高效地完成。
   - **程序负责：** 在接收到文件描述符就绪的通知后，**执行实际的 I/O 操作。**

五种I/O模型。

**非阻塞I/O模型和异步I/O模型的区别：**

- **非阻塞 I/O模型：**主要是**程序主动检查 I/O 状态**，需要结合多路复用机制处理高并发。

- **异步 I/O模型：** 将 I/O 操作完全交给内核管理，**程序仅需处理完成通知**，是更高效的模型，但实现复杂度更高。

在 Unix 和类 Unix 操作系统中，**文件描述符（File descriptor，FD）是用于访问文件或者其他 I/O 资源的抽象句柄**，例如：管道或者网络套接字。而**不同的 I/O 模型会使用不同的方式操作文件描述符。**

### 阻塞I/O模型

当发出**IO读写的系统调用**时，应用程序被阻塞。是最常见的I/O模型。

系统调用syscall，IO操作，**需要与操作系统交换**（文件，网络属于操作系统资源），这类操作称为系统调用。

如图：

<img src="https://fynotefile.oss-cn-zhangjiakou.aliyuncs.com/fynote/fyfile/13080/1684143948010/81ad8124b7294ca4a5cd46a959fcf035.png" alt="image.png" style="zoom: 50%;" />

**编码时，常用到的也是阻塞I/O：**

网络阻塞I/O示例：

```go
// 网络IO(使用系统调用syscall的IO)的阻塞
func BIONet() {
	addr := "127.0.0.1:5678"
	wg := sync.WaitGroup{}

	// 1.模拟读，体会读的阻塞状态
	wg.Add(1)
	go func(wg *sync.WaitGroup) {
		defer wg.Done()
		conn, _ := net.Dial("tcp", addr)
		defer conn.Close()
		buf := make([]byte, 1024)
		// 注意：两次时间的间隔
		log.Println("start read.", time.Now().Format("03:04:05.000"))
		n, _ := conn.Read(buf)
		log.Println("content:", string(buf[:n]), time.Now().Format("03:04:05.000"))
	}(&wg)

	// 2.模拟写
	wg.Add(1)
	go func(wg *sync.WaitGroup) {
		defer wg.Done()

		l, _ := net.Listen("tcp", addr)
		defer l.Close()

		for {
			conn, _ := l.Accept()
			go func(conn net.Conn) {
				defer conn.Close()
				log.Println("connected.")

				// 阻塞时长
				time.Sleep(3 * time.Second)
				conn.Write([]byte("Blocking I/O"))
			}(conn)
		}
	}(&wg)

	wg.Wait()
}
```

测试：

```shell
> go test -run BIONet
2023/05/15 19:42:24 connected.
2023/05/15 19:42:24 start read. 07:42:24.150
2023/05/15 19:42:27 content: Blocking I/O 07:42:27.160
```

阻塞时长3s.

Channel的阻塞I/O示例：

```go
// Channel(Go的自管理的IO）的阻塞
func BIOChannel() {
	// 0初始化数据
	wg := sync.WaitGroup{}
	// IO channel
	ch := make(chan struct{})

	// 1模拟读，体会读的阻塞状态
	wg.Add(1)
	go func(wg *sync.WaitGroup) {
		defer wg.Done()
		log.Println("start read.", time.Now().Format("03:04:05.000"))

		content := <-ch // IO Read, Receive

		log.Println("content:", content, time.Now().Format("03:04:05.000"))
	}(&wg)

	// 2模拟写
	wg.Add(1)
	go func(wg *sync.WaitGroup) {
		defer wg.Done()
		// 阻塞时长
		time.Sleep(3 * time.Second)
		ch <- struct{}{} // Write, Send
	}(&wg)
	wg.Wait()
}
```

测试

```shell
> go test -run BIOChannel
2023/05/15 19:48:56 start read. 07:48:56.855
2023/05/15 19:48:59 content: {} 07:48:59.857
```

阻塞时长3s.

### 非阻塞I/O模型

Non-Blocking IO Model，当FD为非阻塞时，IO操作，会立即返回，不会在未就绪时阻塞（未就绪返回错误）。当然资源就绪时是可以正确返回的。

如图：

![image.png](https://fynotefile.oss-cn-zhangjiakou.aliyuncs.com/fynote/fyfile/13080/1684143948010/2088c97a39b54064ad2ec180b73b12d2.png)

非阻塞IO，**在并发编程时非常常用**，可以将原本应该阻塞的goroutine（或者线程），来处理其他事件。

网络非阻塞编码示例：

conn.Read() 或 conn.Write IO 操作，目前不具备非阻塞操作。但可以通过设置截止时间来完成。

```go
// 网络IO(使用系统调用syscall的IO)的非阻塞
func NIONet() {
	addr := "127.0.0.1:5678"
	wg := sync.WaitGroup{}

	// 1模拟读，体会读的阻塞状态
	wg.Add(1)
	go func(wg *sync.WaitGroup) {
		defer wg.Done()

		conn, _ := net.Dial("tcp", addr)
		defer conn.Close()

		buf := make([]byte, 1024)
		// 注意：两次时间的间隔
		log.Println("start read.", time.Now().Format("03:04:05.000"))
		// 通过设置截止时间来实现非阻塞
		conn.SetReadDeadline(time.Now().Add(400 * time.Millisecond))
		n, _ := conn.Read(buf)
		log.Println("content:", string(buf[:n]), time.Now().Format("03:04:05.000"))
	}(&wg)

	// 2模拟写
	wg.Add(1)
	go func(wg *sync.WaitGroup) {
		defer wg.Done()

		l, _ := net.Listen("tcp", addr)
		defer l.Close()

		for {
			conn, _ := l.Accept()
			go func(conn net.Conn) {
				defer conn.Close()
				log.Println("connected.")

				// 阻塞时长
				time.Sleep(3 * time.Second)
				conn.Write([]byte("Blocking I/O"))
			}(conn)
		}
	}(&wg)

	wg.Wait()
}
```

**Channel配合Select的default子句**，可以完成非阻塞操作：

```go
// 网络IO(使用系统调用syscall的IO)的非阻塞
func NIONet() {
	addr := "127.0.0.1:5678"
	wg := sync.WaitGroup{}

	// 1模拟读，体会读的阻塞状态
	wg.Add(1)
	go func(wg *sync.WaitGroup) {
		defer wg.Done()

		conn, _ := net.Dial("tcp", addr)
		defer conn.Close()

		buf := make([]byte, 1024)
		// 注意：两次时间的间隔
		log.Println("start read.", time.Now().Format("03:04:05.000"))
		// 设置截止时间
		conn.SetReadDeadline(time.Now().Add(400 * time.Millisecond))
		n, _ := conn.Read(buf)
		log.Println("content:", string(buf[:n]), time.Now().Format("03:04:05.000"))
	}(&wg)

	// 2模拟写
	wg.Add(1)
	go func(wg *sync.WaitGroup) {
		defer wg.Done()

		l, _ := net.Listen("tcp", addr)
		defer l.Close()

		for {
			conn, _ := l.Accept()
			go func(conn net.Conn) {
				defer conn.Close()
				log.Println("connected.")

				// 阻塞时长
				time.Sleep(3 * time.Second)
				conn.Write([]byte("Blocking I/O"))
			}(conn)
		}
	}(&wg)

	wg.Wait()
}
```

conn的IO操作，**配合Channel，也可以完成非阻塞操作**

```go
// Channel(Go的自管理的IO）的非阻塞
func NIOChannel() {
	// 0初始化数据
	wg := sync.WaitGroup{}
	// IO channel
	ch := make(chan struct{ id uint })

	// 1模拟读，体会读的阻塞状态
	wg.Add(1)
	go func(wg *sync.WaitGroup) {
		defer wg.Done()
		log.Println("start read.", time.Now().Format("03:04:05.000"))

		content := struct{ id uint }{}
		select {
		case content = <-ch: // IO Read, Receive
		default:
		}

		log.Println("content:", content, time.Now().Format("03:04:05.000"))
	}(&wg)

	// 2模拟写
	wg.Add(1)
	go func(wg *sync.WaitGroup) {
		defer wg.Done()

		// 阻塞时长
		time.Sleep(3 * time.Second)
		ch <- struct{ id uint }{42} //匿名结构体，Write, Send
	}(&wg)

	wg.Wait()
}

```

网络的IO也可以配合select channel完成非阻塞的操作，IO操作，将内容发送到Channel，外层select处理channel，在网络编程中，非常常用，示例代码如下。

```go
// 网络IO(使用系统调用syscall的IO)的非阻塞
func NIONetChannel() {
	addr := "127.0.0.1:5678"
	wg := sync.WaitGroup{}

	// 1模拟读，体会读的阻塞状态
	wg.Add(1)
	go func(wg *sync.WaitGroup) {
		defer wg.Done()

		conn, _ := net.Dial("tcp", addr)
		defer conn.Close()

		// 注意：两次时间的间隔
		log.Println("start read.", time.Now().Format("03:04:05.000"))

		// 独立的goroutine，完成Read操作，将结果发送到channel中
		wgwg := sync.WaitGroup{}
		chRead := make(chan []byte)
		wgwg.Add(1)
		go func() {
			defer wgwg.Done()
			buf := make([]byte, 1024)
			n, _ := conn.Read(buf)
			chRead <- buf[:n]
		}()

		//time.Sleep(100 * time.Millisecond)

		// 使用select + default实现非阻塞操作
		var data []byte
		select {
		case data = <-chRead:
		default:
		}

		log.Println("content:", string(data), time.Now().Format("03:04:05.000"))
		wgwg.Wait()
	}(&wg)

	// 2模拟写
	wg.Add(1)
	go func(wg *sync.WaitGroup) {
		defer wg.Done()

		l, _ := net.Listen("tcp", addr)
		defer l.Close()

		for {
			conn, _ := l.Accept()
			go func(conn net.Conn) {
				defer conn.Close()
				log.Println("connected.")

				// 阻塞时长
				time.Sleep(3 * time.Second)
				conn.Write([]byte("Blocking I/O"))
			}(conn)
		}
	}(&wg)

	wg.Wait()
}

```

强调，以上的Go的程序的例子，都不是通过设置FD的属性实现的。而是通过外部技术实现的（相当于只是模拟一下），截止时间和Select语句。某些语言是支持在FD上设置属性的，C语言：

```c
int flags = fcntl(fd, F_GETFL, 0);
fcntl(fd, F_SETFL, flags | O_NONBLOCK);
```

### 信号驱动I/O模型

信号驱动I/O(signal-driven I/O)，就是**预先告知系统内核，当某个FD准备发生某件事情的时候，让内核发送一个信号通知应用进程**。

如图：

<img src="https://fynotefile.oss-cn-zhangjiakou.aliyuncs.com/fynote/fyfile/13080/1684143948010/649e4f86c90b4f93af59756df788fea1.png" alt="image.png" style="zoom:80%;" />

信号驱动I/O模型，在等待信号的期间，**应用程序不阻塞**，这意味着如果 I/O 操作不能立即完成，系统调用会返回错误（如 `EAGAIN`），而不是阻塞程序。

### 异步I/O模型

应用程序告知内核启动某个操作，**并让内核在整个操作完成之后，通知应用程序**，整个I/O操作由应用程序完成。这种模型与信号驱动模型的主要区别在于，**信号驱动I/O只是由内核通知我们文件描述符的状态变化**（合适就可以开始下一个IO操作 ），而**异步IO模型是由内核通知我们操作什么时候完成**，而整个I/O过程是由内核完成。

**信号驱动I/O模型和异步I/O模型的区别：**

| 特性         | 信号驱动 I/O 模型                         | 异步 I/O 模型                  |
| ------------ | ----------------------------------------- | ------------------------------ |
| **控制权**   | 应用程序负责完成 I/O 操作（读取或写入）。 | 内核负责完成整个 I/O 操作。    |
| **通知时机** | **文件描述符状态就绪时通知。**            | I/O 操作完成后通知应用程序。   |
| **I/O 责任** | 应用程序需主动调用 I/O 操作。             | 内核直接将结果传递给应用程序。 |
| **适用场景** | 适合中等复杂度的异步任务。                | 适合高性能、大规模异步任务。   |

### 多路复用I/O模型

多路复用，Multiplexing，**指的是监听一组FD，当FD的状态发生变化时（变为可读或可写），内核通知应用程序，应用程序完成对FD的操作。**

I/O 多路复用的核心在于：

- 应用程序不直接阻塞在某个 I/O 操作上，而是通过系统调用（如 `select`、`poll`、`epoll`）等待**多个文件描述符**中的一个或多个变为就绪状态（如可读、可写、发生错误等）。
- 当某个文件描述符就绪时，**程序执行相应的 I/O 操作。**

如图所示：

![image.png](https://fynotefile.oss-cn-zhangjiakou.aliyuncs.com/fynote/fyfile/13080/1684143948010/f6a8449a0e534844987fecf2e865b882.png)

我们使用**net包，完成TCP的服务器端，就是典型的多路复用：**

- 一个goroutine在持续监听是否有连接建立
- 一旦建立连接，就获取了一个FD，交由一个独立的goroutine处理
- 核心方法：net.Listen, net.Accept, net.Dial, conn.Read, conn.Write

```go
func() {
    l, _ := net.Listen("tcp", addr)
    defer l.Close()

    for {
        conn, _ := l.Accept()
        go func(conn net.Conn) {
            defer conn.Close()
            log.Println("connected.")

            // 阻塞时长
            time.Sleep(3 * time.Second)
            conn.Write([]byte("Blocking I/O"))
        }(conn)
    }
}
```

Go语言中对I/O多路复用做了封装，底层同样是基于OS操作系统的多路复用特性，不同操作系统使用的多路复用实现不同，Go支持的列表如下：

```go
// linux
src/runtime/netpoll_epoll.go
// macOS
src/runtime/netpoll_kqueue.go
src/runtime/netpoll_solaris.go
// windows
src/runtime/netpoll_windows.go
src/runtime/netpoll_aix.go
src/runtime/netpoll_fake.go
```

**Go会基于当前的操作系统，选择对应的实现**。所以我们经常说**，go的net网络编程是基于epoll实现的，可见这主要针对于Linux来说的。**

系统不同的I/O多路复用模型介绍：

- select

  - **存在最大描述符数量的限制，通常为1024**（底层实现为文件描述符**数组**）

  - 需要在内核开辟空间存储文件描述符
  - 良好的跨平台性能
  - **会将所有的文件描述符都返回，需要程序自己去遍历和区分是哪一个文件描述符**

- poll

  - **底层采用链表的方式实现，没有数量上的限制**

    - 每一个节点上的node都是一个pollfd

    - 包含文件描述符、发生的事件

  - 同样需要程序自己去遍历获取发生变化的文件描述符

- epoll

  - **同样没有描述符数量的限制**

  - 使用**红黑树**的形式管理文件描述符以及对应的监听事件
  - 当触发对应事件的时候，进行回调，放入双向链表节点
  - 程序获取发生的变化的文件描述符的时候，只需要去检查双向链表节点中有没有即可，有的话，将事件以及事件数量返回给用户

**模型对比表**

| **特性**           | **阻塞 I/O**                     | **非阻塞 I/O**                   | **I/O 多路复用**                                  | **信号驱动 I/O**                    | **异步 I/O**                                    |
| ------------------ | -------------------------------- | -------------------------------- | ------------------------------------------------- | ----------------------------------- | ----------------------------------------------- |
| **I/O 操作流程**   | 阻塞等待 I/O 就绪 + 数据传输完成 | 轮询检查 I/O 就绪 + 数据传输完成 | 通过 `select` 或 `poll` 等等待就绪后执行 I/O 操作 | 信号通知 I/O 就绪后主动执行数据传输 | 应用程序发起 I/O 请求后，内核完成所有操作并通知 |
| **应用程序参与度** | 全程参与 I/O 操作                | 主动轮询文件描述符               | 等待就绪时阻塞，I/O 操作由应用程序完成            | 收到信号后执行数据传输              | 发起请求后无需再参与，结果直接返回              |
| **内核通知方式**   | 无                               | 无                               | 就绪事件通知                                      | 信号通知                            | 操作完成通知                                    |
| **是否阻塞**       | 阻塞                             | 非阻塞                           | 阻塞（等待就绪）                                  | 非阻塞                              | 非阻塞                                          |
| **CPU 使用效率**   | 较低（等待期间 CPU 空闲）        | 较低（轮询耗费 CPU 时间）        | 较高（等待时释放 CPU 资源）                       | 较高（信号触发避免轮询）            | 最高（操作完全由内核负责）                      |
| **适用场景**       | 简单 I/O 操作                    | 简单但低效的非阻塞操作           | 高并发连接场景                                    | 中小型异步通知场景                  | 高性能、高并发的 I/O 场景                       |

## 网络轮询器（epoll原理需要再看）

**Go的网络轮询器是对 I/O 多路复用技术的封装**，配合Goroutine的GMP并发调度，实现Go语言层面的I/O多路复用。应用在文件 I/O、网络 I/O 以及计时器操作中。

网络轮询器的核心操作有：

1. 网络轮询器的**初始化**；
2. 如何向网络轮询器**加入**待监听的FD；
3. 如何从网络轮询器**获取**触发的事件；

三个核心操作。

### Epoll核心概念

`epoll` 是 Linux 内核提供的一种高效的 **I/O 多路复用机制**，解决了传统 `select` 和 `poll` 在监控大量文件描述符时的性能瓶颈问题。

源码分析会使用 Linux 操作系统上的 `epoll` 实现作介绍，其他 I/O 多路复用模块的实现大同小异。在学习网络轮询器时，注意如何封装epoll或其他实现。

### **Epoll 工作流程**

初始化网络轮循器的操作如下

1. **创建 epoll 文件描述符：**使用 `epoll_create` 创建一个 `epoll` 文件描述符。（接着会创建一个channel）
2. **添加文件描述符：**使用 `epoll_ctl` **添加需要监控的文件描述符**，并设置监控的事件（如可读、可写）。
3. **等待事件发生：**使用 `epoll_wait` 等待文件描述符的状态变化，当有事件发生时返回相应的事件。
4. **处理事件：**遍历 `epoll_wait` 返回的事件集合，执行相应的 I/O 操作。
5. **循环处理：**重复步骤 3 和 4，直到程序退出。

核心概念：

- **epoll在Linux内核中构建了一个文件系统，该文件系统采用红黑树来构建**。因为数据结构红黑树在增加和删除上面的效率高。
- epoll提供了两种触发模式，水平触发(LT， Level Triggered)和边沿触发(ET， Edge Triggered)
  - 水平触发，文件描述符状态满足条件时持续通知
  - 边沿触发，文件描述符状态变化时仅通知一次
- epoll的工作流程
  - epoll_create，epoll初始化
  - epoll_ctl，epoll操作
  - epoll_wait，事件就绪等待

### 初始化

当使用文件 I/O、网络 I/O 以及计时器时：

- `internal/poll.pollDesc.init`
  - net.netFD.init， 初始化网络 I/O
  - os.newFile， 初始化文件 I/O
- runtime.doaddtimer， 增加新的计时器

以上调用会通过：`netpollGenericInit` 函数完成初始化：

```go
// src/runtime/netpoll.go
func netpollGenericInit() {
	if atomic.Load(&netpollInited) == 0 {
		lockInit(&netpollInitLock, lockRankNetpollInit)
		lock(&netpollInitLock)
		if netpollInited == 0 {
            // 初始化网络轮询器
			netpollinit()
			atomic.Store(&netpollInited, 1)
		}
		unlock(&netpollInitLock)
	}
}
```

注意函数：`netpollinit()` 会基于特定OS平台的多路复用技术来说实现。本例中，Linux上就是epoll。

```go
// src/runtime/netpoll_epoll.go 

var (
	epfd int32 = -1 // epoll descriptor
	netpollBreakRd, netpollBreakWr uintptr // for netpollBreak
	netpollWakeSig atomic.Uint32 // used to avoid duplicate calls of netpollBreak
)

func netpollinit() {
	var errno uintptr
    // 创建Epoll文件描述符FD
	epfd, errno = syscall.EpollCreate1(syscall.EPOLL_CLOEXEC)
	if errno != 0 {
		println("runtime: epollcreate failed with", errno)
		throw("runtime: netpollinit failed")
	}
    // 创建非阻塞管道
	r, w, errpipe := nonblockingPipe()
	if errpipe != 0 {
		println("runtime: pipe failed with", -errpipe)
		throw("runtime: pipe failed")
	}
	ev := syscall.EpollEvent{
		Events: syscall.EPOLLIN,
	}
	*(**uintptr)(unsafe.Pointer(&ev.Data)) = &netpollBreakRd
    // 将epoll FD加入事件监听
	errno = syscall.EpollCtl(epfd, syscall.EPOLL_CTL_ADD, r, &ev)
	if errno != 0 {
		println("runtime: epollctl failed with", errno)
		throw("runtime: epollctl failed")
	}
	netpollBreakRd = uintptr(r)
	netpollBreakWr = uintptr(w)
}
```

初始化的核心功能在以上函数中实现：

- 是调用 `syscall.EpollCreate1` 创建一个新的 `epoll` 文件描述符，这个文件描述符会在整个程序的生命周期中使用
- 通过 `runtime.nonblockingPipe` 创建一个用于通信的管道
- 使用 `epollctl` 将用于读取数据的文件描述符打包成 `epollevent` 事件加入监听

### 轮询事件

初始化时，使用 `epollctl` 将用于读取数据的文件描述符打包成 `epollevent` 事件加入监听。

同时也需要将监听的文件描述符的可读和可写状态事件，加入到全局轮询的文件描述符epfd中。

也就是，在调用 `internal/poll.pollDesc.init` 初始化时，同时完成初始化全局文件描述符，和注册所监听文件描述的可读、可写事件监听：

```go
// $GOROOT/src/internal/poll/fd_poll_runtime.go

func (pd *pollDesc) init(fd *FD) error {
    // 初始化全局文件描述符 epfd
	serverInit.Do(runtime_pollServerInit)
    // 所监听文件描述的可读、可写事件监听
	ctx, errno := runtime_pollOpen(uintptr(fd.Sysfd))
	if errno != 0 {
		return errnoErr(syscall.Errno(errno))
	}
	pd.runtimeCtx = ctx
	return nil
}
```

在epoll的pollOpen实现中，使用EpollCtl，完成了注册具体事件监听的操作：

```go
// $GOROOT/src/runtime/netpoll_epoll.go 
func netpollopen(fd uintptr, pd *pollDesc) uintptr {
	var ev syscall.EpollEvent
    // 事件：输入、输出、Close、ET边沿触发
	ev.Events = syscall.EPOLLIN | syscall.EPOLLOUT | syscall.EPOLLRDHUP | syscall.EPOLLET
	*(**pollDesc)(unsafe.Pointer(&ev.Data)) = pd
	return syscall.EpollCtl(epfd, syscall.EPOLL_CTL_ADD, int32(fd), &ev)
}
```

### 事件循环

初始化epoll（操作系统的多路复用）后，就需要操作系统中的 I/O 多路复用机制和 Go 语言的运行时联系起来，进而达到基于OS的Multiplexing来实现Go语言层面的IO多路复用编程的目的。

存在两个核心过程：

- Go运行时调度Goroutine让出线程并等待读写事件
- 当多路复用读写事件触发时唤醒Goroutine执行

该过程，称为事件循环（EventLoop）。

**Goroutine让出线程并等待读写事件**

当我们在FD上执行读写操作时，如果FD不可用，当前 Goroutine 会等待FD的可读或者可写：

```go
// $GOROOT/src/runtime/netpoll.go
// go:linkname poll_runtime_pollWait internal/poll.runtime_pollWait
func poll_runtime_pollWait(pd *pollDesc, mode int) int {
	errcode := netpollcheckerr(pd, int32(mode))
	if errcode != pollNoError {
		return errcode
	}
	// As for now only Solaris, illumos, and AIX use level-triggered IO.
	if GOOS == "solaris" || GOOS == "illumos" || GOOS == "aix" {
		netpollarm(pd, mode)
	}
    // epoll 阻塞
	for !netpollblock(pd, int32(mode), false) {
		errcode = netpollcheckerr(pd, int32(mode))
		if errcode != pollNoError {
			return errcode
		}
		// Can happen if timeout has fired and unblocked us,
		// but before we had a chance to run, timeout has been reset.
		// Pretend it has not happened and retry.
	}
	return pollNoError
}

```

其中netpollblock就是Block阻塞。会调用 runtime.gopart()让出当前线程M，将Goroutine转为阻塞休眠状态。

**轮询网络**

Go 的runtime会在调度中轮询网络，轮询网络的核心过程是：

- 计算 `epoll` 系统调用需要等待的时间
- 调用 `epollwait` 等待可读或者可写事件的发生；
- 循环处理 `epollevent` 事件；

轮询网络：

```go
// $GOROOT/src/runtime/netpoll_epoll.go

func netpoll(delay int64) gList {
	if epfd == -1 {
		return gList{}
	}
    // 计算系统调用等待时间
	var waitms int32
	if delay < 0 {
		waitms = -1
	} else if delay == 0 {
		waitms = 0
	} else if delay < 1e6 {
		waitms = 1
	} else if delay < 1e15 {
		waitms = int32(delay / 1e6)
	} else {
		// An arbitrary cap on how long to wait for a timer.
		// 1e9 ms == ~11.5 days.
		waitms = 1e9
	}
	var events [128]epollevent
  
    // 等待可读或者可写事件的发生
retry:
	n := epollwait(epfd, &events[0], int32(len(events)), waitms)
	if n < 0 {
		if n != -_EINTR {
			println("runtime: epollwait on fd", epfd, "failed with", -n)
			throw("runtime: netpoll failed")
		}
		// If a timed sleep was interrupted, just return to
		// recalculate how long we should sleep now.
		if waitms > 0 {
			return gList{}
		}
		goto retry
	}
	var toRun gList
  
    // 循环处理 `epollevent` 事件
	for i := int32(0); i < n; i++ {
		ev := &events[i]
		if ev.events == 0 {
			continue
		}

		if *(**uintptr)(unsafe.Pointer(&ev.data)) == &netpollBreakRd {
			if ev.events != _EPOLLIN {
				println("runtime: netpoll: break fd ready for", ev.events)
				throw("runtime: netpoll: break fd ready for something unexpected")
			}
			if delay != 0 {
				// netpollBreak could be picked up by a
				// nonblocking poll. Only read the byte
				// if blocking.
				var tmp [16]byte
				read(int32(netpollBreakRd), noescape(unsafe.Pointer(&tmp[0])), int32(len(tmp)))
				atomic.Store(&netpollWakeSig, 0)
			}
			continue
		}

		var mode int32
		if ev.events&(_EPOLLIN|_EPOLLRDHUP|_EPOLLHUP|_EPOLLERR) != 0 {
			mode += 'r'
		}
		if ev.events&(_EPOLLOUT|_EPOLLHUP|_EPOLLERR) != 0 {
			mode += 'w'
		}
		if mode != 0 {
			pd := *(**pollDesc)(unsafe.Pointer(&ev.data))
			pd.everr = false
			if ev.events == _EPOLLERR {
				pd.everr = true
			}
            // 处理Ready可用事件
			netpollready(&toRun, pd, mode)
		}
	}
	return toRun
}
```

其中，在处理事件时，利用netpollready来处理IO事件。

```go
// $GOROOT/src/runtime/netpoll.go
//go:nowritebarrier
func netpollready(toRun *gList, pd *pollDesc, mode int32) {
	var rg, wg *g
	if mode == 'r' || mode == 'r'+'w' {
        // 唤醒Goroutine的执行
		rg = netpollunblock(pd, 'r', true)
	}
	if mode == 'w' || mode == 'r'+'w' {
        // 唤醒Goroutine的执行
		wg = netpollunblock(pd, 'w', true)
	}
	if rg != nil {
		toRun.push(rg)
	}
	if wg != nil {
		toRun.push(wg)
	}
}
```

其中，netpollunblock，会在读写事件发生时，取消Goroutine的阻塞，唤醒Goroutine的执行。

也就是netpollblock的逆向操作。


## 小结

- IO模型
  - 阻塞I/O模型，Blocking I/O Model
  - 非阻塞I/O模型，Non-Blocking I/O Model
  - 信号驱动I/O模型，Signal-Driven I/O Model
  - 异步I/O模型 ，Asynchronous I/O Model
  - I/O多路复用模型,，Multiplexing I/O Model
- Go的网络轮询器实现Go层面的多路复用
