# WebSocket程序设计

## 协议说明

**WebSocket 是一种在单个TCP连接上进行全双工通信的协议**。WebSocket 使得客户端和服务器之间的数据交换变得更加简单，**允许服务端主动向客户端推送数据**。Websocket主要用在B/S架构的应用程序中，在 WebSocket API 中，浏览器和服务器**只需要完成一次握手**，两者之间就直接可以创建持久性的连接， 并进行双向数据传输。它的最大特点就是，**服务器可以主动向客户端推送信息，客户端也可以主动向服务器发送信息**，是**真正的双向平等对话**，属于服务器推送技术的一种。

相比之下，HTTP协议是**请求-响应模型**，客户端与服务端并不平等。

WebSocket 协议在2008年诞生，2011年成为国际标准。现在最新版本浏览器都已经支持了。

WebSocket 是一种**应用层协议**。

WebSocket 的典型特点：

- **基于 TCP 协议（运输层是TCP）**的应用层协议，实现相对简单
- 单个TCP连接上进行全双工通信
- **兼容 HTTP 协议**，默认端口也是80和443
- **握手阶段采用 HTTP 协议**，能通过各种 HTTP 代理服务器
- 数据格式比较轻量，性能开销小，通信高效
- 可以发送文本和二进制数据
- 没有浏览器的同源限制
- **协议标识符是 `ws`或 `wss`**，网址就是 URL，例如：ws://mashibing.com:80/some/path

websocket的典型场景：

- 即时通信
- 协同编辑/编辑
- 实时数据流的拉取与推送

## websocket推送和浏览器端轮询

在BS开发领域，若需要浏览器B即时得到服务器的状态更新，常使用两个方案：

1. 浏览器端轮询
2. 服务器端推送

浏览器轮询：浏览器端，当需要获取最新数据状态时，利用脚本程序循环向服务端发送请求。

服务器推送，服务器端，当状态改变时，将数据发送到浏览器端。

如图所示：

![image.png](https://fynotefile.oss-cn-zhangjiakou.aliyuncs.com/fynote/fyfile/13080/1688106596061/2fdce5d4a24f498f9d0ab697e36d0d71.png)

如果需要服务器端推送，则需要使用websocket协议。当然HTTP/2版本，也支持服务器端推送，但实现上以推送静态资源为主，不能基于业务逻辑推送特定的消息，因此当前的普及使用率websocket还是主流。

## websocket 与 http 的对比

WebSocket通常和HTTP进行对比，如图：

![image.png](https://fynotefile.oss-cn-zhangjiakou.aliyuncs.com/fynote/fyfile/13080/1688106596061/954a23824b1e4cefa5292812230ac3d2.png)

**WebSocket和HTTP 的相同点:**

- 应用层协议
- B/S 架构中使用
- **基于TCP协议**
- 端口默认都是：80和443

**WebSocket和HTTP 的不同点:**

|              | WebSocket | HTTP                |
| ------------ | --------- | ------------------- |
| 通信模式     | 双向      | 单向                |
| 握手         | 双方协商  | 浏览器发起          |
| 服务器端推送 | 支持      | 不支持。H/2支持部分 |

**WebSocket和HTTP 的联系:**

websocket是在http基础上握手升级得到的。

## WebSocket握手过程

通过HTTP请求响应，中的头信息，完成websocket握手（只需一次握手即可）

- 客户端发起握手请求
- 服务器响应握手请求

如图：

<img src="https://fynotefile.oss-cn-zhangjiakou.aliyuncs.com/fynote/fyfile/13080/1688106596061/7f62abad049e4e1fa0437bc215cf6a06.png" alt="image.png" style="zoom:80%;" />

**请求头如下：**

```http
GET /chat HTTP/1.1
# 指定目标服务器的主机
Host: server.chenpeng.com
# 表明客户端希望升级协议为 WebSocket
Upgrade: websocket
# 表明当前连接支持协议升级。
Connection: Upgrade
# 随机生成的Base64编码密钥，用于服务器验证客户端。也是验证服务器端是否支持websocket
Sec-WebSocket-Key: x4JJHMbDL22zLk1GBhXDw==
# 握手阶段协商子协议，可以视为不同业务逻辑的频道
Sec-WebSocket-Protocol: chat
# 指定WebSocket协议版本（常见值为 13，几乎不需要更改）。
Sec-WebSocket-Version: 13
# 发起请求的源
Origin: http://chenpeng.com
```

基于以上请求头，服务器端就知道需要将协议升级为websocket协议，并提供一些验证信息。

**响应头如下：**状态码101表明服务器同意协议升级

```http
HTTP/1.1 101 Switching Protocols
# 确认协议已升级为 WebSocket
Upgrade: websocket
# 连接状态，用于同步客户端请求，标识协议升级。
Connection: Upgrade
# 服务器根据客户端发送的 Sec-WebSocket-Key 和一个固定的 GUID (258EAFA5-E914-47DA-95CA-C5AB0DC85B11) 计算生成的 Base64 编码字符串，用于验证握手请求的合法性
Sec-WebSocket-Accept: s3pPLMBiTxaQ9kYGzzhZRbK+xOo=
# 握手阶段协商子协议（Subprotocol）
Sec-WebSocket-Protocol: chat
```

基于以上响应头，浏览器端就知道服务器端升级成功，并通过了验证。

至此，B/S端可以基于该连接，完成Websocket双向通信了。

## WebSocket的状态码

| 状态码     | 名称                 | 描述                                                                                                |
| :--------- | :------------------- | :-------------------------------------------------------------------------------------------------- |
| 0–999     | -                    | 保留段, 未使用。                                                                                    |
| 1000       | CLOSE_NORMAL         | 正常关闭; 无论为何目的而创建, 该链接都已成功完成任务。                                              |
| 1001       | CLOSE_GOING_AWAY     | 终端离开, 可能因为服务端错误, 也可能因为浏览器正从打开连接的页面跳转离开。                          |
| 1002       | CLOSE_PROTOCOL_ERROR | 由于协议错误而中断连接。                                                                            |
| 1003       | CLOSE_UNSUPPORTED    | 由于接收到不允许的数据类型而断开连接 (如仅接收文本数据的终端接收到了二进制数据)。                   |
| 1004       | -                    | 保留。 其意义可能会在未来定义。                                                                     |
| 1005       | CLOSE_NO_STATUS      | 保留。 表示没有收到预期的状态码。                                                                   |
| 1006       | CLOSE_ABNORMAL       | 保留。 用于期望收到状态码时连接非正常关闭 (也就是说, 没有发送关闭帧)。                              |
| 1007       | Unsupported Data     | 由于收到了格式不符的数据而断开连接 (如文本消息中包含了非 UTF-8 数据)。                              |
| 1008       | Policy Violation     | 由于收到不符合约定的数据而断开连接。 这是一个通用状态码, 用于不适合使用 1003 和 1009 状态码的场景。 |
| 1009       | CLOSE_TOO_LARGE      | 由于收到过大的数据帧而断开连接。                                                                    |
| 1010       | Missing Extension    | 客户端期望服务器商定一个或多个拓展, 但服务器没有处理, 因此客户端断开连接。                          |
| 1011       | Internal Error       | 客户端由于遇到没有预料的情况阻止其完成请求, 因此服务端断开连接。                                    |
| 1012       | Service Restart      | 服务器由于重启而断开连接。 [Ref]                                                                    |
| 1013       | Try Again Later      | 服务器由于临时原因断开连接, 如服务器过载因此断开一部分客户端连接。 [Ref]                            |
| 1014       | -                    | 由 WebSocket                                                                                        |
| 1015       | TLS Handshake        | 保留。 表示连接由于无法完成 TLS 握手而关闭 (例如无法验证服务器证书)。                               |
| 1016–1999 | -                    | 由 WebSocket 标准保留以便未来使用。                                                                 |
| 2000–2999 | -                    | 由 WebSocket 拓展保留使用。                                                                         |
| 3000–3999 | -                    | 可以由库或框架使用。 不应由应用使用。 可以在 IANA 注册, 先到先得。                                  |
| 4000–4999 | -                    | 可以由应用使用。                                                                                    |

## 服务端编码

需要：

- HTTP服务器，net/http 或者 gin（或其他HTTP框架）
- 处理WebSocket协议的包，https://github.com/gorilla/websocket

其中：https://github.com/gorilla/websocket 是github上Go语言Star数最高的websocket包，推荐使用。

安装gorilla/websocket：

```shell
go get github.com/gorilla/websocket
```

实现流程：

1. 创建HTTP服务器
2. 提供特定路由处理websocket协议
3. 升级为websocket协议
4. 处理Websocket消息
   1. 发送消息
   2. 接收消息

编码实现：略

测试：略

