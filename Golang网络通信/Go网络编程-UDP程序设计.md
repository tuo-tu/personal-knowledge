# UDP程序设计

## UDP协议概述

UDP，User Datagram Protocol，用户数据报协议，是一个简单的面向数据报(package-oriented)的传输层协议，规范为：[RFC 768](https://datatracker.ietf.org/doc/html/rfc768)。

UDP提供**数据的不可靠传输**，它一旦把应用程序发给网络层的数据发送出去，就不保留数据备份。缺乏可靠性，缺乏拥塞控制（congestion control）。

## 基本示例

由于UDP是“无连接”的，所以服务器端不需要创建监听套接字，**只需要监听地址，等待客户端与之建立连接，即可通信。**

net包支持的典型UDP函数：

```go
// 解析UDPAddr
func ResolveUDPAddr(network, address string) (*UDPAddr, error)
// 监听UDP地址
func ListenUDP(network string, laddr *UDPAddr) (*UDPConn, error)
// 连接UDP服务器
func DialUDP(network string, laddr, raddr *UDPAddr) (*UDPConn, error)
// UDP读
func (c *UDPConn) ReadFromUDP(b []byte) (n int, addr *UDPAddr, err error)
// UDP写
func (c *UDPConn) WriteToUDP(b []byte, addr *UDPAddr) (int, error)
```

编写示例，一次读写操作：

服务端流程：

- 解析UDP地址
- 监听UDP
- 读内容
- 写内容

```go
func UDPServerBasic() {
	// 1.解析地址
	laddr, err := net.ResolveUDPAddr("udp", ":9876")
	if err != nil {
		log.Fatalln(err)
	}

	// 2.监听
	udpConn, err := net.ListenUDP("udp", laddr)
	if err != nil {
		log.Fatalln(err)
	}
	log.Printf("%s server is listening on %s\n", "UDP", udpConn.LocalAddr().String())
	defer udpConn.Close()

	// 3.读
	buf := make([]byte, 1024)
	rn, raddr, err := udpConn.ReadFromUDP(buf)
	if err != nil {
		log.Fatalln(err)
	}
	log.Printf("received %s from %s\n", string(buf[:rn]), raddr.String())

	// 4.写
	data := []byte("received:" + string(buf[:rn]))
	wn, err := udpConn.WriteToUDP(data, raddr)
	if err != nil {
		log.Fatalln(err)
	}
	log.Printf("send %s(%d) to %s\n", string(data), wn, raddr.String())
}

```

客户端流程：

- 建立连接
- 写操作
- 读操作

```go
func UDPClientBasic() {
	// 1.解析地址
	raddr, err := net.ResolveUDPAddr("udp", "127.0.0.1:9876")
	if err != nil {
		log.Fatalln(err)
	}
    // 2.拨号，建立连接
	udpConn, err := net.DialUDP("udp", nil, raddr)
	if err != nil {
		log.Fatalln(err)
	}
	log.Println(udpConn)

	// 3.写
	data := []byte("Go UDP program")
	wn, err := udpConn.Write(data) // WriteToUDP(data, raddr)
	if err != nil {
		log.Fatalln(err)
	}
	log.Printf("send %s(%d) to %s\n", string(data), wn, raddr.String())

	// 3.读
	buf := make([]byte, 1024)
	rn, raddr, err := udpConn.ReadFromUDP(buf)
	if err != nil {
		log.Fatalln(err)
	}
	log.Printf("received %s from %s\n", string(buf[:rn]), raddr.String())
}
```

测试：

```shell
> go test -run UDPServerBasic
2023/05/25 17:26:34 UDP server is listening on [::]:9876
2023/05/25 17:29:22 received Go UDP program from 127.0.0.1:58657
2023/05/25 17:29:24 send received:Go UDP program(23) to 127.0.0.1:58657
```

```shell
> go test -run UDPClientBasic
2023/05/25 17:29:22 &{{0xc000108f00}}
2023/05/25 17:29:22 send Go UDP program(14) to 127.0.0.1:9876
2023/05/25 17:29:24 received received:Go UDP program from 127.0.0.1:9876
```

## connected和unconnected的UDPConn

UDP的连接分为：

- 已连接（客户端拨号的），connected, 使用方法 DialUDP建立的连接，称为已连接，pre-connected
- 未连接（服务端监听的），unconnected，使用方法 ListenUDP 获得的连接，称为未连接，not connected

如果 `*UDPConn`是 `connected`,读写方法是 `Read`（用`ReadFromUDP`也行）和 `Write`。
如果 `*UDPConn`是 `unconnected`,读写方法是 `ReadFromUDP`和 `WriteToUDP`

主要的差异在写操作上。读操作如果使用混乱，不会影响读操作本身，但一些参数细节上要注意：

**示例：获取远程地址，conn.RemoteAddr()**

unconnected，ListenUDP

```go
func UDPServerConnect() {
	// 1.解析地址
	laddr, err := net.ResolveUDPAddr("udp", ":9876")
	if err != nil {
		log.Fatalln(err)
	}

	// 2.监听
	udpConn, err := net.ListenUDP("udp", laddr)
	if err != nil {
		log.Fatalln(err)
	}
	log.Printf("%s server is listening on %s\n", "UDP", udpConn.LocalAddr().String())
	defer udpConn.Close()

	// 测试输出远程地址
	log.Println(udpConn.RemoteAddr())

	// 3.读
	buf := make([]byte, 1024)
	rn, raddr, err := udpConn.ReadFromUDP(buf)
	if err != nil {
		log.Fatalln(err)
	}
	log.Printf("received %s from %s\n", string(buf[:rn]), raddr.String())

	// 测试输出远程地址
	log.Println(udpConn.RemoteAddr())

	// 4.写
	data := []byte("received:" + string(buf[:rn]))
	wn, err := udpConn.WriteToUDP(data, raddr)
	if err != nil {
		log.Fatalln(err)
	}
	log.Printf("send %s(%d) to %s\n", string(data), wn, raddr.String())

	// 测试输出远程地址
	log.Println(udpConn.RemoteAddr())
}
```

测试：

```shell
> go test -run UDPServerConnect
2023/05/25 18:24:19 UDP server is listening on [::]:9876
2023/05/25 18:24:19 <nil>
2023/05/25 18:24:32 received Go UDP program from 127.0.0.1:63583
2023/05/25 18:24:35 <nil>
2023/05/25 18:24:35 send received:Go UDP program(23) to 127.0.0.1:63583
2023/05/25 18:24:35 <nil>
```

connected，DialUDP

```go
func UDPClientConnect() {
	// 1.建立连接
	raddr, err := net.ResolveUDPAddr("udp", "127.0.0.1:9876")
	if err != nil {
		log.Fatalln(err)
	}
	udpConn, err := net.DialUDP("udp", nil, raddr)
	if err != nil {
		log.Fatalln(err)
	}

	// 测试输出远程地址
	log.Println(udpConn.RemoteAddr())

	// 2.写
	data := []byte("Go UDP program")
	wn, err := udpConn.Write(data) // WriteToUDP(data, raddr)
	if err != nil {
		log.Fatalln(err)
	}
	log.Printf("send %s(%d) to %s\n", string(data), wn, raddr.String())

	// 测试输出远程地址
	log.Println(udpConn.RemoteAddr())

	// 3.读
	buf := make([]byte, 1024)
	rn, raddr, err := udpConn.ReadFromUDP(buf)
	if err != nil {
		log.Fatalln(err)
	}
	log.Printf("received %s from %s\n", string(buf[:rn]), raddr.String())

	// 测试输出远程地址
	log.Println(udpConn.RemoteAddr())
}
```

测试：

```shell
> go test -run UDPClientConnect
2023/05/25 18:24:32 127.0.0.1:9876
2023/05/25 18:24:32 send Go UDP program(14) to 127.0.0.1:9876
2023/05/25 18:24:32 127.0.0.1:9876
2023/05/25 18:24:35 received received:Go UDP program from 127.0.0.1:9876
2023/05/25 18:24:38 127.0.0.1:9876
```

**示例：connected+WriteToUDP错误**：

```go
udpConn, err := net.DialUDP("udp", nil, raddr)
wn, err := udpConn.WriteToUDP(data, raddr)
```

测试：

```shell
> go test -run UDPClientConnect
2023/05/25 18:27:41 127.0.0.1:9876
2023/05/25 18:27:41 write udp 127.0.0.1:52787->127.0.0.1:9876: use of WriteTo with pre-connected connection
```

**示例：unconnected+Write错误**：

```go
udpConn, err := net.ListenUDP("udp", laddr)
wn, err := udpConn.Write(data)
```

测试：

```shell
write udp [::]:9876: wsasend: A request to send or receive data was disallowed because the socket is not connected and (when sending on a datagram socket using a sendto call) no address was supplied.
```

Read的使用尽量遵循原则，但语法上可以混用，内部有兼容处理。

## 对等服务端和客户端

函数

```go
func ListenUDP(network string, laddr *UDPAddr) (*UDPConn, error)
```

可以直接返回UDPConn，是unconnected状态。在编程时，我们的**客户端和服务端都可以使用该函数建立UDP连接**，而不是一定要在客户端使用DialUDP函数。

**这样创建的客户端和服务端时对等的关系**。适用于例如**广播**类的网络通信项目中。

示例代码：

server：

```go
func UDPServerPeer() {
	// 1.解析地址
	laddr, err := net.ResolveUDPAddr("udp", "127.0.0.1:9876")
	if err != nil {
		log.Fatalln(err)
	}

	// 2.监听
	udpConn, err := net.ListenUDP("udp", laddr)
	if err != nil {
		log.Fatalln(err)
	}
	log.Printf("%s server is listening on %s\n", "UDP", udpConn.LocalAddr().String())
	defer udpConn.Close()

	// 远程地址
	raddr, err := net.ResolveUDPAddr("udp", "127.0.0.1:6789")
	if err != nil {
		log.Fatalln(err)
	}

	// 3.读
	buf := make([]byte, 1024)
	rn, raddr, err := udpConn.ReadFromUDP(buf)
	if err != nil {
		log.Fatalln(err)
	}
	log.Printf("received %s from %s\n", string(buf[:rn]), raddr.String())

	// 4.写
	data := []byte("received:" + string(buf[:rn]))
	wn, err := udpConn.WriteToUDP(data, raddr)
	if err != nil {
		log.Fatalln(err)
	}
	log.Printf("send %s(%d) to %s\n", string(data), wn, raddr.String())
}
```

client：

```go
func UDPClientPeer() {
	// 1.解析地址
	laddr, err := net.ResolveUDPAddr("udp", "127.0.0.1:6789")
	if err != nil {
		log.Fatalln(err)
	}
	// 2.监听
	udpConn, err := net.ListenUDP("udp", laddr)
	if err != nil {
		log.Fatalln(err)
	}
	log.Printf("%s server is listening on %s\n", "UDP", udpConn.LocalAddr().String())
	defer udpConn.Close()

	// 远程地址
	raddr, err := net.ResolveUDPAddr("udp", "127.0.0.1:9876")
	if err != nil {
		log.Fatalln(err)
	}

	// 2.写
	data := []byte("Go UDP program")
	wn, err := udpConn.WriteToUDP(data, raddr)
	if err != nil {
		log.Fatalln(err)
	}
	log.Printf("send %s(%d) to %s\n", string(data), wn, raddr.String())

	// 3.读
	buf := make([]byte, 1024)
	rn, raddr, err := udpConn.ReadFromUDP(buf)
	if err != nil {
		log.Fatalln(err)
	}
	log.Printf("received %s from %s\n", string(buf[:rn]), raddr.String())
}
```

测试：

```shell
> go test -run UDPServerPeer
2023/05/25 19:08:34 UDP server is listening on 127.0.0.1:9876
2023/05/25 19:08:46 received Go UDP program from 127.0.0.1:6789
2023/05/25 19:08:46 send received:Go UDP program(23) to 127.0.0.1:6789
```

```shell
> go test -run UDPClientPeer
2023/05/25 19:08:46 UDP server is listening on 127.0.0.1:6789
2023/05/25 19:08:46 send Go UDP program(14) to 127.0.0.1:9876
2023/05/25 19:08:46 received received:Go UDP program from 127.0.0.1:9876
```

## 多播编程（同一组）

多播（Multicast）方式的数据传输是基于 UDP 完成的。与 UDP 服务器端/客户端的单播方式不同，区别是，单播数据传输以单一目标进行，而多播数据同时传递到加入（注册）特定组的大量主机。换言之，采用多播方式时，可以同时向多个主机传递数据。

多播的特点如下：

- 多播发送端针对特定多播组
- 发送端发送 1 次数据，但**该组内的所有接收端都会接收数据**
- 多播组数可以在 IP 地址范围内任意增加

如图所示：

<img src="https://fynotefile.oss-cn-zhangjiakou.aliyuncs.com/fynote/fyfile/13080/1685002732071/f0a4a2e3d7824f3e9d39fc1740cfae2b.png" alt="image.png" style="zoom: 67%;" />

多播组是 D 类IP地址（224.0.0.0~239.255.255.255）：

- 224.0.0.0～224.0.0.255为预留的组播地址（永久组地址），地址224.0.0.0保留不做分配，其它地址供路由协议使用；
- 224.0.1.0～224.0.1.255是公用组播地址，Internetwork Control Block；
- 224.0.2.0～238.255.255.255为用户可用的组播地址（临时组地址），全网范围内有效；
- 239.0.0.0～239.255.255.255为本地管理组播地址，仅在特定的本地范围内有效

Go的标准库net支持多播编程，主要的函数：

```go
func ListenMulticastUDP(network string, ifi *Interface, gaddr *UDPAddr) (*UDPConn, error)
```

用于创建多播的UDP连接。

示例：多播通信

接收端

```go
// 多播接收端
func UDPReceiverMulticast() {
	// 1.多播监听地址
	address := "224.1.1.2:6789"
	gaddr, err := net.ResolveUDPAddr("udp", address)
	if err != nil {
		log.Fatalln(err)
	}

	// 2.多播监听
	udpConn, err := net.ListenMulticastUDP("udp", nil, gaddr)
	if err != nil {
		log.Fatalln(err)
	}
	defer udpConn.Close()
	log.Printf("%s server is listening on %s\n", "UDP", udpConn.LocalAddr().String())

	// 3.接受数据，循环接收
	buf := make([]byte, 1024)
	for {
		rn, raddr, err := udpConn.ReadFromUDP(buf)
		if err != nil {
			log.Println(err)
		}
		log.Printf("received \"%s\" from %s\n", string(buf[:rn]), raddr.String())
	}
}
```

发送端：

```go
// 多播的发送端
func UDPSenderMulticast() {
	// 1.建立UDP多播组连接
	address := "224.1.1.2:6789"
	raddr, err := net.ResolveUDPAddr("udp", address)
	if err != nil {
		log.Fatalln(err)
	}
    // 注意这里用DialUDP，因此用远程地址raddr
	udpConn, err := net.DialUDP("udp", nil, raddr)
	if err != nil {
		log.Fatalln(err)
	}
	defer udpConn.Close()
    
	// 2.发送内容，循环发送
	for {
		data := fmt.Sprintf("[%s]: %s", time.Now().Format("03:04:05.000"), "hello!")
		wn, err := udpConn.Write([]byte(data))
		if err != nil {
			log.Println(err)
		}        
		log.Printf("send \"%s\"(%d) to %s\n", data, wn, raddr.String())

		time.Sleep(time.Second)
	}
}
```

测试：

启动发送端：

```
# go test -run UDPSenderMulticast
2023/05/26 16:36:43 send "[04:36:43.976]: hello!"(22) to 224.1.1.2:6789
2023/05/26 16:36:44 send "[04:36:44.977]: hello!"(22) to 224.1.1.2:6789
2023/05/26 16:36:45 send "[04:36:45.979]: hello!"(22) to 224.1.1.2:6789
2023/05/26 16:36:46 send "[04:36:46.980]: hello!"(22) to 224.1.1.2:6789

```

启动多个接收端，也可以在过程中继续启动：

```
# go test -run UDPReceiverMulticast
2023/05/26 16:36:00 UDP server is listening on 0.0.0.0:6789
2023/05/26 16:36:00 received "[04:36:43.499]: hello!" from 192.168.50.130:56777
2023/05/26 16:36:01 received "[04:36:43.500]: hello!" from 192.168.50.130:56777
2023/05/26 16:36:02 received "[04:36:43.500]: hello!" from 192.168.50.130:56777

```

```
# go test -run UDPReceiverMulticast
2023/05/26 16:36:00 UDP server is listening on 0.0.0.0:6789
2023/05/26 16:36:00 received "[04:36:43.499]: hello!" from 192.168.50.130:56777
2023/05/26 16:36:01 received "[04:36:44.500]: hello!" from 192.168.50.130:56777
2023/05/26 16:36:02 received "[04:36:45.500]: hello!" from 192.168.50.130:56777
```

## 附：Goland远程开发步骤截图：

* 建立ssh连接
* 打开项目

<img src="https://fynotefile.oss-cn-zhangjiakou.aliyuncs.com/fynote/fyfile/13080/1685002732071/394ba849879447da82c6e569f901e350.png" alt="image.png" style="zoom:80%;" />

<img src="https://fynotefile.oss-cn-zhangjiakou.aliyuncs.com/fynote/fyfile/13080/1685002732071/0a534c1a1c9345dab7976f3b444c9ca1.png" alt="image.png" style="zoom: 67%;" />

<img src="https://fynotefile.oss-cn-zhangjiakou.aliyuncs.com/fynote/fyfile/13080/1685002732071/bcad00e09ed047a7864e529f82139450.png" alt="image.png" style="zoom: 67%;" />

<img src="https://fynotefile.oss-cn-zhangjiakou.aliyuncs.com/fynote/fyfile/13080/1685002732071/1e7f7404b9fe4f08b9d0f1b08867d948.png" alt="image.png" style="zoom:67%;" />

<img src="https://fynotefile.oss-cn-zhangjiakou.aliyuncs.com/fynote/fyfile/13080/1685002732071/fa4173713cde495a84dfd54d898b9075.png" alt="image.png" style="zoom:67%;" />

<img src="https://fynotefile.oss-cn-zhangjiakou.aliyuncs.com/fynote/fyfile/13080/1685002732071/6178a55afc4144b79e75fda5831f44aa.png" alt="image.png" style="zoom:67%;" />

<img src="https://fynotefile.oss-cn-zhangjiakou.aliyuncs.com/fynote/fyfile/13080/1685002732071/bbded7a7d5fc493ebb96904e33a421e9.png" alt="image.png" style="zoom: 67%;" />

<img src="https://fynotefile.oss-cn-zhangjiakou.aliyuncs.com/fynote/fyfile/13080/1685002732071/176c243c0a87431580d864f88e1b0a2f.png" alt="image.png" style="zoom:67%;" />

<img src="https://fynotefile.oss-cn-zhangjiakou.aliyuncs.com/fynote/fyfile/13080/1685002732071/e55e0ff4237c4da595faa5ead78291cc.png" alt="image.png" style="zoom:67%;" />

<img src="https://fynotefile.oss-cn-zhangjiakou.aliyuncs.com/fynote/fyfile/13080/1685002732071/9d4879f9d1e84e2186b0e8f07237d44f.png" alt="image.png" style="zoom:67%;" />

<img src="https://fynotefile.oss-cn-zhangjiakou.aliyuncs.com/fynote/fyfile/13080/1685002732071/99a190e01d9c44d1aa0d91682eb58b1f.png" alt="image.png" style="zoom: 80%;" />

## 广播编程（同一网络）

广播地址，Broadcast，指的是将消息发送到在**同一广播网络上的每个主机**。

例如对于网络：

```shell
# ip a
ens33: <BROADCAST,MULTICAST,UP,LOWER_UP>
inet 192.168.50.130/24 brd 192.168.50.255
```

来说，IP ADDR 就是 192.168.50.130/24， **广播地址就是 192.168.50.255**。

意味着，只要发送数据包的目标地址（接收地址）为192.168.50.255时，那么该数据会发送给该网段上的所有计算机。

如图：

<img src="https://fynotefile.oss-cn-zhangjiakou.aliyuncs.com/fynote/fyfile/13080/1685002732071/23bd74af0240497db7f32babbb805a7a.png" alt="image.png" style="zoom: 67%;" />

编码实现：

编码时数据发的发送端，同样使用 `ListenUDP` 方法建立UDP连接，调用WriteToUDP完成数据的发送。就是上面的对等服务端和客户端模式。

接收端：

```go
// 广播接收端
func UDPReceiverBroadcast() {
	// 1.广播监听地址
	laddr, err := net.ResolveUDPAddr("udp", ":6789")
	if err != nil {
		log.Fatalln(err)
	}

	// 2.广播监听
	udpConn, err := net.ListenUDP("udp", laddr)
	if err != nil {
		log.Fatalln(err)
	}
	defer udpConn.Close()
	log.Printf("%s server is listening on %s\n", "UDP", udpConn.LocalAddr().String())

	// 3.接收数据
	// 4.处理数据
	buf := make([]byte, 1024)
	for {
		rn, raddr, err := udpConn.ReadFromUDP(buf)
		if err != nil {
			log.Println(err)
		}
		log.Printf("received \"%s\" from %s\n", string(buf[:rn]), raddr.String())
	}
}
```

发送端：

```go
// 广播发送端
func UDPSenderBroadcast() {
	// 1.监听地址
	// 2.建立连接
	laddr, err := net.ResolveUDPAddr("udp", ":9876")
	if err != nil {
		log.Fatalln(err)
	}
	udpConn, err := net.ListenUDP("udp", laddr)
	if err != nil {
		log.Fatalln(err)
	}
	defer udpConn.Close()
	log.Printf("%s server is listening on %s\n", "UDP", udpConn.LocalAddr().String())

	// 3.发送数据
	// 广播地址
	rAddress := "192.168.50.255:6789"
	raddr, err := net.ResolveUDPAddr("udp", rAddress)
	if err != nil {
		log.Fatalln(err)
	}
	for {
		data := fmt.Sprintf("[%s]: %s", time.Now().Format("03:04:05.000"), "hello!")
		// 广播发送,将字节数据通过 UDP 发送到目标广播地址（raddr）。
		wn, err := udpConn.WriteToUDP([]byte(data), raddr)
		if err != nil {
			log.Println(err)
		}
		log.Printf("send \"%s\"(%d) to %s\n", data, wn, raddr.String())

		time.Sleep(time.Second)
	}
}
```

测试：

接收端：

```shell
# go test -run UDPReceiverBroadcast
2023/06/01 17:13:27 UDP server is listening on [::]:6789
2023/06/01 17:13:34 received "[05:13:34.707]: hello!" from 192.168.50.130:9876
2023/06/01 17:13:35 received "[05:13:35.709]: hello!" from 192.168.50.130:9876

```

发送端：

```go
# go test -run UDPSenderBroadcast
2023/06/01 17:13:34 UDP server is listening on [::]:9876
2023/06/01 17:13:34 send "[05:13:34.707]: hello!"(22) to 192.168.50.255:6789
2023/06/01 17:13:35 send "[05:13:35.709]: hello!"(22) to 192.168.50.255:6789

```

## 文件传输案例

### 案例说明

UDP协议在传输数据时，由于不能保证稳定性传输，因此比较适合**多媒体通信领域**，例如直播、视频、音频即时播放，即时通信等领域。

本案例使用文件传输为例。

客户端设计：

* 发送文件mp3（任意类型都ok）
* 发送文件名
* 发送文件内容

服务端设计：

* 接收文件
* 存储为客户端发送的名字
* 接收文件内容
* 写入到具体文件中

### 编码实现

客户端：

```go
// 文件传输（上传）
func UDPFileClient() {
	// 1.获取文件信息
	filename := "./data/Beyond.mp3"
	// 打开文件
	file, err := os.Open(filename)
	if err != nil {
		log.Fatalln(err)
	}
	// 关闭文件
	defer file.Close()
	// 获取文件信息
	fileinfo, err := file.Stat()
	if err != nil {
		log.Fatalln(err)
	}
	//fileinfo.Size(), fileinfo.Name()
	log.Println("send file size:", fileinfo.Size())

	// 2.连接服务器
	raddress := "192.168.50.131:5678"
	raddr, err := net.ResolveUDPAddr("udp", raddress)
	if err != nil {
		log.Fatalln(err)
	}
    
	udpConn, err := net.DialUDP("udp", nil, raddr)
	if err != nil {
		log.Fatalln(err)
	}
	defer udpConn.Close()

	// 3.发送文件名
	if _, err := udpConn.Write([]byte(fileinfo.Name())); err != nil {
		log.Fatalln(err)
	}

	// 4.服务端确认
	buf := make([]byte, 4*1024)
	rn, err := udpConn.Read(buf)
	if err != nil {
		log.Fatalln(err)
	}
	// 判断是否为文件名正确接收响应
	if "filename ok" != string(buf[:rn]) {
		log.Fatalln(errors.New("server not ready"))
	}

	// 5.发送文件内容
	// 读取文件内容，利用连接发送到服务端
	// file.Read()
	i := 0
	for {
		// 读取文件内容
		rn, err := file.Read(buf)
		if err != nil {
			// io.EOF 错误表示文件读取完毕
			if err == io.EOF {
				break
			}
			log.Fatalln(err)
		}

		// 发送到服务端
		if _, err := udpConn.Write(buf[:rn]); err != nil {
			log.Fatalln(err)
		}
		i++
	}
	log.Println(i)
	// 文件发送完成。
	log.Println("file send complete.")

	// 等待的测试
	time.Sleep(2 * time.Second)
}
```

服务端：

```go
// UDP文件传输
func UDPFileServer() {
	// 1.建立UDP连接
	laddress := ":5678"
	laddr, err := net.ResolveUDPAddr("udp", laddress)
	if err != nil {
		log.Fatalln(err)
	}
	udpConn, err := net.ListenUDP("udp", laddr)
	if err != nil {
		log.Fatalln(err)
	}
	defer udpConn.Close()
	log.Printf("%s server is listening on %s\n", "UDP", udpConn.LocalAddr().String())

	// 2.接收文件名，并确认
	buf := make([]byte, 4*1024)
	rn, raddr, err := udpConn.ReadFromUDP(buf)
	if err != nil {
		log.Fatalln(err)
	}
	filename := string(buf[:rn])
	if _, err := udpConn.WriteToUDP([]byte("filename ok"), raddr); err != nil {
		log.Fatalln(err)
	}

	// 3.接收文件内容，并写入文件
	// 打开文件（创建）
	file, err := os.Create(filename)
	if err != nil {
		log.Fatalln(err)
	}
	defer file.Close()

	// 网络读取
	i := 0
	for {
		// 一次读取
		rn, _, err := udpConn.ReadFromUDP(buf)
		if err != nil {
			log.Fatalln(err)
		}

		// 写入文件
		if _, err := file.Write(buf[:rn]); err != nil {
			log.Fatalln(err)
		}
		i++
		log.Println("file write some content", i)
	}
}
```

测试，将文件从win传输到linux（centos）中。

上传成功，但文件内容未完整接收，注意这个**UDP内容传输的特点（劣势）**。

```shell
# ll
total 16344
-rw-r--r--. 1 root root 9954453 Jun  2 18:08 Beyond.mp3

# ll
total 16344
-rw-r--r--. 1 root root 9757845 Jun  2 18:14 Beyond.mp3

```

对比源文件大小：

```shell
> go test -run UDPFileClient
2023/06/02 18:14:54 send file size: 10409109
```


## 小结

- UDP，User Datagram Protocol，用户数据报协议，是一个简单的面向数据报(package-oriented)的传输层协议
- 单播，点对点
- 多播，组内，使用多播（组播）地址
- 广播，网段内，使用广播地址
- udp连接
  - connected, net.DialUDP, Read, Write
  - unconnected, net.ListenUDP, ReadFromUDP, WriteToUDP
- 场景：
  - 实时性要求高
  - 完整性要求不高
