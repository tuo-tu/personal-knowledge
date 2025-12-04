# Channel通信

## Channel概述

> **不要通过共享内存的方式进行通信，而是应该通过通信的方式共享内存**

这是Go语言最核心的设计模式之一。

在很多主流的编程语言中，多个线程传递数据的方式一般都是共享内存，而Go语言中多Goroutine通信的主要方案是Channel。Go语言也可以使用共享内存的方式支持Goroutine通信。

![image.png](https://fynotefile.oss-cn-zhangjiakou.aliyuncs.com/fynote/fyfile/13080/1672136903072/cceb79db926449819974f11cae5a8d0f.png)

Go语言实现了**CSP**通信模式，CSP是Communicating Sequential Processes的缩写，**通信顺序进程**。Goroutine和Channel分别对应CSP中的实体和传递信息的媒介。CSP是Tony Hoare于1977年提出。

Channel提供可接收和发送特定类型值的用于并发函数(Goroutine)通信的数据类型，是满足FIFO（先进先出）原则的队列类型，先进先出不仅体现在数据类型上，也体现在操作上：

- channel类型的元素是先进先出的，先发送到channel的value会先被receive
- 先向Channel发送数据的Goroutine会先执行
- 先从Channel接收数据的Goroutine会先执行

如图：

![image.png](https://fynotefile.oss-cn-zhangjiakou.aliyuncs.com/fynote/fyfile/13080/1672136903072/6a24b9d6e254423dbfb04b186c9fa1d3.png)

## Channel操作语法

- 一个关键字
  - chan
  - chan <-
  - <- chan
- 两个函数
  - make
  - close
- 一个语句
  - 发送语句 ch <- expression
- 一个操作符
  - 接收操作符 <- ch

### Channel类型

channel是Go语言中的数据类型，类型声明如下：

```
ChannelType = ( "chan" | "chan" "<-" | "<-" "chan" ) ElementType .
```

其中：

- chan channel类型关键字
- <- 操作符，用于Channel收发，定义Channel类型时，用于表示Channel的方向：
  - 默认是双向，可收可发。还可以定义为定向的，仅收仅发。
  - chan<- 仅发送Channel
  - <-chan 仅接收Channel
- ElementType Channel中元素类型

可见，Channel类型是用于接收或发送特定类型元素的Go数据类型。

示例：

```go
chan int
chan struct{}
chan<- int
<-chan int
```

### 初始化Channel值

内建函数make()可用于初始化Channel值。支持两个参数：

```go
make(ChannelType, Capacity)
```

其中：

- ChannelType是channel类型
- Capacity是缓冲容量。可以省略或为0，表示无缓冲Channel

channel是引用类型，类似于map和slice。

示例：

```go
ch := make(chan int)
var ch = make(chan int)
ch := make(chan int, 10)
ch := make(<-chan int)
ch := make(chan<- int, 10)
```

未使用make()初始化的channel为nil。nil channel不能执行收发通信操作，例如：

```go
var ch chan int
```

ch就是nil channel。

### Send语句和接收操作符

- Send语句用于向Channel发送值
- 接收操作符用于从Channel中接收值

Send语句语法：

```
SendStmt = Channel "<-" Expression .
Channel  = Expression .
```

```go
ch <- Expression
ch <- 42
ch <- f()
```

接收操作符语法：

```
<-ch
v1 := <-ch // 声明
v = <-ch // 赋值
f(<-ch) // 函数调用
<-strobe // 等待接收
```

### 关闭channel

内置函数close()用于关闭channel。

关闭Channel的意思是记录该Channel不能再被发送任何元素了，而不是销毁该Channel的意思。也就意味着关闭的Channel是可以继续接收值的。因此：

- 向已关闭的Channel发送回引发runtime panic
- 关闭nil Channel会引发runtime panic
- 不能关闭仅接收Channel
- 不能关闭已经关闭的Channel，否则会引发runtime panic

当从已关闭的Channel接收时：

- 可以接收关闭前发送的全部值
- 若没有已发送的值会返回类型的零值，不会被阻塞

使用接收操作符的多值返回结构，可以判断Channel是否已经关闭：

```go
var x, ok = <-ch
x, ok := <-ch
```

- ok为true，channel未关闭
- ok为false，channel已关闭

### for range channel

for语句的range子句可以持续地从Channel中接收元素，语法如下：

```go
for e := range ch {
    // e是ch中元素值
}
```

持续接收操作与接收操作<-行为一致：

- 若ch为nil channel会阻塞
- 若ch没有已发送元素会阻塞

for会持续执行到channel被关闭，关闭后，若channel中存在已发送元素，for会全部读取完毕。

示例：

```go
func ChannelFor() {
    // 一，初始化部分数据
    ch := make(chan int) // 无缓冲的channel
    wg := sync.WaitGroup{}

    // 二，持续发送
    wg.Add(1)
    go func() {
        defer wg.Done()
        for i := 0; i < 5; i++ {
            // random send value
            ch <- rand.Intn(10)
        }
        // 关闭
        close(ch)
    }()

    // 三，持续接收，for range
    wg.Add(1)
    go func() {
        defer wg.Done()
        // 持续接收
        for e := range ch {
            println("received from ch, element is ", e)
        }
    }()

    wg.Wait()
}
```

## 缓冲与无冲 channel

Channel区别于是否存在缓冲区，分为：

- 缓冲Channel，make(chan T, cap)，cap是大于0的值。
- 无缓冲Channel, make(chan T), make(chan T, 0)

### 无缓冲channel

![image.png](https://fynotefile.oss-cn-zhangjiakou.aliyuncs.com/fynote/fyfile/13080/1672136903072/2a70e2fd29a04678ad97a79d5a365829.png)

也称为同步Channel，只有当发送方和接收方都准备就绪时，通信才会成功。

同步操作示例：

```go
func ChannelSync() {
    // 初始化数据
    ch := make(chan int)
    wg := sync.WaitGroup{}

    // 间隔发送
    wg.Add(1)
    go func() {
        defer wg.Done()
        for i := 0; i < 5; i++ {
            ch <- i
            println("Send ", i, ".\tNow:", time.Now().Format("15:04:05.999999999"))
            // 间隔时间
            time.Sleep(1 * time.Second)
        }
        close(ch)
    }()

    // 间隔接收
    wg.Add(1)
    go func() {
        defer wg.Done()
        for v := range ch {
            println("Received ", v, ".\tNow:", time.Now().Format("15:04:05.999999999"))
            // 间隔时间，注意与send的间隔时间不同
            time.Sleep(3 * time.Second)
        }
    }()

    wg.Wait()
}
```

代码中，采用同步channel，使用两个goroutine完成发送和接收。每次发送和接收的时间间隔不同。我们分别打印发送和接收的值和时间。注意结果：

- 发送和接收时间一致
- 间隔以长的为准，可见发送和接收操作为同步操作

同步Channel适合在gotoutine间做同步信号！

### 缓冲Channel

![image.png](https://fynotefile.oss-cn-zhangjiakou.aliyuncs.com/fynote/fyfile/13080/1672136903072/91c9919010274c759c3f1120ac32c208.png)

缓冲Channel也称为异步Channel，接收和发送方不用等待双方就绪即可成功。缓冲Channel会存在一个容量为cap的缓冲空间。当使用缓冲Channel通信时，接收和发送操作是在操作Channel的Buffer：

- 接收时，从缓冲中接收元素，只要缓冲不为空，不会阻塞。反之，缓冲为空，会阻塞，goroutine挂起
- 发送时，向缓冲中发送元素，只要缓冲未满，不会阻塞。反之，缓冲满了，会阻塞，goroutine挂起

是典型的队列操作。

缓冲channel操作示例：

```go
func ChannelASync() {
    // 初始化数据
    ch := make(chan int, 5)
    wg := sync.WaitGroup{}

    // 间隔发送
    wg.Add(1)
    go func() {
        defer wg.Done()
        for i := 0; i < 5; i++ {
            ch <- i
            println("Send ", i, ".\tNow:", time.Now().Format("15:04:05.999999999"))
            // 间隔时间
            time.Sleep(1 * time.Second)
        }
    }()

    // 间隔接收
    wg.Add(1)
    go func() {
        defer wg.Done()
        for v := range ch {
            println("Received ", v, ".\tNow:", time.Now().Format("15:04:05.999999999"))
            // 间隔时间，注意与send的间隔时间不同
            time.Sleep(3 * time.Second)
        }
    }()

    wg.Wait()
}
```

代码中，与同步channel一致，只是采用了容量为5的缓冲channel，使用两个goroutine完成发送和接收。每次发送和接收的时间间隔不同。我们分别打印发送和接收的值和时间。注意结果：

- 发送和接收时间不同
- 发送和接收操作不会阻塞，可见发送和接收操作为异步操作

缓冲channel非常适合做goroutine的数据通信了。

### 长度和容量，len()和cap()

内置函数 len() 和 cap() 可以分别获取：

- len()长度，缓冲中元素个数。
- cap()容量，缓冲的总大小。cap()返回0，意味着是无缓冲通道

## 单向Channel

![image.png](https://fynotefile.oss-cn-zhangjiakou.aliyuncs.com/fynote/fyfile/13080/1672136903072/27c1c8ed7a6b4bf895c4ea5c92638ed5.png)

单向Channel，指的是仅支持接收或仅支持发送操作的Channel。语法上：

- `chan<- T` 仅发送Channel
- `<-chan T` 仅接收Channel

单向Channel的意义在于约束Channel的使用方式。

仅使用单向Channel通常没有实际意义，单向Channel最典型的使用方式是：

**使用单向通道约束双向通道的操作。**

语法上来说，就是我们会将双向Channel转换为单向Channel来使用。典型使用在函数参数或返回值类型中。

示例代码：

```go
func ChannelDirectional() {
    // 初始化数据
    ch := make(chan int)
    wg := &sync.WaitGroup{}

    // send and receive
    wg.Add(2)
    go setElement(ch, 42, wg)
    go getElement(ch, wg)

    wg.Wait()
}

// only receive channel
func getElement(ch <-chan int, wg *sync.WaitGroup) {
    defer wg.Done()

    println("received from ch, element is ", <-ch)
}

// only send channel
func setElement(ch chan<- int, v int, wg *sync.WaitGroup) {
    defer wg.Done()

    ch <- v
    println("send to ch, element is ", v)
}
```

函数getElement和setElement，分别使用了单向的接收和发送channel，在语义上表示只能接收和只能发送操作，同时程序上限定了操作。

典型的单向Channel的标准库例子：

```go
// signal.Notify()
func Notify(c chan<- os.Signal, sig ...os.Signal)

// time.After
func After(d Duration) <-chan Time
```

以上两个示例分别展示了单向Channel作为函数参数和函数返回值的语法。

## Channel结构

### Channel定义结构

![image.png](https://fynotefile.oss-cn-zhangjiakou.aliyuncs.com/fynote/fyfile/13080/1672136903072/74ff7e2ca61c4331a5d2bb3915924622.png)

Channel的结构定义为 `runtime.hchan`：

```go
// GOROOT/src/runtime/chan.go
type hchan struct {
    qcount   uint           // 元素个数。len()
    dataqsiz uint           // 缓冲队列的长度。cap()
    buf      unsafe.Pointer // 缓冲队列指针，无缓冲队列为nil
    elemsize uint16 // 元素大小
    closed   uint32
    elemtype *_type // 元素类型
    sendx    uint   // send index
    recvx    uint   // receive index
    recvq    waitq  // list of recv waiters
    sendq    waitq  // list of send waiters

    // lock protects all fields in hchan, as well as several
    // fields in sudogs blocked on this channel.
    //
    // Do not change another G's status while holding this lock
    // (in particular, do not ready a G), as this can deadlock
    // with stack shrinking.
    lock mutex
}
```

其中：

- 存储空间分为channel和channel.buf两块
- channel上记录channel的属性，**长度、容量、元素类型、元素大小，接收发送索引、接收发送等待队列**
- channel.buf为elemtype类型的array
- 若为无缓冲channel，不分配channel.buf空间
- make()初始化的核心操作就是分配内存空间

### 缓冲数组

![image.png](https://fynotefile.oss-cn-zhangjiakou.aliyuncs.com/fynote/fyfile/13080/1672136903072/3efd58b1945f49a48aad6aae0f42e045.png)

![image.png](https://fynotefile.oss-cn-zhangjiakou.aliyuncs.com/fynote/fyfile/13080/1672136903072/4b85a56eccd847848b773df63151f797.png)

缓冲为数组结构，channel记录发送和接收元素的索引:

```go
 sendx    uint   // 发送索引
 recvx    uint   // 接收索引
```

缓冲数组是循环使用的，也就是若数组的最后一个元素存储了元素，那么下一次会尝试存储在第一个元素位置。

### Channel与Goroutine的关系

![image.png](https://fynotefile.oss-cn-zhangjiakou.aliyuncs.com/fynote/fyfile/13080/1672136903072/3f0bf05adaa24010a25fd2c0db521c3b.png)

Channel记录两个属性，由于记录等待接收和发送的goroutine队列：

```go
recvq    waitq  // 等待接收goroutine队列
sendq    waitq  // 等待发送goroutine队列
```

当基于某channel的接收或发送的goroutine无法理解执行时，也就是需要阻塞时，会被记录到Channel的等待队列中。当channel可以完成相应的接收或发送操作时，从等待队列中唤醒goroutine进行操作。

其中等待队列是 runtime.waitq 类型，是一个双向链表结构，具体的某个链表节点存在两个指针，指向前后节点：

```go
// GOROOT/src/runtime/chan.go
type waitq struct {
    first *sudog
    last  *sudog
}
```

其中 *sudog 可以理解为一个挂起的goroutine。

### 初始化channel流程

make()初始化channel时，会根据是否存在缓冲，选择：

- 存在缓冲，为channel和buffer分别分配内存，同时channel.buf指向buffer地址
- 不存在缓冲，仅为channel分配内存，channel.buf为nil。
- 初始化channel中其他属性

![image.png](https://fynotefile.oss-cn-zhangjiakou.aliyuncs.com/fynote/fyfile/13080/1672136903072/f6dc4a8966934f2a9f26d9bd5574ba9b.png)

### 向channel发送流程

语句 ch <- element 向channel发送元素时，大体的执行流程如下：

- 直接发送：当channel存在等待接受者时，channel.recvq，直接将元素拷贝给等待接受者，并唤醒等待接受者goroutine将其放在M的runnext位置，下次调度立即执行
- 直接写缓冲区，当缓冲区存在空间时，将发送元素直接写入缓冲区，调整channel.sendx的位置
- 阻塞发送，当缓冲区已满或无缓冲区时，发送goroutine进入channel.sendq队列，转为阻塞状态，等待其他goroutine从channel中接收元素，进而唤醒发送goroutine

![image.png](https://fynotefile.oss-cn-zhangjiakou.aliyuncs.com/fynote/fyfile/13080/1672136903072/73f96aca61da483692faf36672a20ef8.png)

### 从channel接收流程

操作符 <- ch 从channel中接收元素，大体流程如下：

- 当存在等待发送者时，channel.sendq
  - 若无缓冲区，直接将元素从发送者拷贝到接受者，并唤醒发送者gorutine，进入runnext下次调度执行
  - 若存在缓冲区，此时缓冲区是满的，从缓冲区获取元素，并将等待发送者发送元素拷贝到缓冲区，唤醒发送者goroutine。调整channel的recvx和sendx索引位置
- 当缓冲区有元素时（无等待发送者），直接从缓冲区读取元素
- 如果缓冲区不存在或缓冲区没有元素时，接收者goroutine进入阻塞状态，进入channel.recvq接受者队列，等待发送者发送数据唤醒。

![image.png](https://fynotefile.oss-cn-zhangjiakou.aliyuncs.com/fynote/fyfile/13080/1672136903072/89f73cffcd37448cb1085ed6de144b9e.png)

### 关闭channel流程

close(ch)关闭channel，主要工作是：

- 取消channel关联的sendq和recvq队列
- 调度阻塞在sendq和recvq中的goroutine

## select 语句

`select` 语句能够从多个可读或者可写的Channel中选择一个继续执行 ，若没有Channel发生读写操作，`select` 会一直阻塞当前Goroutine。

![image.png](https://fynotefile.oss-cn-zhangjiakou.aliyuncs.com/fynote/fyfile/13080/1672136903072/9d4a496c12a1483f8527c00c2c721360.png)

### select语法

```go
SelectStmt = "select" "{" { CommClause } "}" .
CommClause = CommCase ":" StatementList .
CommCase   = "case" ( SendStmt | RecvStmt ) | "default" .
RecvStmt   = [ ExpressionList "=" | IdentifierList ":=" ] RecvExpr .
RecvExpr   = Expression .
```

语法结构与 switch 类似，但case都要涉及channel操作，示例：

```go
func SelectStmt() {
    // 声明需要的变量
    var a [4]int
    var c1, c2, c3, c4 = make(chan int), make(chan int), make(chan int), make(chan int)
    var i1, i2 int

    // 用于操作channel的goroutine
    go func() {
        c1 <- 10
    }()
    go func() {
        <-c2
    }()
    go func() {
        close(c3)
    }()
    go func() {
        c4 <- 40
    }()

    // 用于select的goroutine
    go func() {
        select {
        case i1 = <-c1:
            println("received ", i1, " from c1")
        case c2 <- i2:
            println("sent ", i2, " to c2")
        case i3, ok := <-c3:
            if ok {
                println("received ", i3, " from c3")
            } else {
                println("c3 is closed")
            }
        case a[f()] = <-c4:
            println("received ", a[f()], " from c4")
        default:
            println("no communication")
        }
    }()

    // 简单sleep测试
    time.Sleep(100 * time.Millisecond)
}

func f() int {
    print("f() was run")
    return 2
}
```

测试：

```go
func TestSelectStmt(t *testing.T) {
    for i := 0; i < 10; i++ {
        println(i, ":")
        SelectStmt()
    }
}
```

### 执行流程

![image.png](https://fynotefile.oss-cn-zhangjiakou.aliyuncs.com/fynote/fyfile/13080/1672136903072/ae21b1bafe63437895c4ef8642165afb.png)

select语句的执行分为几个步骤：

1. 对于全部的 case，receive操作的channel操作数、send语句的channel和右表达式在进入select语句时只会基于源码顺序计算一次。计算结果是一组要从中接收或者发送到的channel，以及要发送的相应值。RecvStmt的左侧带有短变量声明或赋值的表达式尚未计算。
2. 如果一个或多个通信可以继续，通过伪随机数选择其中一个继续执行。否则，如果存在default case，选择该case。如果没有default case，select 语句会阻塞直到至少一个通信操作可以继续。
3. 除非选择了default case，那么相应的通信操作会被执行。
4. 如果选择的case是带有短变量声明或赋值的RecvStmt，左侧表达式会被计算，并分配接收的值（或多个值）。
5. 执行所选择的case的语句列表。

### for + select

select 匹配到可操作的case或者是defaultcase后，就执行完毕了。实操时，我们通常需要持续监听某些channel的操作，因此典型的select使用会配合for完成。

例如：持续从某个ch内获取数据

```go
func SelectFor() {
    ch := make(chan int)
    // send to channel
    go func() {
        for {
            // 模拟演示数据来自于随机数
            // 实操时，数据可以来自各种I/O，例如网络、缓存、数据库等
            ch <- rand.Intn(100)
            time.Sleep(200 * time.Millisecond)
        }
    }()
    // select receive from channel
    go func() {
        for {
            select {
            case v := <-ch:
                println("received value: ", v)
            }
        }

    }()

    time.Sleep(3 * time.Second)
}
```

### 阻塞select

以下典型的情况会直接导致阻塞goroutine：

- 不存在任何case的
- case监听都是nil channel

示例：

```go
func SelectBlock() {
    // 空select阻塞
    println("before select")
    select {}
    println("after select")

    // nil select阻塞
    var ch chan int
    go func() {
        ch <- 1024
    }()
    println("before select")
    select {
    case <-ch:
    case ch <- 42:
    }
    println("after select")
}
```

go test 测试时，会一直阻塞。若上面的代码出现在常规执行流程中，会导致 deadlock。

### nil channel的case

nil channel 不能读写，因此通过将channel设置为nil，可以控制某个case不再被执行。

例如，3秒后，不再接受ch的数据：

```go
func SelectNilChannel() {
    ch := make(chan int)
    // 写channel
    go func() {
        // 随机写入int
        rand.Seed(time.Now().Unix())
        for {
            ch <- rand.Intn(10)
            time.Sleep(400 * time.Millisecond)
        }
    }()

    // 读channel
    go func() {
        sum := 0
        t := time.After(3 * time.Second)
        for {
            select {
            case v := <-ch:
                println("received value: ", v)
                sum += v
            case <-t:
                // 将channel设置为nil，不再读写
                ch = nil
                println("ch was set nil, sum is ", sum)
            }
        }

    }()

    // sleep 5 秒
    time.Sleep(5 * time.Second)
}
```

### 带有default的select，非阻塞收发

当select语句存在default case时：

- 若没有可操作的channel，会执行default case
- 若有可操作的channel，会执行对应的case

这样select语句不会进入block状态，称之为非阻塞（non-block）的收发（channel 的接收和发送）。

示例：多人猜数字游戏，我们在乎是否有人猜中数字：

```go
func SelectNonBlock() {
    // 初始化数据
    counter := 10 // 参与人数
    max := 20     // [0, 19] // 最大范围
    rand.Seed(time.Now().UnixMilli())
    answer := rand.Intn(max) // 随机答案
    println("The answer is ", answer)
    println("------------------------------")

    // 正确答案channel
    bingoCh := make(chan int, counter)
    // wg
    wg := sync.WaitGroup{}
    wg.Add(counter)
    for i := 0; i < counter; i++ {
        // 每个goroutine代表一个猜数字的人
        go func() {
            defer wg.Done()
            result := rand.Intn(max)
            println("someone guess ", result)
            // 答案争取，写入channel
            if result == answer {
                bingoCh <- result
            }
        }()
    }
    wg.Wait()

    println("------------------------------")
    // 是否有人发送了正确结果
    // 可以是0或多个人
    // 核心问题是是否有人猜中，而不是几个人
    select {
    case result := <-bingoCh:
        println("some one hint the answer ", result)
    default:
        println("no one hint the answer")
    }
}
```

特别的情况是存在两个case，其中一个是default，另一个是channel case，那么go的优化器会优化内部这个select。内部会以if结构完成处理。因为这种情况，不用考虑随机性的问题。类似于：

```go
select {
    case result := <-bingoCh:
    println("some one hint the answer ", result)
    default:
    // 非阻塞的保证，存在default case
    println("no one hint the answer")
}

// 优化伪代码
if selectnbrecv(bingoCh) {
    println("some one hint the answer ", result)
} else {
    println("no one hint the answer")
}
```

### Race模式

![image.png](https://fynotefile.oss-cn-zhangjiakou.aliyuncs.com/fynote/fyfile/13080/1672136903072/fcceb74a358c4bc29468ad3f8a5fea80.png)

Race模式，典型的并发执行模式之一，多路同时操作资源，哪路先操作成功，优先使用，同时放弃其他路的等待。简而言之，从多个操作中选择一个最快的。核心工作：

- 选择最快的
- 停止其他未完成的

示例代码，示例从多个查询器同时读取数据，使用最先反返回结果的，其他查询器结束：

```go
func SelectRace() {
    // 一，初始化数据
    // 模拟查询结果，需要与具体的querier建立联系
    type Rows struct {
        // 数据字段

        // 索引标识
        Index int
    }
    // 模拟的querier数量
    const QuerierNum = 8
    // 用于通信的channel，数据，停止信号
    ch := make(chan Rows, 1)
    stopChs := [QuerierNum]chan struct{}{}
    for i := range stopChs {
        stopChs[i] = make(chan struct{})
    }
    // wg,rand
    wg := sync.WaitGroup{}
    rand.Seed(time.Now().UnixMilli())

    // 二，模拟querier查询，每个查询持续不同的时间
    wg.Add(QuerierNum)
    for i := 0; i < QuerierNum; i++ {
        // 每一个 querier
        go func(i int) {
            defer wg.Done()
            // 模拟执行时间
            randD := rand.Intn(1000)
            println("querier ", i, " start fetch data, need duration is ", randD, " ms.")
            // 查询结果的channel
            chRst := make(chan Rows, 1)

            // 执行查询工作
            go func() {
                // 模拟时长
                time.Sleep(time.Duration(randD) * time.Millisecond)
                chRst <- Rows{
                    Index: i,
                }
            }()

            // 监听查询结果和停止信号channel
            select {
            // 查询结果
            case rows := <-chRst:
                println("querier ", i, " get result.")
                // 保证没有其他结果写入，才写入结果
                if len(ch) == 0 {
                    ch <- rows
                }
            // stop信号
            case <-stopChs[i]:
                println("querier ", i, " is stopping.")
                return
            }

        }(i)
    }

    // 三，等待第一个查询结果的反馈
    wg.Add(1)
    go func() {
        defer wg.Done()
        // 等待ch中传递的结果
        select {
        // 等待第一个查询结果
        case rows := <-ch:
            println("get first result from ", rows.Index, ". stop other querier.")
            // 循环结构，全部通知querier结束
            for i := range stopChs {
                // 当前返回结果的goroutine不需要了，因为已经结束
                if i == rows.Index {
                    continue
                }
                stopChs[i] <- struct{}{}
            }

        // 计划一个超时时间
        case <-time.After(5 * time.Second):
            println("all querier timeout.")
            // 循环结构，全部通知querier结束
            for i := range stopChs {
                stopChs[i] <- struct{}{}
            }
        }
    }()

    wg.Wait()
}
```

其中核心点：

- 获取了结果，通知结束
- 通过多个无缓冲channel通知goroutine结束
- 通过缓冲channel传递结果

执行结果示例：

```
querier  2  start fetch data, Need duration is  674  ms.
querier  6  start fetch data, Need duration is  695  ms.
querier  1  start fetch data, Need duration is  484  ms.
querier  4  start fetch data, Need duration is  544  ms.
querier  0  start fetch data, Need duration is  101  ms.
querier  7  start fetch data, Need duration is  233  ms.
querier  5  start fetch data, Need duration is  721  ms.
querier  3  start fetch data, Need duration is  727  ms.
querier  0  get result.
get first result from  0 . stop other querier.
querier  7  is stopping.
querier  2  is stopping.
querier  4  is stopping.
querier  6  is stopping.
querier  5  is stopping.
querier  1  is stopping.
querier  3  is stopping.
```

### All 模式

![image.png](https://fynotefile.oss-cn-zhangjiakou.aliyuncs.com/fynote/fyfile/13080/1672136903072/68cfba1c99534b52bc7fb2ace50d42db.png)

Race模式是多个Goroutine获取相同的结果，优先使用快速响应的。

而All模式是多个Goroutine分别获取结果的各个部分，全部获取完毕后，组合成完整的数据，要保证全部的Goroutine都响应后，继续执行。

示例代码，核心逻辑：

- 一个整体内容Content，分为三个goroutine分别处理subject、tags、views三个部分
- 3个goroutine要全部执行完毕，数据才会整体获取
- 不会一直等待，设置超时时间。

本例中，使用具体的每个goroutine的标识方案来识别goroutine。对比Race方案使用的是索引号的方案来识别goroutine。

判定是否全部结束的方案，也是基于具体的标志key。

其中某次执行结果为：

```
start fetch  tags  data, need duration is  396  ms.
start fetch  views  data, need duration is  693  ms.
start fetch  subject  data, need duration is  597  ms.
querier  tags  get result.
received some part  tags
querier timeout. Content is incomplete.
querier  subject  is stopping.
querier  views  is stopping.
received content  { [go Goroutine Channel select] 0 }
```

### 无缓冲Channel+关闭作典型同步信号

基于：

- 无缓冲Channel是同步的
- closed 的channel是可以接收内容的

以上两点原因，经常使用关闭无缓冲channel的方案来作为信号传递使用。前提是，信号纯粹是信号，没有其他含义，比如关闭时间等。

示例代码：

```go
func SelectChannelCloseSignal() {
	wg := sync.WaitGroup{}
	// 定义无缓冲channel
	// 作为一个终止信号使用（啥功能的信号都可以，信号本身不分功能）
	ch := make(chan struct{})

	// goroutine，用来close, 表示
发出信号
	wg.Add(1)
	go func() {
		defer wg.Done()
		time.Sleep(2 * time.Second)
		fmt.Println("发出信号, close(ch)")
		close(ch)
	}()

	// goroutine，接收ch，表示接收信号
	wg.Add(1)
	go func() {
		defer wg.Done()
		// 先正常处理，等待ch的信号到来
		for {
			select {
			case <-ch:
				fmt.Println("收到信号, <-ch")
				return
			default:

			}
			// 正常的业务逻辑
			fmt.Println("业务逻辑处理中....")
			time.Sleep(300 * time.Millise
cond)
		}
	}()

	wg.Wait()
}

// ====
> go test -run TestSelectChannelCloseSignal
业务逻辑处理中....
业务逻辑处理中....
业务逻辑处理中....
业务逻辑处理中....
业务逻辑处理中....
业务逻辑处理中....
业务逻辑处理中....
发出信号, close(ch)
收到信号, <-ch
PASS
ok      goConcurrency   2.168s
```

### signal.Notify 信号通知监控

系统信号也是通过channel与应用程序交互，例如典型的 ctrl+c 中断程序， `os.Interrupt`，若不监控系统信号，ctrl+c后程序会直接终止，而如果监控了信号，那么可以在ctrl+c后，执行一系列的关闭处理，例如：

```go
func SelectSignal() {
    // 一：模拟一段长时间运行的goroutine
    go func() {
        for {
            fmt.Println(time.Now().Format(".15.04.05.000"))
            time.Sleep(300 * time.Millisecond)
        }
    }()

    // 要求主goroutine等待上面的goroutine，方案：
    // 1. wg.Wait()
    // 2. time.Sleep()
    // 3. select{}

    // 持久阻塞
    //select {}

    // 二，监控系统的中断信号,interrupt
    // 1 创建channel，用于传递信号
    chSignal := make(chan os.Signal, 1)
    // 2 设置该channel可以监控哪些信号
    signal.Notify(chSignal, os.Interrupt)
    //signal.Notify(chSignal, os.Interrupt, os.Kill)
    //signal.Notify(chSignal) // 全部类型的信号都可以使用该channel
    // 3 监控channel
    select {
    case <-chSignal:
        fmt.Println("received os signal: Interrupt")
    }
}
```

## 定时器与断续器，Timer&Ticker

> Timer&Ticker是Go标准包time中定义的类型，通过Channel与程序进行通信。

time包中两个与Channel紧密关联的结构：

```go
// 定时器
time.Timer
// 断续器
time.Ticker
```

- 定时器Timer类似于一次性闹钟
- 断续器Ticker类似于重复性闹钟，也成循环定时器

无论是一次性还是重复性计时器，都是通过Channel与应用程序交互的。我们通过监控Timer和Ticker返回的Channel，来确定是否到时的需求。

### 定时器

![image.png](https://fynotefile.oss-cn-zhangjiakou.aliyuncs.com/fynote/fyfile/13080/1672136903072/164a606a60c841129b548fb05f074df0.png)

使用语法：

```go
// time.NewTimer
func NewTimer(d Duration) *Timer
```

创建定时器。参数是Duration时间。返回为 `*Timer`。`*Timer.C`是用来接收到期通知的单向Channel。

```go
type Timer struct {
    C <-chan Time
}
```

因此我们只要可从 `*Timer.C`上接收数据，就意味着定时器时间到。接收到的元素是 `time.Time` 类型数据，为到时时间。

示例：

```go
func TimerA() {
    t := time.NewTimer(time.Second)
    println("Set the timer, \ttime is ", time.Now().String())

    now := <-t.C
    println("The time is up, time is ", now.String())
}
```

Timer除了C之外，还有两个方法：

```go
// 停止计时器
// 返回值bool类型，返回false，表示该定时器早已经停止，返回true表示由本次调用停止
func (t *Timer) Stop() bool

// 重置定时器
// 返回值bool类型，返回false，表示该定时器早已经停止，返回true表示由本次调用重置
func (t *Timer) Reset(d Duration) bool
```

使用这两个方法，可以完整定时器的业务逻辑。

示例代码，简单的猜数字游戏，共猜5次，每次有超时时间3秒钟：

```go
func TimerB() {
    ch := make(chan int)

    // 写channel
    go func() {
        // 随机写入int
        for {
            ch <- rand.Intn(10)
            time.Sleep(400 * time.Millisecond)
        }
    }()

    // 每局时间
    t := time.NewTimer(time.Second * 3)
    hint, miss := 0, 0
    // 统计结果，共玩5次
    for i := 0; i < 5; i++ {
    guess:
        for {
            select {
            case v := <-ch:
                println("guess value: ", v)
                if v == 4 {
                    println("Bingo! some one hint the answer.")
                    // 新游戏，重置定时器
                    t.Reset(time.Second * 3)
                    hint++
                    break guess
                }
            case <-t.C:
                println("The time is up, no one hint.")
                miss++
                // 重新创建定时器
                t = time.NewTimer(time.Second * 3)
                break guess
            }
        }
    }
    println("Game Over! Hint ", hint, ", Miss ", miss)
}
```

代码在猜中或者时间到时，要重置或新建定时器。

如果不需要定时器的关闭和重置操作，可以使用函数：

```go
func After(d Duration) <-chan Time
```

直接返回定时器到期的通知Channel。

```go
func TimerC() {
    ch := time.After(time.Second)
    println("Set the timer, \ttime is ", time.Now().String())

    now := <-ch
    println("The time is up, time is ", now.String())
}
```

如果希望在定时器到期时执行特定函数，可以使用如下函数：

```go
func AfterFunc(d Duration, f func()) *Timer
```

该函数返回*Timer用于控制定时器，例如Stop或Reset.

### 断续器

![image.png](https://fynotefile.oss-cn-zhangjiakou.aliyuncs.com/fynote/fyfile/13080/1672136903072/d7601031b918471a915cdef202bf1a54.png)

也叫循环定时器。

使用语法：

```go
func NewTicker(d Duration) *Ticker
```

创建断续器。参数是Duration时间。返回为 `*Ticker`。`*Ticker.C`是用来接收到期通知的单向Channel。

```go
type Ticker struct {
    C <-chan Time // The channel on which the ticks are delivered.
}
```

因此我们只要可从 `*Ticker.C`上接收数据，就意味着断续器时间到。接收到的元素是 `time.Time` 类型数据，为到时时间。当接收到到期时间后，间隔下一个Duration还会再次接收到到期时间。

`*Ticker`也有方法：

```go
// 停止断续器
func (t *Ticker) Stop()
// 重置断续器间隔时间
func (t *Ticker) Reset(d Duration)
```

示例：

```go
func TickerA() {
    // 断续器
    ticker := time.NewTicker(time.Second)

    // 定时器
    timer := time.After(5 * time.Second)
loop: // 持续心跳
    for now := range ticker.C {
        println("now is ", now.String())
        // heart beat
        println("http.Get(\"/ping\")")

        // 非阻塞读timer，到时结束断续器
        select {
        case <-timer:
            ticker.Stop()
            break loop
        default:
        }
    }
}
```

代码模拟了一个心跳程序，间隔1秒，发送ping操作。整体到时，运行结束。

## 小结

Channel的分类

- nil channel
- 缓冲Channel
- 无缓冲Channel
- 单向Channel

Channel的操作

- 初始化，make(channel type[, cap])
- 发送，ch <- expression
- 接收, v, ok := <- ch
- 遍历接收，for e := range ch {}
- 关闭， close(ch)

select语句

- channel的多路复用
- 执行第一个可以操作channel的case
- 若同时多个channel可操作随机选择case避免饥饿case的出现
- 增加default case可以达到非阻塞channel操作的目的
- 经常配合for select使用循环多路监听
- 典型的多路模式有：Race和All

timer和ticker

- 定时器，到时执行一次，可以在到时前，重置或提前结束
- 断续器，配置间隔重复执行，重复定时器，可以重置间隔时间和提前结束
