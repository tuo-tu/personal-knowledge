package router

import (
	"mini-program/api/v1"
	"mini-program/middleware"

	"github.com/gin-gonic/gin"
)

// 轮播图相关路由
func InitBannerRouter(Router *gin.RouterGroup) {
	BannerRouter := Router.Group("banner")
	{
		// 不需要认证的路由 - 供小程序端使用
		BannerRouter.GET("list", v1.GetBannerList) // 获取启用的轮播图列表

		// 需要认证的路由 - 供管理后台使用
		authBannerRouter := BannerRouter.Group("").Use(middleware.JWTAuth())
		{
			authBannerRouter.POST("create", v1.CreateBanner)  // 创建轮播图
			authBannerRouter.PUT("update", v1.UpdateBanner)   // 更新轮播图
			authBannerRouter.DELETE("delete/:id", v1.DeleteBanner) // 删除轮播图
			authBannerRouter.GET("detail/:id", v1.GetBanner)  // 获取轮播图详情
			authBannerRouter.GET("all", v1.GetAllBanner)      // 获取所有轮播图(包括禁用的)
		}
	}
}
