
测试 HttpServerRun：
    http://127.0.0.1:8080/ping
    "pong"

测试 HttpsServerRun：
    https://127.0.0.1:4433/ping
     "pong"

测试 TcpServerRun：
    telnet 127.0.0.1 8011
    set msb mashibing
    +OK
    get msb
    "mashibing"

测试 GrpcServerRun：
    运行测试程序：gateway/proxy/grpc_proxy/grpc_server_client/grpc_client.go
    addr: 127.0.0.1:8012
