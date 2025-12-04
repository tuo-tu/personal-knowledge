### gRPC 通信

#### RPC 介绍

RPC, Remote Procedure Call，远程过程调用。与 HTTP 一致，也是应用层协议。该协议的目标是实现：调用远程过程（方法、函数）就如调用本地方法一致。

如图所示：

![image.png](https://fynotefile.oss-cn-zhangjiakou.aliyuncs.com/fynote/fyfile/13080/1657590789084/41b39fe5425f4e8ba0b83075520afab5.png)

说明：

- ServiceA 需要调用 ServiceB 的 FuncOnB 函数，对于 ServiceA 来说 FuncOnB 就是远程过程
- RPC 的目的是让 ServiceA 可以像调用 ServiceA 本地的函数一样调用远程函数 FuncOnB，也就是 ServieA 上代码层面使用：`serviceB.FuncOnB()` 即可完成调用
- **RPC 是 C/S 模式，调用方为 Client，远程方为 Server**
- RPC 把整体的调用过程，数据打包、网络请求等，封装完毕，存储在 C、S 两端的 Stub 中。**Stub（代码存根）**
- 调用流程如下，**共6步：**
  1. ServiceA 将调回需求告知 Client Sub
  2. Client Sub 将调用目标（Call ID）、参数数据（params）等调用信息进行打包（序列化），并将打包好的调用信息**通过网络传输给 Server Sub**
  3. Server Sub 将根据调用信息，调用相应过程。期间涉及到数据的**拆包（反序列化）**等操作。
  4. 远程过程 FuncOnB 运行，并得到结果，将结果告知 Server Sub
  5. Server Sub 将结果打包，并传输回给 Client Sub
  6. Client Sub 将结果拆包，把最终函数调用的结果告知 ServiceA

以上就是典型 RPC 的流程。

RPC 协议没有对网络层做规范，那也就意味着具体的 RPC 实现可以基于 TCP，也可以基于其他协议，例如 HTTP，UDP 等。RPC 也没有对数据传输格式做规范，也就是逻辑层面，传输 JSON、Text、protobuf 都可以。这些都要看具体的 RPC 产品的实现。广泛使用的 RPC 产品有 gRPC，Thrift 等。

#### gRPC 介绍

gPRC 官网（https://grpc.io/）上的 Slogan 是：A high performance, open source universal RPC framework。就是：一个高性能、开源的通用 RPC 框架。

支持多数主流语言：C#、C++、Dart、**Go**、Java、Kotlin、Node、Objective-C、PHP、Python、Ruby。其中 Go 支持 Windows, Linux, Mac 上的 Go 1.13+ 版本。

gRPC 是一个 Google 开源的高性能远程过程调用 (RPC) 框架，可以在任何环境中运行。它可以通过对负载平衡、跟踪、健康检查和身份验证的可插拔支持有效地连接数据中心内和跨数据中心的服务。它也适用于分布式计算的最后一步，将设备、移动应用程序和浏览器与后端服务接。

![Concept Diagram](https://grpc.io/img/landing-2.svg)

在 gRPC 中，客户端应用程序可以直接调用不同机器上的服务器应用程序的方法，就像它是本地对象一样，使您更容易创建分布式应用程序和服务。与许多 RPC 系统一样，gRPC 基于定义服务的思想，指定可以远程调用的方法及其参数和返回类型。在服务端，服务端实现这个接口并运行一个 gRPC 服务器来处理客户端调用。在客户端，客户端有一个存根（在某些语言中仅称为客户端），它提供与服务器相同的方法。

技术上，gRPC 基于 HTTP/2 通信，采用 Protocol Buffers 作数据序列化。

#### 准备 gRPC 环境

使用 gRPC 需要：

- Go
- Protocol Buffer 编译器，`protoc`，推荐版本3
- Go Plugin，用于 Protocol Buffer 编译器

**安装 protoc：**

可以使用 yum 或 apt 包管理器安装，但通常版本会比较滞后。因此更建议使用预编译的二进制安装。

下载地址：

```url
https://github.com/protocolbuffers/protobuf/releases
```

基于系统和版本找到合适的二进制下载并安装。

CentOS  演示：

```shell
# 下载特定版本，当前（2022年08月）最新 21.4
$ curl -LO https://github.com/protocolbuffers/protobuf/releases/download/v21.4/protoc-21.4-linux-x86_64.zip
# 解压到特定目录
$ sudo unzip protoc-21.4-linux-x86_64.zip -d /usr/local
# 如果特定目录中的bin不在环境变量 path 中，手动加入 path

# 测试安装结果，注意版本应该是 3.x
$ protoc --version
libprotoc 3.21.4
```

Win 演示，下载，解压到指定目录，在 CMD 中运行：

```powershell
# 解压到指定目录即可，要保证 protoc/bin 位于环境变量 path 中，可以随处调用
> protoc.exe --version
libprotoc 3.21.4

```

**安装 Go Plugin：**

```powershell
# 下载特定版本，当前（2022年08月）最新 v1.28.1
> go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
# 下载特定版本，当前（2022年08月）最新 v1.2.0
> go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
# 安装完毕后，要保证 $GOPATH/bin 位于环境变量 path 中

# 测试安装结果
> protoc-gen-go --version
protoc-gen-go.exe v1.28.1
> protoc-gen-go-grpc --version
protoc-gen-go-grpc 1.2.0
```

#### Protocol Buffer 的基础使用

默认情况下，gRPC 使用 Protocol Buffers，这是 Google 用于序列化结构化数据的成熟开源机制（尽管它可以与 JSON 等其他数据格式一起使用）。

> Protocol Buffers 的文档：https://developers.google.com/protocol-buffers/docs/overview

使用 Protocol Buffers 的基本步骤是：

1. 使用 protocol buffers 语法定义消息，消息是用于传递的数据
2. 使用 protocol buffers 语法定义服务，服务是 RPC 方法的集合，来使用消息
3. 使用 Protocol Buffer编 译工具 `protoc` 来编译，生成对应语言的代码，例如 Go 的代码

使用 Protocol Buffers 的第一步是在 `.proto` 文件中定义序列化的数据的结构，.proto 文件是普通的文本文件。Protocol Buffers 数据被结构化为消息，其中每条消息都是一个小的信息逻辑记录，包含一系列称为字段的 name-value 对。

除了核心内容外，`.proto` 文件还需要指定语法版本，目前主流的也是最新的 proto3 版本。在 `.proto` 文件的开头指定。

一个简单的产品信息示例：

product.proto

```protobuf
syntax = "proto3";

// 定义 Product 消息
message Product {
  string name = 1;
  int64 id = 2;
  bool is_sale = 3;
}
```

第二步是在 `.proto` 文件中定义 gRPC 服务，将 RPC 方法参数和返回类型指定为 Protocol Buffers 消息，继续编辑 product.proto :

```protobuf
syntax = "proto3";

// 为了生成 go 代码，需要增加 go_package 属性，表示代码所在的包。protoc 会基于包构建目录
option go_package = "./proto-codes";

// 定义 ProductInfo 消息
message ProductInfoResponse {
  string name = 1;
  int64 int64 = 2;
  bool is_sale = 3;
}

// rpc 方法 ProductInfo 需要的参数消息
message ProductInfoRequest {
  int64 int64 = 1;
}

// 定义 Product 服务
service Product {
  // 获取产品信息
  rpc ProductInfo (ProductInfoRequest) returns (ProductInfoResponse) {}
}
```

**第三步是使用 `protoc` 工具**将 `.proto` 定义的消息和包含 rpc 方法的服务编译为目标语言的代码，我们选择 Go 代码。

```shell
$ protoc --go_out=. --go-grpc_out=. product.proto
# --go_out *.pb.go 目录
# --go-grpc_out *_grpc.pb.go 目录
```

其中：

- `*.pb.go` **包含消息类型的定义和操作的相关代码**
- `*_grpc.pb.go` **包含客户端和服务端的相关代码**

生成的代码主要是结构上的封装，在继续使用时，还需要继续充实业务逻辑。

#### 基于 gRPC 的服务间通信示例

示例说明，存在两个服务，订单服务和产品服务。其中：

- 订单服务提供 HTTP 接口，用于完成订单查询。订单中包含产品信息，要利用 grpc 从产品服务获取产品信息
- 产品服务提供 grpc 接口，用于响应微服务内部产品信息查询

本例中，对于 grpc 来说，产品服务为服务端、订单服务为客户端。

同时不考虑其他业务逻辑，例如产品服务也需要对外提供 http 接口等，仅在乎 grpc 的通信示例。同时不考虑服务发现和网关等。

![image.png](https://fynotefile.oss-cn-zhangjiakou.aliyuncs.com/fynote/fyfile/13080/1657590789084/55e0c7d4c1b042cf80aee8d03bd6d20e.png)

编码实现：

**一：基于之前定义的 .proto 文件生成 pb.go 文件**

注意，客户端和服务端，都需要使用生成的 pb.go 文件

**二：实现订单服务**

orderService/httpService.go

```go
package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"log"
	"net/http"
	"orderService/protos/codes"
	"time"
)

var (
	// 目标 grpc 服务器地址
	gRPCAddr = flag.String("grpc", "localhost:50051", "the address to connect to")
	// http 命令行参数
	addr = flag.String("addr", "127.0.0.1", "The Address for listen. Default is 127.0.0.1")
	port = flag.Int("port", 8080, "The Port for listen. Default is 8080.")
)

func main() {
	flag.Parse()
	// 连接 grpc 服务器
	conn, err := grpc.Dial(*gRPCAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	// 实例化 grpc 客户端
	c := codes.NewProductClient(conn)

	// 定义业务逻辑服务，假设为产品服务
	service := http.NewServeMux()
	service.HandleFunc("/orders", func(writer http.ResponseWriter, request *http.Request) {
		// 调用 grpc 方法，完成对服务器资源请求
		ctx, cancel := context.WithTimeout(context.Background(), time.Second)
		defer cancel()
		r, err := c.ProductInfo(ctx, &codes.ProductInfoRequest{
			Int64: 42,
		})
		if err != nil {
			log.Fatalln(err)
		}

		resp := struct {
			ID       int                          `json:"id"`
			Quantity int                          `json:"quantity"`
			Products []*codes.ProductInfoResponse `json:"products"`
		}{
			9527, 1,
			[]*codes.ProductInfoResponse{
				r,
			},
		}
		respJson, err := json.Marshal(resp)
		if err != nil {
			log.Fatalln(err)
		}
        writer.Header().Set("Content-Type", "application/json")
		_, err = fmt.Fprintf(writer, "%s", string(respJson))
		if err != nil {
			log.Fatalln(err)
		}
	})

	// 启动监听
	address := fmt.Sprintf("%s:%d", *addr, *port)
	fmt.Printf("Order service is listening on %s.\n", address)
	log.Fatalln(http.ListenAndServe(address, service))
}
```

**三，实现产品服务**

productService/grpcService.go

```go
package main

import (
	"context"
	"flag"
	"fmt"
	"google.golang.org/grpc"
	"log"
	"net"
	"productService/protos/compiles"
)

//grpc 监听端口
var port = flag.Int("port", 50051, "The server port")

// ProductServer 实现 UnimplementedProductServer
type ProductServer struct {
	compiles.UnimplementedProductServer
}

func (ProductServer) ProductInfo(ctx context.Context, pr *compiles.ProductInfoRequest) (*compiles.ProductInfoResponse, error) {
	return &compiles.ProductInfoResponse{
		Name:   "马士兵 Go 云原生",
		Int64:  42,
		IsSale: true,
	}, nil
}

func main() {
	flag.Parse()
	//设置 tcp 监听器
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", *port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	// 新建 grpc Server
	s := grpc.NewServer()
	// 将 ProductServer 注册到 grpc Server 中
	compiles.RegisterProductServer(s, ProductServer{})
	log.Printf("server listening at %v", lis.Addr())
	// 启动监听
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}

```

测试，访问 order 的 http 接口。获取订单信息中，包含产品信息。

#### gRPC 核心概念
