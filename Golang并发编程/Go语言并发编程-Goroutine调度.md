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
    // 用于保证不会被拷贝
    noCopy noCopy
    // 当前状态，存储计数器，存储等待的goroutine
    state1 uint64
    state2 uint32
}
```

状态 32bit和64bit的计算机不同，以64bit为例：

- 高32 bits是计数器
- 低32 bits是等待者

Add() 和 Done() 是用来操作计数器，操作计数器的操作是原子操作，保证并发安全性。

Wait()操作，在计数器为0时，结束阻塞状态。

核心代码示例：

```go
func (wg *WaitGroup) Add(delta int) {
    // 原子操作，累加计数器
    state := atomic.AddUint64(statep, uint64(delta)<<32)
}

func (wg *WaitGroup) Done() {
    wg.Add(-1)
}

func (wg *WaitGroup) Wait() {

    for {
        state := atomic.LoadUint64(statep)
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

### goroutine的最小为2KB

之所以支持百万级的goroutine并发，核心原因是因为每个goroutine的**初始栈内存为2KB**，用于保持goroutine中的执行数据，例如局部变量等。相对来说，线程线程的栈内存通常为2MB。除了比较小的初始栈内存外，**goroutine的栈内存可扩容的**，也就是说支持按需增大或缩小，一个goroutine最大的栈内存当前限制为1GB。

### goroutine内部资源竞争溢出

在goroutine内增加，fmt.Println() 测试：

```shell
panic: too many concurrent operations on a single file or socket (max 1048575)
```

### 控制并发数量

实际开发时，要根据系统资源和每个goroutine锁消耗的资源来控制并发规模。

典型方案 goroutine pool，典型的包：

- Jeffail/tunny
- panjf2000/ants
- go-playground/pool

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

![image.png](https://fynotefile.oss-cn-zhangjiakou.aliyuncs.com/fynote/fyfile/13080/1672136903072/0dca3eb676b04825b0ef480379915ce0.png)

当前的计算机都支持多线程，因此现代语言实现的协程调度器通常都是多对多的调度方案：

![image.png](https://fynotefile.oss-cn-zhangjiakou.aliyuncs.com/fynote/fyfile/13080/1672136903072/94443c4e3e5749eda7c594ca3f5a69d6.png)

### GMP模型结构

Goroutine 就是Go语言实现的协程模型。其核心结构有三个，称为GMP，也叫GMP模型。分别是：

- G，Goroutine，我们使用关键字go调用的函数。存储于P的本地队列或者是全局队列中。
- M，**M**achine，就是 Work Thread，就是传统意义的线程，用于执行Goroutine，G。**只有在M与具体的P绑定后，才能执行P中的G。**
- P，Processor，**处理器**，主要用于协调G和M之间的关系，存储需要执行的G队列，与特定的M绑定后，执行Go程序，也就是G。

GMP整体结构逻辑图：

![image.png](https://fynotefile.oss-cn-zhangjiakou.aliyuncs.com/fynote/fyfile/13080/1672136903072/ec55e59f84304b9787e3c2706868b318.png)

图例说明：

- M程序线程，由OS负责调度交由具体的CPU核心中执行
- 待执行的G可能存储于全局队列，和某个P的本地队列中。**P的本地队列当前最多存储256个G。**
- 若要执行P中的G，则P必须要于对应的M建立联系。建立联系后，就可以执行P中的G了。
- M与P不是强绑定的关系，若一个M阻塞，那么P就会选择一个新的M执行后续的G。该过程由Go调度器调度。

GMP的关系

- G 是独立的运行单元
- M是执行任务的线程
- P是G和M的关联纽带
- G要在M中执行，P的任务就是合理的将G分配给M

### P的数量

**P的数量通常是固定的，当程序启动时由 `$GOMAXPROCS`环境变量决定创建P的数量。默认的值为当前CPU的核心数所有的 P 都在程序启动时创建。**

这意味着程序的执行过程中，最多同时有$GOMAXPROCS个Goroutine同时运行，默认与CPU核数保持一致，可以最大程度利用多核CPU并行执行的能力。

程序运行时，`runtime.GOMAXPROCS()`函数可以动态改变P的数量，但通常不建议修改，或者即使修改也不建议数量超过CPU的核数。调动该函数的典型场景是控制程序的并行规模，例如：

```go
// 最多利用一半的CPU
runtime.GOMAXPROCS(runtime.NumCPU() / 2)

// 获取当前CPU的核数
runtime.NumCPU()
```

我们知道Go没有限定G的数量，那M的数量呢？

- Go对M的数量做了一个上限，**10000个**，但通常不会到达这个规模，因为操作系统很难支持这么多的线程。
- M的数量是由P决定的。
- 当P需要执行时，会去找可用的M，若没有可用的M，就会创建新的M，这就会导致M的数量不断增加
- 当M线程长时间空闲，也就是长时间没有新任务时，GC会将线程资源回收，这会导致M的数量减少
- 整体来说，M的数量一定会多于P的数量，取决于空闲（没有G可执行的）的，和完成其他任务（例如CGO操作，GC操作等）的M的数量多少

### P与G关联的流程

go 创建的待执行的Goroutine与P建立关联的核心流程：

- 新创建的G会**优先保持在P的本地队列中**。例如A函数中执行了 go B()，那么B这个Goroutine会优先保存在A所属的P的本地队列中。
- 若G加入P的本地队列时**本地队列已满**，那么G会被加入到全局G队列中。新G加入全局队列时，会把**P本地队列中一半的G**也同时移动到全局队列中（是乱序入队列），以保证P的本地队列可以继续加入新的G。
- 当P要执行G时
  - 会从P的本地队列查找G。
  - 若本地队列中没有G，则会尝试从**其他的P中**偷取（Steal）G来执行，通常会偷取一半的G。
  - 若无法从其他的P中偷取G，则从全局G队列中获取G，会一次获取多个G。
  - **整体：本地G队列->其他P的本地G队列->全局G队列**
- 当全局运行队列中有待执行的 G 时，还会有固定几率（每61个调度时钟周期 schedtick）会从全局的运行队列中查找对应的 G，为了保证全局G队列一定可以被调度。

核心流程图例：

A 中调用了 go B()， P的本地队列未满时：

![image.png](https://fynotefile.oss-cn-zhangjiakou.aliyuncs.com/fynote/fyfile/13080/1672136903072/d2296ed925264f3cbc624d27de06637e.png)

A 中调用了 go B()， P的本地队列已满时：

![image.png](https://fynotefile.oss-cn-zhangjiakou.aliyuncs.com/fynote/fyfile/13080/1672136903072/de05e69f33664d2295fd528eea5e2d47.png)

当P要执行G时：

![image.png](https://fynotefile.oss-cn-zhangjiakou.aliyuncs.com/fynote/fyfile/13080/1672136903072/615654aa1f5f45ff89c4d7cd90835b2b.png)

### P与M关联的流程

P中关联了大量待执行的G，若需要执行G，P要去找可用的M。P不能执行G，只有M才能真正执行G。

P与M建立关联的核心过程：

- 当P需要执行时，P要寻找可用的M，优先从**空闲M池中找**，若没有空闲的，则新建M来执行
- 在创建G时，G会尝试唤醒空闲的M
- 当M的执行因为G进行了**系统调用时**，M会释放与之绑定的P，把P转移给其他的M去执行。称为**P抢占。**
- 当M执行完的系统调用阻塞的G后，M会尝试获取新的空闲P，同时将G放入P的本地队列执行。若没有空闲的P，则将G放入全局G队列，M进入休眠，等待被唤醒或被垃圾回收

如图所示：

P要寻找可用的M：

![image.png](https://fynotefile.oss-cn-zhangjiakou.aliyuncs.com/fynote/fyfile/13080/1672136903072/1094a65c5a534e7a8ff9118597cd6fbf.png)

G执行了系统调用，M与P解绑（释放），P转移到新的M上执行：

![image.png](https://fynotefile.oss-cn-zhangjiakou.aliyuncs.com/fynote/fyfile/13080/1672136903072/e71e6b03e1414dcd99b538633ebec8d4.png)

M1执行完G的系统调用后，G不一定结束，还要继续执行，则M1会尝试获取空闲的P（没有与M绑定的P），若没有空闲的P可用，将M1执行的G放入全局G队列，M1进入空闲状态：

![image.png](https://fynotefile.oss-cn-zhangjiakou.aliyuncs.com/fynote/fyfile/13080/1672136903072/7402e28b89194a23be7814f2b766a86e.png)

### G的调度流程总结

整体上看

1. go func() 创建Goroutine
2. 将Goroutine放入队列

- 放入本地队列
- 本地队列满，放入全局队列

3. M通过P获取G运行

- 从本地队列获取G
- 从其他P的本地队列获取G
- 从全局队列获取G

4. M执行G

- 调度周期循环执行G
- G主动让出
- G执行系统调用

5. G执行系统调用

- 解绑G和P
- P抢占其他的M继续执行
- 系统调用的G结束，将G放入其他P队列执行，M空闲

6. 若G执行完毕，释放

### M0 和 G0

- M0, 启动程序后的编号为 0 的主线程，负责执行初始化操作和启动第一个 G，也就是 main Goroutine。之后与其他M一样调度。
- G0，每个 M 创建的第一个 Goroutine。G0 仅用于负责调度的 G，G0 不指向任何可执行的函数，每个 M 都会有一个自己的 G0。在调度或系统调用时会使用 G0 的栈空间。

如图：

M的G0：

![image.png](https://fynotefile.oss-cn-zhangjiakou.aliyuncs.com/fynote/fyfile/13080/1672136903072/98e97bf2f1d5401189882a781949f27e.png)

### 协作和抢占调度

当某个 G 执行时间过长，其他的 G 如何调度。通常有两种方案：

- 协作式，主动让出执行权，让其他G执行。通过runtime.Gosched()可以让出执行权。
- 抢占式，被动让出执行权，也就是调度器将G的执行权取消，分配给其他的G。Go目前默认的方式。在Go中一个G最多可以执行10ms，超时就会被让出调度权。

函数：

```go
runtime.Gosched()
```

方法可以要求Go主动调度该goroutine，去执行其他的goroutine。这种是典型的协作调度模式，类似于 py 的 yield。

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

我们采用1个P来进行模拟，看主动让出交替执行的情况。

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

  - G，Goroutine，独立并发执行的代码段
  - M，mechine, 系统线程
  - P，Processor，逻辑处理器，用于联系G和M。
  - G存在与P的本地队列或全局队列中
  - M要与P绑定，P中的G才会执行
  - M执行G中的系统调用时，会解绑M和P，P会找到新的M执行
