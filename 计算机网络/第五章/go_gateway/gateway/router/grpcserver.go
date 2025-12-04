package router

import (
	"fmt"
	"go_gateway/bussiness/mvc/dao"
	_ "go_gateway/common"
	"go_gateway/gateway/middleware/grpc_mid"
	"go_gateway/gateway/proxy"
	"google.golang.org/grpc"
	"log"
	"net"
)

var grpcServerList = []*wrapGrpcServer{}

type wrapGrpcServer struct {
	Addr string
	*grpc.Server
}

func GrpcServerRun() {
	serviceList := dao.ServiceManagerHandler.GetGrpcServiceList()
	for _, serviceItem := range serviceList {
		tempItem := serviceItem
		go func(serviceDetail *dao.ServiceDetail) {
			addr := fmt.Sprintf(":%d", serviceDetail.GRPCRule.Port)
			rb, err := dao.LoadBalancerHandler.GetLoadBalancer(serviceDetail)
			if err != nil {
				log.Fatalf(" [INFO] GetTcpLoadBalancer %v err:%v\n", addr, err)
				return
			}
			lis, err := net.Listen("tcp", addr)
			if err != nil {
				log.Fatalf(" [INFO] GrpcListen %v err:%v\n", addr, err)
			}
			grpcHandler := proxy.NewGrpcLoadBalanceHandler(rb)
			s := grpc.NewServer(
				grpc.ChainStreamInterceptor(
					grpc_mid.GrpcFlowCountMiddleware(serviceDetail),
					grpc_mid.GrpcFlowLimitMiddleware(serviceDetail),
					//grpc_mid.GrpcJwtAuthTokenMiddleware(serviceDetail),
					//grpc_mid.GrpcJwtFlowCountMiddleware(serviceDetail),
					//grpc_mid.GrpcJwtFlowLimitMiddleware(serviceDetail),
					grpc_mid.GrpcWhiteListMiddleware(serviceDetail),
					grpc_mid.GrpcBlackListMiddleware(serviceDetail),
					grpc_mid.GrpcHeaderTransferMiddleware(serviceDetail),
				),
				grpc.UnknownServiceHandler(grpcHandler))

			grpcServerList = append(grpcServerList, &wrapGrpcServer{
				Addr:   addr,
				Server: s,
			})
			log.Printf(" [INFO] grpc_proxy_run %v\n", addr)
			if err := s.Serve(lis); err != nil {
				log.Printf(" [INFO] grpc_proxy_run %v err:%v\n", addr, err)
			}
		}(tempItem)
	}
}

func GrpcServerStop() {
	for _, grpcServer := range grpcServerList {
		grpcServer.GracefulStop()
		log.Printf(" [INFO] grpc_proxy_stop %v stopped\n", grpcServer.Addr)
	}
}
