package router

import (
	"github.com/gin-gonic/gin"
	"go_gateway/gateway/middleware"
	"go_gateway/gateway/middleware/http_mid"
)

func InitRouter(middlewares ...gin.HandlerFunc) *gin.Engine {
	router := gin.New()
	router.Use(middlewares...)
	router.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})

	oauth := router.Group("/oauth")
	oauth.Use(middleware.TranslationMiddleware())
	{
		middleware.OAuthRegister(oauth)
	}

	router.Use(
		http_mid.HTTPAccessModeMiddleware(),
		http_mid.HTTPFlowCountMiddleware(),
		http_mid.HTTPFlowLimitMiddleware(),
		http_mid.HTTPJwtAuthTokenMiddleware(),
		http_mid.HTTPJwtFlowCountMiddleware(),
		http_mid.HTTPJwtFlowLimitMiddleware(),
		http_mid.HTTPWhiteListMiddleware(),
		http_mid.HTTPBlackListMiddleware(),
		http_mid.HTTPHeaderTransferMiddleware(),
		http_mid.HTTPStripUriMiddleware(),
		http_mid.HTTPUrlRewriteMiddleware(),
		http_mid.HTTPReverseProxyMiddleware())

	return router
}
