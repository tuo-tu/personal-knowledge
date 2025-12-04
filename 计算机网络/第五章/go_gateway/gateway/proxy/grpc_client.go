package proxy

import (
	"context"
	"fmt"
	"go_gateway/common"
	"go_gateway/gateway/proxy/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"io"
	"log"
	"time"
)

var (
	msg = "this is client data "
	//AccessToken = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE5ODM0MTc4MzAsImlzcyI6ImFwcF9pZF9hIn0.NbDJ81fJN-3T3g2bE52wJWySz4AVHKR9a2r9w4Jpwb0"
	AccessToken = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE5ODM0MTc5NDcsImlzcyI6ImFwcF9pZF9hIn0.pi4y3qY5RWVam90fLduyps7Sn2Jyp4Etw-MB_Boj9Xs"
)

func unaryEchoWithMetadata(c proto.EchoClient, msg string) {
	fmt.Println("---- UnaryEcho Client -----")

	md := metadata.Pairs("timestamp", time.Now().Format(time.StampNano))
	//md.Append("authorization", "Bearer some-secret-token")
	md.Append("authorization", "Bearer "+AccessToken)
	ctx := metadata.NewOutgoingContext(context.Background(), md)

	resp, err := c.UnaryEcho(ctx, &proto.EchoRequest{Message: msg},
		grpc.CallContentSubtype(common.Codec().Name()))
	if err != nil {
		log.Fatalf("failed to call UnaryEcho method error:%v", err)
	} else {
		fmt.Printf("response:%v\n", resp.Message)
	}
}

func serverStreamingWithMetadata(c proto.EchoClient, msg string) {
	fmt.Println("---- ServerStreaming Client -----")

	md := metadata.Pairs("timestamp", time.Now().Format(time.StampNano))
	md.Append("authorization", "Bearer "+AccessToken)
	ctx := metadata.NewOutgoingContext(context.Background(), md)

	stream, err := c.ServerStreamingEcho(ctx, &proto.EchoRequest{Message: msg},
		grpc.CallContentSubtype(common.Codec().Name()))
	if err != nil {
		log.Fatalf("failed to call ServerStreamingEcho method error:%v", err)
	}

	var rpcError error
	for {
		// err 读取到流末尾，err = io.EOF
		resp, err := stream.Recv()
		if err != nil {
			rpcError = err
			break
		}
		fmt.Printf("response is :%s\n", resp.Message)
	}
	if rpcError != io.EOF {
		log.Fatalf("failed to finish ServerStreaming:%v", rpcError)
	}
}

func clientStreamingWithMetadata(c proto.EchoClient, msg string) {
	fmt.Println("---- ClientStreaming Client -----")

	md := metadata.Pairs("timestamp", time.Now().Format(time.StampNano))
	md.Append("authorization", "Bearer "+AccessToken)
	ctx := metadata.NewOutgoingContext(context.Background(), md)

	stream, err := c.ClientStreamingEcho(ctx, grpc.CallContentSubtype(common.Codec().Name()))
	if err != nil {
		log.Fatalf("failed to call ClientStreamingEcho method error:%v", err)
	}

	for i := 0; i < 5; i++ {
		err := stream.Send(&proto.EchoRequest{Message: msg})
		if err != nil {
			log.Fatalf("Failed to send:%v", err.Error())
		}
	}

	// 获取响应
	resp, err := stream.CloseAndRecv()
	if err != nil {
		log.Fatalf("failed to finish ClientStreaming:%v", err)
	}
	// 处理服务端响应
	fmt.Printf("response:%v\n", resp.Message)
}

func bidirectionalStreamingWithMetadata(c proto.EchoClient, msg string) {
	fmt.Println("---- bidirectionalStreamingWithMetadata Client -----")

	md := metadata.Pairs("timestamp", time.Now().Format(time.StampNano))
	md.Append("authorization", "Bearer "+AccessToken)
	ctx := metadata.NewOutgoingContext(context.Background(), md)

	stream, err := c.BidirectionalStreamingEcho(ctx, grpc.CallContentSubtype(common.Codec().Name()))
	if err != nil {
		log.Fatalf("failed to call ClientStreamingEcho method error:%v", err)
	}

	// 新建协程发送消息
	go func() {
		for i := 0; i < 5; i++ {
			err := stream.Send(&proto.EchoRequest{Message: msg})
			if err != nil {
				log.Fatalf("Failed to send:%v", err.Error())
			}
		}
		//stream.CloseSend()
	}()

	// 获取响应
	var rpcError error
	for {
		// err 读取到流末尾，err = io.EOF
		resp, err := stream.Recv()
		if err != nil {
			rpcError = err
			break
		}
		fmt.Printf("response is :%s\n", resp.Message)
	}
	if rpcError != io.EOF {
		log.Fatalf("failed to finish ServerStreaming:%v", rpcError)
	}
}
