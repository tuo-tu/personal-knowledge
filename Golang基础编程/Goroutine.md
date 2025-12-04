# goroutine

## 概念

在Go中，每个并发执行的单元称为goroutine。通常称为Go协程。

## go 关键字启动goroutine

go中使用关键字 go 即可启动新的goroutine。

示例代码：

两个函数分别输出奇数和偶数。采用常规调用顺序执行，和采用go并发调用，通过结果了解并发执行：

```go
func GoroutineGo() {
    // 定义输出奇数的函数
    printOdd := func() {
        for i := 1; i <= 10; i += 2 {
            fmt.Println(i)
            time.Sleep(100 * time.Millisecond)
        }
    }

    // 定义输出偶数的函数
    printEven := func() {
        for i := 2; i <= 10; i += 2 {
            fmt.Println(i)
            time.Sleep(100 * time.Millisecond)
        }
    }

    // 顺序调用
    //printOdd()
    //printEven()

    // 在 main goroutine 中，开启新的goroutine
    // 并发调用
    go printOdd()
    go printEven()

    // 典型的go
    //go func() {}()
    //func() {}()

    // main goroutine 运行结束
    // 内部调用的goroutine也就结束
    time.Sleep(time.Second)
}

// 测试时，需要定义对应的测试文件，例如goroutine_test.go
// 增加单元测试函数：
//file:goroutine_test.go
//package goConcurrency
//
//import "testing"
//
//func TestGoroutineGo(t *testing.T) {
//    GoroutineGo()
//}

// 输出测试结果
goConcurrency> go test -run TestGoroutineGo
1
2
4
3
5
6
8
7
9
10
PASS
ok      goConcurrency   1.052s
```

注意：`time.Sleep(time.Second)`的目的是主goroutine要等待内部goroutine运行结束才能结束，否则主goroutine结束了，内部调用的goroutine也会随之结束。因此我们简单的增加了time.Sleep的方式，进行等待，当然还有很多其他办法。

执行流程如图所示：

![image.png](https://fynotefile.oss-cn-zhangjiakou.aliyuncs.com/fynote/fyfile/13080/1672136903072/1e2f358481964afc94ac275d21b1e8f2.png)

## 使用sync.WaitGroup实现协同调度

WaitGroup用于等待一组goroutine完成。等待思路是计数器方案：

- 调用等待goroutine时，调用Add()增加等待的goroutine的数量
- 当具体的goroutine运行结束后，Done()用来减去计数。
- 主goroutine可以使用Wait来阻塞，直到所有goroutine都完成（计数器归零）。

示例代码：

```go
func GoroutineWG() {
    // 1. 初始化 WaitGroup
    wg := sync.WaitGroup{}
    // 定义输出奇数的函数
    printOdd := func() {
        // 3.并发执行结束后，计数器-1
        defer wg.Done()
        for i := 1; i <= 10; i += 2 {
            fmt.Println(i)
            time.Sleep(100 * time.Millisecond)
        }
    }

    // 定义输出偶数的函数
    printEven := func() {
        // 3.并发执行结束后，计数器-1
        defer wg.Done()
        for i := 2; i <= 10; i += 2 {
            fmt.Println(i)
            time.Sleep(100 * time.Millisecond)
        }
    }
    // 在 main goroutine 中，开启新的goroutine
    // 并发调用
    // 2, 累加WG的计数器
    wg.Add(2)
    go printOdd()
    go printEven()

    // main goroutine 运行结束
    // 内部调用的goroutine也就结束
    // 4. 主goroutine等待
    wg.Wait()
    fmt.Println("after main wait")
}
```

**WaitGroup() 适用于主goroutine需要等待其他goroutine全部运行结束后，才结束的情况。不适用于，主goroutine需要结束，而通知其他goroutine结束的情景**。

如图：

![image.png](https://fynotefile.oss-cn-zhangjiakou.aliyuncs.com/fynote/fyfile/13080/1672136903072/4f5508ffdc3c43178925e8c36ed7f55a.png)

注意，不得复制WaitGroup。因为内部维护的计数器不能被意外修改。

可以同时存在多个goroutine进行等待。

### WaitGroup的基本实现原理

WaitGroup 结构：

```go
type WaitGroup struct {
    // 用于保证WaitGroup不会被复制
    noCopy noCopy
    // 当前状态，存储计数器，存储等待的goroutine
    state1 uint64
    state2 uint32
}
```

 `state` 的高 32 位和低 32 位可能被拆分用于存储不同的信息，状态 32bit和64bit的计算机不同，以64bit为例：

- 高32 bits是计数器的值（`delta` 的累加）
- 低32 bits是等待者（等待中的 goroutine 数量）

Add() 和 Done() 是用来操作计数器，操作计数器的操作是原子操作，保证并发安全性。

Wait()操作，在计数器为0时，结束阻塞状态。

核心代码示例：

```go
func (wg *WaitGroup) Add(delta int) {
    // 原子操作，累加计数器，用于对无符号 64 位整数进行线程安全的加法。
    state := atomic.AddUint64(statep, uint64(delta)<<32)
}

// 相当于加（-1）
func (wg *WaitGroup) Done() {
    wg.Add(-1)
}


// wait函数用一个循环结构来实现，直到计数器为0
func (wg *WaitGroup) Wait() {

    for {
        state := atomic.LoadUint64(statep)
        // 提取 state 的高 32 位(右移相当于丢掉数据的低位)
        v := int32(state >> 32)
        w := uint32(state)
        // 如果计数器为0，则不需要等待
        if v == 0 {
            // Counter is 0, no need to wait.
            if race.Enabled {
                race.Enable()
                race.Acquire(unsafe.Pointer(wg))
            }
            return
        }
        // Increment waiters count.
}
```

## 调度的随机性

我们不能期望goroutine的执行顺序依照源代码的顺序执行，如下代码，会随机输出0-9：

```go
func GoroutineRandom() {
    wg := sync.WaitGroup{}
    wg.Add(10)
    for i := 0; i < 10; i++ {
        go func(n int) {
            defer wg.Done()
            fmt.Println(n)
        }(i)
    }
    wg.Wait()
}
```

## goroutine的并发规模

> Goroutine 的并发数量有上限吗？
>
> - 受goroutine占用的栈内存限制
> - 受内部操作资源限制
> - goroutine本身无上限

函数 `runtime.NumGoroutine()` 可以获取当前存在的Goroutine数量。

示例，大量执行耗时的gorutine，并统计goroutine的数量：

```go
func GoroutineNum() {
    // 1. 统计当前存在的goroutine的数量
    go func() {
        for {
            fmt.Println("NumGoroutine:", runtime.NumGoroutine())
            time.Sleep(500 * time.Millisecond)
        }
    }()

    // 2. 启动大量的goroutine
    for {
        go func() {
            time.Sleep(100 * time.Second)
        }()
    }

}
```

内存溢出的运行错误

```shell
NumGoroutine: 3000000

runtime: VirtualAlloc of 32768 bytes failed with errno=1455
fatal error: out of memory
```

### goroutine最小为2KB

之所以支持百万级的goroutine并发，核心原因是因为每个goroutine的**初始栈内存为2KB**，用于保持goroutine中的执行数据，例如局部变量等。相对来说，线程的栈内存通常为2MB。除了比较小的初始栈内存外，**goroutine的栈内存可扩容的**，也就是说支持按需增大或缩小，**一个goroutine最大的栈内存当前限制为1GB。**

### goroutine内部资源竞争溢出

在goroutine内增加，fmt.Println() 测试：

```shell
panic: too many concurrent operations on a single file or socket (max 1048575)
```

### 控制并发数量

实际开发时，要根据系统资源和每个goroutine锁消耗的资源来控制并发规模。

典型方案 goroutine pool，典型的包：

1. Jeffail/tunny：性能一般，轻量级的，没有专门的优化机制。通过复用 goroutine，可以减少频繁创建和销毁 goroutine 的开销。
2. panjf2000/ants：性能领先，支持**任务内存复用**，减少了垃圾回收的开销。
3. go-playground/pool：性能中规中矩，适合对性能要求不高的小型项目。

| 功能特性       | **Jeffail/tunny**         | **panjf2000/ants**             | **go-playground/pool** |
| -------------- | ------------------------- | ------------------------------ | ---------------------- |
| **任务类型**   | 任意类型（`interface{}`） | 固定任务函数或支持任务参数传递 | 固定任务函数           |
| **池大小调整** | 固定池大小                | 动态调整池大小                 | 固定池大小             |
| **任务优先级** | 不支持                    | 支持延迟任务                   | 不支持                 |
| **内存复用**   | 不支持                    | **支持**                       | 不支持                 |
| **性能优化**   | 一般                      | **强优化（高吞吐、低延迟）**   | 一般                   |
| **任务结果**   | 支持返回任务处理结果      | 不直接支持（需自行处理）       | 不支持任务结果         |

以 ants 为例，展示goroutine池的使用：

```go
$ go get -u github.com/panjf2000/ants/v2

func GoroutineAnts() {
    // 1. 统计当前存在的goroutine的数量
    go func() {
        for {
            fmt.Println("NumGoroutine:", runtime.NumGoroutine())
            time.Sleep(500 * time.Millisecond)
        }
    }()

    // 2. 初始化协程池，goroutine pool
    size := 1024
    pool, err := ants.NewPool(size)
    if err != nil {
        log.Fatalln(err)
    }
    // 保证pool被关闭
    defer pool.Release()

    // 3. 利用 pool，调度需要并发的大量goroutine
    for {
        // 向pool中提交一个执行的goroutine
        err := pool.Submit(func() {
            time.Sleep(100 * time.Second)
        })
        if err != nil {
            log.Fatalln(err)
        }
    }
}

// ======
> go test -run TestGoroutineAnts
runtime.NumGoroutine():  8
runtime.NumGoroutine():  1031
runtime.NumGoroutine():  1031
runtime.NumGoroutine():  1031
runtime.NumGoroutine():  1031
```

其中：

- ants.NewPool() 创建池
- pool.Release() 释放池
- pool.Submit() 提交goroutine操作

除了使用特定包之外，还可自定义计数器等方案实现，请参考Channel的示例：《使用channel控制并发数量示例》

## 并发调度

### 多对多的协程调度

Goroutine 的概念来自于协程 Coroutine，协程，又称微线程，纤程。通过将多个协程映射到特定线程，提高程序（函数）的并发执行能力。

在典型的语言中，函数（子程序）的调用都是通过栈（先进后出）实现的，通常都是层级调用。例如：

```
- A
  - B
      - C
      - D
  - E
  - F
```

以上调用，A函数分别调用了B，E，F，而B调用了C，D。通常来说A必须要等到B，C，D执行完才能返回；E要等到B执行完才继续执行，而B要等到C，D执行完才能返回。整体来说就是一个函数一旦执行，不能被打断去执行别的函数。

而协程Coroutine的设计就是在某个函数执行的过程中，可以主动（Python的yield）和被动（go 的goroutine）的被终止执行，转而去执行其他函数。也就是，上例子中，A调用了B、E、F，可以做到，执行一会B，再去执行E、E没有执行完毕，又暂停去执行B，或F。直到全部执行完毕。

以单线程为例，多个协程通过应用程序自己维护的协程调度器完成多个协程的调度：

<img src="https://fynotefile.oss-cn-zhangjiakou.aliyuncs.com/fynote/fyfile/13080/1672136903072/0dca3eb676b04825b0ef480379915ce0.png" alt="image.png" style="zoom: 67%;" />

**当前的计算机都支持多线程**，因此现代语言实现的协程调度器通常都是多对多的调度方案：**指的是多协程对多线程。**

<img src="https://fynotefile.oss-cn-zhangjiakou.aliyuncs.com/fynote/fyfile/13080/1672136903072/94443c4e3e5749eda7c594ca3f5a69d6.png" alt="image.png" style="zoom:67%;" />

### GMP模型结构

Goroutine 就是Go语言实现的协程模型。其核心结构有三个，称为GMP，也叫GMP模型。分别是：

- **G**，Goroutine，**用户级线程**，是Go 程序中的最小执行单元，使用关键字go创建。**G会存储在P的本地队列或者是全局队列中。**特别注意，`G0`是M 的调度专用 Goroutine（仅用于执行调度逻辑）。
- **M**，Machine，**内核线程**，与操作系统的线程一一对应，就是Work Thread，就是传统意义的线程（系统线程），**用于执行Goroutine**。M只有在与具体的P绑定后，才能执行P中的G。
- **P**，Processor，**处理器**，或称逻辑处理器，主要**用于协调G和M之间的关系**，每个 P 都有自己的任务队列（本地队列），存储待执行的G队列，P与特定的M绑定后，M执行P队列中的G。P 的数量由 `GOMAXPROCS` 控制，默认值为 CPU 核心数。调度器按照 FIFO（先进先出）策略从本地队列取任务。

GMP整体结构逻辑图：

<img src="https://fynotefile.oss-cn-zhangjiakou.aliyuncs.com/fynote/fyfile/13080/1672136903072/ec55e59f84304b9787e3c2706868b318.png" alt="image.png" style="zoom: 50%;" />

图例说明：

- M程序线程，由OS负责调度交由具体的CPU核心中执行
- 待执行的G可能存储于全局队列，和某个P的本地队列中。**P的本地队列当前最多存储256个G。**
- 若要执行P中的G，则P必须要于对应的M建立联系。建立联系后，就可以执行P中的G了。
- M与P不是强绑定的关系，**若一个M阻塞，那么P就会选择一个新的M执行后续的G**。该过程由Go调度器调度。

GMP的关系

- G 是独立的运行单元
- M是执行任务的线程
- P是G和M的关联纽带
- G要在M中执行，P的任务就是合理的将G分配给M

### P的数量

P的数量通常是固定的，当程序启动时**由 `$GOMAXPROCS`环境变量决定创建P的数量。默认的值为当前CPU的核数，所有的 P 都在程序启动时创建。**这意味着程序的执行过程中，最多同时有$GOMAXPROCS个Goroutine同时运行，默认与CPU核数保持一致，可以最大程度利用多核CPU并行执行的能力。

程序运行时，`runtime.GOMAXPROCS()`函数可以动态改变P的数量，但通常不建议修改，或者即使修改也不建议数量超过CPU的核数。调动该函数的典型场景是控制程序的并行规模，例如：

```go
// 最多利用一半的CPU
runtime.GOMAXPROCS(runtime.NumCPU() / 2)

// 获取当前CPU的核数
runtime.NumCPU()
```

我们知道Go没有限定G的数量，那M的数量呢？

- Go对M的数量做了一个上限，10000个，但通常不会到达这个规模，因为操作系统很难支持这么多的线程。
- M的数量是由P决定的。
- 当P需要执行时，会去找可用的M，若没有可用的M，就会创建新的M，这就会导致M的数量不断增加
- 当M线程长时间空闲，也就是长时间没有新任务时，GC会将线程资源回收，这会导致M的数量减少
- 整体来说，M的数量一定会多于P的数量，取决于空闲（没有G可执行的）的，和完成其他任务（例如CGO操作，GC操作等）的M的数量多少。

### P与G关联的流程

go 创建的待执行的Goroutine与P建立关联的核心流程：

- 新创建的G会优先保持在P的本地队列中。例如A函数中执行了 go B()，那么B这个Goroutine会优先保存在A所属的P的本地队列中。
- 若G加入P的本地队列时本地队列已满，那么G会被加入到全局G队列中。**新G加入全局队列时，会把P本地队列中一半的G也同时移动到全局队列中（是乱序入队列）**，以保证P的本地队列可以继续加入新的G。
- 当P要执行G时
  - 会从P的本地队列查找G。
  - 若本地队列中没有G，则会尝试从其他的P中偷取（Steal）G来执行，**通常会偷取一半的G。**
  - 若无法从其他的P中偷取G，则**从全局G队列中获取G**，会一次获取多个G。
  - 整体：**本地G队列->其他P的本地G队列->全局G队列**
- 当全局运行队列中有待执行的 G 时，还会有固定几率（每61个调度时钟周期 schedtick）会从全局的运行队列中查找对应的 G，**为了保证全局G队列一定可以被调度。**

核心流程图例：

A 中调用了 go B()， P的本地队列未满时：

<img src="https://fynotefile.oss-cn-zhangjiakou.aliyuncs.com/fynote/fyfile/13080/1672136903072/d2296ed925264f3cbc624d27de06637e.png" alt="image.png" style="zoom:50%;" />

A 中调用了 go B()， P的本地队列已满时：

<img src="https://fynotefile.oss-cn-zhangjiakou.aliyuncs.com/fynote/fyfile/13080/1672136903072/de05e69f33664d2295fd528eea5e2d47.png" alt="image.png" style="zoom: 50%;" />



当P要执行G时：

<img src="https://fynotefile.oss-cn-zhangjiakou.aliyuncs.com/fynote/fyfile/13080/1672136903072/615654aa1f5f45ff89c4d7cd90835b2b.png" alt="image.png" style="zoom:50%;" />

### P与M关联的流程

P中关联了大量待执行的G，若需要执行G，P要去找可用的M。**P不能执行G，只有M才能真正执行G。**

P与M建立关联的核心过程：

- 当P需要执行时，P要寻找可用的M，**优先从空闲M池中找M，**若没有空闲的，则新建M来执行
- 在创建G时，**G会尝试唤醒空闲的M**
- 当**G进行了系统调用时，M会阻塞，并释放与之绑定的P**，把P转移给其他的M去执行。称为P抢占。
- 当G完成系统调用后，M也不阻塞了，M会尝试获取新的空闲P，同时将G放入P的本地队列执行。**若没有空闲的P，则将G放入全局G队列，M进入休眠，等待被唤醒或被垃圾回收**

如图所示：

P要寻找可用的M：

<img src="https://fynotefile.oss-cn-zhangjiakou.aliyuncs.com/fynote/fyfile/13080/1672136903072/1094a65c5a534e7a8ff9118597cd6fbf.png" alt="image.png" style="zoom: 67%;" />

G执行了系统调用，M与P解绑（释放），P转移到新的M上执行：

<img src="https://fynotefile.oss-cn-zhangjiakou.aliyuncs.com/fynote/fyfile/13080/1672136903072/e71e6b03e1414dcd99b538633ebec8d4.png" alt="image.png" style="zoom:50%;" />

M1执行完G的系统调用后，G不一定结束，还要继续执行，则M1会尝试获取空闲的P（没有与M绑定的P），若没有空闲的P可用，将M1执行的G放入全局G队列，M1进入空闲状态，进入空闲M池：

![image.png](https://fynotefile.oss-cn-zhangjiakou.aliyuncs.com/fynote/fyfile/13080/1672136903072/7402e28b89194a23be7814f2b766a86e.png)

### G的调度流程总结

整体上看

1. go func() 创建Goroutine
2. 将Goroutine放入队列

- 放入本地队列
- 本地队列满，放入全局队列

3. M通过P获取G并执行

- 从本地队列获取G
- 从其他P的本地队列获取G
- 从全局队列获取G

4. M执行G

- 调度周期循环执行G
- G主动让出
- G执行系统调用

5. G执行系统调用

- 解绑G和P（M和P吧？）
- P抢占其他的M继续执行
- 系统调用的G结束，将G放入其他P队列执行，M空闲

6. 若G执行完毕，释放

### M0 和 G0

- M0, 启动程序后的编号为 0 的主线程，负责执行初始化操作和启动第一个G，也就是 main Goroutine。之后与其他M一样调度。
- G0，每个 M 创建的第一个 Goroutine。**G0 是M 的调度专用协程，仅用于负责调度G **，G0 不指向任何可执行的函数，**每个 M 都会有一个自己的 G0**。在调度或系统调用时会使用 G0 的栈空间。

如图：

M的G0：

![image.png](https://fynotefile.oss-cn-zhangjiakou.aliyuncs.com/fynote/fyfile/13080/1672136903072/98e97bf2f1d5401189882a781949f27e.png)

### 协作和抢占调度

当某个 **G 执行时间过长**，其他的 G 如何调度。通常有两种方案：

- 协作式，主动让出执行权，让其他G执行。通过runtime.Gosched()可以让出执行权。
- 抢占式，被动让出执行权，也就是**调度器将G的执行权取消**，分配给其他的G。抢占式是Go目前默认的方式。在Go中一个Goroutine**最多可以执行10ms**，超时就会被让出调度权。

函数：

```go
runtime.Gosched()
```

此方法可以要求Go主动调度该goroutine，去执行其他的goroutine。这种是典型的协作调度模式，类似于 py 的 yield。

示例：

```go
func GoroutineSched() {
    runtime.GOMAXPROCS(1)
    wg := sync.WaitGroup{}
    wg.Add(1)
    max := 100
    go func() {
        defer wg.Done()
        for i := 1; i <= max; i += 2 {
            fmt.Print(i, " ")
            runtime.Gosched()
            //time.Sleep(time.Millisecond)
        }
    }()

    wg.Add(1)
    go func() {
        defer wg.Done()
        for i := 2; i <= max; i += 2 {
            fmt.Print(i, " ")
            runtime.Gosched()
            //time.Sleep(time.Millisecond)
        }
    }()

    wg.Wait()
}
```

我们采**用1个P来进行模拟**，看主动让出交替执行的情况。

上面代码中，若goroutine中，没有runtime.GoSched，则会先执行完一个，再执行另一个。若存在runtime.GoSched，则会交替执行。这就是协作式。

除此之外，增加sleep时间1ms，不增加runtime.GoSched调用，也会出现交替执行的情况，这种情况就是调度器主动调度Goroutine了，是抢占式。

## 小结

- Goroutine：Go语言中实现的协程。
- go 关键字：使用go 关键字调用函数，可以让函数独立运行在Goroutine中。
- main 函数也是运行在Goroutine中
- 通常 main 函数需要等待其他Goroutine运行结束
- 典型的并发等待使用 sync.WaitGroup 类型。
- 并发Goroutine的调度在应用层面可以认为是随机的
- 支持海量gouroutine的特点：

  - goroutine语法层面没有限制，但使用时通常要限制，避免并发的goroutine过多，资源占用过大
  - 更小的goroutine栈内存
  - 强大的GMP调度
- GMP

  - G，Goroutine，用户级线程，独立并发执行的代码段
  - M，Machine, 系统线程，用于执行G
  - P，Processor，逻辑处理器，用于协调G和M。
  - G存在与P的本地队列或全局队列中
  - P要与M绑定，P中的G才会被执行
  - M执行G中的系统调用时，会解绑M和P，P会找到新的M执行

# Channel通信

## Channel概述

> **不要通过共享内存的方式进行通信，而是应该通过通信的方式共享内存**

这是Go语言最核心的设计模式之一。

在很多主流的编程语言中，多个线程传递数据的方式一般都是共享内存，而**Go语言中多Goroutine通信的主要方案是Channel**。Go语言也可以使用共享内存的方式支持Goroutine通信。

<img src="https://fynotefile.oss-cn-zhangjiakou.aliyuncs.com/fynote/fyfile/13080/1672136903072/cceb79db926449819974f11cae5a8d0f.png" alt="image.png" style="zoom:80%;" />

Go语言实现了**CSP**通信模式，CSP是Communicating Sequential Processes的缩写，**通信顺序进程**。Goroutine和Channel分别对应CSP中的实体和传递信息的媒介。CSP是Tony Hoare于1977年提出。

Channel提供可接收和发送特定类型值的用于并发函数(Goroutine)通信的数据类型，是满足FIFO（先进先出）原则的队列类型，先进先出不仅体现在数据类型上，也体现在操作上：

- **channel类型的元素是先进先出的**，先发送到channel的value会先被receive
- 先向Channel发送数据的Goroutine会先执行
- 先从Channel接收数据的Goroutine会先执行

如图：

<img src="https://fynotefile.oss-cn-zhangjiakou.aliyuncs.com/fynote/fyfile/13080/1672136903072/6a24b9d6e254423dbfb04b186c9fa1d3.png" alt="image.png" style="zoom:80%;" />

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

关闭Channel的意思是记录该Channel**不能再被发送任何元素了，而不是销毁该Channel**的意思。也就意味着**关闭的Channel是可以继续接收值的**。因此：

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

**也称为同步Channel**，只有当发送方和接收方都准备就绪时，通信才会成功。

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

**同步Channel适合在goroutine间做同步信号！**

### 缓冲Channel

![image.png](https://fynotefile.oss-cn-zhangjiakou.aliyuncs.com/fynote/fyfile/13080/1672136903072/91c9919010274c759c3f1120ac32c208.png)

**缓冲Channel也称为异步Channel**，接收和发送方不用等待双方就绪即可成功。缓冲Channel会存在一个容量为cap的缓冲空间。当使用缓冲Channel通信时，接收和发送操作是在操作Channel的Buffer：

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

**缓冲channel非常适合做goroutine的数据通信了。**

### 长度和容量，len()和cap()

内置函数 len() 和 cap() 可以分别获取：

- len()长度，**缓冲中元素个数**。
- cap()容量，**缓冲的总大小。cap()返回0，意味着是无缓冲通道**

### 示例：使用channel控制并发数量

核心思路是使用缓冲channel作为goroutine计数器来使用：

```go
func ChannelGoroutineNumCtl() {
    // 1 独立的goroutine输出goroutine数量
    go func() {
        for {
            fmt.Println("NumGoroutine:", runtime.NumGoroutine())
            time.Sleep(500 * time.Millisecond)
        }
    }()

    // 2 初始化channel，设置缓冲大小（并发规模）
    const size = 1024
    ch := make(chan struct{}, size)

    // 3 并发的goroutine
    for {
        // 一，启动goroutine前，执行 ch send
        // 当ch的缓冲已满时，阻塞
        ch <- struct{}{}
        go func() {
            time.Sleep(10 * time.Second)
            // 二，goroutine结束时，接收一个ch中的元素
            <-ch
        }()
    }
}

// ======
> go test -run TestChannelGoroutineNumCtl
NumGoroutine: 7
NumGoroutine: 1029
NumGoroutine: 1029
NumGoroutine: 1029
NumGoroutine: 1029
NumGoroutine: 1029
NumGoroutine: 1029
NumGoroutine: 1029
NumGoroutine: 1029
NumGoroutine: 1029
```

其中：

- 当要开启goroutine时，先执行ch的send操作，若ch的缓冲已满，则阻塞在send操作，不会开启新的goroutine
- 当goroutine执行结束时，从ch中接收一个元素，减少ch的缓冲元素数量

## 单向Channel

![image.png](https://fynotefile.oss-cn-zhangjiakou.aliyuncs.com/fynote/fyfile/13080/1672136903072/27c1c8ed7a6b4bf895c4ea5c92638ed5.png)

单向Channel，指的是仅支持接收或仅支持发送操作的Channel。语法上：

- `chan<- T` 仅发送Channel
- `<-chan T` 仅接收Channel

单向Channel的意义在于约束Channel的使用方式。

**仅使用单向Channel通常没有实际意义**，单向Channel最典型的使用方式是：**使用单向通道约束双向通道的操作。**

**语法上来说，就是我们会将双向Channel转换为单向Channel来使用**。典型使用在函数参数或返回值类型中。双向通道（实参）可以作为参数传递给单向通道（形参），具体来说，**双向通道**可以被隐式地转换为**单向通道**，以限制通道在函数内部的操作权限。

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

// only send channel
func setElement(ch chan<- int, v int, wg *sync.WaitGroup) {
    defer wg.Done()
    ch <- v
    println("send to ch, element is ", v)
}

// only receive channel
func getElement(ch <-chan int, wg *sync.WaitGroup) {
    defer wg.Done()
    println("received from ch, element is ", <-ch)
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
    qcount   uint           // 元素个数len()
    dataqsiz uint           // 缓冲队列的长度cap()
    buf      unsafe.Pointer // 缓冲队列指针，无缓冲队列为nil
    elemsize uint16 // 元素大小
    closed   uint32
    elemtype *_type // 元素类型
    sendx    uint   // send index 发送索引
    recvx    uint   // receive index 接收索引
    recvq    waitq  // list of recv waiters 等待接收队列
    sendq    waitq  // list of send waiters 等待发送队列

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

- channel存储空间分为buf和其他属性两块
- channel（其他属性）上记录channel的属性，长度、容量、元素类型、元素大小，接收/发送索引、接收/发送等待队列
- channel.buf为elemtype类型的array
- 若为**无缓冲channel，不分配channel.buf空间**
- **make()初始化的核心操作就是分配内存空间**

### 缓冲数组

<img src="https://fynotefile.oss-cn-zhangjiakou.aliyuncs.com/fynote/fyfile/13080/1672136903072/3efd58b1945f49a48aad6aae0f42e045.png" alt="image.png" style="zoom:80%;" />

<img src="https://fynotefile.oss-cn-zhangjiakou.aliyuncs.com/fynote/fyfile/13080/1672136903072/4b85a56eccd847848b773df63151f797.png" alt="image.png" style="zoom:80%;" />

buf（即缓冲）为数组结构，channel还记录了**buf的发送和接收元素的索引**：

```go
 sendx  uint  // 发送索引
 recvx  uint  // 接收索引
```

**缓冲数组是循环使用的**，也就是**若数组的最后一个元素存储了元素，那么下一次会尝试存储在第一个元素位置。**

### Channel与Goroutine的关系(感觉还挺重要)

![image.png](https://fynotefile.oss-cn-zhangjiakou.aliyuncs.com/fynote/fyfile/13080/1672136903072/3f0bf05adaa24010a25fd2c0db521c3b.png)

Channel有两个属性，用于**记录等待接收或发送的goroutine队列**：

```go
recvq  waitq  // 等待接收goroutine队列
sendq  waitq  // 等待发送goroutine队列
```

**当基于某channel的接收或发送的goroutine无法执行时，也就是需要阻塞时，会被记录到Channel的等待队列中**。当channel可以完成相应的接收或发送操作时，**从等待队列中唤醒goroutine进行操作。**

其中**等待队列是 runtime.waitq 类型，是一个双向链表结构**，具体的某个链表节点存在两个指针，指向前后节点：

```go
// GOROOT/src/runtime/chan.go
type waitq struct {
    first *sudog
    last  *sudog
}
```

其中***sudog 可以理解为一个挂起的goroutine**（即待唤醒的G）。

### 初始化channel流程

make()初始化channel时，会根据是否存在缓冲，选择：

- 存在缓冲，为channel和buffer分别分配内存，同时channel.buf指向buffer地址
- **不存在缓冲，仅为channel分配内存，channel.buf为nil。**
- 初始化channel中其他属性

![image.png](https://fynotefile.oss-cn-zhangjiakou.aliyuncs.com/fynote/fyfile/13080/1672136903072/f6dc4a8966934f2a9f26d9bd5574ba9b.png)

### 向channel发送流程

语句 ch <- element 向channel发送元素时，大体的执行流程如下：

1. 直接发送：**当channel存在等待接受者（channel.recvq）时，直接将元素拷贝给等待接受者**，并**唤醒等待接受者**goroutine将其放在M的runnext位置，下次调度立即执行，
2. 若没有recvq，并且存在缓冲区，将发送元素直接写入缓冲区（缓冲区存在空间时），调整channel.sendx的位置
3. **当缓冲区已满或无缓冲区时，发送goroutine进入channel.sendq队列，转为阻塞状态**，等待其他goroutine从channel中接收元素，进而唤醒发送goroutine

<img src="https://fynotefile.oss-cn-zhangjiakou.aliyuncs.com/fynote/fyfile/13080/1672136903072/73f96aca61da483692faf36672a20ef8.png" alt="image.png" style="zoom:80%;" />

### 从channel接收流程

操作符 <- ch 从channel中接收元素，大体流程如下：

- 当**存在等待发送者**，channel.sendq
  - **若无缓冲区，直接将元素从发送者拷贝到接受者**，并唤醒发送者gorutine，进入runnext下次调度执行
  - **若有缓冲区**，此时缓冲区是满的，**从缓冲区获取元素**，**并将等待发送者发送元素拷贝到缓冲区**，唤醒发送者goroutine。调整channel的recvx和sendx索引位置
- 当**不存在等待发送者，缓冲区有元素**时，直接从缓冲区读取元素
- 当**不存在等待发送者，缓冲区不存在或缓冲区无元素**时，**接收者goroutine进入阻塞状态**，进入channel.recvq（等待接受者队列），等待发送者发送数据唤醒。

<img src="https://fynotefile.oss-cn-zhangjiakou.aliyuncs.com/fynote/fyfile/13080/1672136903072/89f73cffcd37448cb1085ed6de144b9e.png" alt="image.png" style="zoom:80%;" />

### 关闭channel流程

close(ch)关闭channel，主要工作是：

- 标记 `channel` 为已关闭

- 解除关联队列：将channel关联的sendq和recvq队列统统解除
- 唤醒goroutine：唤醒阻塞在sendq和recvq中的goroutine

## select 语句

`select` 语句能够从**多个可读或者可写的Channel**中选择一个继续执行 ，**若没有Channel发生读写操作，`select` 会一直阻塞当前Goroutine。**

![image.png](https://fynotefile.oss-cn-zhangjiakou.aliyuncs.com/fynote/fyfile/13080/1672136903072/9d4a496c12a1483f8527c00c2c721360.png)

### select语法

```go
SelectStmt = "select" "{" { CommClause } "}" .
CommClause = CommCase ":" StatementList .
CommCase   = "case" ( SendStmt | RecvStmt ) | "default" .
RecvStmt   = [ ExpressionList "=" | IdentifierList ":=" ] RecvExpr .
RecvExpr   = Expression .
```

语法结构与 switch 类似，**但case都要涉及channel操作**，示例：

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

1. **计算channel和值**：在 `select` 语句开始执行时，所有 case 中的 channel 表达式会根据源码顺序计算一次，并不会重复计算。如果有 `send` 操作，还会计算要发送的值。

   - 解释：**所有** `case` 中涉及的 `channel` 表达式（如接收操作的 `channel`，发送操作的 `channel` 以及要发送的值）在进入 `select` 时就会被**一次性**计算完成，而不会在 `select` 的执行过程中重复计算；

   - 注意：对于接收操作`RecvStmt`，此过程中左侧带有短变量声明或赋值的表达式尚未计算，如果带有短变量声明（如 `v := <-ch`），**只有在对应 case 被选择时，才会计算左侧表达式并赋值**。

2. **伪随机选择：**如果一个或多个通信可以执行，`select` 会随机选择其中一个执行，避免因固定选择顺序而产生的偏向问题。如果没有可执行的channel，往下继续。
3. **Default case：**如果存在default case，选择该case（`default` case 是非阻塞的）。如果没有default case，**select 语句会阻塞**直到至少一个通信操作可以被执行。

4. **执行RecvStmt**：如果选择的case是带有短变量声明或赋值的RecvStmt，左侧表达式会被计算，并分配接收的值（或多个值）。

5. 按顺序执行所选择的case的语句列表。

### for + select

select 匹配到可操作的case或者是default case后，就执行完毕了。实操时，我们**通常需要持续监听某些channel的操作**，因此典型的select使用会配合for完成。

例如：**持续从某个ch内获取数据**

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

- **不存在任何case的**（挺有用）
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

go test 测试时，会一直阻塞。若上面的代码出现在常规执行流程中，会导致 deadlock（死锁）。

### nil channel的case

nil channel 不能读写，**因此通过将channel设置为nil，可以控制某个case不再被执行。**

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

特别的情况是**如果存在两个case，其中一个是default，另一个是channel case，那么go的优化器会优化内部这个select。内部会以if结构完成处理**。因为这种情况，不用考虑随机性的问题。类似于：

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

**Goroutine的两种常见的任务处理模式**，**Race模式** 和 **All模式**（用于处理多个并行任务）

### Race模式（选择最快完成任务的结果）

<img src="https://fynotefile.oss-cn-zhangjiakou.aliyuncs.com/fynote/fyfile/13080/1672136903072/fcceb74a358c4bc29468ad3f8a5fea80.png" alt="image.png" style="zoom:80%;" />

Race模式，典型的并发执行模式之一，**多路同时操作资源，哪路先操作成功，优先使用，同时放弃其他路的等待**。简而言之，**从多个操作中选择一个最快的**。多个Goroutine**获取相同的结果**，优先使用快速响应的。核心工作：

- 选择最快的
- 停止其他未完成的

示例代码，示例从多个查询器同时查询数据，使用最先反返回结果的，其他查询器结束：

```go
func SelectRace() {
    // 一，初始化数据
    // 模拟查询返回的结果，包含唯一标识 Index。需要与具体的querier（查询者）建立联系
    type Rows struct {
        Index int
    }
    // 并行的查询者数量
    const QuerierNum = 8
    // 用于通信的channel，数据，停止信号？？
    // 一个带缓冲的通道，用来传递第一个返回的查询结果。
    ch := make(chan Rows, 1)
    // 每个查询任务都有一个stopCh，用于接收停止信号，以序号进行区分
    stopChs := [QuerierNum]chan struct{}{}
    for i := range stopChs {
        stopChs[i] = make(chan struct{})
    }
   
    wg := sync.WaitGroup{}
    rand.Seed(time.Now().UnixMilli())
    // 二，模拟查询，每个查询持续不同的时间
    wg.Add(QuerierNum)
    for i := 0; i < QuerierNum; i++ {
        // 每一次循环表示，每一个查询者，每个查询会随机模拟一个耗时，并尝试获取数据。
        go func(i int) {
            defer wg.Done()
            // 模拟执行时间
            randD := rand.Intn(1000)
            println("querier ", i, " start fetch data, need duration is ", randD, " ms.")
            // 每个查询者有一个独立的 chRst，用于传递模拟的查询结果。
            chRst := make(chan Rows, 1)
            // 执行查询工作，将查询结果放入chRst
            go func() {
                // 模拟查询的随机耗时，
                time.Sleep(time.Duration(randD) * time.Millisecond)
                // 生成一个 Rows 结果并传递到 chRst 中。
                chRst <- Rows{ 
                    Index: i,
                }
            }()

            // 监听自己的查询结果和停止信号channel
            select {
            // chRst虽然有多个，但是在下面这行代码在每个查询者的内部，因此内部的每个select语句可以监听自己独有的 chRst，从而确保只有自己的结果被处理。
            case rows := <-chRst:
                println("querier ", i, " get result.")
                // 保证没有其他结果写入，才写入结果。这里实现了最快的一个查询被写入了ch中
                if len(ch) == 0 {
                    ch <- rows
                }
            // 同时监听stop信号，收到停止信号，直接结束本函数（本goroutine）
            case <-stopChs[i]:
                println("querier ", i, " is stopping.")
                return
            }
        }(i)
    }

    // 三、启动一个goroutine监听第一个查询结果的反馈，或进行超时处理
    wg.Add(1)
    go func() {
        defer wg.Done()
        select {
        // 监听第一个查询结果
        case rows := <-ch:
            println("get first result from ", rows.Index, ". stop other querier.")
            // 循环结构，通知其余的查询者，结束查询
            for i := range stopChs {
                // 当前返回结果的goroutine不需要通知，因为已经结束
                if i == rows.Index {
                    continue
                }
                stopChs[i] <- struct{}{}
            }

        // 设置一个5s定时器，若超时，通知所有查询者停止。
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

### All 模式（也称收集模式，组合所有任务的结果）

<img src="https://fynotefile.oss-cn-zhangjiakou.aliyuncs.com/fynote/fyfile/13080/1672136903072/68cfba1c99534b52bc7fb2ace50d42db.png" alt="image.png" style="zoom:80%;" />

All模式是**多个Goroutine分别获取结果的各个部分**，全部获取完毕后，**组合成完整的数据**，要保证全部的Goroutine都响应后，继续执行。

示例代码，核心逻辑：

- 一个整体内容Content，分为三个goroutine分别处理subject、tags、views三个部分
- 3个goroutine要全部执行完毕，数据才会整体获取
- 不会一直等待，设置超时时间。

本例中，使用具体的每个goroutine的标识方案来识别goroutine。对比Race方案使用的是索引号的方案来识别goroutine。

判定是否全部结束的方案，也是基于具体的标志key。

```go
func SelectAll() {
	// 定义资源类型数据结构
	type Content struct {
		Subject string   // 主题
		Tags    []string // 标签
		Views   int      // 浏览量
		Part    string   // 标识当前数据来自哪个部分
	}

	// 定义三个用于表示不同部分的常量
	const (
		PartSubject = "subject"
		PartTags    = "tags"
		PartViews   = "views"
	)
	// 初始化一个的map结构（包括其中的channel），值表示content各部分的停止信号Channel，此时channel没有值
	stopChs := map[string]chan struct{}{
		PartSubject: make(chan struct{}),
		PartTags:    make(chan struct{}),
		PartViews:   make(chan struct{}),
	}
	// 用于存储从各个部分的 goroutine 收集的结果，缓冲区大小为部分的数量
	ch := make(chan Content, len(stopChs))
	// 初始化超时Channel
	to := time.After(100 * time.Millisecond)
	wg := sync.WaitGroup{}
	// 对content的每个部分（subject、tags、views），启动一个 goroutine 负责获取数据。
	for part := range stopChs { // 遍历stopChs得到key
		wg.Add(1)
		go func(part string) {
			defer wg.Done()
			// 模拟每个数据获取操作所需时间
			randD := rand.Intn(1000)
			fmt.Println("querier", part, ", 开始取数据, 需要持续时间为:", randD, "ms")
			// 初始化每个part的查询结果Channel
			chRst := make(chan Content, 1)

			go func() { // 启动一个goroutine，将获取到的content放入chRst中。
				// 模拟获取资源的执行时间
				time.Sleep(time.Duration(randD) * time.Millisecond)
				context := Content{ // 标识这个循环代表了哪一部分
					Part: part,
				}
				
				// 基于不同的 part，完成不同属性的设置
				switch part {
				case PartSubject:
					context.Subject = "Subject of content"
				case PartTags:
					context.Tags = []string{"Go", "Goroutine", "Channel", "Select"}
				case PartViews:
					context.Views = 1024
				}
				// 将本部分的数据发送到chRst
				chRst <- context
			}()

			// 监控获取成功 or 超时
			select {
			// 此case表示从本部分获取到结果
			case rst := <-chRst:
				fmt.Println("querier", part, ", 得到结果")
				ch <- rst
			// 用于超时退出，停止工作
			case <-stopChs[part]:
				fmt.Println("querier", part, ", 已经停止")
				return
			}
		}(part)
	}

	// 单独启动一个goroutine，接收和整合数据
	wg.Add(1)
	go func() {
		defer wg.Done()
		// 初始化content用于整合数据。
		content := Content{}
		// received是本函数的全局变量，用于标记（信号）记录已完成接收的部分
		received := map[string]struct{}{}
		// 等待接收或者超时退出，超时需通知未完成的goroutine结束
		// 未到超时时间，将结果整合到一起，并判断是否需要继续等待
	loopReceive:
		for { // 一直循环，直到数据接收完毕
			select {
			// 监控ch，不停从里面接收数据（ch最多有3个元素）
			case rst := <-ch:
				fmt.Println("收到数据的部分元素:", rst.Part)
				// 根据不同的part，更新整体content字段
				// 同时记录，哪个part已经完成接收
				switch rst.Part {
				case PartSubject:
					content.Subject = rst.Subject
					received[PartSubject] = struct{}{}
				case PartTags:
					content.Tags = rst.Tags
					received[PartTags] = struct{}{}
				case PartViews:
					content.Views = rst.Views
					received[PartViews] = struct{}{}
				}
				// 判定是否已经接收完毕，是否需要继续等待
				finish := true
				// 确认是否都接收了
				for part := range stopChs {
					if _, exists := received[part]; !exists {
						// 不存在，说明存在处理完成但还未写入的任务，还不能结束
						finish = false
						break
					}
				}
				// 判定已经全部处理完成
				if finish {
					// 说明全部处理完毕，结束
					fmt.Println("所有查询器都完成, 数据是完整的")
					close(ch)
					break loopReceive
				}
			// 超时
			case <-to:
				fmt.Println("查询器超时, Content 是不完整的")
				// 超时时间到要通知未完成的goroutine结束
				// 遍历stopChs，判断是否存在received中
				for part := range stopChs {
					if _, exists := received[part]; !exists {
						// 不存在，说明未处理完成的任务已超时，应该结束
						stopChs[part] <- struct{}{}
					}
				}
				// 关闭
				close(ch)
				// 不再继续监听，结束
				break loopReceive
			}
		}
		// 输出结果（无论数据是否完整）
		fmt.Println("Content:", content)
	}()
	wg.Wait()
}
```

```
func TestSelectAll(t *testing.T) {
	SelectAll()
}
```

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

### 无缓冲Channel + 关闭作典型同步信号

基于：

- 无缓冲Channel是同步的
- **关闭的（无缓冲）channel是可以接收内容的（接收到的是零值）**

以上两点原因，经常使用关闭无缓冲channel的方案来作为信号传递使用。前提是，信号纯粹是信号，没有其他含义，比如关闭时间等。

示例代码：

```go
func SelectChannelCloseSignal() {
	wg := sync.WaitGroup{}
	// 定义无缓冲channel
	// 作为一个终止信号使用（啥功能的信号都可以，信号本身不分功能）
	ch := make(chan struct{})

	// goroutine，用来close, 表示
	// 发出信号
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
			time.Sleep(300 * time.Millisecond)
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

系统信号也是通过channel与应用程序交互，例如典型的 ctrl+c 中断程序， `os.Interrupt`。若不监控系统信号，ctrl+c后程序会直接终止，而**如果监控了信号，那么可以在ctrl+c后，执行一系列的关闭处理**，例如：

#### **`SIGINT` 信号的来源**

`SIGINT` 是 POSIX 标准定义的信号之一，表示 **中断信号**（Signal Interrupt）。它通常由用户通过键盘输入 `Ctrl+C` 触发，**用于请求终止正在运行的程序。**

#### **`os.Interrupt` 的用途**

在 Go 中，`os.Interrupt` （**中断信号**）被用于**信号通知机制**，它可以被 `os/signal` 包监听和捕获。程序可以通过捕获 `os.Interrupt` 信号，实现优雅退出（如清理资源或保存状态）。

#### 小结

- **`os.Interrupt` 代表中断信号 `SIGINT`**，通常由用户通过键盘输入 `Ctrl+C` 触发。
- 它被广泛用于捕获和处理用户终止操作，实现程序的优雅退出。

```go
func SelectSignal() {
    // 一、模拟一段长时间运行的goroutine
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

- 定时器**Timer类似于一次性闹钟**
- 断续器**Ticker类似于重复性闹钟，也成循环定时器**

无论是一次性还是重复性计时器，都是通过Channel与应用程序交互的。我们通过监控Timer和Ticker返回的Channel，来确定是否到时的需求。

### 定时器

![image.png](https://fynotefile.oss-cn-zhangjiakou.aliyuncs.com/fynote/fyfile/13080/1672136903072/164a606a60c841129b548fb05f074df0.png)

使用语法：

```go
// time.NewTimer
func NewTimer(d Duration) *Timer
```

创建定时器，参数是Duration时间。返回为 `*Timer`。`*Timer.C`是用来接收到期通知的单向Channel。

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

如果希望在**定时器到期时执行特定函数**，可以使用如下函数：

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
