# Context上下文

## goroutine的通信方式

主要就是两种：

- **channel**：推荐使用，因为它符合 Go 的并发设计理念：**不要通过共享内存来通信，而是通过通信来共享内存**。

- context：本节内容

## Context概述

Go 1.7 标准库引入 context，译作“上下文”，**准确说它是 goroutine 的上下文（实际中貌似不局限于G）**，包含 goroutine 的**运行状态、环境、现场等信息**。

**context 主要用来在 goroutine 之间传递上下文信息**，包括：取消信号、超时时间、截止时间、k-v 等。

随着 context 包的引入，标准库中很多接口因此加上了 context 参数，例如 database/sql 包。**context 几乎成为了并发控制和超时控制的标准做法。**

**在一组goroutine 之间传递共享的值、取消信号、deadline是Context的作用**。

以典型的HTTPServer为例：

![image.png](https://fynotefile.oss-cn-zhangjiakou.aliyuncs.com/fynote/fyfile/13080/1672136903072/f6704747c32d4736b7e63103fa112d61.png)

我们以 Context II为例，若没有上下文信号，当其中一个goroutine出现问题时，其他的goroutine不知道，还会继续工作。这样的无效的goroutine积攒起来，就会导致goroutine雪崩，进而导致服务宕机！

没有同步信号：

<img src="https://fynotefile.oss-cn-zhangjiakou.aliyuncs.com/fynote/fyfile/13080/1672136903072/91e4d12dec0c4a07887c210d302ab8d5.png" alt="image.png" style="zoom:80%;" />

增加同步信号：

<img src="https://fynotefile.oss-cn-zhangjiakou.aliyuncs.com/fynote/fyfile/13080/1672136903072/0e9e2717de4a4f00b68aff69b187fe44.png" alt="image.png" style="zoom:90%;" />

参考：Context传递取消信号小结。

**Context 的作用**：

- **超时控制**：可以为操作设置超时时间。
- **取消信号**：可以通过取消一个 `context` 终止一组协程。
- **携带数据**：可在上下文中传递轻量级、只读的数据（例如请求 ID、用户认证信息）。

## Context 核心结构

`context.Context` 是 Go 语言在 1.7 版本中引入标准库的接口，该接口定义了四个需要实现的方法：

```go
type Context interface {
    // 返回上下文何时超时的时间点（若未设置，则 ok 返回 false）。
	Deadline() (deadline time.Time, ok bool)
    // 返回用于通知Context完结的channel
    // 当context被取消时，会关闭此channel
    // 在子goroutine里读这个channel，除非被关闭，否则读不出来任何东西（这点挺重要）
	Done() <-chan struct{}
    // 返回Context取消的错误
    Err() error
    // 返回key对应的value
	Value(key any) any
}
```

重点：

- 当 `Context` 被取消或超时时，`Done` 通道会关闭。

- 通过监听 `Done()` 通道，可以优雅地结束 goroutine 的工作。

除了Context接口，还存在一个canceler接口，用于**实现Context可以被取消**：

```go
type canceler interface {
	cancel(removeFromParent bool, err error)
	Done() <-chan struct{}
}
```

除了以上两个接口，还有4个**预定义的Context类型**：

```go
// 1.空Context
type emptyCtx int

// 2.取消Context
type cancelCtx struct {
	Context
	mu       sync.Mutex            // protects following fields
	done     atomic.Value          // of chan struct{}, created lazily, closed by first cancel call
	children map[canceler]struct{} // set to nil by the first cancel call
	err      error                 // set to non-nil by the first cancel call
}

// 3.定时取消Context
type timerCtx struct {
	cancelCtx
	timer *time.Timer // Under cancelCtx.mu.

	deadline time.Time
}

// 4.KV值Context
type valueCtx struct {
	Context
	key, val any
}

```

## 默认(空)Context的使用

context 包中最常用的方法是 `context.Background`、`context.TODO`，这两个方法都会返回**预先初始化好的**私有变量 background 和 todo，它们会在同一个 Go 程序中被复用：

- `context.Background`： 是上下文的默认值，**所有其他的上下文都应该从它衍生出来**，在多数情况下，如果当前函数没有上下文作为入参，我们都会使用 `context.Background` 作为**起始的**上下文向下传递。
- context.TODO，是一个备用，一个context占位，通常用在并不知道传递什么 context的情形。

使用示例，`database/sql`包中的执行：

```sql
func (db *DB) PingContext(ctx context.Context) error
func (db *DB) ExecContext(ctx context.Context, query string, args ...any) (Result, error)
func (db *DB) QueryContext(ctx context.Context, query string, args ...any) (*Rows, error)
func (db *DB) QueryRowContext(ctx context.Context, query string, args ...any) *Row
```

方法，其中第一个参数就是context.Context。

例如：操作时：

```go
db, _ := sql.Open("", "")
query := "DELETE FROM `table_name` WHERE `id` = ?"
db.ExecContext(context.Background(), query, 42)
```

当然，单独 `database.sql`包中，也支持不传递context.Context的方法。功能一致，但缺失了context.Context相关功能。

```go
func (db *DB) Exec(query string, args ...any) (Result, error)
```

context.Background 和 context.TODO 返回的都是预定义好的 emptyCtx 类型数据，其结构如下：

```go
// 创建方法
func Background() Context {
    return background
}
func TODO() Context {
    return todo
}

// 预定义变量
var (
    background = new(emptyCtx)
    todo       = new(emptyCtx)
)

// emptyCtx 定义
type emptyCtx int

func (*emptyCtx) Deadline() (deadline time.Time, ok bool) {
    return
}

func (*emptyCtx) Done() <-chan struct{} {
    return nil
}

func (*emptyCtx) Err() error {
    return nil
}

func (*emptyCtx) Value(key any) any {
    return nil
}

func (e *emptyCtx) String() string {
    switch e {
    case background:
        return "context.Background"
    case todo:
        return "context.TODO"
    }
    return "unknown empty Context"
}
```

可见，**emptyCtx 是不具备取消、KV值和Deadline的相关功能的，称为空Context**，没有任何功能。

## Context传递取消信号

**`context.WithCancel`** 函数能够从 context.Context 中衍生出一个新的子上下文并返回用于取消该上下文的函数。**一旦我们执行返回的取消函数，当前上下文以及它的子上下文都会被取消，所有的 Goroutine 都会同步收到这一取消信号**。取消操作通常分为**主动取消，定时取消**两类。

### 主动取消

需要的操作为：

- 创建带有cancel函数的Context，

  ```go
  func WithCancel(parent Context) (ctx Context, cancel CancelFunc)
  ```

- 接收cancel的Channel，ctx.Done()；

- **主动取消**：主动Cancel的函数，cancel CancelFunc；

![image.png](https://fynotefile.oss-cn-zhangjiakou.aliyuncs.com/fynote/fyfile/13080/1672136903072/ec8961c6d6e94c3ab5e4da6633704fe6.png)

示例代码：这是cancel的典型用法，3s后主goroutine主动取消context，因此**当前上下文以及它的子上下文都会被取消，只需监控`ctx.Done()`即可**；

```go
func ContextCancelCall() {
    // 1. 创建cancelContext
    ctx, cancel := context.WithCancel(context.Background())

    wg := sync.WaitGroup{}
    wg.Add(4)
    // 2. 启动goroutine，携带cancelCtx
    for i := 0; i < 4; i++ {
        // 启动goroutine，携带ctx参数
        go func(c context.Context, n int) {
            defer wg.Done()
            // 监听context的取消完成channel，来确定是否执行了主动cancel操作
            for {
                select {
                // 子goroutine等待接收c.Done()这个channel
                case <- c.Done():
                    fmt.Println("Cancel")
                    return
                default:
                }
                fmt.Println(strings.Repeat("  ", n), n)
                time.Sleep(300 * time.Millisecond)
            }
        }(ctx, i)
    }

    select {
         // 3. 3s后主动取消context，这是cancel的典型用法。当前上下文以及它的子上下文都会被取消，只需监控ctx.Done()即可；
    case <-time.NewTimer(2 * time.Second).C:
        cancel() // ctx.Done() <- struct{}
    }

    select {
    case <-ctx.Done():
        fmt.Println("main Cancel")
    }
    
    wg.Wait()
}

// ======
> go test -run TestContextCancelCall
       3
   1  
 0  
     2
   1
       3
     2  
 0  
 0
   1  
       3
     2  
     2
   1
       3
 0
 0
   1
       3
     2
     2
   1
 0
       3
       3
 0
   1
     2
main Cancel
Cancel
Cancel
Cancel
Cancel
PASS
ok      goConcurrency   2.219s

```

**当调用cancel()时，全部的goroutine会从 ctx.Done() 接收到内容（即包括子goroutine），进而完成后续控制操作。**

WithCancel函数返回context和一个`CancelFunc`。

```go
func WithCancel(parent Context) (ctx Context, cancel CancelFunc) {
	c := withCancel(parent)
	return c, func() { c.cancel(true, Canceled, nil) }
}
```

withCancel函数：

```go
func withCancel(parent Context) *cancelCtx {
	if parent == nil {
		panic("cannot create context from nil parent")
	}
	c := &cancelCtx{}
	c.propagateCancel(parent, c)
	return c
}
```

withCancel函数返回的Context是 `context.cancelCtx` 结构体对象

其中 `context.cancelCtx` 结构如下：。。。。

```go
// A cancelCtx can be canceled. When canceled, it also cancels any children
// that implement canceler.
type cancelCtx struct {
	Context

	mu       sync.Mutex            // protects following fields
	done     atomic.Value          // of chan struct{}, created lazily, closed by first cancel call
	children map[canceler]struct{} // set to nil by the first cancel call
	err      error                 // set to non-nil by the first cancel call
	cause    error                 // set to non-nil by the first cancel call
}
```

其中：

- Context，上级Context对象
- mu， 互斥锁
- done，用于处理cancel通知信号的channel。懒惰模式创建，调用cancel时关闭。
- children，**以该context为parent的可cancel的context们**
- err，标准化的错误类型（例如被取消或超时）。
- cause，存储更具体的取消原因或上下文相关的错误信息。

### Deadline和Timeout定时取消

与主动调用 CancelFunc 的差异在于，定时取消，增加了一个到时自动取消的机制：

- **Deadline，某个时间点后**，使用 `func WithDeadline(parent Context, d time.Time) (Context, CancelFunc)`创建
- **Timeout，某个时间段后**，使用 `func WithTimeout(parent Context, timeout time.Duration) (Context, CancelFunc)` 创建

示例代码如下，与主动cancel的代码类似：

```go
// 1s后cancel
ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)

// 每天 20:30 cancel
curr := time.Now()
t := time.Date(curr.Year(), curr.Month(), curr.Day(), 20, 30, 0, 0, time.Local)
ctx, cancel := context.WithDeadline(context.Background(), t)
```

其他代码一致，**当时间到时，ctx.Done() 可以接收内容**，进而控制goroutine停止。

不论WithDeadline和WithTimeout都会构建 `*timerCtx` 类型的Context，结构如下：

```go
// A timerCtx carries a timer and a deadline. It embeds a cancelCtx to
// implement Done and Err. It implements cancel by stopping its timer then
// delegating to cancelCtx.cancel.
type timerCtx struct {
   cancelCtx
   timer *time.Timer // Under cancelCtx.mu.

   deadline time.Time
}
```

其中：

- cancelCtx，基于parent构建的cancelCtx
- deadline，cancel时间
- **timer，定时器，用于自动cancel**

### Cancel操作的向下传递

当父上下文被取消时，子上下文也会被取消。Context 结构如下：

```
ctxOne
  |    \
ctxTwo    ctxThree
  |
ctxFour
```

示例代码：

```go
func ContextCancelDeep() {
    ctxOne, cancel := context.WithCancel(context.Background())
    ctxTwo, _ := context.WithCancel(ctxOne)
    ctxThree, _ := context.WithCancel(ctxOne)
    ctxFour, _ := context.WithCancel(ctxTwo)

    // 带有timeout的cancel
    //ctxOne, _ := context.WithTimeout(context.Background(), 1*time.Second)
    //ctxTwo, cancel := context.WithTimeout(ctxOne, 1*time.Second)
    //ctxThree, _ := context.WithTimeout(ctxOne, 1*time.Second)
    //ctxFour, _ := context.WithTimeout(ctxTwo, 1*time.Second)

    cancel()
    wg := sync.WaitGroup{}
    wg.Add(4)
    go func() {
        defer wg.Done()
        select {
        case <-ctxOne.Done():
            fmt.Println("one cancel")
        }
    }()
    go func() {
        defer wg.Done()
        select {
        case <-ctxTwo.Done():
            fmt.Println("two cancel")
        }
    }()
    go func() {
        defer wg.Done()
        select {
        case <-ctxThree.Done():
            fmt.Println("three cancel")
        }
    }()
    go func() {
        defer wg.Done()
        select {
        case <-ctxFour.Done():
            fmt.Println("four cancel")
        }
    }()
    wg.Wait()
}
```

我们调用 ctxOne 的 cancel, 其后续的context都会接收到取消的信号。

如果调用了其他的cancel，例如ctxTwo，那么ctxOne和ctxThree是不会接收到信号的。

### 取消操作流程

#### cancelCtx创建流程

重点是理解cancel传播的流程。

使用 `context.WithCancel`, `context.WithDeadlime`, `context.WithTimeout` 创建cancelCtx或timerCtx的核心过程基本一致，以 `context.WithCancel` 为例：

```go
func WithCancel(parent Context) (ctx Context, cancel CancelFunc) {
    if parent == nil {
        panic("cannot create context from nil parent")
    }
    // 1.构建cancelCtx对象
    c := newCancelCtx(parent)
    // 2.传播Cancel操作
    propagateCancel(parent, &c)
    // 3.返回值，注意第二个cancel函数的实现
    return &c, func() { c.cancel(true, Canceled) }
}

func newCancelCtx(parent Context) cancelCtx {
    return cancelCtx{Context: parent}
}
```

由此可见，核心过程有两个：

1. newCancelCtx， 使用 `parent` 构建 cancelCtx
2. **propagateCancel， cancel传播机制**，用来构建父子Context的关联，用于保证在父级Context取消时可以同步取消子级Context

核心的propagateCancel 的实现如下：

```go
// propagateCancel arranges for child to be canceled when parent is.
func propagateCancel(parent Context, child canceler) {
    // parent不会触发cancel操作
    done := parent.Done()
    if done == nil {
        return // parent is never canceled
    }

    // parent已经触发了cancel操作
    select {
    case <-done:
        // parent is already canceled
        child.cancel(false, parent.Err())
        return
    default:
    }

    // parent还没有触发cancel操作
    if p, ok := parentCancelCtx(parent); ok {
        // 内置cancelCtx类型
        p.mu.Lock()
        if p.err != nil {
            // parent has already been canceled
            child.cancel(false, p.err)
        } else {
            if p.children == nil {
                p.children = make(map[canceler]struct{})
            }
            // 将当前context放入parent.children中
            p.children[child] = struct{}{}
        }
        p.mu.Unlock()
    } else {
        // 非内置cancelCtx类型
        atomic.AddInt32(&goroutines, +1)
        go func() {
            select {
            case <-parent.Done():
                child.cancel(false, parent.Err())
            case <-child.Done():
            }
        }()
    }
}
```

以上代码在建立child和parent的cancelCtx联系时，处理了下面情况：

- parent不会触发cancel操作，不做任何操作，直接返回
- parent已经触发了cancel操作，执行child的cancel操作，返回
- parent还没有触发cancel操作，**`child` 会被加入 `parent` 的 `children` 列表中**，等待 `parent` 释放取消信号
- 如果是自定义Context实现了可用的Done()，那么开启goroutine来监听parent.Done()和child.Done()，同样在parent.Done()时取消child。

**propagateCancel的增强版（新版）：**

```go
// propagateCancel arranges for child to be canceled when parent is.
// It sets the parent context of cancelCtx.
func (c *cancelCtx) propagateCancel(parent Context, child canceler) {
	// 将父 Context 赋值给当前 cancelCtx，以便建立关联。
    c.Context = parent
    // 1.父ctx未取消：检查父 Ctx是否有取消信号（Done() 是否为 nil），如果没有，则父 Context 不会被取消，直接返回。
	done := parent.Done()
	if done == nil {
		return // parent is never canceled
	}
    
    // 2.如果父 Context 已取消，直接取消子 Context。
    // 使用 Cause(parent) 来获取父 Context 取消的原因。
	select {
	case <-done:
		// parent is already canceled
		child.cancel(false, parent.Err(), Cause(parent))
		return
	default:
	}

    // 3.检查父上下文是不是CancelCtx
	if p, ok := parentCancelCtx(parent); ok {
		// parent is a *cancelCtx, or derives from one.
        // 翻译：如果父上下文是 cancelCtx 或派生自它
		p.mu.Lock()
		if p.err != nil {
			// 3.1 如果父已取消（p.err != nil），子直接取消
			child.cancel(false, p.err, p.cause)
		} else {
			if p.children == nil {
				p.children = make(map[canceler]struct{})
			}
            // 3.2 如果父未取消，将子添加到父的 children 列表中，以便后续取消时一起通知子。
			p.children[child] = struct{}{}
		}
		p.mu.Unlock()
		return
	}
    
    // 4.如果父上下文实现了 AfterFunc 方法（如超时或定时取消），设置一个回调函数来取消子 Context。
	if a, ok := parent.(afterFuncer); ok {
		// parent implements an AfterFunc method.
		c.mu.Lock()
        // 设置一个回调函数来取消子 Context。
		stop := a.AfterFunc(func() {
			child.cancel(false, parent.Err(), Cause(parent))
		})
        // 使用 stopCtx 包装父上下文，以便在取消完成后停止回调。
		c.Context = stopCtx{
			Context: parent,
			stop:    stop,
		}
		c.mu.Unlock()
		return
	}

    // 5.如果父 Context 既不是 cancelCtx，也不支持 AfterFunc，则启动一个新的 Goroutine 来监听父 Context 的取消信号。
	goroutines.Add(1)
	go func() {
		select {
		case <-parent.Done():
            // 如果父 Context 被取消，取消子 Context。
			child.cancel(false, parent.Err(), Cause(parent))
		case <-child.Done():
		}
	}()
}
```

`propagateCancel` 主要作用是确保：

- 当父 `Context` 被取消时，子 `Context` 也会相应地被取消。
- 如果父 `Context` 已经被取消，则子 `Context` 直接处理取消逻辑。
- 如果父 `Context` 是 `cancelCtx` 类型或支持其他特性（如 `AfterFunc`），它会合理地利用这些特性。

如果**是WithDeadline构建的timerCtx，构建的过程多了两步**：

- **对截止时间的判定**，判定是否已经截止
- **如果截止时间未过，再设置定时器**

示例代码：

```go
func WithDeadline(parent Context, d time.Time) (Context, CancelFunc) {
    // 确保 parent 不为空，否则直接触发 panic
    if parent == nil {
        panic("cannot create context from nil parent")
    }
    // 如果父上下文已经有一个截止时间，并且在新设置的截止时间之前(说明老截止时间已经到了)
    if cur, ok := parent.Deadline(); ok && cur.Before(d) {
        // The current deadline is already sooner than the new one.
        // 直接创建一个取消上下文
        return WithCancel(parent)
    }
    
    // 否则创建一个timerCtx，这是一个带有截止时间的上下文，它继承自 cancelCtx。
    c := &timerCtx{
        cancelCtx: newCancelCtx(parent),
        deadline:  d,
    }
    // 传播cancel操作
    propagateCancel(parent, c)

    dur := time.Until(d)
    
    // 1.截止时间已经过了，则立即取消上下文，并返回
    if dur <= 0 {
        c.cancel(true, DeadlineExceeded) // deadline has already passed
        return c, func() { c.cancel(false, Canceled) }
    }
    c.mu.Lock()
    defer c.mu.Unlock()
    
    // 2.如果未过截止时间，则设置一个定时器，定时器到期后自动触发取消。
    if c.err == nil {
        // 注意传入了dur作为参数，表示离截至时间还有多久
        // time.AfterFunc用于在指定的时间间隔后执行一个函数。
        c.timer = time.AfterFunc(dur, func() {
            c.cancel(true, DeadlineExceeded)
        })
    }
    return c, func() { c.cancel(true, Canceled) } // 返回上下文和取消函数：
}
```

#### ctx.Done() 操作流程

初始信号channel流程，以 cancelCtx 为例，返回一个 `<-chan struct{}`

```go
func (c *cancelCtx) Done() <-chan struct{} {
    // 1.如果已经存在一个通道，直接返回
    d := c.done.Load()
    if d != nil {
        return d.(chan struct{})
    }
    c.mu.Lock()
    defer c.mu.Unlock()

    // 2.如果通道尚未初始化，创建一个新的通道并存储
    d = c.done.Load()
    if d == nil {
        d = make(chan struct{})
        c.done.Store(d)
    }
    return d.(chan struct{})
}
```

其中两个步骤：

1. 先尝试加载已经存在的
2. 不存在再初始化新的

核心要点是，当调用Done()时，初始化chan struct{}， 而不是在上下文cancelCtx创建时，就初始化完成了。称为懒惰初始化（懒加载）。

**懒加载**（Lazy Loading）是一种设计模式，指的是在**需要时才初始化或加载资源**，而不是在程序启动或某功能初始化时立即加载所有资源。其核心思想是**延迟加载**，避免不必要的计算和资源消耗，从而提高程序的性能和资源利用率。

#### cancel()操作流程

取消流程，我们以 cancelCtx 的主动取消函数cancel的实现为例：通过关闭通道、通知子上下文，以及从父上下文中移除自身，来完成取消操作。

```go
// cancel closes c.done, cancels each of c's children, and, if
// removeFromParent is true, removes c from its parent's children.
func (c *cancelCtx) cancel(removeFromParent bool, err error) {
    // 1.必须you确保取消时必须提供一个错误原因,如果没有提供错误，直接触发 panic。
    if err == nil {
        panic("context: internal error: missing cancel error")
    }
    c.mu.Lock()
    // 检查是否已经取消，如果 c.err 已经被设置，说明上下文已经取消，直接返回，避免重复取消。
    if c.err != nil {
        c.mu.Unlock()
        return // already canceled
    }
    // 设置取消错误err
    c.err = err
    // 2.关闭自己的done通道（done channel只有关闭了才能读取）
    d, _ := c.done.Load().(chan struct{})
    if d == nil {
        // 如果通道尚未初始化，则将其设置为一个已关闭的共享通道 closedchan。
        c.done.Store(closedchan)
    } else {
        // 如果通道已存在，则调用 close(d) 关闭通道。
        close(d)
    }
    
    // 3.取消子上下文：遍历全部的子上下文，调用它们的 cancel 方法，传播cancel信号，全部取消。
    for child := range c.children {
        // NOTE: acquiring the child's lock while holding parent's lock.
        child.cancel(false, err)
    }
    
    // 遍历结束后，将自己的 children 清空
    c.children = nil
    c.mu.Unlock()

    // 4.从父上下文中的children列表中移除自己
    if removeFromParent {
        removeChild(c.Context, c)
    }
}
```

以上流程的核心操作：

- 关闭done channel，用来通知全部使用该ctx的goroutine
- 遍历全部可取消的子context，执行child的取消操作
- 从parent的children删除自己

## Context传值

![image.png](https://fynotefile.oss-cn-zhangjiakou.aliyuncs.com/fynote/fyfile/13080/1672136903072/2fd01066703042469b2869fbdd3c2868.png)

**若希望在使用context时，携带额外的Key-Value数据**，可以使用 `context.WithValue` 方法，构建带有值的context。并使用 `Value(key any) any` 方法获取值。带有值

对应方法的签名如下：

```go
func WithValue(parent Context, key, val any) Context

type Context interface {
    Value(key any) any
}
```

需要三个参数：

- 上级 Context
- **key 要求是comparable的（可比较的）**，实操时，**推荐使用特定的Key类型**，避免直接使用string或其他内置类型而带来package之间的冲突。
- val any

示例代码

```go
type MyContextKey string

func ContextValue() {
    wg := sync.WaitGroup{}

    ctx := context.WithValue(context.Background(), MyContextKey("title"), "Go")

    wg.Add(1)
    go func(c context.Context) {
        defer wg.Done()
        if v := c.Value(MyContextKey("title")); v != nil {
            fmt.Println("found value:", v)
            return
        }
        fmt.Println("key not found:", MyContextKey("title"))
    }(ctx)

    wg.Wait()
}
```

使用`MyContextKey`（一个自定义类型）代替普通字符串，确保键不会意外与其他使用 `context` 的代码发生冲突。（避免键的意外冲突）

`context.WithValue` 方法返回 `context.valueCtx` 结构体类型。`context.valueCtx` 结构体**包含了上级Context和key、value**：

```go
// A valueCtx carries a key-value pair. It implements Value for that key and
// delegates all other calls to the embedded Context.
type valueCtx struct {
    Context
    key, val any
}


func (c *valueCtx) Value(key any) any {
    if c.key == key {
        return c.val
    }
    return value(c.Context, key)
}
```

也就是除了 value 功能，其他Contenxt功能都由parent Context实现。

如果 [`context.valueCtx.Value`](https://draveness.me/golang/tree/context.valueCtx.Value) 方法查询的 key 不存在于当前 valueCtx 中，就会**从父上下文中查找该键对应的值**直到某个父上下文中返回 `nil` 或者查找到对应的值。

**`context.valueCtx.Value` 的工作机制**：

1. **在当前上下文中查找**：
   - 如果当前 `valueCtx` 包含与 `key` 匹配的值，则直接返回这个值。

2. **从父上下文中查找**：
   - 如果当前上下文中没有匹配的 `key`，则会向上递归查询父上下文，直到：
     - 找到对应的 `key`，返回其值。
     - 或者遍历到 `context.Background()` 或 `context.TODO()` 这样的根上下文，返回 `nil`，表示未找到对应的值。

例如：

```go
func ContextValueDeep() {
    wgOne := sync.WaitGroup{}

    ctxOne := context.WithValue(context.Background(), MyContextKey("title"), "One")
    ctxTwo := context.WithValue(ctxOne, MyContextKey("key"), "Value")
    ctxThree := context.WithValue(ctxTwo, MyContextKey("key"), "Value")

    wgOne.Add(1)
    go func(c context.Context) {
        defer wgOne.Done()
        if v := c.Value(MyContextKey("title")); v != nil {
            fmt.Println("found value:", v)
            return
        }
        fmt.Println("key not found:", MyContextKey("title"))
    }(ctxThree)

    wgOne.Wait()
}
```

## 小结

特定的结构体类型：

- emptyCtx，函数 context.Background, context.TODO
- cancelCtx，函数 context.WithCancel
- timerCtx， 函数 context.WithDeadline, context.WithTimeout
- valueCtx，函数 context.WithValue

官方博客对Context使用的建议：

- 直接将 Context 类型作为函数的第一参数，而且一般都命名为 ctx。
- 如果你实在不知道传什么，标准库给你准备好了一个  context.TODO。
- context 存储的应该是一些goroutine共同的数据。
- **context 是并发安全的（例如cancelCtx）。**
