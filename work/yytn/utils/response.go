package utils

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

// Response 基础响应结构
type Response struct {
	Code int         `json:"code"`
	Data interface{} `json:"data"`
	Msg  string      `json:"msg"`
}

// PageResult 分页响应结构
type PageResult struct {
	List  interface{} `json:"list"`
	Total int64       `json:"total"`
	Page  int         `json:"page"`
	Size  int         `json:"size"`
}

// Ok 返回成功响应
func Ok(c *gin.Context) {
	c.JSON(http.StatusOK, Response{
		Code: 200,
		Data: nil,
		Msg:  "success",
	})
}

// OkWithMessage 返回带消息的成功响应
func OkWithMessage(message string, c *gin.Context) {
	c.JSON(http.StatusOK, Response{
		Code: 200,
		Data: nil,
		Msg:  message,
	})
}

// OkWithData 返回带数据的成功响应
func OkWithData(data interface{}, c *gin.Context) {
	c.JSON(http.StatusOK, Response{
		Code: 200,
		Data: data,
		Msg:  "success",
	})
}

// OkWithDetailed 返回详细的成功响应
func OkWithDetailed(data interface{}, message string, c *gin.Context) {
	c.JSON(http.StatusOK, Response{
		Code: 200,
		Data: data,
		Msg:  message,
	})
}

// Fail 返回失败响应
func Fail(c *gin.Context) {
	c.JSON(http.StatusOK, Response{
		Code: 500,
		Data: nil,
		Msg:  "fail",
	})
}

// FailWithMessage 返回带消息的失败响应
func FailWithMessage(message string, c *gin.Context) {
	c.JSON(http.StatusOK, Response{
		Code: 500,
		Data: nil,
		Msg:  message,
	})
}
