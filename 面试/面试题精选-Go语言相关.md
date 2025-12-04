# Go语言相关面试题

## 结构体字段tag的作用

**tag可以为结构体字段提供属性**。

特性：

* tag的声明格式为跟随字段的字符串字面量
* tag字符串可以是任意格式
* tag需要通过反射解析，获取字符串后，根据特定格式解析
* 通用格式：`key1:"value1" key2:"value2"`

常见的tag示例：

* json: json序列化或反序列化时字段的名称
* gorm: gorm模型定义属性配置
* form: gin框架中对应的前端的数据字段名

示例：

```go
type User struct {
	Name  string `json:"showName" form:"" gorm:""`
	Email string `json:"email"`
}
```

## 如何解析结构体字段tag

核心类型：reflect.StructTag

```go
// 字段tag的类型定义
type StructTag string

// 获取tag中，对应的key的值
func (tag StructTag) Get(key string) string
// 查找，当key不存在时，ok返回false
func (tag StructTag) Lookup(key string) (value string, ok bool)
```

示例：

```go
func GetFieldTag() {
	type User struct {
		Name  string `json:"showName" form:"userName" gorm:"realName"`
		Email string `json:"email"`
	}

	// reflect
	user := User{}
	userType := reflect.TypeOf(user)
	// 获取字段的tag，全部字段的tag字符串
	for i, l := 0, userType.NumField(); i < l; i++ {
		fieldType := userType.Field(i)
		fmt.Println(fieldType.Tag)
	}

	// 获取tag中，具体的key
	nameFieldType, exists := userType.FieldByName("Name")
	if !exists {
		log.Fatal("Name field not exists")
	}
	fmt.Println(nameFieldType.Tag.Get("json"))
	fmt.Println(nameFieldType.Tag.Get("form"))
	fmt.Println(nameFieldType.Tag.Get("gorm"))
}

```

```
> go test -run GetFieldTag
json:"showName" form:"userName" gorm:"realName"
json:"email"
showName  
userName  
realName  
PASS
ok      question        0.020s

```

## init函数何时执行？

**简单答案：在main.main函数执行前执行，在var声明变量后执行。**

**init函数的作用**：

完成应用程序初始化操作。

特点：

* init函数是可选的
* init函数可以定义多个
* init函数自动执行
* 首字母小写，一定是 init
* 函数签名：`func init()`
  * 没有参数
  * 没有返回值

**执行顺序：**

在main.main()函数前执行。

当存在多个init函数时，执行顺序是什么？

分成：

* 单个源文件：依据定义顺序执行
* 单个包：依据**包中文件名的字典顺序**，依次执行其中的init函数
* 多个包不存在导入关系：依据import导入优化后的顺序执行，也就是**包名字典排序顺序**
* 存在导入关系的多个包：依据导入依赖关系，优先执行**被导入**包的init函数。例如：
  * main import a， 那么a包init先于main包的init执行
  * main import a, a import b, 那就是：b.init , a.init, main.init

测试：

单个源文件

```go
// main.go
package main

import "fmt"

func init() {
	fmt.Println("a")
}

func init() {
	fmt.Println("b")
}

func init() {
	fmt.Println("c")
}

func main() {
	fmt.Println("main")
}

```

```
> go run .\main.go
a
b
c
main
```

单个包：

```go
// main.go
package main

import "fmt"

func main() {
	fmt.Println("main")
}


// a.go
package main

import "fmt"

func init() {
	fmt.Println("z")
}

func init() {
	fmt.Println("zz")
}

func init() {
	fmt.Println("zzz")
}

// b.go
package main

import "fmt"

func init() {
	fmt.Println("y")
}


// c.go
package main

import "fmt"

func init() {
	fmt.Println("x")
}

```

```
> go run .
z
zz
zzz
y
x
main

```

多个包不存在导入关系：

main import a

main import b

main import c

```go
// main.go
package main

import "fmt"
import _ "question/init/multi/b"
import _ "question/init/multi/c"
import _ "question/init/multi/a"

func main() {
	fmt.Println("main")
}


// a/init.go
package a

import "fmt"

func init() {
	fmt.Println("a")
}

// b/init.go
package b

import "fmt"

func init() {
	fmt.Println("b")
}


// c/init.go
package c

import "fmt"

func init() {
	fmt.Println("c")
}

```

测试时，任意调整main中的导入顺序：

```
> go run .
a
b
c
main

```

存在导入关系的多个包：

在以上源码的main包中，增加init

```go
// main.go
package main

import "fmt"
import _ "question/init/multi/b"
import _ "question/init/multi/c"
import _ "question/init/multi/a"

func main() {
	fmt.Println("main")
}

func init() {
	fmt.Println("init main")
}

```

```
> go run .
a
b  
c  
init main
main

```

## 如何判断 map 中是否包含某个 key

map类型的下标语法的返回值，支持两种形式：

* 单个值，只返回下标key对应的value
* 两个值，返回下标key对应的value，及该key是否存在于map中

因此可以通过第二个返回值来判断map中的key是否存在。

示例代码：

```go
func TestMapKeyExists(t *testing.T) {
	m := map[string]int{
		"go":   42,
		"java": 365,
	}

	if v, exists := m["cpp"]; exists {
		log.Println("key exists. value is ", v)
	}
}
```

## `=` 和 `:=` 的区别

两个不同的操作符：

* `=` 赋值运算符，为**已经存在的变量**赋值
* `:=` **短声明运算符**，类比var来说，省略var关键字，通过**类型推断**的方式，来声明变量。左值变量列表一定要存在未声明的变量

示例：

```go
package question

import "testing"

func TestVarShort(t *testing.T) {
	v1, v2 = 10, 20  //Unresolved reference 'v1'
	v3, v4 := 10, 20 // Unused variable 'v3'
}

```

## 使用过context吗？context 有哪些使用场景？

context 标准包。

用过。

context：在goroutine间，传递信息。

```go
func() {
	ctx := context.Backgroud()
	go func(context.Context) {}(ctx)
}()
```

context 有哪些使用场景?(在goroutine间传递信息，具体是哪种信息)

* **取消信号**
  * 主动取消
  * 超时取消
  * 时间截至取消
* **值的信息**

示例：

```go
func TestContextDemo(t *testing.T) {
	// 主动取消
	wg := sync.WaitGroup{}
	// 1. 构建context
	ctx, cancel := context.WithCancel(context.Background())

	wg.Add(1)
	go func(ctx context.Context) {
		defer wg.Done()
		for {
			select {
			case <-ctx.Done():
				log.Println("received cancel")
				return
			default:
				log.Println("inner goroutine")
				time.Sleep(300 * time.Millisecond)
			}

		}

	}(ctx)

	// 2. 外层终止，将信号传递到其他goroutine
	time.Sleep(2 * time.Second)
	cancel() // 取消信号
	wg.Wait()
}

```

## 如何在Go语言中实现类型断言？

强类型。

**接口类型**，可以存储（表示）任意实现该接口的类型数据。空接口类型interface{}, any，可以存储任意类型的数据。

类型断言，用于检查和获取接口中实际存储的数据类型的值。

```go
// 断言, 将i断言为T类型的数据
v := i.(T)
v, ok := i.(T)
```

另外的**类型switch语法**，也可以完成类型的判断：

```go
switch i.type {
	case string:
	case int:
	default:
}
```

测试：

```go
func TestTypeAssert(t *testing.T) {
	var interfaceVar any
	intVar := 32
	interfaceVar = intVar
	v := interfaceVar.(int)
	log.Println(v)

	stringVar := "GoLang"
	interfaceVar = stringVar
	vs, ok := interfaceVar.(string)
	log.Println(vs, ok)

	switch v := interfaceVar.(type) {
	case int:
		log.Println("int", v)
	case string:
		log.Println("string", v)
	default:
		log.Println("no match any type")
	}
}
```

## 如何在Go语言中实现一个接口

实现。implement。

* 某类型实现了某接口声明的全部方法，即认为该类型实现了该接口
* 不需要（没有）类型和接口间通过类似implement关键字来强行绑定
* method set，接口声明的方法，声明了某种类型的方法集

```
// 接口类型
interface
// 定义类型，实现接口通常是struct结构体类型
type 

```

示例：

```go
// 定义接口
type Study interface {
	DoHomework() (bool, error)
	HaveWeekend() error
}

// 实现接口
type Student struct {
}

// 实现接口，定义Student的方法集合
func (Student) DoHomework() (bool, error) {
	return false, nil
}

func (Student) HaveWeekend() error {
	return nil
}
func (Student) Sleep() error {
	return nil
}

func DoStudy(user Study) {

}

func TestImplementInterface(t *testing.T) {
	s := Student{}
	DoStudy(s)
}
```

## Go的并发原语有哪些

原语，primitive。

```
goroutine // 协程
channel // 通信
sync // 同步控制
sync.WaitGroup // 用于等待一组goroutine完成
sync.Cond // 条件变量，用于goroutine之间的事件通知
sync.Once // 确保函数只执行一次，常用于单例模式
sync.Mutex // 互斥锁
sync/atomic // 原子操作
context // goroutine间信息传递

```

gmp调度？

## Go的错误处理机制

特点：

* **不使用异常式**的处理机制，把错误和功能拆分开
* **推荐显式的错误**处理机制，体现就是错误应该是函数的返回值的一部分，**将错误和函数功能合并**
* 认为错误应该被检查和处理，而不是被抛出和捕获

```go
func() (result any, err error)

if result, err := func(); err != nil {
	// 函数执行没有错误
}
```

## Go语言中的切片和数组的区别

切片和数组都是典型的**列表式结构**，逻辑上连续的一段类型的相同的存储数据。

差异比较：

| 特性           | 切片slice                 | 数组array                            |
| -------------- | ------------------------- | ------------------------------------ |
| 长度           | 元素的个数                | 元素的个数                           |
| **容量** | **支持扩容**        | **不支持扩容，没有容量的概念** |
| 类型字面量     | []T                       | [len]T                               |
| 引用类型       | 是                        | 否（数组是**值类型**）                  |
| 构造           | 字面量，make(T, len, cap) | 字面量                               |

```
s1 := []int{10, 20, 30}
s2 = s1
s1[0] = 100
s1[0] == s2[0] == 100
```

## 简要概述Go语言的并发模型是什么（重要）

**并发模型的核心基础：**

Goroutine是Go语言实现的协程调度。

协程 coroutine，可以中断执行的函数，称为协程。

```
go func1()
go func2()
```

**程序开发的角度，并发模型：**

* goroutine, 可以独立并发执行的程序单元，语法上就是go 调用的函数。`go func()`
  * go 关键字，让函数以独立的goroutine来运行
  
  * GMP的调度模型
    * M对应操作系统的线程
    
    * G独立的Goroutine
    
    * P用来将具体的某个G（goroutine），绑定到某个M上执行。每个P有本地的G队列。整体还有全局的G队列
    
      P（processor）代表**处理器。**
    
  * go 语言调度器自己实现的GMP的调度模型，用户层面（不是系统层面）的调度
  
* channel，goroutine间的通信
  * 通道（信道）
  * CSP，通信顺序进程。channel实现的通信模式，就是CSP。并发实体（goroutine)通过共享的通信管道(channel)完成通信。
  * 数据通信，缓冲channel，make(chan, 3)
  * 信号通信，非缓存channel， make(chan struct{})
  * 操作
    * <-chan, 接收，receive，读取
    * chan<-，发送，send，写入
  
* context：用于在goroutine间传递信息，主要完成取消信号的传递

* sync，同步控制包

## Go语言中map类型的使用

map映射类型。**键值对集合**类型，内部是基于HashTable实现。

关键字map来构建map类型，需要提供键类型和值类型：

```go
// map类型的字面量
map[keyType]ValueType
```

keyType，必须是可比较的。

```go
// comparable is an interface that is implemented by all comparable types
// (booleans, numbers, strings, pointers, channels, arrays of comparable types,
// structs whose fields are all comparable types).
// The comparable interface may only be used as a type parameter constraint,
// not as the type of a variable.
type comparable interface{ comparable }
```

valueType，任意。

map支持的操作，效率极高：

* [key] 操作
  * 添加，更新
  * 查找，支持第二个返回值判断key是否存在
* delete()内置函数，删除key

**map类型不是并发安全的**。并发执行的goroutine操作同一个map，导致数据不可期待。

应该使用加锁，或者channel的方案来保证并发安全。或者使用**并发安全的map类型**，例如sync.Map。

## Go语言中的defer语句的作用

defer，**延迟执行**。将函数的调用延迟到defer所在函数结束前执行。

```
defer funcName()
```

当存在多个defer，栈顺序，先进后出的顺序。先defer的后执行。

```go
func TestDeferCall(t *testing.T) {
	defer func() {
		log.Println("1 defer")
	}()
	defer func() {
		log.Println("2 defer")
	}()
	defer func() {
		log.Println("3 defer")
	}()
}
```

```
> go test -run Defer
2024/07/20 16:26:29 3 defer
2024/07/20 16:26:29 2 defer   
2024/07/20 16:26:29 1 defer   
PASS
```

defer 和 return 的执行顺序？

return语句涉及的表达式要执行eval，然后去执行defer函数，函数返回。**不要期望在defer中影响返回值**（影响核心业务），这是说defer中对值的操作和最终的return 值无关。

示例：

```go
func TestDeferCall2(t *testing.T) {
	log.Println(deferReturn())
}

func deferReturn() int {
	v := 0
	defer func() { // 2. 执行defer函数
		v = v + 1
	}()
	return v + 1 // 1. 执行v+1 = 1
	// 3. 执行 return 1
}

```

```
> go test -run DeferCall2
2024/07/20 16:29:37 1   
PASS 
```

**defer的作用：释放函数占用的资源**

```go
defer file.Close()
defer lock.Unlock()
defer wg.Done()
```

## 如何并发安全的使用map

map类型不是并发安全的。

示例：

```go
func TestMapSafe(t *testing.T) {
	m := map[string]int{
		"counter": 0,
	}
	wg := sync.WaitGroup{}
	// 两个（或多个）goroutine同时操作m
	wg.Add(3)
	go func() {
		defer wg.Done()
		for range 1000 {
			m["counter"]++
		}
	}()
	go func() {
		defer wg.Done()
		for range 1000 {
			m["counter"]++
		}
	}()
	go func() {
		defer wg.Done()
		for range 1000 {
			m["counter"]++
		}
	}()

	wg.Wait()
	log.Println(m, m["counter"])
}
```

```
> go test -run MapSafe
fatal error: concurrent map writes
```

如何处理？

* 使用内置的安全map结构
* 手动加锁

示例：

```go
func TestMapSafe1(t *testing.T) {
	m := map[string]int{
		"counter": 0,
	}
	wg := sync.WaitGroup{}
	// 两个（或多个）goroutine同时操作m
	wg.Add(2)
	// 手动加锁
	lock := sync.Mutex{}
	go func() {
		defer wg.Done()
		for range 100000 {
			lock.Lock()
			m["counter"]++
			lock.Unlock()
		}
	}()
	go func() {
		defer wg.Done()
		for range 100000 {
			lock.Lock()
			m["counter"]++
			lock.Unlock()
		}
	}()
	wg.Wait()
	log.Println(m, m["counter"])
}
```

## 切片Slice的扩容策略是什么

扩容，切片的容量增加。

扩容时机：切片追加元素时（append())，若**容量不足以存储全部元素，则会发生扩容。**

扩容策略，指的是每次增加多少容量的策略：

* 容量大小
* 内存对齐

核心扩容策略：

* 若期望容量大于当前容量的两倍，则使用期望容量
* 若当前容量小于1024，通过**翻倍**的方式扩容
* 容当前容量大于1024，每次**增加当前容量的1/4**来扩容，直到满足需求

新的包方法可以主动扩容：

```go
package slices
func Grow[S ~[]E, E any](s S, n int) S {
```

## 在Go语言中如何交叉编译

交叉编译：编译**不同平台**（操作系统）的应用程序。

```
go build
```

通过设置Go的**环境变量**来实现：

* CGO_ENABLED, 是否启用CGO。是否使用C语言版本的Go编译器。**默认为1**，会允许在Go中调用C代码。
* GOOS，目标操作系统，例如 Mac，Linux，Win
* GOARCH，目标系统架构，32bit, 64bit, arm

> CGO 是 Go 语言与 C 语言互操作的工具，通过 CGO，Go 程序可以调用 C 语言代码或 C 库函数，扩展 Go 的功能。
>
> 设置 `CGO_ENABLED=0` 可以确保生成的 Go 二进制文件是纯静态的，这样编译后的程序可以在目标平台上运行，而不依赖于动态链接的 C 库。
>
> - `CGO_ENABLED=1`（启用 CGO，默认设置）
>   - 在支持 CGO 的环境中，Go 编译器可以编译包含 C 语言代码或使用 C 库的 Go 程序。
>   - 编译出的二进制文件可能会**动态链接**到系统的 C 语言运行时或其他 C 库。
>
> - **`CGO_ENABLED=0`（禁用 CGO）**
>   - 禁用 CGO，编译器生成的二进制文件是纯 Go 的静态链接文件，不依赖于任何 C 语言运行时或库。

示例：

win下编译：

```
# Mac 64 位系统
SET CGO_ENABLED=0
SET GOOS=darwin
SET GOARCH=amd64
go build

# Linux 通用 64 位系统(Centos、Ubuntu等)
SET CGO_ENABLED=0
SET GOOS=linux
SET GOARCH=amd64
go build

```

mac下编译：

```
# Linux 通用 64 位系统(Centos、Ubuntu等)
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build

# Windows 64 位系统
CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build 

```

Linux下编译：

```
# Mac 64 位系统
CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 go build 

# Windows 64 位系统
CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build 

```

## Go语言的泛型及用途是什么

泛型：generic type

Go语法层面，成为：类型参数，type parameter。**数据类型作为参数来使用。**

go1.18后，引入泛型。用于处理不同类型的相同逻辑的工作。

示例，go的标注库slices:

```go
func Contains[S ~[]E, E comparable](s S, v E) bool
```

泛型定义：允许程序中的函数、类型、接口中使用类型参数，来增加对**多个类型（相同结构）**的支持。

典型的用途：

* 简化相同结构的不同类型的编码操作
* 减少接口，反射，类型断言等操作
* 约束类型范围。相对interface{}，**可以限制类型范围**

## 在Go语言中如何编写和运行测试

```
go test
```

命令，go中完成测试的命令。

支持：

1. 单元测试，unit test
2. 基准测试，benchmark test
3. Fuzz测试，覆盖测试

编码上，通过 _test.go 后缀的go文件。编写测试函数即可：

```
func TestFunc(*testing.T)
```

测试函数的特点：

* Test测试函数名前缀
* 测试参数
  * *testing.T，单元测试
  * *testing.B，基准测试
  * *testing.F，覆盖测试
* _test.go 测试文件名后缀

## 解释Go语言中的闭包

闭包：Closure。闭包一种**语法现象。**

指的是在**内层作用域的函数，使用外层函数作用域的变量**，当外层函数运行结束后，内部函数在调用时，还可以操作到外层函数作用域的变量，该现象成为闭包。

* 函数**作用域的嵌套**，语法上涉及到**匿名函数。**
* 变量是外层函数的局部变量
* 内部匿名函数中使用该局部变量
* 外层函数执行完毕后，调用内层函数

示例：

```go
func Outer() (func(), func() int) {
	counter := 0

	incr := func() {
		counter++
	}
	get := func() int {
		return counter
	}
	return incr, get
}

func TestClosure(t *testing.T) {
	incr, get := Outer() // counter = 0，注意这里只是返回2个函数，并不执行函数内的操作。
	incr() // counter = 1
	incr() // counter = 2
	log.Println(get()) // 输出2
	incr() // counter = 3
	log.Println(get()) // 输出3
}
```

```
> go test -run Closure
2024/07/22 11:12:25 2
2024/07/22 11:12:25 3
PASS
ok      question        0.071s

```

以上例子就是两个内部incr,get和外部函数中的变量counter形成了2个闭包。

## 如何在Go语言中封装错误

自定义错误。

为什么要自定义错误？go中内置的error结构很简单，看error接口：

```go
type error interface {
	Error() string
}
```

同时，errors包提供的方法，功能也简单。

error功能不足时，可以通过自定义error来扩展错误的功能。

代码上，**通过定义类型，实现error接口，即可完成自定义错误。**

示例：

```go
type CustomError struct {
	Code    string
	Message string
}

func (e *CustomError) Error() string {
	return e.Code + " " + e.Message
}

func NewCustomError(code, msg string) *CustomError {
	return &CustomError{Code: code, Message: msg}
}

func DoSomeError() *CustomError {
	return NewCustomError("404", "not found")
}

func TestCustomError(t *testing.T) {
	if err := DoSomeError(); err != nil {
		log.Println(err)
	}
}
```

```
> go test -run CustomError
2024/07/22 20:57:13 404 not found
PASS
ok      question        0.051s

```

## 解释Go语言中的空接口的使用方法

interface{}是不包含任何函数声明（方法集合）的接口，方法集合为空的接口。

意味着**全部类型都实现了该接口**。使用该类型表示任意类型。

示例：使用 `接口类型.(type)` 来匹配类型，便于后续操作。

```go
type HeroX struct {
	Name string
}

func TestEmptyInterface(t *testing.T) {
	saySomeInterface("golang")
	saySomeInterface(42)
	saySomeInterface(HeroX{Name: "金钢狼"})

}

func saySomeInterface(content interface{}) {
	//log.Println(content.(int) + 20) // Invalid operation: content + 20 (mismatched types interface{} and untyped int)
	switch v := content.(type) { 
	case string:
		log.Println("Hi " + v)
	case int:
		log.Println(v + 20)
	case HeroX:
		log.Println("I am ", v.Name)
	}
}
```

```
> go test -run EmptyInterface
2024/07/22 21:08:40 Hi golang
2024/07/22 21:08:40 62
2024/07/22 21:08:40 I am  金钢狼
PASS
ok      question        0.045s
```

当前 any 是 interface{}的别名：

```
type any = interface{}
```

劣势：需要大量的断言或type switch语法来识别具体的数据类型。

很多地方都可以用类型参数，**泛型来代替空接口的使用。**

## Go语言中字符串是如何实现的

string类型，字符的集合。

数据结构，由两部分构成：

1. 底层存储：一段**连续的存储空间**。Go中该存储空间是**只读的。**
2. 上层字符串头信息。

如图：

![image.png](https://fynotefile.oss-cn-zhangjiakou.aliyuncs.com/fynote/fyfile/13080/1721653091029/e368d74d0dc14a739c63952298236485.png)

StringHeader的源码：

```go
type StringHeader struct {
	Data uintptr
	Len  int
}
```

Go中该存储空间是只读的，go的字符串是只读的。

通常与[]byte可以进行转换。

## 如何在Go语言中使用环境变量

函数：

```go
package os
func Getenv(key string) string
func Setenv(key, value string) error
```

示例：

```go
func TestUseEnv(t *testing.T) {
	log.Println(os.Getenv("path"))
	//os.Setenv("", "")
}
```

```
> go test -run UseEnv
2024/07/22 21:23:47 D:\devel\go\bin;D:\devel\go\bin;D:\apps\gopath\bin;D:\devel\go\bin;D:\apps\gopath\bin;C:\Win
dows\system32;C:\Windows;C:\Windows\System32\Wbem;C:\Windows\System32\WindowsPowerShell\v1.0\;C:\Windows\System3
2\OpenSSH\;C:\Program Files\Intel\WiFi\bin\;C:\Program Files\Common Files\Intel\WirelessCommon\;C:\Program Files
\Docker\Docker\resources\bin;D:\devel\Git\cmd;D:\apps\gopath\bin;D:\devel\go\bin;D:\devel\nodejs\;D:\devel\proto
c-25.3-win64\bin;C:\Users\54009\AppData\Local\Microsoft\WindowsApps;C:\Users\54009\go\bin;C:\Users\54009\AppData
\Roaming\npm;C:\Program Files\JetBrains\GoLand 2023.3.4\bin;;C:\Users\54009\AppData\Local\Programs\Microsoft VS 
Code\bin
PASS
ok      question        0.074s
```

## 解释Go语言中的原子操作及其用途

原子操作由：sync/atomic 包提供支持。

原子操作：一组不可分割的操作，**在并发编程中可以安全的执行**，主要是操作共享资源（变量）。

示例：

```go
func TestAtomicCounter(t *testing.T) {
	wg := sync.WaitGroup{}

	var counter int64
	wg.Add(1)
	go func() {
		defer wg.Done()
		for range 10000 {
			//counter++
			atomic.AddInt64(&counter, 1)
		}
	}()
	wg.Add(1)
	go func() {
		defer wg.Done()
		for range 10000 {
			//counter++
			atomic.AddInt64(&counter, 1)
		}
	}()
	wg.Wait()

	log.Println(counter)
}
```

```
> go test -run AtomicCounter
2024/07/22 21:37:54 20000
PASS
ok      question        0.047s

```

使用++测试，不会是20000，因为++不是原子操作，执行的过程会被打断。

## Go语言中什么是嵌入类型

嵌入类型，嵌入字段，带有嵌入字段的结构体类型。

核心：embeded field。

**嵌入字段：**在声明struct类型时，某个字段如果仅给出了类型但未提供字段名，该字段称为嵌入字段。

**嵌入字段的字段名是类型的名字**，unqualified name， 就是类型的最后的一部分。例如嵌入的字段类型为：package.Type, 此时名字是Type。

> 一个更严谨的解释是：如果一个字段只含有字段类型而没有指定字段的名字，那么这个字段就是一个嵌入类型字段。

引入嵌入字段的目的：

* 将其他类型的功能组合到当前结构体类型中。（类似于其他语言的继承）
* 可以直接使用嵌入类型的字段，不通过嵌入字段名来使用。称为**字段提升**，promoted。

注意事项：

* 允许多个嵌入
* 需要注意嵌入类型的**字段名冲突**
* T和 `*T`的方法集合。T就是嵌入字段的结构体类型。
  * 如果T嵌入的是 `*F`，那么T和 `*T` 都包含 F和 `*F`的方法集
  * 如果T嵌入的是 `F`，那么T和 `*T` 都包含F方法集，但是只有`*T`包含 `*F`的方法集合，`T`不包含`*F`的方法集。

## go中编译标签的用途是什么

> 编译标签（Build Tags）是 Go 编程语言中的一种机制，用于**控制代码在编译时是否被包含在最终的程序中。**它们常用于条件编译，帮助开发者根据目标操作系统、架构或其他条件选择不同的代码逻辑。
>
> **基本语法**
>
> - 编译标签写在文件的头部，紧跟在注释 `//` 后面，格式为：
>
>   ```
>   // +build <tag1> [tag2 ...]
>   ```
>
> - 示例：
>
>   ```
>   // +build linux darwin
>   ```
>
>   这表示该文件只会在 `linux` 或 `darwin`（macOS）平台上被编译。
>
> ------
>
> **规则和约束**
>
> 1. **位置**：
>
>    - 编译标签必须在文件的最顶部，且与 `package` 声明之间不能有任何空行。
>
>    - 可以与版权声明或其他注释共存：
>
>      ```
>      // +build linux
>      // This file is for Linux only.
>      package main
>      ```
>
> 2. **逻辑运算符（自定义标签）**：
>
>    - 编译标签之间可以使用布尔逻辑：
>
>      - **空格（逻辑或 OR）**：`linux darwin` 表示 Linux 或 macOS。
>      - **逗号（逻辑与 AND）**：`linux,amd64` 表示 Linux 且是 AMD64 架构。
>
>    - 布尔表达式支持加括号组合：
>
>      ```
>      // +build (linux darwin) !windows
>      ```
>
> 3. **否定符**：
>
>    - 使用 `!` 表示非，例如 `!windows` 表示非 Windows 系统。
>
> ------
>
> **常见编译标签**
>
> 1. **预定义标签**（内置）：
>
>    - **操作系统**：
>
>      - `linux`、`windows`、`darwin`、`freebsd`、`netbsd`、`openbsd`、`solaris`、`aix`、`android` 等。
>
>    - **架构**：
>
>      - `amd64`、`386`、`arm`、`arm64`、`ppc64`、`ppc64le`、`mips`、`mipsle` 等。
>
>    - 示例：
>
>      ```
>      // +build linux,amd64
>      ```
>
>      仅在 Linux 且架构为 AMD64 时编译。
>
> 2. **自定义标签**：
>
>    - 开发者可以通过命令行传递自定义标签。
>
>    - 示例：
>
>      ```
>      // +build mytag
>      ```
>
>      在编译时通过 `-tags` 参数启用：
>
>      ```
>      go build -tags mytag
>      ```
>
> 

编译标签：build tag。

编译标签写在go语言源文件的第一行，用来提供编译的一些配置信息。指导go build工具工作的。

```go
// +build dev, prod
```

```bash
go build -tags "dev"
```

典型的例子，部分源文件在开发环境使用，部分源文件在生产环境使用。

通过提供 dev 和 prod（都是自定义的）编译标签，在go build -tags "" 匹配，进而实现**选择特定的源码文件进行编译**的目的。

编译标签的逻辑关系：

* tag1 tag2， 或 ||
* tag1, tag2，与 &&
* !tag1，非

> 总结：编译标签分为 **预定义标签** 和 **自定义标签**：
>
> 1.**预定义标签（如操作系统、架构）**：
>
> - 如果是使用 Go 内置的 **预定义标签**（如 `linux`, `windows`, `amd64` 等），编译器会自动根据目标平台和架构选择对应的文件，**无需**显式使用 `go build -tags` 命令。
>
> 2.**自定义标签**：
>
> - 如果是使用自定义的标签（如 `mytag`），则需要通过 `go build -tags` 显式启用，否则这些文件将被忽略。

## 在Go语言中如何进行编译优化

通过go build的参数，来控制go build的部分行为，进行优化go build操作。

典型的选项：

* -gcflags 与内存管理的相关的参数

  > 用于控制 **Go编译器（gc，Go Compiler）** 的行为，主要作用在编译阶段。

* -ldflags 与程序连接的相关参数

  > 用于调整 **链接器（linker）** 的行为，主要作用在链接阶段。

```
go build -gcflags="" -ldflags=""
```

**gcflags规则**

```
-gcflags="pattern=arg list"
```

pattern，匹配哪些包（module）需要控制

* main main所在的顶级包，自己写的应用程序
* all 全部用到的包
* std

arg list 具体的参数

* -N 禁止优化
* -l 关闭内联
* -c N 编译过程的并发数量
* -o 编译的优化级别

例如：

```
-gcflags="all=-N -l"
```

**ldflags规则**

* -w 不生成DWARF信息，调试信息的一种方案
* -s 关闭符号表

典型的：

```
-ldflags="-w -s"
```

**使用场景对比**

| 参数       | 阶段     | 作用             | 常用场景                                     |
| ---------- | -------- | ---------------- | -------------------------------------------- |
| `-gcflags` | 编译阶段 | 调整代码生成行为 | 调试代码（关闭优化、查看逃逸分析）           |
| `-ldflags` | 链接阶段 | 调整链接器行为   | 设置版本信息、减小二进制文件大小、静态链接等 |

## Go的运行时 (Runtime)提供了哪些功能

运行时（runtime），在Go程序的运行时期内，对应用程序进行管理的组件。

主要功能：

* 垃圾回收，GC
* 内存分配
* 并发调度
* 信号处理
* 调式分析工具
  * pprof
  * trace

## 什么是类型常量和无类型常量

类型常量，typed constant，确定了具体类型的常量

无类型常量，untyped constant，未确定具体类型的常量

常量不同于变量，声明时必须准确确定类型。

例如：

```
42
```

是什么类型？

```
const C = 42
```

C 又是什么类型？

C常量就是未确定类型常量。对应：

```
const D uint = 42
```

D就是类型常量。

使用差异：

* 类型常量，只能在匹配具体类型的环境中使用
* 无类型常量，可以用在满足字面量需求的环境中

```
func Name1(p uint)
func Name2(p int)
```

C 可以用在两个函数中，D 只可以用在Name1中。

## Go有哪些方式可以安全地共享变量

共享变量，在并发中，在 goroutine 中共享变量。

> Go中的名言，期望以通信的方式共享内存，而不期望以共享内存的方式进行通信。
>
> goroutine中共享数据，**强烈推荐channel**，尽量不要使用共享变量的方案。

如果使用共享变量，可以安全处理的话：

* 加锁，sync.Mutex, sync.RWMutex
* 原子操作，sync/atomic
* **使用channel**（数据放在channel中，可获取，可写入）

## Go语言中的select语句如何使用

select 多路复用语句。多路指的是**多个channel通路。**

select与channel配合使用的语句，用于监听具体某个channel上发生操作，进而进行处理。

语法结构与switch类似。

```go
select {
case:
case:
default:
}

for {
	select {
	}
}
```

与switch不同的执行流程是，switch是按照上下顺序进行case判断的，而select不依据语法顺序，是哪个case上优先出现了channel的send或receive操作，先执行哪个case。

* 若同时触发的多个case，那么随机选择其中一个case执行
* 若某时刻没有任何的case触发，那么如果有default就执行default分支，此时select称为非阻塞select，没有default时select语句会阻塞，直到某个case被触发
* 不论执行了任何的分支（case和default），select语句都会执行完毕。
* **select通常在for结构中**，进行持续监控channel的行为。

示例：

```go
func TestSelectStmt(t *testing.T) {
	c := make(chan int, 1)
	for range 10 {
		// select
		select {
		case c <- 1:
		case c <- 2:
		case c <- 3:
		case c <- 4:
		case c <- 5:
		case c <- 6:
		}
		log.Println(<-c)
	}
}
```

```
> go test -run SelectStmt
2024/07/24 10:40:25 2
2024/07/24 10:40:25 6
2024/07/24 10:40:25 4
2024/07/24 10:40:25 3
2024/07/24 10:40:25 1
2024/07/24 10:40:25 2
2024/07/24 10:40:25 2
2024/07/24 10:40:25 3
2024/07/24 10:40:25 6
2024/07/24 10:40:25 4
PASS
ok      question        0.042s
```

随机获得 1-6.

## Go语言中的time.Tick有何作用

Ticker 断续器，本质可以按照**时间周期**定时向信道channel发送信号，协助程序实现周期性的逻辑。

```go
func Tick(d Duration) <-chan Time
```

示例：

```go
func TestTimeTicker(t *testing.T) {
	for t := range time.Tick(1 * time.Second) {
		log.Println(t.Format("15:04:05"))
	}
}
```

```
> go test -run TimeTicker
2024/07/24 10:58:54 10:58:54
2024/07/24 10:58:55 10:58:55
2024/07/24 10:58:56 10:58:56
2024/07/24 10:58:57 10:58:57
2024/07/24 10:58:58 10:58:58
```

每个1秒，输出当前时间

若需要**控制断续器停止**，使用Ticker类型，示例：

```go
func TestTimeTicker(t *testing.T) {
	//for t := range time.Tick(1 * time.Second) {
	//	log.Println(t.Format("15:04:05"))
	//}
    
	ticker := time.NewTicker(1 * time.Second)
	for t := range ticker.C {
		log.Println(t.Format("15:04:05"))
		if t.Second() == 20 {
			ticker.Stop()
			break
		}
	}
}
```

上面的逻辑，当秒针到达20时，断续器终止！

## Go语言中的time.After有何作用

After(), 得到是 Timer，**定时器**，当时间达到时执行。Timer是一次执行。ticker是循环执行。

方法：

```go
func After(d Duration) <-chan Time
```

示例：

```go
func TestTimeTimer(t *testing.T) {
	log.Println(time.Now().Format("04:05"))
	select {
	// 3秒后执行某个操作
	case t := <-time.After(3 * time.Second):
		log.Println(t.Format("04:05"))
	}
}
```

```
> go test -run TimeTimer
2024/07/24 11:05:10 05:10
2024/07/24 11:05:13 05:13
PASS  
ok      question        3.056s
```

加强控制，可以使用timer类型：

```go
func TestTimeTimer(t *testing.T) {
	log.Println(time.Now().Format("04:05"))
	//select {
	//// 3秒后执行某个操作
	//case t := <-time.After(3 * time.Second):
	//	log.Println(t.Format("04:05"))
	//}
	wg := sync.WaitGroup{}
	timer := time.NewTimer(3 * time.Second)
	stopC := make(chan struct{})
	wg.Add(1)
	go func() {
		defer wg.Done()
		select {
		// 3秒后执行某个操作
		case <-time.After(1 * time.Second):
			timer.Stop()
			stopC <- struct{}{}
			log.Println("timer stop")
		}
	}()
	wg.Add(1)
	go func() {
		defer wg.Done()
		select {
		// 3秒后执行某个操作
		case t := <-timer.C:
			log.Println("time up", t.Format("04:05"))
		case <-stopC:
			return
		}
	}()
	wg.Wait()
}
```

```
> go test -run TimeTimer
2024/07/24 11:11:07 11:07
2024/07/24 11:11:08 timer stop
PASS

```

上面的例子，在定时器到时前，提前进行了关闭。

## Go语言中的GOPATH和GOMOD的区别

go管理项目（包）的方式：

* GOPATH，一个特定目录。
* GOMOD，没有这个环境变量，指的是**以module的方式管理项目（包）。**

**重要：直接使用mod的项目管理方案即可**。

GOPATH，早期的go管理项目包模块依赖的方案。要求将包，项目，都放在GOPATH/src目录下。

目前，go1.11加入对mod的支持，可以在任意的位置管理项目，项目通过项目目录中的go.mod进行依赖模块的组织管理。通过 go mod系列命令完成模块项目管理即可。不再需要（强烈不推荐）将项目代码放在GOPATH/src下。此时GOPATH主要用于管理非自己的源码，第三方的依赖。

## 什么情况下defer会修改返回值

defer 和 return 的执行顺序。

* return 的操作是最后的，函数结束了
* 但是 return 涉及的表达式，return exp, **exp的运算eval会先于defer执行**

```go
func deferReturn() int {
	v := 0
	defer func() { // 2. 执行defer函数
		v = v + 1
	}()
	return v + 1 // 1. 执行v+1 = 1
	// 3. 执行 return 1
}
```

在如下情况下，defer会影响返回值：

* return 后没表达式，使用了具名返回变量
* 函数返回指针类型

```go
func TestDeferReturnEffected(t *testing.T) {
	log.Println(deferReturnNoEffected())
	log.Println(deferReturnEffected())
	log.Println(*deferReturnEffectedToo())
}
func deferReturnNoEffected() int {
	v := 0
	defer func() {
		v = v + 1
	}()
	v = v + 1
	return v
}
func deferReturnEffected() (v int) {
	defer func() { // 2. 执行defer函数
		v = v + 1
	}()
	v = v + 1
	return
}

func deferReturnEffectedToo() *int {
	v := 0
	defer func() {
		v = v + 1
	}()
	v = v + 1
	return &v
}

```

```
> go test -run DeferReturnEffected
2024/07/24 16:43:13 1
2024/07/24 16:43:13 2
2024/07/24 16:43:13 2
   
```

## Go中make和new内置函数的区别

make和new，两个功能的内置函数。

这道问题的原因，由于make和new是否在做内存分配，因此会拿来作比较。

函数签名：

```go
// The make built-in function allocates and initializes an object of type
// slice, map, or chan (only). Like new, the first argument is a type, not a
// value. Unlike new, make's return type is the same as the type of its
// argument, not a pointer to it. The specification of the result depends on
// the type:
//
//	Slice: The size specifies the length. The capacity of the slice is
//	equal to its length. A second integer argument may be provided to
//	specify a different capacity; it must be no smaller than the
//	length. For example, make([]int, 0, 10) allocates an underlying array
//	of size 10 and returns a slice of length 0 and capacity 10 that is
//	backed by this underlying array.
//	Map: An empty map is allocated with enough space to hold the
//	specified number of elements. The size may be omitted, in which case
//	a small starting size is allocated.
//	Channel: The channel's buffer is initialized with the specified
//	buffer capacity. If zero, or the size is omitted, the channel is
//	unbuffered.
func make(t Type, size ...IntegerType) Type

// The new built-in function allocates memory. The first argument is a type,
// not a value, and the value returned is a pointer to a newly
// allocated zero value of that type.
func new(Type) *Type
```

区别一览：

* **功能不同**
  * make，为slice, map, or chan (only)**分配内存和初始化**，返回类型数据。以上三种类型，都有自己的容量的概念，slice底层数组容量，map,容量，chan缓冲容量
  * new，分配内存，为任意类型**分配内存**，返回类型指针
* 语法不同
  * 参数不同
    * make，支持2个或以上参数
    * new, 支持一个参数
  * 返回值不同
    * make，返回**类型数据**
    * new，返回**类型指针**

示例：

```go
func TestMakeVsNew(t *testing.T) {
	//make()
	//new()
	s, m, c := make([]int, 10, 20), make(map[string]int, 20), make(chan int, 20)
	pi, pf, ps := new(int), new(float32), new(string)
	psl := new([]int) // *[]int var psl *[]int
	var plsvar *[]int
	log.Println(s, m, c)
	log.Println(pi, pf, ps)
	log.Println(psl, plsvar)
}

```

```
> go test -run MakeVsNew
2024/07/24 17:00:01 [0 0 0 0 0 0 0 0 0 0] map[] 0xc000176000
2024/07/24 17:00:02 0xc00011c180 0xc00011c188 0xc00012e290
2024/07/24 17:00:02 &[] <nil>
PASS

```

## for循环结构的循环变量是同一个吗

**新版本中不是同一个，是在go1.21后带来的改变**。早期的版本，循环变量是同一个。

```go
for i:=0; i<10; i++ {}
for i := range 10 {}
for k, v := range []string{} {}
```

i，在每次循环中，是否为同一个变量？

示例：

```go
func TestForVar(t *testing.T) {
	for i := 0; i < 3; i++ {
		log.Printf("%p", &i)
	}
	log.Println()
	for i := range 3 {
		log.Printf("%p", &i)
	}
	log.Println()
	for k, v := range []string{"a", "b", "c"} {
		log.Printf("%p %p", &k, &v)
	}
}
```

```
> go test -run ForVar
2024/07/24 17:06:48 0xc00010a0a0
2024/07/24 17:06:48 0xc00010a5d0   
2024/07/24 17:06:48 0xc00010a5d8   
2024/07/24 17:06:48                
2024/07/24 17:06:48 0xc00010a5e8   
2024/07/24 17:06:48 0xc00010a5f0   
2024/07/24 17:06:48 0xc00010a5f8   
2024/07/24 17:06:48
2024/07/24 17:06:48 0xc00010a600 0xc000110240
2024/07/24 17:06:48 0xc00010a608 0xc000110250
2024/07/24 17:06:48 0xc00010a610 0xc000110260

```

## 双引号单引号反引号的区别

```go
// 双引号，常规字符串字面量定义语法，支持转义字符
"字符串，\n\t"
// 反引号，原生字符串字面量定义语法，不支持转义，除了反引号本身\`
`\n\t字符串`

// 单引号, rune类型字面量定义语法
'G'
```

不同的引号，有不同的语法含义（语义）。

## Go语言中函数和方法的区别

本质没有区别，函数和方法都是一段封装好的可执行性代码。

方法：**带有接收器的函数称为方法**。函数接收器的作用就是表示该方法属于哪种类型的方法集合中。方法默认比函数多了一个参数，就是接收器参数。

接收器：reveiver，语法：

差异在：

* 逻辑性：
  * 函数，function，独立的一段功能代码
  * 方法，method，属于某种类型的一段功能代码
* 语法上：
  * 函数，没有接收器，直接调用
  * 方法，带有接收器，使用选择器Selector（t.Method()）的方式来调用

```go
type Stu struct {
}

func (s Stu) Study() error {
	return nil
}
```

## Go语言中如何比较两个map是否相等

==，!= 比较相等运算符，在map类型上未实现。

```go
m1, m2 := map[string]int{}, map[string]int{}
// Invalid operation: m1 == m2 (the operator == is not defined on map[string]int)
log.Println(m1 == m2)
```

类似：slice, map, func 类型是未定义==，!=运算符。

但是，可以直接于nil进行比较，用来判断是否未初始化：

```go
log.Println(m1 == nil, m2 == nil)
```

如果需要比较：

* **maps.Equal()，推荐的，1.20+增加的**
* reflect.DeepEqual()，不仅仅支持map，还支持slice，map，等
* 自定义比较函数

示例：

```go
func TestMapEqual(t *testing.T) {

	m1, m2 := map[string]int{"go": 2008}, map[string]int{"go": 2008}
	// Invalid operation: m1 == m2 (the operator == is not defined on map[string]int)
	//log.Println(m1 == m2)
	log.Println(m1 == nil, m2 == nil)
	log.Println(maps.Equal(m1, m2))
	log.Println(reflect.DeepEqual(m1, m2))
}
```

```
> go test -run MapEqual
2024/07/24 17:54:43 false false
2024/07/24 17:54:43 false
2024/07/24 17:54:43 false
PASS

```

map相等的条件：

* 包含相同的key/value对
* 同时为nil

## Go语言如何高效的拼接字符串

**选择+或strings.Builder即可。如果大量拼凑，strings.Builder的方案会好些。**

几种不同的方案拼接字符串：

1. 运算符+
2. strings.Builder
3. fmt.Sprintf 格式化字符串
4. strings.Join() 将string的切片连接起来，已经存在[]string时才使用
5. bytes.Buffer，类似的[]byte的方案，类似于append，在存在[]bytes时，可以考虑使用。

从功能，从最合适的业务逻辑来选方案。直接相关的就是+和strings.Builder。效率都是ok的，如果大量拼凑，Builder的方案会好些，小量+和Builder方案类似。

示例：

```go
func TestConcatString(t *testing.T) {
	s1, s2, s3, s4, s5 := "Go", "is", "the", "best", "language"

	log.Println(s1 + s2 + s3 + s4 + s5)

	builder := strings.Builder{}
	builder.WriteString(s1)
	builder.WriteString(s2)
	builder.WriteString(s3)
	builder.WriteString(s4)
	builder.WriteString(s5)
	log.Println(builder.String())
}
```

## Go语言中的interface间可以比较吗

可以比较。==运算符支持interface。

相等的逻辑：

* 动态类型一致
* 动态值相等

上面两个条件满足，意味着相等。

都是 nil 也意味着相等。

reflect.DeepEqual() 也支持interface的比较

示例：

```go
func TestInterfaceEqual(t *testing.T) {
	var i1, i2 interface{}
	i1 = 42
	i2 = uint(42)
	log.Println(i1 == i2)

	var i3, i4 interface{}
	i3 = 42
	i4 = 42
	log.Println(i3 == i4)

	var i5, i6 interface{}
	log.Println(i5 == nil, i6 == nil, i5 == i6)

	log.Println(reflect.DeepEqual(i3, i4))
}
```

```
> go test -run InterfaceEqual
2024/07/24 18:33:07 false
2024/07/24 18:33:07 true
2024/07/24 18:33:07 true true true
2024/07/24 18:33:07 true
PASS

```

## Go语言中Map元素可以取地址吗

不可以。

原因是：Go中Map的底层是Hash表，Key/Value存储在Bucket中，因此不能取地址。

示例：

```go
func TestMapElementAddress(t *testing.T) {
	m := map[string]int{
		"go": 2008,
	}
	log.Println(&m)
	//log.Println(&m["go"]) // Cannot take the address of 'm["go"]'
}
```

## Go的Map的Key为什么是无序的

主要体现是遍历map时，顺序会不一致。

```go
func TestMapRange(t *testing.T) {
	m := map[string]int{
		"go":   2008,
		"cpp":  1998,
		"java": 1234,
	}
	m["js"] = 123
	m["ts"] = 44
	m["rust"] = 897
	for k, v := range m {
		log.Println(k, v)
	}
}
```

```
> go test -run MapRange
2024/07/24 19:44:52 ts 44
2024/07/24 19:44:52 rust 897
2024/07/24 19:44:52 go 2008
2024/07/24 19:44:52 cpp 1998
2024/07/24 19:44:52 java 1234
2024/07/24 19:44:52 js 123
PASS
ok      question        0.051s
> go test -run MapRange
2024/07/24 19:44:56 go 2008
2024/07/24 19:44:56 cpp 1998
2024/07/24 19:44:56 java 1234
2024/07/24 19:44:56 js 123
2024/07/24 19:44:56 ts 44
2024/07/24 19:44:56 rust 897
PASS
```

顺序无序的原因：

map的底层Hash表，**key/value落在具体的bucket上。**

当map的容量发生改变时，key就会搬迁。不能保证原来的顺序。

遍历的时候，通常会按照bucket遍历，bucket中的key遍历。

for range map时，有时会（取决于语言）随机选择bucket进行遍历。操作逻辑上就是无序的。

## Go如何导入本地自定义的包

本地自定义包：

* 项目内
* 项目外

重点讨论项目外的。

**首先不推荐这种方案！**

示例：

```go
abc/
	import_test.go
	go.mod
xyz/
	funcs.go
	go.mod
```

需要 abc/import_test.go中使用xyz/funcs.go中的函数。

abc, xyz 是两个独立的模块。

错误：直接通过路径的方式导入，会出现如下错：

```go
import "../xyz"

func TestImportOuter(t *testing.T) {
	xyz.SayWhere()
}

// import_test.go:4:8: local import "../xyz" in non-local package

```

方案：

修改go.mod来解决：

abc/go.mod

```go
module abc

go 1.22.0

require xyz v0.0.0
replace xyz => ../xyz

```

abc/import_test.go

```go
package abc

import "testing"
import "xyz"

func TestImportOuter(t *testing.T) {
	xyz.SayWhere()
}

```

```
> go test -run ImportOuter
2024/07/25 11:12:59 in xyz module
PASS

```

## 阐述Go语言中的Map扩容机制

Hash表的扩容机制。

哈希表的实现，通常有两种设计：

1. 开放寻址法。数组
2. 拉链法，最典型的方案。**数组+链表来实现（有些会增加红黑树）**

拉链法，也是Go的Map整体结构：

![image.png](https://fynotefile.oss-cn-zhangjiakou.aliyuncs.com/fynote/fyfile/13080/1721814440012/3d6a47950e5a4c8392ce5cd88347e6b7.png)

**何时需要扩容?(扩容的条件）容量不足时需要扩容：**

1. 元素数量超过了过载因子（6.5)时。逻辑上，元素过多。
2. 溢出桶的数量过多时，多的标准指的是，溢出桶大概与常规桶的数量差不多时，大约是数组元素数量的一半。逻辑上，表示哈希函数不够均匀或桶中存在大量的空洞。

过载因子：Load Factor，指的是元素数量与桶数量：

```
过载因子 = 元素数量/桶数量
```

溢出桶，与桶（常规桶）相对，当常规桶满了时后，map（hashtable）为了避免频繁的扩容操作，允许再常规桶后，增加额外的桶来存储键值对，这些称为溢出桶。

**扩容算法：**

根据不同的扩容原因，会有两种扩容方案：

1. 元素数量超过了过载因子，扩容算法为：双倍容量。指的是重新创建原来长度的2倍的数组，来进行桶的管理。
2. 溢出桶的数量过多时，扩容的算法为：等量扩容。指的创建相同长度的数组扩容，扩容的操作本质就是整理内存。

**何时触发扩容？**

在对map做写操作时，会触发扩容。删除，更新，添加时触发。

**扩容策略：**

不是全量扩容，而是**增量扩容**。

* 全量扩容，一次性把扩容操作全部完成
* 增量扩容，一点点的将扩容操作完成，一次移动一个桶，一次清理一个桶，释放一个桶

## Go语言中的Map的数据结构

Map就是HashTable，Go采用的拉链法实现，因此底层结构就是**数组+链表。**

如图所示：

![image.png](https://fynotefile.oss-cn-zhangjiakou.aliyuncs.com/fynote/fyfile/13080/1721814440012/1157c5662f36472396f9465f768e1a25.png)

对应的代码描述：

```go
type hmap struct {
	count     int // 元素数量
	flags     uint8 // 操作的标志位
	B         uint8 // 桶的数量，存储桶数量对数，len(buckets) == 2^B
	noverflow uint16 // 溢出桶的数量
	hash0     uint32 // 哈希函数的种子，增加随机性

	buckets    unsafe.Pointer // 桶地址
	oldbuckets unsafe.Pointer // 扩容时，旧桶地址
	nevacuate  uintptr // 扩容进度

	extra *mapextra // 扩展操作的信息
}

type mapextra struct {
	overflow    *[]*bmap // 溢出桶地址集合
	oldoverflow *[]*bmap // 旧的溢出桶的地址集合
	nextOverflow *bmap // 处理的溢出桶
}
```

## 为什么说Go语言是鸭子类型

鸭子类型，Duck typing，If it looks like a duck, swims like a duck, and quacks like a duck, then it probably is a duck.

如果看起来像鸭子，游泳像鸭子，叫起来像鸭子，那么就认为是鸭子。

是**动态编程**语言在做类型推测的一种策略，在乎的功能本身，而不是语法绑定。

动态编程语言，很多操作，包括类型推测是在运行期间完成的，而不是在编译期间完成。

而我们的**Go语言是典型的静态语言**，会在编译期间，确定很多事情。

因为Go新兴语言，作为静态语言，引入了鸭子类型的特征和语法便利。

**主要体现在接口的实现上，不要求显式的实现接口，没有implements关键字将类型和接口绑定，只要实现了某个接口要求的方法集合，那么就视为实现这个接口。因此Go是鸭子类型。**

## 接口中iface和eface的实现有何不同

都是接口类型的底层结构：

* iface，常规接口，具有方法集合定义的接口
* eface, 空接口， interface{}， 没有方法集合定义的接口

eface的结构：

```go
type eface struct {
    _type *_type
    data  unsafe.Pointer
}
```

* _type 实体值类型
* data 实体值

iface的结构：

```go
type iface struct {
	tab  *itab
	data unsafe.Pointer
}

type itab struct {
	inter  *interfacetype
	_type  *_type
	link   *itab
	hash   uint32 // copy of _type.hash. Used for type switches.
	bad    bool   // type does not implement interface
	inhash bool   // has this itab been added to hash?
	unused [2]byte
	fun    [1]uintptr // variable sized
}
```

* tab itab表示，**接口本身类型**及接口**动态实体类型**
  * inter 接口本身类型
  * _type 实体值类型
  * link 相关的itab
  * hash 拷贝_type.hash，用在 type switch中
  * bad 类型是否未实现该接口
  * inhash 当前的itab是否在hash中
  * unused 未用标志
  * fun 变量来确认，该接口具体类型的方法地址
* data 实体值

## 切片Map是值传递还是引用传递

**总结：语法上类似引用传递，但实际是值传递，由于拷贝后使用的是同一个底层存储，修改时会互相影响。**

（类似，slice类型，go中本质都是值传递）

如下代码：

```go
func TestMapTrans(t *testing.T) {
	m1 := map[string]int{
		"go": 2008,
	}
	m2 := m1

	m1["go"] = 42
	log.Println(m1)
	log.Println(m2)
	fmt.Printf("%p, %p\n", &m1, &m2)
}
```

```
> go test -run MapTrans
2024/07/26 10:32:01 map[go:42]
2024/07/26 10:32:01 map[go:42]
0xc000058060, 0xc000058068
PASS
```

从结果上看，修改m1影响m2说明map是引用传递。

但是，map的结构，是由两部分构成：

* map本身，语法上hmap，存储了map的元素个数等信息，其中有桶地址的字段，用来存储map的key/value数据的桶的地址。
* 桶，存储数据的单元（一大堆链表）

当拷贝（复制）map的时候，**拷贝的是map本身的部分，而没有拷贝桶的部分。**这就导致拷贝得到的新map，与之前的map使用的是相同的底层存储结构。

基于这个原因，m2 := m1，**其实是值传递**，但是由于存储结构的原因，导致m1与m2共享hash表数据。

## 如何在Go中拷贝map类型数据

目前的方案：

* maps.Clone()，**推荐的方案**，标准库maps的方案。需要go的新版本
* reflect.Copy(), **基于反射的拷贝**
* 自定义拷贝方法，为需要的map类型定义拷贝函数，因为要考虑map的value是否继续需要深度拷贝

示例：

```go
func TestMapCopy(t *testing.T) {
	m1 := map[string]int{
		"go": 2008,
	}
	m2 := maps.Clone(m1)
	m1["go"] = 42
	log.Println(m1)
	log.Println(m2)
}
```

```
> go test -run MapCopy 
2024/07/26 10:41:00 map[go:42]
2024/07/26 10:41:00 map[go:2008]
PASS

```

类似的slice也有方法：

```go
slices.Clone()
func Clone[S ~[]E, E any](s S) S
```

## 两个nil会不相等吗

结论：**两个nil不会不相等。特例是当空接口interface{}的类型确定但值为nil时，该接口不等于nil**。

```go
var p *int = nil // 指针类型变量，值为 nil
var i interface{} = p // 空接口的动态类型是 *int，值是 nil
fmt.Println(i == nil) // false
```

注意：nil 不能直接比较，==没有在nil上实现。

```go
nil == nil
```

是不会出现的。

可以使用值与nil比较，判断是否==nil

```go
if err != nil {}
```

**如果两个值同时为nil，那么就是相等的。**

注意特殊的地方法：在inteface{}空接口类型的处理上。

考点：

* 如何比较两个interface{}是否相等？
* 如何比较interface{}和值是否相等？
* 如何比较interface{}和nil是否相等？

通用的示例：

```go
func TestNilEqual(t *testing.T) {
	// p 变量为nil，*int类型
	var p *int = nil
	// i 变量，interface{}空接口， 接口的值是p
	var i interface{} = p

	// p是nil
	fmt.Println(p == nil) // true
	// 判断i和p相等
	fmt.Println(i == p) // true

	// 1. p是nil 2. p和i相等 3. 推导出i是nil
	// 测试：i不等于nil
	// 结论：nil 不等于 nil。两个nil可能会不等
	fmt.Println(i == nil) // false
}
```

用来说明两个nil可能会不相等。

以上示例的逻辑漏洞：为什么i == p? 也就是**如何比较interface{}和值是否相等**？

当接口类型和值类型做比较时，将值类型转换为接口类型做比较。因此上面的比较是将p转为interface{}再比较，也就是：

```go
i == p
i == interface{}(p)
```

此时，就要了解**如何比较两个interface{}是否相等**？

空接口interface{}内部的实现叫eface, 有两个字段，T(_type)和V(data)，类型和值。

比较的方式，**先比较T，再比较V**，任何不等都认为两个接口不相等。

最后 i != nil 呢？思考**如何比较interface{}和nil是否相等？**

只有T和V都是nil的interface{}才是nil。

```go
var ii interface{}
ii == nil // true
```

```go
func TestNilEqual(t *testing.T) {
	// p 变量为nil，*int类型
	var p *int = nil
	// i 变量，interface{}空接口， 接口的值是p
	var i interface{} = p

	// p是nil
	log.Println(p == nil) // true
	// 判断i和p相等
	log.Println(i == p) // true

	// 1. p是nil 2. p和i相等 3. 推导出i是nil
	// 测试：i不等于nil
	// 结论：nil 不等于 nil。两个nil可能会不等
	log.Println(i == nil) // false

	var ii interface{}
	log.Println(ii == nil) // true
}
```

```
> go test -run NilEqual
2024/07/26 11:24:26 true
2024/07/26 11:24:26 true
2024/07/26 11:24:26 false
2024/07/26 11:24:26 true
PASS

```

## Map的delete操作删除key后内存会立即释放吗

不会的。

- 虽然键值对被删除，但 `map` 的底层数据结构（如分配的存储空间）可能仍然存在，等待被后续操作重用。
- Go 的垃圾回收机制会在适当的时候回收未使用的内存。

内存优化方案：

* 什么都不做，等待map在检测到溢出桶过多时，进行等量扩容操作，进行内存的优化整理
* 手动做，可尽量将不用的key删除或设置为nil（内存上优势不大，操作效率上会有所提升）

## Go中nilslice和空slice有何区别?

```go
// nil slice
var nilSlice []int

// empty slice
var emptySlice []int = []int{}
emptySlice := []int{}
```

- nilSlice == nil, 其实是没有分配内存的slice数据。

- emptySlice 是分配了内存的Slice，长度为0.


类似的：

nil map 和 空map

```go
var nilMap map[string]int
emptyMap := map[string]int{}
```

## 如何理解Go中的Rune类型

Rune，字符类型，**表示单字符的类型。**

语法本质：int32的别名，与int32是一致的，4字节的整数类型。

```go
// rune is an alias for int32 and is equivalent to int32 in all ways. It is
// used, by convention, to distinguish character values from integer values.
type rune = int32
```

rune类型的字面量，是使用单引号包裹的单字符，或者单转义字符，或单字符编码。

```go
'符'
'a'
'ä'
'本'
'\t'
'\000'
'\007'
'\377'
'\x07'
'\xff'
'\u12e4'
'\U00101234'
```

rune类型存储的是字符的Unicode码值（码点）（Unicode Code Point)。unicode4字节，utf8mb4。

Rune类型用于表示存储单个字符的unicode码值，本质上是int32类型，逻辑上表现为单个字符。

## Go语言中Struct类型是否可以比较

== 是否支持结构体类型。

常规情况下是可以直接对于struct类型做==运算的，可以比较。

**前提是，struct类型的全部字段要是可比较的**。例如当字段为map类型时，会导致整个struct不能比较。

示例：

```go
func TestStructCmp(t *testing.T) {
	type X struct {
		F1 string
		F2 int
	}
	x1, x2 := X{}, X{}
	log.Println(x1 == x2)

	type Y struct { // Y 不是可比较类型
		F1 string
		F2 map[string]int // F2 不是可比较类型
	}
	//y1, y2 := Y{}, Y{}
	//Invalid operation: y1 == y2 (the operator == is not defined on Y)
	//log.Println(y1 == y2)
}
```

上面的例子，Y就是不可比较的结构体，理由是Y的F2是不可比较类型。而X是可比较的结构体，理由是X的两个（全部）字段都是可比较的。

## 函数返回局部变量指针是否安全

答案：安全。

考点：

* 函数的局部变量，大多分配到**函数的运行栈空间**
* 函数的运行栈空间，会在函数运行结束后释放

技术点：

栈逃逸，当Go检测到函数的局部变量（栈变量）需要继续使用时，就会将该变量分配到堆中，称为栈逃逸。从栈中逃逸到堆中。

## 概述Go中chan类型的底层结构

chan channel 通道类型的数据结构。

```go
// 构建元素类型为int，缓冲容量为3的chan类型
c := make(chan int, 3)
```

chan的结构通过源码：runtime.hchan:

```go
type hchan struct {
	qcount   uint           // total data in the queue。channel中的元素个数
	dataqsiz uint           // size of the circular queue。channel。循环队列的长度
	buf      unsafe.Pointer // points to an array of dataqsiz elements。缓冲区指针
	elemsize uint16 // 缓冲元素大小
	closed   uint32 // 是否关闭
	elemtype *_type // element type 。缓冲元素类型
	sendx    uint   // send index // 发送操作处理的位置
	recvx    uint   // receive index // 接收操作处理的位置
	recvq    waitq  // list of recv waiters // receive操作阻塞的goroutine队列
	sendq    waitq  // list of send waiters // send操作阻塞的goroutine队列

	// lock protects all fields in hchan, as well as several
	// fields in sudogs blocked on this channel.
	//
	// Do not change another G's status while holding this lock
	// (in particular, do not ready a G), as this can deadlock
	// with stack shrinking.
	lock mutex // 锁
}

// 等待队列， 是双向链表
type waitq struct {
	first *sudog
	last  *sudog
}
```

每个字段的参考字段注释！

## Go语言中的多值返回是如何实现的

总结答案：

基于**SP + Offset** 的方案，来确定返回值的位置，进而进行多值返回。

Go中将返回变量，和局部变量，参数变量是统一管理的。如下代码中的：p1, p2, x, y 以及两个匿名的返回值变量是统一管理。统一以 Stack Pointer堆栈指针 + Offset偏移量的方案管理。确定变量基于SP的偏移位置，来确定变量地址，返回变量也是如此。

找了个函数的固定位置，来存储返回值。

示例代码：

```go
func Test(p1, p2 int) (int, int) {
	x := p1 + p2
	y := p1 - p2
	return x, y
}
```

基于以上示例代码，生成汇编代码，可以直观的看到SP+Offset的方案。

创建 main.go 存储上面的代码。

运行go tool compile 命令进行编译：

```
>go tool compile -N -l -S main.go > main.s
```

-N, -l 禁用编译器优化，禁用连接

-S 生成汇编

查看生成的main.s

```go
main.Test STEXT nosplit size=74 args=0x10 locals=0x28 funcid=0x0 align=0x0
	0x0000 00000 (D:/apps/mashibing/question/func/main.go:3)	TEXT	main.Test(SB), NOSPLIT|ABIInternal, $40-16
	0x0000 00000 (D:/apps/mashibing/question/func/main.go:3)	PUSHQ	BP
	0x0001 00001 (D:/apps/mashibing/question/func/main.go:3)	MOVQ	SP, BP
	0x0004 00004 (D:/apps/mashibing/question/func/main.go:3)	SUBQ	$32, SP
	0x0008 00008 (D:/apps/mashibing/question/func/main.go:3)	FUNCDATA	$0, gclocals·g2BeySu+wFnoycgXfElmcg==(SB)
	0x0008 00008 (D:/apps/mashibing/question/func/main.go:3)	FUNCDATA	$1, gclocals·g2BeySu+wFnoycgXfElmcg==(SB)
	0x0008 00008 (D:/apps/mashibing/question/func/main.go:3)	FUNCDATA	$5, main.Test.arginfo1(SB)
	// p1, p2 参数的位置 MOVQ（移动8个字节）, 赋值。
	0x0008 00008 (D:/apps/mashibing/question/func/main.go:3)	MOVQ	AX, main.p1+48(SP)
	0x000d 00013 (D:/apps/mashibing/question/func/main.go:3)	MOVQ	BX, main.p2+56(SP)
	// r0, r1 返回值的位置
	0x0012 00018 (D:/apps/mashibing/question/func/main.go:3)	MOVQ	$0, main.~r0+8(SP)
	0x001b 00027 (D:/apps/mashibing/question/func/main.go:3)	MOVQ	$0, main.~r1(SP)
	0x0023 00035 (D:/apps/mashibing/question/func/main.go:4)	ADDQ	BX, AX
	0x0026 00038 (D:/apps/mashibing/question/func/main.go:4)	MOVQ	AX, main.x+24(SP)
	0x002b 00043 (D:/apps/mashibing/question/func/main.go:5)	MOVQ	main.p1+48(SP), CX
	0x0030 00048 (D:/apps/mashibing/question/func/main.go:5)	SUBQ	BX, CX
	0x0033 00051 (D:/apps/mashibing/question/func/main.go:5)	MOVQ	CX, main.y+16(SP)
	// 将计算结果，赋值给返回值变量位置
	0x0038 00056 (D:/apps/mashibing/question/func/main.go:6)	MOVQ	AX, main.~r0+8(SP)
	0x003d 00061 (D:/apps/mashibing/question/func/main.go:6)	MOVQ	CX, main.~r1(SP)
	0x0041 00065 (D:/apps/mashibing/question/func/main.go:6)	MOVQ	CX, BX
	0x0044 00068 (D:/apps/mashibing/question/func/main.go:6)	ADDQ	$32, SP
	0x0048 00072 (D:/apps/mashibing/question/func/main.go:6)	POPQ	BP
	0x0049 00073 (D:/apps/mashibing/question/func/main.go:6)	RET
	0x0000 55 48 89 e5 48 83 ec 20 48 89 44 24 30 48 89 5c  UH..H.. H.D$0H.\
	0x0010 24 38 48 c7 44 24 08 00 00 00 00 48 c7 04 24 00  $8H.D$.....H..$.
	0x0020 00 00 00 48 01 d8 48 89 44 24 18 48 8b 4c 24 30  ...H..H.D$.H.L$0
	0x0030 48 29 d9 48 89 4c 24 10 48 89 44 24 08 48 89 0c  H).H.L$.H.D$.H..
	0x0040 24 48 89 cb 48 83 c4 20 5d c3                    $H..H.. ].
go:cuinfo.producer.<unlinkable> SDWARFCUINFO dupok size=0
	0x0000 2d 4e 20 2d 6c 20 72 65 67 61 62 69              -N -l regabi
go:cuinfo.packagename.main SDWARFCUINFO dupok size=0
	0x0000 6d 61 69 6e                                      main
main..inittask SNOPTRDATA size=8
	0x0000 00 00 00 00 00 00 00 00                          ........
gclocals·g2BeySu+wFnoycgXfElmcg== SRODATA dupok size=8
	0x0000 01 00 00 00 00 00 00 00                          ........
main.Test.arginfo1 SRODATA static dupok size=5
	0x0000 00 08 08 08 ff                                   .....

```

上面加注释的部分，用于描述返回值变量初始化，和运算完毕赋值的过程。

注意的是返回值变量的表示方案：**SP+Offset的方案。**

## Go中的类型指针与unsafe.Pointer和uintptr有何区别

> 在 Go 中，类型指针、`unsafe.Pointer` 和 `uintptr` 是三种不同类型的指针，它们的用途和特性也有所不同。以下是详细的区别：
>
> ------
>
> ### 1. **类型指针（Type Pointer）**
>
> #### 特点：
>
> - 类型指针是指向特定类型的变量的指针，例如 `*int`、`*string` 等。
> - 编译器会对类型指针进行严格的类型检查。
> - 类型指针可以安全地解引用，访问或修改其指向的值。
>
> #### 使用场景：
>
> - 用于常规指针操作，比如引用变量、传递指针参数等。
>
> #### 示例：
>
> ```go
> var i int = 42
> var p *int = &i
> *p = 43  // 修改指针指向的值
> fmt.Println(i) // 输出 43
> ```
>
> ------
>
> ### 2. **`unsafe.Pointer`**
>
> #### 特点：
>
> - `unsafe.Pointer` 是一种**通用指针类型，可以与任何指针类型互相转换。**
> - 它没有类型信息，不支持直接解引用。
> - `unsafe.Pointer` 的主要目的是用于绕过 Go 的类型安全系统（Unsafe 操作）。
>
> #### 使用场景：
>
> - 用于实现底层内存操作和**高性能优化**，比如与 C 库交互或操作结构体内存布局。
> - 通常与 `unsafe` 包中的其他操作配合使用，比如类型转换、强制对齐等。
>
> #### 示例：
>
> ```go
> import "unsafe"
> 
> var i int = 42
> var p = unsafe.Pointer(&i)  // 转换为 unsafe.Pointer
> var pi = (*int)(p)          // 再转换回 *int
> *pi = 43
> fmt.Println(i) // 输出 43
> ```
>
> ------
>
> ### 3. **`uintptr`**
>
> #### 特点：
>
> - `uintptr` 是一种用于表示指针地址的整型类型。
> - 它可以将指针转换为整数，表示内存地址。
> - 与 `unsafe.Pointer` 不同，`uintptr` 的主要用途是进行数值运算，比如**偏移操作。**
>
> #### 注意：
>
> - `uintptr` 不是一个指针，垃圾回收器不会将其视为指针。如果你通过 `uintptr` 操作地址，可能会引入问题，比如指针悬空。
>
> #### 使用场景：
>
> - 用于与内存地址相关的计算或偏移操作。
>
> #### 示例：
>
> ```go
> import "unsafe"
> 
> type MyStruct struct {
>     a int
>     b int
> }
> 
> var s = MyStruct{a: 1, b: 2}
> ptr := uintptr(unsafe.Pointer(&s)) // 获取结构体地址
> offset := unsafe.Offsetof(s.b)     // 获取字段偏移量
> bPtr := (*int)(unsafe.Pointer(ptr + offset))
> *bPtr = 42
> fmt.Println(s.b) // 输出 42
> ```
>
> ------
>
> ### **总结**
>
> | 特性           | 类型指针     | `unsafe.Pointer`         | `uintptr`          |
> | -------------- | ------------ | ------------------------ | ------------------ |
> | 类型安全       | 是           | 否                       | 否                 |
> | 可解引用       | 是           | 否                       | 否                 |
> | 可用于地址运算 | 否           | 否                       | 是                 |
> | 垃圾回收支持   | 是           | 是                       | 否                 |
> | 常用场景       | 普通指针操作 | 底层内存操作，跨语言交互 | 地址计算，偏移操作 |
>
> 这三者分别适用于不同场景，使用时需注意其特性和可能的风险（特别是 `unsafe.Pointer` 和 `uintptr`）。



Go中三种指针类型：

* 类型指针（Type Pointer），某种数据类型的指针类型，*Type， *int, *string。**不能运算，仅用于传递**。
* unsafe.Pointer，通用类型的指针。任意类的数据类型构建指针，**核心功能是转换不同的指针类型，但不能执行运算**。
* uintptr，**用于指针的计算，主要就是偏移**。配合unsafe.Pointer使用。

目的：

在程序员可控的情况下（对程序员的水平有要求），unsafe.Pointer + uintptr 增加了程序对指针的运算能力。

---

### unsafe.Pointer

```go
type ArbitraryType int
type Pointer *ArbitraryType
```

ArbitraryType 逻辑上用于表示任意类型的文档类型（用于文档描述类型）。

类型指针间的转换示例：

```go
func TestPointer(t *testing.T) {
	v1 := int(32)
	v2 := int64(24)
	p := &v1 // p *int

	// 使用p，找到v2
	// 将v2的指针*int64类型转换为 *int类型。
	//p = (*int)(&v2)
	p = (*int)(unsafe.Pointer(&v2)) // *int64 to unsafe.Pointer to *int
	log.Println(*p)
}
```

```
> go test -run Pointer
2024/07/28 09:13:19 24
PASS
```

官方例子：

```go
package math
// Float64bits returns the IEEE 754 binary representation of f,
// with the sign bit of f and the result in the same bit position,
// and Float64bits(Float64frombits(x)) == x.
func Float64bits(f float64) uint64 { return *(*uint64)(unsafe.Pointer(&f)) }
```

返回浮点数的IEEE754的二进制表示，使用uint64的方案。

```go
log.Println(math.Float64bits(3.14))
4614253070214989087
```

### uintptr

```go
// uintptr is an integer type that is large enough to hold the bit pattern of
// any pointer.
type uintptr uintptr
```

用于计算的。

典型的场景，利用uintptr获取数组的元素偏移或结构体的字段偏移：

示例：

```go
func TestPointerUintptr(t *testing.T) {
	a := [...]int{10, 20, 30}
	// 通过指针运算，获取a的第二个元素
	ap := unsafe.Pointer(&a)
	offset := unsafe.Sizeof(a[0]) // 元素大小，偏移位置
	// ap + offset
	v := *(*int)(unsafe.Pointer(uintptr(ap) + offset))
	log.Println(v)
}
```

```
> go test -run PointerUintptr
2024/07/28 09:26:09 20
PASS

```

tip：结构体操作时，要注意结构体的字段内存对齐问题。

## Go语言空结构体的作用

```
struct{}
```

逻辑上，不使用来设计某类型的。

主要是用于做**有无类信号数据**。

struct{}的特点：

* 内存占用为0，作为值来说，效率极高
* 通常多个struct{}可能会使用同一个地址

例如：

```go
// channel信号时
ch := make(chan struct{})

// Set 类型，不重复的key的集合
set := map[string]struct{}
```

类似的还有 [0]int ，空数组也具备类似特性。数组要指定元素类型，而空结构体只有struct关键字，语义更清晰简洁。

## 如何在不使用第三个变量下交换两个数值变量的值

```
a, b := 10, 20
```

示例：

```go
func TestSwapAB(t *testing.T) {
	a, b := 10, 20
	a, b = b, a // 先计算右侧表达式的值，再给左侧变量赋值
	log.Println(a, b)

	c := 10
	d := 20
	c = c + d // c == sum
	d = c - d // sum - d == c
	c = c - d // sum - c == d
	log.Println(c, d)
}
```

```
> go test -run SwapAB
2024/07/28 09:47:51 20 10
2024/07/28 09:47:51 20 10
PASS
```

## string类型的值可以修改吗

不可以，string类型是只读的。指的是字符串的下标操作只能读取某个字节，不能修改某个字节。

```go
s := "GoLang"
log.Println(s[2])
// s[2] = X // 不能操作
```

通常与[]byte转换，来实现字符串内容的修改：

```go
s := "GoLang"
b := []byte(s)
b[2] = 97
s = string(b)
```

## switch语句如何执行下一个case

fallthrough 语句即可。

go中的switch case，每个case执行完毕会break。（有些语言case执行完了，会执行下面的case，需要强行break）

通常建议，将多个case与default写在后边，通过fallthrough完成一致的操作：

```go
switch {
	case case1, case2: fallthrough
    default:
		some operate
}
```

## Go语言支持哪些数据类型

* 布尔型 bool
* 数值型 number
  * 整数型
    * 固定位宽的无符号整数：uint8, uint16, uint32, uint64
    * 固定位宽的符号整数：int8, int16, int32, int64
    * 架构位宽的整数型：uint，int（默认类型）
    * 别名
      * byte uint8的别名
      * rune int32的别名
    * uintptr 存储地址的整数型，用于地址计算
  * 浮点数型
    * float32 单精度浮点数
    * float64 双精度浮点数（默认类型）
  * 复数型
    * complex64
    * complex128（默认类型）
* 字符串类型string
* 数组类型array
* 切片类型slice
* 映射表类型map
* 结构体类型struct
* 指针类型Pointer
* 接口类型interface
* 通道类型Chan
* 函数类型Func

支持使用go关键字，定义类型。包括类型别名和类型定义。

```go
// 类型别名
type byte as uint8
// 类型定义
type Counter int64
```

> ### 1. **类型定义（Type Definition）**
>
> #### 语法：
>
> ```go
> type NewType ExistingType
> ```
>
> #### 特点：
>
> - **新类型与原类型是完全不同的类型**，即便它们的底层表示相同。
>
> ### 2. **类型别名（Type Alias）**
>
> #### 语法：
>
> ```go
> type AliasName = ExistingType
> ```
>
> #### 特点：
>
> - **类型别名与原类型是同一种类型**，它们在编译时会被视为同一类型。

## 如何从panic中恢复

使用recover()内置函数恢复。

panic：go中**运行时期的错误**，称为panic。相对编译时期的错误来说。

例如，索引过界，就是运行时的错误。

分为：

* 程序panic，程序执行时发现的
* 用户panic，用户通过内置函数panic发出的

**recover()要在panic后执行**，通常会配合函数的defer机制完成。因为一旦程序中出现panic，意味着程序执行终止，相当于该函数运行结束。

```go
func mayPanic() {
	defer func() {
		if err := recover(); err != nil {
			log.Println("after panic")
		}
	}()
	a := [...]int{10, 20, 30}
	i := 3
	log.Println(a[i])
}

func TestErrorRecover(t *testing.T) {
	mayPanic()
}
```

```
> go test -run ErrorRecover
2024/07/30 12:07:21 after panic
PASS
```


自己通过panic也可以生成被recover捕获的错误：

```go
func mayPanic() {
	defer func() {
		if err := recover(); err != nil {
			log.Println("after panic", err)
		}
	}()
	a := [...]int{10, 20, 30}
	i := 3
	if i < 0 || i >= 3 {
		panic("out of index")
	}
	log.Println(a[i])
}

func TestErrorRecover(t *testing.T) {
	mayPanic()
}
```
