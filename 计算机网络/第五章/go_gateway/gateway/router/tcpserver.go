package router

import (
	"context"
	"fmt"
	"go_gateway/bussiness/mvc/dao"
	"go_gateway/gateway/middleware/tcp_mid"
	"go_gateway/gateway/proxy"
	"log"
	"net"
)

var tcpServerList = []*proxy.TCPServer{}

type tcpHandler struct {
}

func (t *tcpHandler) ServeTCP(ctx context.Context, src net.Conn) {
	src.Write([]byte("tcpHandler\n"))
}

func TcpServerRun() {
	serviceList := dao.ServiceManagerHandler.GetTcpServiceList()
	for _, serviceItem := range serviceList {
		tempItem := serviceItem
		go func(serviceDetail *dao.ServiceDetail) {
			addr := fmt.Sprintf(":%d", serviceDetail.TCPRule.Port)
			rb, err := dao.LoadBalancerHandler.GetLoadBalancer(serviceDetail)
			if err != nil {
				log.Fatalf(" [INFO] GetTcpLoadBalancer %v err:%v\n", addr, err)
				return
			}

			// 构建路由及设置中间件
			router := tcp_mid.NewTcpSliceRouter()
			router.Group("/").Use(
				//tcp_mid.TCPRecoveryMiddleware(),
				tcp_mid.TCPFlowCountMiddleware(),
				tcp_mid.TCPFlowLimitMiddleware(),
				tcp_mid.TCPWhiteListMiddleware(),
				tcp_mid.TCPBlackListMiddleware(),
			)

			// 构建回调handler
			routerHandler := tcp_mid.NewTcpSliceRouterHandler(
				func(c *tcp_mid.TcpSliceRouterContext) proxy.TCPHandler {
					return proxy.NewTcpLoadBalanceReverseProxy(c.Ctx, rb)
				}, router)

			baseCtx := context.WithValue(context.Background(), "service", serviceDetail)
			tcpServer := &proxy.TCPServer{
				Addr:    addr,
				Handler: routerHandler,
				BaseCtx: baseCtx,
			}
			tcpServerList = append(tcpServerList, tcpServer)
			log.Printf(" [INFO] tcp_proxy_run %v\n", addr)
			// 启动TCP服务，并处理服务异常
			// proxy.ErrServerClosed: 不处理服务关闭异常
			if err := tcpServer.ListenAndServe(); err != nil && err != proxy.ErrServerClosed {
				log.Printf(" [INFO] tcp_proxy_run %v err:%v\n", addr, err)
			}
		}(tempItem)
	}
}

func TcpServerStop() {
	for _, tcpServer := range tcpServerList {
		tcpServer.Close()
		log.Printf(" [INFO] tcp_proxy_stop %v stopped\n", tcpServer.Addr)
	}
}
