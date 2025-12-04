package router

import (
	"context"
	"github.com/gin-gonic/gin"
	"go_gateway/common"
	"go_gateway/gateway/middleware"
	"log"
	"net/http"
	"time"
)

var (
	HttpSrvHandler  *http.Server
	HttpsSrvHandler *http.Server
)

func HttpServerRun() {
	gin.SetMode(common.GetStringConf("proxy.base.debug_mode"))
	r := InitRouter(middleware.RecoveryMiddleware(), middleware.RequestLog())
	HttpSrvHandler = &http.Server{
		Addr:           common.GetStringConf("proxy.http.addr"),
		Handler:        r,
		ReadTimeout:    time.Duration(common.GetIntConf("proxy.http.read_timeout")) * time.Second,
		WriteTimeout:   time.Duration(common.GetIntConf("proxy.http.write_timeout")) * time.Second,
		MaxHeaderBytes: 1 << uint(common.GetIntConf("proxy.http.max_header_bytes")),
	}
	log.Printf(" [INFO] http_proxy_run %s\n", common.GetStringConf("proxy.http.addr"))
	if err := HttpSrvHandler.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Printf(" [ERROR] http_proxy_run %s err:%v\n", common.GetStringConf("proxy.http.addr"), err)
	}
}

func HttpsServerRun() {
	gin.SetMode(common.GetStringConf("proxy.base.debug_mode"))
	r := InitRouter(middleware.RecoveryMiddleware(),
		middleware.RequestLog())
	HttpsSrvHandler = &http.Server{
		Addr:           common.GetStringConf("proxy.https.addr"),
		Handler:        r,
		ReadTimeout:    time.Duration(common.GetIntConf("proxy.https.read_timeout")) * time.Second,
		WriteTimeout:   time.Duration(common.GetIntConf("proxy.https.write_timeout")) * time.Second,
		MaxHeaderBytes: 1 << uint(common.GetIntConf("proxy.https.max_header_bytes")),
	}
	log.Printf(" [INFO] https_proxy_run %s\n", common.GetStringConf("proxy.https.addr"))
	if err := HttpsSrvHandler.ListenAndServeTLS(
		// 以下命令只在编译机有效，如果交叉编译需要单独设置路径
		"./certfile/server.crt", "./certfile/server.key"); err != nil && err != http.ErrServerClosed {
		log.Printf(" [ERROR] https_proxy_run %s err:%v\n", common.GetStringConf("proxy.https.addr"), err)
	}
}

func HttpServerStop() {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := HttpSrvHandler.Shutdown(ctx); err != nil {
		log.Printf(" [ERROR] http_proxy_stop err:%v\n", err)
	}
	log.Printf(" [INFO] http_proxy_stop %v stopped\n", common.GetStringConf("proxy.http.addr"))
}

func HttpsServerStop() {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := HttpsSrvHandler.Shutdown(ctx); err != nil {
		log.Fatalf(" [ERROR] https_proxy_stop err:%v\n", err)
	}
	log.Printf(" [INFO] https_proxy_stop %v stopped\n", common.GetStringConf("proxy.https.addr"))
}
