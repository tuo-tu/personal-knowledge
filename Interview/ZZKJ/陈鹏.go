package main

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/binary"
	"encoding/hex"
	"errors"
	"fmt"
	"hash"
	"io"
	"net"
	"sync"
)

const (
	Key       byte = iota // key 本身
	KeyData               // key 对应的 data
	End                   // 结束标志
	HeaderLen = 9
	ChunkSize = 8 * 1024 // 数据帧分割长度/读缓冲区大小
)

type Header struct {
	MessageType byte // 帧类型
	L           int  // 数据长度（每次）
}

func (h *Header) Marshal() []byte {
	buf := make([]byte, HeaderLen)
	buf[0] = h.MessageType
	binary.BigEndian.PutUint64(buf[1:], uint64(h.L))
	return buf
}

func (h *Header) Unmarshal(data []byte) {
	h.MessageType = data[0]
	h.L = int(binary.BigEndian.Uint64(data[1:]))
}

// Conn 是你需要实现的一种连接类型，它支持下面描述的若干接口；
// 为了实现这些接口，你需要设计一个基于 TCP 的简单协议；
type Conn struct {
	tcpConn net.Conn
}

// Send 传入一个 key 表示发送者将要传输的数据对应的标识；
// 返回 writer 可供发送者分多次写入大量该 key 对应的数据；
// 当发送者已将该 key 对应的所有数据写入后，调用 writer.Close 告知接收者：该 key 的数据已经完全写入；
func (conn *Conn) Send(key string) (writer io.WriteCloser, err error) {
	header := &Header{
		MessageType: Key,
		L:           len(key),
	}

	data := append(header.Marshal(), []byte(key)...)
	if _, err = conn.tcpConn.Write(data); err != nil {
		return nil, err
	}
	writer = &Writer{conn: conn.tcpConn}
	return
}

type Writer struct {
	conn net.Conn
}

// 当此发送了多少数据
func (w *Writer) Write(p []byte) (n int, err error) {
	// 按照接口的定义，这里最多返回len(p)
	// It returns the number of bytes written from p (0 <= n <= len(p))
	header := &Header{
		MessageType: KeyData,
		L:           len(p),
	}

	data := append(header.Marshal(), p...)
	/*  多余的一段
	n = 0
	for n < len(data) {
		end := n + ChunkSize
		if end > len(data) {
			end = len(data)
		}
		wn, werr := w.conn.Write(data[n:end])
		n += wn
		if werr != nil {
			return n - HeaderLen, werr
		}
	}
	*/
	n, err = w.conn.Write(data)
	return n - HeaderLen, nil
}

// 写入结束帧
func (w *Writer) Close() error {
	header := &Header{
		MessageType: End,
		L:           0,
	}
	data := header.Marshal()
	_, err := w.conn.Write(data)
	return err
}

// Receive 返回一个 key 表示接收者将要接收到的数据对应的标识；
// 返回的 reader 可供接收者多次读取该 key 对应的数据；
// 当 reader 返回 io.EOF 错误时，表示接收者已经完整接收该 key 对应的数据；
func (conn *Conn) Receive() (key string, reader io.Reader, err error) {
	// 首先处理 Key 帧
	headerData := make([]byte, HeaderLen)
	if _, err = conn.tcpConn.Read(headerData); err != nil {
		return "", nil, err
	}

	header := &Header{}
	header.Unmarshal(headerData)

	if header.MessageType != Key {
		return "", nil, errors.New("不是Key帧")
	}

	// 读取 Key 帧
	l := header.L
	keyData := make([]byte, l)
	total := 0
	for total < l {
		toRead := ChunkSize
		if l-total < toRead {
			toRead = l - total
		}
		n, err := conn.tcpConn.Read(keyData[total : total+toRead])
		if err != nil {
			return "", nil, err
		}
		total += n
	}

	key = string(keyData)
	reader = &Reader{
		conn:   conn.tcpConn,
		buffer: []byte{}, // 省略也可以。但是保留着可读性更好
	}
	return
}

type Reader struct {
	conn   net.Conn
	buffer []byte
	isEnd  bool
}

// p实际上也是一个buf
func (r *Reader) Read(p []byte) (n int, err error) {
	if len(p) == 0 {
		return 0, nil
	}

	for {
		// 调换位置，前两个if一开始不会触发，一触发就结束，因此放到前面
		// 读到结束帧
		if r.isEnd {
			if len(r.buffer) > 0 {
				n = copy(p, r.buffer)
				r.buffer = r.buffer[n:]
				return n, nil
			}
			return 0, io.EOF
		}

		if len(r.buffer) > 0 {
			n = copy(p, r.buffer)
			// n 不会超出buf的边界
			r.buffer = r.buffer[n:]
			return n, nil
		}

		// 以上不满足，从conn里面读
		headerData := make([]byte, HeaderLen)
		if _, err = r.conn.Read(headerData); err != nil {
			if err == io.EOF {
				r.isEnd = true
				// 继续读缓冲区
				continue
			}
			// 注意这里是读取头信息，还没有读取body，所以返回0，err是对的，而头信息在receive里面处理了。
			return 0, err
		}

		header := &Header{}
		header.Unmarshal(headerData)
		if header.MessageType == End {
			r.isEnd = true
			continue
		}

		// 注意，要有错误才结束，没错误会继续循环
		if err = r.readBuffer(header.L); err != nil {
			// 读到缓冲的时候就出错，当然返回0，err
			return 0, err
		}
	}
}

func (r *Reader) readBuffer(L int) error {
	bufSize := L
	if L > ChunkSize {
		bufSize = ChunkSize
	}

	tmp := make([]byte, bufSize)
	readTotal := 0
	for readTotal < L {
		n, err := r.conn.Read(tmp[:bufSize])
		if err != nil {
			return err
		}
		r.buffer = append(r.buffer, tmp[:n]...)
		readTotal += n

		remaining := L - readTotal
		if remaining < bufSize {
			bufSize = remaining
		}
	}
	return nil
}

// Close 关闭你实现的连接对象及其底层的 TCP 连接
func (conn *Conn) Close() {
	_ = conn.tcpConn.Close()
}

// NewConn 从一个 TCP 连接得到一个你实现的连接对象
func NewConn(conn net.Conn) *Conn {
	return &Conn{
		tcpConn: conn,
	}
}

// 除了上面规定的接口，你还可以自行定义新的类型，变量和函数以满足实现需求

//////////////////////////////////////////////
///////// 接下来的代码为测试代码，请勿修改 /////////
//////////////////////////////////////////////

// 连接到测试服务器，获得一个你实现的连接对象
func dial(serverAddr string) *Conn {
	conn, err := net.Dial("tcp", serverAddr)
	if err != nil {
		panic(err)
	}
	return NewConn(conn)
}

// 启动测试服务器
func startServer(handle func(*Conn)) net.Listener {
	ln, err := net.Listen("tcp", ":0")
	if err != nil {
		panic(err)
	}
	go func() {
		for {
			conn, err := ln.Accept()
			if err != nil {
				fmt.Println("[WARNING] ln.Accept", err)
				return
			}
			go handle(NewConn(conn))
		}
	}()
	return ln
}

// 简单断言
func assertEqual[T comparable](actual T, expected T) {
	if actual != expected {
		panic(fmt.Sprintf("actual:%v expected:%v\n", actual, expected))
	}
}

// 简单 case：单连接，双向传输少量数据
func testCase0() {
	const (
		key  = "Bible"
		data = `Then I heard the voice of the Lord saying, “Whom shall I send? And who will go for us?”
And I said, “Here am I. Send me!”
Isaiah 6:8`
	)
	ln := startServer(func(conn *Conn) {
		// 服务端等待客户端进行传输
		_key, reader, err := conn.Receive()
		if err != nil {
			panic(err)
		}
		assertEqual(_key, key)
		dataB, err := io.ReadAll(reader)
		if err != nil {
			panic(err)
		}
		assertEqual(string(dataB), data)

		// 服务端向客户端进行传输
		writer, err := conn.Send(key)
		if err != nil {
			panic(err)
		}
		n, err := writer.Write([]byte(data))
		if err != nil {
			panic(err)
		}
		if n != len(data) {
			panic(n)
		}
		conn.Close()
	})
	//goland:noinspection GoUnhandledErrorResult
	defer ln.Close()

	conn := dial(ln.Addr().String())
	// 客户端向服务端传输
	writer, err := conn.Send(key)
	if err != nil {
		panic(err)
	}
	n, err := writer.Write([]byte(data))
	if n != len(data) {
		panic(n)
	}
	err = writer.Close()
	if err != nil {
		panic(err)
	}
	// 客户端等待服务端传输
	_key, reader, err := conn.Receive()
	if err != nil {
		panic(err)
	}
	assertEqual(_key, key)
	dataB, err := io.ReadAll(reader)
	if err != nil {
		panic(err)
	}
	assertEqual(string(dataB), data)
	fmt.Println(_key, string(dataB))
	conn.Close()
}

// 生成一个随机 key
func newRandomKey() string {
	buf := make([]byte, 8)
	_, err := rand.Read(buf)
	if err != nil {
		panic(err)
	}
	return hex.EncodeToString(buf)
}

// 读取随机数据，并返回随机数据的校验和：用于验证数据是否完整传输
func readRandomData(reader io.Reader, hash hash.Hash) (checksum string) {
	hash.Reset()
	var buf = make([]byte, 23<<20) //调用者读取时的 buf 大小不是固定的，你的实现中不可假定 buf 为固定值
	for {
		n, err := reader.Read(buf)
		if err == io.EOF {
			break
		}
		if err != nil {
			panic(err)
		}
		_, err = hash.Write(buf[:n])
		if err != nil {
			panic(err)
		}
	}
	checksum = hex.EncodeToString(hash.Sum(nil))
	return checksum
}

// 写入随机数据，并返回随机数据的校验和：用于验证数据是否完整传输
func writeRandomData(writer io.Writer, hash hash.Hash) (checksum string) {
	hash.Reset()
	const (
		dataSize = 500 << 20 //一个 key 对应 500MB 随机二进制数据，dataSize 也可以是其他值，你的实现中不可假定 dataSize 为固定值
		bufSize  = 1 << 20   //调用者写入时的 buf 大小不是固定的，你的实现中不可假定 buf 为固定值
	)
	var (
		buf  = make([]byte, bufSize)
		size = 0
	)
	for i := 0; i < dataSize/bufSize; i++ {
		_, err := rand.Read(buf)
		if err != nil {
			panic(err)
		}
		_, err = hash.Write(buf)
		if err != nil {
			panic(err)
		}
		n, err := writer.Write(buf)
		if err != nil {
			panic(err)
		}
		size += n
	}
	if size != dataSize {
		panic(size)
	}
	checksum = hex.EncodeToString(hash.Sum(nil))
	return checksum
}

// 复杂 case：多连接，双向传输，大量数据，多个不同的 key
func testCase1() {
	var (
		mapKeyToChecksum = map[string]string{}
		lock             sync.Mutex
	)
	ln := startServer(func(conn *Conn) {
		// 服务端等待客户端进行传输
		key, reader, err := conn.Receive()
		if err != nil {
			panic(err)
		}
		var (
			h         = sha256.New()
			_checksum = readRandomData(reader, h)
		)
		lock.Lock()
		checksum, keyExist := mapKeyToChecksum[key]
		lock.Unlock()
		if !keyExist {
			panic(fmt.Sprintln(key, "not exist"))
		}
		assertEqual(_checksum, checksum)

		// 服务端向客户端连续进行 2 次传输
		for _, key := range []string{newRandomKey(), newRandomKey()} {
			writer, err := conn.Send(key)
			if err != nil {
				panic(err)
			}
			checksum := writeRandomData(writer, h)
			lock.Lock()
			mapKeyToChecksum[key] = checksum
			lock.Unlock()
			err = writer.Close() //表明该 key 的所有数据已传输完毕
			if err != nil {
				panic(err)
			}
		}
		conn.Close()
	})
	//goland:noinspection GoUnhandledErrorResult
	defer ln.Close()

	conn := dial(ln.Addr().String())
	// 客户端向服务端传输
	var (
		key = newRandomKey()
		h   = sha256.New()
	)
	writer, err := conn.Send(key)
	if err != nil {
		panic(err)
	}
	checksum := writeRandomData(writer, h)
	lock.Lock()
	mapKeyToChecksum[key] = checksum
	lock.Unlock()
	err = writer.Close()
	if err != nil {
		panic(err)
	}

	// 客户端等待服务端的多次传输
	keyCount := 0
	for {
		key, reader, err := conn.Receive()
		if err == io.EOF {
			// 服务端所有的数据均传输完毕，关闭连接
			break
		}
		if err != nil {
			panic(err)
		}
		_checksum := readRandomData(reader, h)
		lock.Lock()
		checksum, keyExist := mapKeyToChecksum[key]
		lock.Unlock()
		if !keyExist {
			panic(fmt.Sprintln(key, "not exist"))
		}
		assertEqual(_checksum, checksum)
		keyCount++
		fmt.Println(key, checksum)
	}
	assertEqual(keyCount, 2)
	//fmt.Println(_key, string(dataB))
	conn.Close()
}

func main() {
	testCase0()
	testCase1()
}
