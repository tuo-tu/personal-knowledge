# 同步和锁

## 概述

**同步是并发编程的基本要素之一**，我们通过channel可以完成多个goroutine间数据和信号的同步。

除了channel外，我们还可以使用go的官方同步包sync，sync/atomic 完成一些基础的同步功能。主要包含同步数据、锁、原子操作等。

一个同步失败的示例：

```go
func SyncErr() {
    wg := sync.WaitGroup{}
    // 计数器
    counter := 0
    // 多个goroutine并发的累加计数器
    gs := 100
    wg.Add(gs)
    for i := 0; i < gs; i++ {
        go func() {
            defer wg.Done()
            // 累加
            for k := 0; k < 100; k++ {
                counter++
                // ++ 操作不是原子的
                // counter = counter + 1
                // 1. 获取当前的counter变量
                // 2. +1
                // 3. 赋值新值到counter
            }
        }()
    }
    // 统计计数结果
    wg.Wait()
    fmt.Println("Counter:", counter)
}
```

Lock解决方案：

```go
func SyncLock() {
    n := 0
    wg := sync.WaitGroup{}

    lk := sync.Mutex{}
    for i := 0; i < 1000; i++ {
        wg.Add(1)
        go func() {
            defer wg.Done()
            for i := 0; i < 100; i++ {
                lk.Lock()
                n++
                lk.Unlock()
            }
        }()
    }
    wg.Wait()
    fmt.Println("n:", n)
}

// run
n: 100000
```

## 互斥锁Mutex的使用

sync包提供了两种锁：

1. 互斥锁，Mutex
2. 读写互斥锁，RWMutex

互斥锁，**同一时刻只能有一个goroutine申请锁定成功**，**不区分读、写操作**。**也称为独占锁、排它锁**。

提供了如下方法完成锁操作：

```go
type Mutex
// 锁定锁m, 若锁m已是锁定状态，调用的goroutine会被阻塞，直到可以锁定
func (m *Mutex) Lock()
// 解锁锁m，若m不是锁定状态，会导致运行时错误
func (m *Mutex) Unlock()
// 尝试是否可以加锁，返回是否成功
func (m *Mutex) TryLock() bool
```

注意：锁与goroutine没有关联，意味着允许一个goroutine加锁，在另一个goroutine中解锁。但这不是最典型的用法。

**典型的锁用法：**

```go
var lck sync.Mutex
func () {
    lck.Lock()
    // 互斥执行的代码
    defer lck.Unlock()
}
```

锁的流程：

![image.png](https://fynotefile.oss-cn-zhangjiakou.aliyuncs.com/fynote/fyfile/13080/1672136903072/f5e31cf6b402490cba3ccf52748d9f37.png)

示例：

```go
func SyncMutex() {
    wg := sync.WaitGroup{}
    var lck sync.Mutex
    for i := 0; i < 4; i++ {
        wg.Add(1)
        go func(n int) {
            defer wg.Done()
            fmt.Println("before lock: ", n)
            lck.Lock()
            fmt.Println("locked: ", n)
            time.Sleep(1 * time.Second)
            lck.Unlock()
            fmt.Println("after lock: ", n)
        }(i)
    }
    wg.Wait()
}
```

某次输出结果：

```
before lock:  3
locked:  3   
before lock:  2
before lock:  1
before lock:  0
after lock:  3
locked:  2
after lock:  2
locked:  1
after lock:  1
locked:  0
after lock:  0
```

可以发现，before lock 都是先执行的，而Lock() 操作，必须要等到其他goroutineUnlock()才能成功。

注意，**如果其他goroutine没有通过相同的锁（1没用锁，2用了其他锁）去操作资源，那么是不受锁限制的**，例如：

```go
func SyncLockAndNo() {
    n := 0
    wg := sync.WaitGroup{}
    lk := sync.Mutex{}
    for i := 0; i < 1000; i++ {
        wg.Add(1)
        go func() {
            defer wg.Done()
            for i := 0; i < 100; i++ {
                lk.Lock()
                n++
                lk.Unlock()
            }
        }()
    }
    wg.Add(1)
    go func() {
        defer wg.Done()
        for i := 0; i < 10000; i++ {
            n++
        }
    }()

    // 其他锁
    //var lk2 sync.Mutex
    //go func() {
    //    defer wg.Done()
    //    for i := 0; i < 10000; i++ {
    //        lk2.Lock()
    //        n++
    //        lk2.Unlock()
    //    }
    //}()

    wg.Wait()
    fmt.Println("n:", n)
}

// 其中一次结果
n: 109876
```

我们在第一个counter的例子上，增加了一个goroutine同去累加计数器counter，但没有使用前面的Mutex（不使用或使其他锁）。可见，出现了资源争用的情况。因此要注意：**如果要限制资源的并发争用，要全部的资源操作都使用同一个锁。**

实操时，**锁除了直接调用外，还经常性出现在结构体中，以某个字段的形式出现**，用于包含struct字段不会被多gorutine同时修改，例如我们 cancelCtx：

```go
type cancelCtx struct {
    Context

    mu       sync.Mutex            // protects following fields
    done     atomic.Value          // of chan struct{}, created lazily, closed by first cancel call
    children map[canceler]struct{} // set to nil by the first cancel call
    err      error                 // set to non-nil by the first cancel call
}
```

我们通常也会这么做，示例：

```go
type Post struct {
    Subject string
    // 赞
    Likes   int
    // 操作锁定
    mu sync.Mutex
}

func (p *Post) IncrLikes() *Post {
    p.mu.Lock()
    defer p.mu.Unlock()
    p.Likes++

    return p
}

func (p *Post) DecrLikes() *Post {
    p.mu.Lock()
    defer p.mu.Unlock()
    p.Likes--

    return p
}
```

## 读写RWMutex的使用

读写互斥锁，将锁操作类型做了区分，分为读锁和写锁，由sync.RWMutex类型实现：

1. 读锁，Read Lock，**共享读**，阻塞写

   解释：允许多个 goroutine 同时读取共享资源，而不阻塞彼此，但会阻止其他 goroutine 进行写操作（也不能获取写锁）。

2. 写锁，Lock，独占操作，阻塞读写

| 并发 | 读     | 写     |
| ---- | ------ | ------ |
| 读   | 支持   | 不支持 |
| 写   | 不支持 | 不支持 |

之所以减小锁的粒度，因为实际操作中读操作的比例要远高于写操作的比例，增加了共享读操作锁后，可以更大程度的提升读的并发能力。

sync.RWMutex 提供了如下方法完成操作：

```go
type RWMutex
// 写锁定
func (rw *RWMutex) Lock()
// 写解锁
func (rw *RWMutex) Unlock()

// 读锁定
func (rw *RWMutex) RLock()
// 读解锁
func (rw *RWMutex) RUnlock()

// 尝试加写锁定
func (rw *RWMutex) TryLock() bool
// 尝试加读锁定
func (rw *RWMutex) TryRLock() bool
```

写锁定，与互斥锁Mutex的语法和操作结果一致，都是保证互斥的独占操作。

**读锁定，可以在已经存在读锁的情况下，加锁成功。**

如图所示：

![image.png](https://fynotefile.oss-cn-zhangjiakou.aliyuncs.com/fynote/fyfile/13080/1672136903072/25437f58ca6447d5a73abe4ae5b4cdd6.png)

读锁示例：

```go
func SyncRLock() {
    wg := sync.WaitGroup{}
    // 模拟多个goroutine
    var rwlck sync.RWMutex
    for i := 0; i < 10; i++ {
        wg.Add(1)
        go func() {
            defer wg.Done()
            //rwlck.Lock()
            rwlck.RLock()
            // 输出一段内容
            fmt.Println(time.Now())
            time.Sleep(1 * time.Second)
            //
            //rwlck.Unlock()
            rwlck.RUnlock()
        }()
    }

    wg.Add(1)
    go func() {
        defer wg.Done()
        //
        rwlck.Lock()
        //rwlck.RLock()
        // 输出一段内容
        fmt.Println(time.Now(), "Lock")
        time.Sleep(1 * time.Second)
        rwlck.Unlock()
        //rwlck.RUnlock()
    }()

    wg.Wait()
}
```

其中，使用读锁，输出操作会全部立即执行，然后集体sleep1s后全部结束。使用写锁，输出和Sleep会间隔1s依次执行。

实操示例：

```go
type Article struct {
    Subject string
    // 赞
    likes int
    // 操作锁定
    mu sync.RWMutex
}

func (a Article) Likes() int {
    a.mu.RLock()
    defer a.mu.RUnlock()
    return a.likes
}

func (a *Article) IncrLikes() *Article {
    a.mu.Lock()
    defer a.mu.Unlock()
    a.likes++
    return a
}
```

## 同步Map sync.Map

**Go中Map不是并发安全的**。并发操作 Map 类型时，会导致 `fatal error: concurrent map read and map write`错误：

```go
func SyncMapErr() {
    m := map[string]int{}
    // 并发map写
    go func() {
        for {
            m["key"] = 0
        }
    }()
    // 并发map读
    go func() {
        for {
            _ = m["key"]
        }
    }()
    // 阻塞
    select {}
}
```

之所以Go不支持Map的并发安全，是**因为Go认为Map的典型使用场景不需要在多个Goroutine间并发安全操作Map。**

并发安全操作Map的方案：

- 锁 + Map，自定义Map操作，增加锁的控制，可以选择 Mutex和RWMutex。
- sync.Map，sync包提供的安全Map.

锁+Map示例，在结构体内嵌入sync.Mutex：

```go
func SyncMapLock() {
    myMap := struct {
        sync.RWMutex
        Data map[string]int
    }{
        Data: map[string]int{},
    }

    // write
    myMap.Lock()
    myMap.Data["key"] = 0
    myMap.Unlock()

    // read
    myMap.RLock()
    _ = myMap.Data["key"]
    myMap.RUnlock()
}
```

sync.Map 的使用：最常用的4个方法：**Store、Range、Delete、Load**

```go
type Map
// 1.存储
func (m *Map) Store(key, value any)
// 2.遍历所有的键值对，如果返回 true，继续遍历；如果返回 false，停止遍历。
func (m *Map) Range(f func(key, value any) bool)
// 3.删除某个key元素
func (m *Map) Delete(key any)
// 4.返回key的值。存在key，返回value，true，不存在返回 nil, false
func (m *Map) Load(key any) (value any, ok bool)

// 若m[key]==old，执行删除。key不存在，返回false
func (m *Map) CompareAndDelete(key, old any) (deleted bool)
// 若m[key]==old，执行交换, m[key] = new
func (m *Map) CompareAndSwap(key, old, new any) bool

// 返回value后删除元素，loaded 表示是否load成功，key不存在，loaded为false
func (m *Map) LoadAndDelete(key any) (value any, loaded bool)
// 加载失败则存储，返回加载或存储的值和是否加载
func (m *Map) LoadOrStore(key, value any) (actual any, loaded bool)

// 存储新值，返回交换前的值，value表示要设置的新值。loaded表示key是否存在
func (m *Map) Swap(key, value any) (previous any, loaded bool)
```

- `Swap()`：将指定键的值直接替换为新值，无需检查当前值是否匹配某个期望值。

- `CompareAndSwap()`：只有在键的值与指定的旧值匹配时，才会将其更新为新值。

sync.Map 不需要类型初始化，即可使用，可以理解为map[comparable]any。

使用示例，不会触发 `fatal error: concurrent map read and map write`：

```go
func SyncSyncMap() {
    var m sync.Map
    go func() {
        for {
            m.Store("key", 0)
        }
    }()
    go func() {
        for {
            _, _ = m.Load("key")
        }
    }()
    select {}
}
```

使用示例：

```go
func SyncSyncMapMethod() {
    wg := sync.WaitGroup{}
    var m sync.Map
    for i := 0; i < 10; i++ {
        wg.Add(1)
        go func(n int) {
            defer wg.Done()
            m.Store(n, fmt.Sprintf("value: %d", n))
        }(i)
    }

    for i := 0; i < 10; i++ {
        wg.Add(1)
        go func(n int) {
            defer wg.Done()
            fmt.Println(m.Load(n))
        }(i)
    }

    wg.Wait()
    m.Range(func(key, value any) bool {
        fmt.Println(key, value)
        return true
    })

    //
    m.Delete(4)
}
```

并发安全操作Map的方案的选择，统计的压测数据显示，相对而言：

- 锁 + Map，写快，读慢
- sync.Map，读快，写慢，删快，适合读多写少的场景

**sync.Map遍历示例代码：**

```go
func main() {
	var m sync.Map

	// 存储键值对
	m.Store("key1", "value1")
	m.Store("key2", "value2")
	m.Store("key3", "value3")

	// 遍历 Map
	m.Range(func(key, value any) bool {
		fmt.Printf("Key: %v, Value: %v\n", key, value)
		return true // 返回 true，继续遍历
	})

	// 停止遍历的示例
	fmt.Println("Stopping early:")
	m.Range(func(key, value any) bool {
		fmt.Printf("Key: %v, Value: %v\n", key, value)
		return false // 返回 false，停止遍历
	})
}
```

**输出示例（具体顺序可能不同）：**

```yaml
Key: key1, Value: value1
Key: key2, Value: value2
Key: key3, Value: value3
Stopping early:
Key: key1, Value: value1
```



## 原子操作 sync/atomic

**原子操作即是进行过程中不能被中断的操作，针对某个值的原子操作在被进行的过程中，CPU绝不会再去进行其他的针对该值的操作**。为了实现这样的严谨性，原子操作仅会由一个独立的CPU指令代表和完成。**原子操作是无锁的**，**常常直接通过CPU指令直接实现。** 事实上，其它同步技术的实现常常依赖于原子操作。

**原子操作是CPU指令级别实现的**，比如在Intel的CPU上主要是使用总线锁的方式，AMD的CPU架构机器上就是使用MESI一致性协议的方式来保证原子操作。

go中 `sync/atomic` 包提供了原子操作的支持，用于同步操作整型（和指针类型）：

- int32
- int64
- uint32
- uint64
- uintptr
- unsafe.Pointer

针对于以上类型，提供了如下操作：Type 是以上的类型之一

```go
// Type 是以上的类型之一
// 比较相等后交换 CAS
func CompareAndSwapType(addr *Type, old, new Type) (swapped bool)
// 交换
func SwapType(addr *Type, new Type) (old Type)
// 累加
func AddType(addr *Type, delta Type) (new Type)
// 获取
func LoadType(addr *Type) (val Type)
// 存储
func StoreType(addr *Type, val Type)
```

除了以上函数，还提供了对应的类型方法操作，以Int32为例：

```go
type Int32
func (x *Int32) Add(delta int32) (new int32)
func (x *Int32) CompareAndSwap(old, new int32) (swapped bool)
func (x *Int32) Load() int32
func (x *Int32) Store(val int32)
func (x *Int32) Swap(new int32) (old int32)
```

除了以上几个整型，bool类型也提供了类型上的原子操作：

```go
type Bool
func (x *Bool) CompareAndSwap(old, new bool) (swapped bool)
func (x *Bool) Load() bool
func (x *Bool) Store(val bool)
func (x *Bool) Swap(new bool) (old bool)
```

示例：

```go
func SyncAtomicAdd() {
    // 并发的过程，没有加锁，Lock
    // var counter int32 = 0
    // type
    // atomic 原子的Int32, counter := 0
    counter := atomic.Int32{}
    wg := sync.WaitGroup{}
    for i := 0; i < 1000; i++ {
        wg.Add(1)
        go func() {
            defer wg.Done()
            for i := 0; i < 100; i++ {
                //atomic.AddInt32(&counter, 1)
                // type
                // 原子累加操作，counter++
                counter.Add(1)
            }
        }()
    }
    wg.Wait()
    //fmt.Println("counter:", atomic.LoadInt32(&counter))
    // type
    fmt.Println("counter:", counter.Load())
}
```

以上示例不会出现不到10000的情况了。

**除了预定义的整型的支持**，还可以使用 **atomic.Value** 类型，完成其他类型的原子操作：

```go
type Value struct {
	v any
}

func (v *Value) CompareAndSwap(old, new any) (swapped bool)
func (v *Value) Load() (val any)
func (v *Value) Store(val any)
func (v *Value) Swap(new any) (old any)
```

使用方法：

```go
func SyncAtomicValue() {

    var loadConfig = func() map[string]string {
        return map[string]string{
            // some config
            "title":   "马士兵Go并发编程",
            "varConf": fmt.Sprintf("%d", rand.Int63()),
        }
    }

    var config atomic.Value

    // 每N秒加载一次配置文件
    go func() {
        for {
            config.Store(loadConfig())
            fmt.Println("latest config was loaded", time.Now().Format("15:04:05.99999999"))
            time.Sleep(time.Second)
        }
    }()

    // 使用配置
    // 不能在加载的过程中使用配置
    for {
        go func() {
            c := config.Load()
            fmt.Println(c, time.Now().Format("15:04:05.99999999"))
        }()

        time.Sleep(400 * time.Millisecond)
    }

    select {}
}
```

## sync.Pool 并发安全池

![image.png](https://fynotefile.oss-cn-zhangjiakou.aliyuncs.com/fynote/fyfile/13080/1672136903072/f5d537e860584db4ae86756dbdf1eee4.png)

**`sync.Pool`用于缓存和复用临时对象。**存储在池中的任何项目可随时自动删除，无需通知。**一个池可以安全地同时被多个goroutine使用。**

典型特征：

- sync.Pool 是并发安全的
- **池中的对象由Go负责删除**，内存由Go自己回收
- 池中元素的数量由Go负责管理，用户无法干预
- **池中元素应该是临时的，不应该是持久的**。例如长连接不适合放入 sync.Pool 中

池的目的是缓存已分配但未使用的项目以供以后**重用**，从而减轻垃圾收集器的压力。也就是说，它使构建高效、线程安全的自由元素变得容易。

池的一个适当用途是**管理一组临时项**，这些临时项在包的并发独立客户端之间默默**共享**，并可能被其重用。池提供了一种在许多客户机上分摊分配开销的方法。

一个很好地使用池的例子是fmt包，它维护了临时输出缓冲区的动态大小存储。

池由 sync.Pool类型实现，具体三个操作：

1. **初始化Pool实例**，需要提供池中缓存元素的New方法。
2. **申请元素**，func (p *Pool) Get() any
3. **交回对象**，func (p *Pool) Put(x any)

`pool`结构体介绍：

```go
type Pool struct {
	noCopy noCopy // 防止 Pool 结构体被复制

	local     unsafe.Pointer // local fixed-size per-P pool, actual type is [P]poolLocal
	localSize uintptr        // 本地数组的大小

	victim     unsafe.Pointer // local from previous cycle
	victimSize uintptr        // victim 数组的大小。

    // 一个可选的用户自定义函数，用于在池中没有可用对象时，生成一个新对象并返回。(破案了)
	New func() any
}
```

pool操作示例：

```go
func SyncPool() {
    // 原子的计数器
    var counter int32 = 0

    // 定义元素的Newer，创建器。定义了一个 atomic 计数器
    elementNewer := func() any {
        // 原子的计数器累加
        atomic.AddInt32(&counter, 1)
        // 池中元素推荐（强烈）是指针类型
        return new(bytes.Buffer)
    }

    // Pool的初始化
    pool := sync.Pool{
        New: elementNewer,
    }

    // 并发的申请和交回元素
    workerNum := 1024 * 1024
    wg := sync.WaitGroup{}
    wg.Add(workerNum)
    for i := 0; i < workerNum; i++ {
        go func() {
            defer wg.Done()
            // 申请元素，通常需要断言为特定类型
            buffer := pool.Get().(*bytes.Buffer)
            // 不用Pool
            // buffer := elementNewer().(*bytes.Buffer)
            // 交回元素
            defer pool.Put(buffer)
            // 使用元素
            _ = buffer.String()
        }()
    }

    wg.Wait()
    // 测试创建元素的次数
    fmt.Println("elements number is :", counter)
}

// elements number is : 12
```

测试的时候，大家可以发现创建的元素数量远远低于goroutine的数量。

## DATA RACE 现象

**数据竞争**（Data Race）是并发编程中的一种现象，指当**多个线程或 goroutine 同时访问同一个共享变量**，**其中至少有一个是写操作**，并且没有使用任何同步机制（如锁、原子操作）来保证访问的顺序时，**可能会导致未定义的行为**。有时在编写代码时很难发现，要经过大量测试才会发现。

可以用 `go run -race`命令 ，通过增加**-race**选项，检测运行时可能出现的竞争问题。

测试之前的计数器累加代码：本例子需要 main来演示，因为是 go run：

```go
func main() {
    wg := sync.WaitGroup{}
    // 计数器
    counter := 0
    // 多个goroutine并发的累加计数器
    gs := 1000
    wg.Add(gs)
    for i := 0; i < gs; i++ {
        go func() {
            defer wg.Done()
            // 累加
            for k := 0; k < 100; k++ {
                counter++
                // ++ 操作不是原子的
                // counter = counter + 1
                // 1. 获取当前的counter变量
                // 2. +1
                // 3. 赋值新值到counter
            }
        }()
    }
    // 统计计数结果
    wg.Wait()
    fmt.Println("Counter:", counter)
}
```

结果：

```shell
# 没有使用 -race
PS D:\apps\goExample\concurrency> go run .\syncRace.go
n: 94077

# 使用 -race
PS D:\apps\goExample\concurrency> go run -race .\syncRace.go
==================
WARNING: DATA RACE  
Read at 0x00c00000e0f8 by goroutine 9:
  main.main.func1()
      D:/apps/goExample/concurrency/syncMain.go:16 +0xa8

Previous write at 0x00c00000e0f8 by goroutine 7:
Goroutine 9 (running) created at:
  main.main()
      D:/apps/goExample/concurrency/syncMain.go:13 +0x84

Goroutine 7 (finished) created at:
  main.main()
      D:/apps/goExample/concurrency/syncMain.go:13 +0x84
==================
n: 98807
Found 1 data race(s)
exit status 66
```

**该选项用于在开发阶段，检测数据竞争情况。**出现 data race情况，可以使用**锁，或原子操作**的来解决。

## sync.Once

若需要保证多个并发goroutine中，**某段代码仅仅执行一次**，就可以使用 sync.Once 结构实现。

例如，在加载配置的时候，往往仅仅需要获取/加载一次，然后去使用。在多个goroutine并发时，要保证能够获取到配置，同时仅获取一次配置，就可以使用sync.Once结构：

```go
func SyncOnce() {

    // 初始化config变量
    config := make(map[string]string)

    // 1.初始化 sync.Once
    once := sync.Once{}

    // 加载配置的函数
    loadConfig := func() {
        // 2.利用 once.Do() 来执行
        once.Do(func() {
            // 保证执行一次
            config = map[string]string{
                "varInt": fmt.Sprintf("%d", rand.Int31()),
            }
            fmt.Println("config loaded")
        })
    }

    // 模拟多个goroutine，多次调用加载配置
    // 测试加载配置操作，执行了几次
    workers := 10
    wg := sync.WaitGroup{}
    wg.Add(workers)
    for i := 0; i < workers; i++ {
        go func() {
            defer wg.Done()
            // 并发的多次加载配置
            loadConfig()
            // 使用配置
            _ = config

        }()
    }
    wg.Wait()
}
```

核心逻辑：

1. 初始化 sync.Once
2. `once.Do(func())` 可以确保函数func()仅仅执行一次

sync.Once 的实现很简单：

```go
type Once struct {
    // 是否已处理，保证一次
    done uint32
    // 锁，保证并发安全
    m   Mutex
}
```

## sync.Cond

sync.Cond是sync包提供的**基于条件（Condition）**的通知结构。

在 Golang 中，`sync.Cond` 是一种同步原语，用于**让 goroutine 等待或通知某个条件的变化**。它**通常用于需要一个或多个 goroutine 等待某个条件满足后才能继续执行的场景。**

该结构提供了4个方法：

```go
// 创建一个与指定锁相关联的 sync.Cond
func NewCond(l Locker) *Cond
// 唤醒所有正在等待该条件的goroutine
func (c *Cond) Broadcast()
// 唤醒一个正在等待该条件的goroutine
func (c *Cond) Signal()
// 阻塞调用的goroutine，直到被通知（Signal 或 Broadcast）
func (c *Cond) Wait()
```

`sync.Cond` 需要与一个 `sync.Locker`（如 `sync.Mutex` 或 `sync.RWMutex`）关联，用于保护条件变量。然后**两个方法用来通知，一个方法用来等待**。

使用逻辑很简单，通常是一个goroutine负责通知，多个goroutine等待处理，如图：

![image.png](https://fynotefile.oss-cn-zhangjiakou.aliyuncs.com/fynote/fyfile/13080/1672136903072/fb09755625274325a6a7fa6d90ca3798.png)

创建Cond，sync.NewCond() 需要提供锁，例如：

```go
cond := sync.NewCond(&sync.Mutex{})
cond := sync.NewCond(&sync.RWMutex{})
```

**关键点：**

- 在调用 `Wait()` 之前，必须先加锁（使用 `sync.Mutex` 或 `sync.RWMutex`），并**在通知或广播后解锁。**
- `Wait()` 会临时释放锁，以允许其他 goroutine 继续运行，并在被唤醒后重新获取锁。

- 而广播（信号）操作的加锁是可选的。
- `sync.Cond` 是处理复杂同步场景（例如**生产者-消费者模式**或协同多个 goroutine）的强大工具。

还有，**cond的广播和信号通知操作是并发安全的，可以重复调用。**

示例代码：

- 一个goroutine负责生产数据（生产者），条件满足后，广播给处理数据的goroutine（消费者）。
- 多个goroutine处理数据（消费者），在条件未满足前，等待广播信号，信号来了，处理数据。

```go
func main() {
	var m sync.Mutex
	cond := sync.NewCond(&m)
	ready := false // 条件变量：是否准备好

	// 消费者 goroutine
	for i := 1; i <= 3; i++ {
		go func(id int) {
			m.Lock()
			for !ready { // 等待条件满足
				fmt.Printf("消费者 %d 正在等待...\n", id)
				cond.Wait() // 等待通知
			}
			fmt.Printf("消费者 %d 收到通知，开始工作！\n", id)
			m.Unlock()
		}(i)
	}

	// 生产者 goroutine
	go func() {
		time.Sleep(2 * time.Second) // 模拟准备过程
		m.Lock()
		ready = true // 条件满足
		fmt.Println("生产者：所有消费者可以开始工作！")
		cond.Broadcast() // 广播通知所有等待的消费者
		m.Unlock()
	}()

	// 等待程序运行
	time.Sleep(5 * time.Second)
}
```

**运行结果示例：**

```cmd
消费者 1 正在等待...
消费者 2 正在等待...
消费者 3 正在等待...
生产者：所有消费者可以开始工作！
消费者 1 收到通知，开始工作！
消费者 2 收到通知，开始工作！
消费者 3 收到通知，开始工作！
```

示例代码要点：

1. **wait所在的goroutine要判定是否需要wait**，所以wait要出现在条件中，因为goroutine调用的关系，不能保证wait在broadcast前面执行
2. **wait必须使用for进行条件判定**，因为在wait返回后，条件不一定成立，因为可能出现 **虚假唤醒**。`Broadcast()` 可能在条件未满足时被调用，导致 `Wait` 过早返回（通常是在其他的goroutine中）。
3. Broadcast() 操作可选的是否加锁解锁，加锁与否取决于应用场景。
4. **Wait() 操作前，一定要加锁**。因为`Wait()`的代码会在调用时临时解锁（以便其他 goroutine 可以获得锁），接收到信号后，再加锁。

### sync.Cond 基本原理

sync.Cond结构：

```go
type Cond struct {
    noCopy noCopy		 // 防止结构体被错误复制，确保线程安全
    L Locker			 // 通常是sync.Mutex或sync.RWMutex
    notify  notifyList 	 // 等待通知的goroutine队列
    checker copyChecker  // 检查Cond结构体是否被复制
}
```

结构上可见，**Cond记录了等待的goroutine列表**，这样就可以做到，广播到全部的等待goroutine。这也是Cond不应该被复制的原因，否则这些goroutine可能会被意外唤醒。

Wait() 操作：

```go
func (c *Cond) Wait() {
   // 1.检查 Cond 是否被非法复制
   c.checker.check()
   // 2.将当前goroutine添加到内部的 notifyList 等待队列中，并返回当前需要插入notifyList的goroutine的编号
   t := runtime_notifyListAdd(&c.notify)
   // 3.临时解锁，这样其他 goroutine 可以获得锁
   c.L.Unlock()
   // 4.将当前的goroutine挂起，直到被Signal()或Broadcast()唤醒
   runtime_notifyListWait(&c.notify, t)
   // 5.唤醒之后，重新加锁，因为阻塞之前解锁了。
   c.L.Lock()
}
```

核心工作就是，**将当前goroutine添加到Cond的notifyList队列中**。之后解锁，挂起，加锁。因此要在Wait()前加锁，后边要解锁。

Broadcast()操作：

```go
func (c *Cond) Broadcast() {
   // 1.检查 sync.Cond 是否被复制了（是否被非法复制）
   c.checker.check()
   // 2.唤醒notifyList中的所有goroutine，让它们从Wait()方法中返回。
   runtime_notifyListNotifyAll(&c.notify)
}
```

核心工作就是**唤醒 notifyList 中全部的 goroutine。**

## 小结

同步类型：

- 数据同步，保证数据操作的原子性
  - sync/atomic
  - sync.Map
  - sync.Mutex, sync.RWMutex
- 操作同步
  - sync.Mutex, sync.RWMutex

锁的类型：

- 互斥锁 sync.Mutex，完全独占
- 读写互斥锁 sync.RWMutex，可以共享读操作

锁的不锁资源，只是锁定申请锁本身的操作。

sync包总结

- 锁：sync.Mutex, sync.RWMutex
- 数据：sync.Map, sync/atomic
- sync.Pool
- sync.Once
- sync.Cond

使用Channel完成数据和信号的同步！
