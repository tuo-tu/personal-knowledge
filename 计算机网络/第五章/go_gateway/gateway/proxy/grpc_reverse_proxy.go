package proxy

import (
	"context"
	"go_gateway/common"
	"go_gateway/gateway/loadbalance"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"log"
)

func NewGrpcLoadBalanceHandler(lb loadbalance.LoadBalance) grpc.StreamHandler {
	return func() grpc.StreamHandler {
		// 定义入口函数：实用负载均衡算法获取下游主机地址
		director := func(ctx context.Context, fullMethodName string) (context.Context, *grpc.ClientConn, error) {
			nextAddr, err := lb.Get(fullMethodName)
			if err != nil {
				log.Fatal("get next address fail")
			}
			c, err := grpc.DialContext(ctx, nextAddr,
				// 自定义编码
				grpc.WithDefaultCallOptions(grpc.CallContentSubtype(common.Codec().Name())),
				// 禁用安全传输
				grpc.WithTransportCredentials(insecure.NewCredentials()))
			//md, _ := metadata.FromIncomingContext(ctx)
			//outCtx, _ := context.WithCancel(ctx)
			//outCtx = metadata.NewOutgoingContext(outCtx, md.Copy())
			return ctx, c, err
		}
		return TransparentHandler(director)
	}()

	//return func() grpc.StreamHandler {
	//	nextAddr, err := lb.Get("")
	//	if err != nil {
	//		log.Fatal("get next addr fail")
	//	}
	//	director := func(ctx context.Context, fullMethodName string) (context.Context, *grpc.ClientConn, error) {
	//		c, err := grpc.DialContext(ctx, nextAddr,
	//			// 自定义编码
	//			grpc.WithDefaultCallOptions(grpc.CallContentSubtype(common.Codec().Name())),
	//			// 禁用安全传输
	//			grpc.WithTransportCredentials(insecure.NewCredentials()))
	//		md, _ := metadata.FromIncomingContext(ctx)
	//		outCtx, _ := context.WithCancel(ctx)
	//		outCtx = metadata.NewOutgoingContext(outCtx, md.Copy())
	//		return outCtx, c, err
	//	}
	//	return TransparentHandler(director)
	//}()
}
