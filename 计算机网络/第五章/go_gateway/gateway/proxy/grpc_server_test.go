package proxy

import (
	"flag"
	"fmt"
	"go_gateway/gateway/proxy/proto"
	"google.golang.org/grpc"
	"log"
	"net"
	"testing"
)

var port = flag.Int("port", 8005, "the port to serve on")

func TestGrpcServer(t *testing.T) {
	flag.Parse()

	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", *port))
	if err != nil {
		log.Fatalf("failed lisenting: %v", err)
	}
	s := grpc.NewServer()
	proto.RegisterEchoServer(s, &server{})
	s.Serve(listener)
}
