package main

import (
	"flag"
	"go_gateway/bussiness/mvc/dao"
	"go_gateway/common"
	"go_gateway/gateway/router"
	"os"
	"os/signal"
	"syscall"
)

//endpoint dashboard后台管理  server代理服务器
//config ./conf/prod/ 对应配置文件夹

var (
	endpoint = flag.String("endpoint", "server", "input endpoint dashboard or server")
	config   = flag.String("config", "./conf/dev/", "input config file like ./conf/dev/")
)

// go run main.go -config=./conf/dev/ -endpoint server
func main() {
	flag.Parse()
	if *endpoint == "" {
		flag.Usage()
		os.Exit(1)
	}
	if *config == "" {
		flag.Usage()
		os.Exit(1)
	}

	common.InitModule(*config)
	defer common.Destroy()
	dao.ServiceManagerHandler.LoadOnce()
	dao.AppManagerHandler.LoadOnce()

	go func() {
		router.HttpServerRun()
	}()
	go func() {
		router.HttpsServerRun()
	}()
	go func() {
		router.TcpServerRun()
	}()
	go func() {
		router.GrpcServerRun()
	}()

	quit := make(chan os.Signal)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	defer func() {
		router.HttpServerStop()
		router.HttpsServerStop()
		router.GrpcServerStop()
		router.TcpServerStop()
	}()
}
