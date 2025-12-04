package util

import (
	"context"
	"github.com/gin-gonic/gin"
	"go_gateway/common/log"
)

// ContextWarning 错误日志
func ContextWarning(c context.Context, dltag string, m map[string]interface{}) {
	v := c.Value("trace")
	traceContext, ok := v.(*log.TraceContext)
	if !ok {
		traceContext = log.NewTrace()
	}
	log.Log.TagWarn(traceContext, dltag, m)
}

// ContextError 错误日志
func ContextError(c context.Context, dltag string, m map[string]interface{}) {
	v := c.Value("trace")
	traceContext, ok := v.(*log.TraceContext)
	if !ok {
		traceContext = log.NewTrace()
	}
	log.Log.TagError(traceContext, dltag, m)
}

// ContextNotice 普通日志
func ContextNotice(c context.Context, dltag string, m map[string]interface{}) {
	v := c.Value("trace")
	traceContext, ok := v.(*log.TraceContext)
	if !ok {
		traceContext = log.NewTrace()
	}
	log.Log.TagInfo(traceContext, dltag, m)
}

// ComLogWarning 错误日志
func ComLogWarning(c *gin.Context, dltag string, m map[string]interface{}) {
	traceContext := GetGinTraceContext(c)
	log.Log.TagError(traceContext, dltag, m)
}

// ComLogNotice 普通日志
func ComLogNotice(c *gin.Context, dltag string, m map[string]interface{}) {
	traceContext := GetGinTraceContext(c)
	log.Log.TagInfo(traceContext, dltag, m)
}

// GetGinTraceContext 从gin的Context中获取数据
func GetGinTraceContext(c *gin.Context) *log.TraceContext {
	if c == nil {
		return log.NewTrace()
	}
	traceContext, exists := c.Get("trace")
	if exists {
		if tc, ok := traceContext.(*log.TraceContext); ok {
			return tc
		}
	}
	return log.NewTrace()
}

// GetTraceContext 从Context中获取数据
func GetTraceContext(c context.Context) *log.TraceContext {
	if c == nil {
		return log.NewTrace()
	}
	traceContext := c.Value("trace")
	if tc, ok := traceContext.(*log.TraceContext); ok {
		return tc
	}
	return log.NewTrace()
}
