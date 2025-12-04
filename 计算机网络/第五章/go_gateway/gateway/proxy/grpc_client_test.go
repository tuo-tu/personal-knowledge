package proxy

import (
	"go_gateway/gateway/proxy/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"log"
	"testing"
	"time"
)

func TestGrpcClient(t *testing.T) {
	//conn, err := grpc.Dial("127.0.0.1:8012")
	conn, err := grpc.Dial("127.0.0.1:8012", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Did not connect:%v", err)
		return
	}
	defer conn.Close()

	c := proto.NewEchoClient(conn)

	// 1.调用一元RPC方法
	unaryEchoWithMetadata(c, msg)
	time.Sleep(1 * time.Second)

	// 2.调用服务端流式处理RPC方法
	serverStreamingWithMetadata(c, msg)
	time.Sleep(1 * time.Second)

	// 3.调用客户端流式处理RPC方法
	clientStreamingWithMetadata(c, msg)
	time.Sleep(1 * time.Second)

	// 4.调用双向流式处理RPC方法
	bidirectionalStreamingWithMetadata(c, msg)
	time.Sleep(1 * time.Second)

}
