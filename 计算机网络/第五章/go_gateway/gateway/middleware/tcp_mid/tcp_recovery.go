package tcp_mid

import (
	"fmt"
	"runtime/debug"
)

func TCPRecoveryMiddleware() func(c *TcpSliceRouterContext) {
	return func(c *TcpSliceRouterContext) {
		defer func() {
			if err := recover(); err != nil {
				fmt.Printf("error:%v", fmt.Sprint(err))
				fmt.Printf("stack:%v", string(debug.Stack()))
			}
		}()
		c.Next()
	}
}
