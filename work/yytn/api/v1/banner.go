package v1

import (
	"mini-program/model"
	"mini-program/service"
	"mini-program/utils/response"

	"github.com/gin-gonic/gin"
)

// CreateBanner 创建轮播图
// @Tags Banner
// @Summary 创建轮播图
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body model.BannerReq true "轮播图信息"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"创建成功"}"
// @Router /banner/create [post]
func CreateBanner(c *gin.Context) {
	var bannerReq model.BannerReq
	if err := c.ShouldBindJSON(&bannerReq); err != nil {
		response.FailWithMessage("参数错误", c)
		return
	}

	if err := service.CreateBanner(bannerReq); err != nil {
		response.FailWithMessage("创建失败", c)
	} else {
		response.OkWithMessage("创建成功", c)
	}
}

// UpdateBanner 更新轮播图
// @Tags Banner
// @Summary 更新轮播图
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body model.BannerReq true "轮播图信息"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"更新成功"}"
// @Router /banner/update [put]
func UpdateBanner(c *gin.Context) {
	var bannerReq model.BannerReq
	if err := c.ShouldBindJSON(&bannerReq); err != nil {
		response.FailWithMessage("参数错误", c)
		return
	}

	if err := service.UpdateBanner(bannerReq); err != nil {
		response.FailWithMessage("更新失败", c)
	} else {
		response.OkWithMessage("更新成功", c)
	}
}

// DeleteBanner 删除轮播图
// @Tags Banner
// @Summary 删除轮播图
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param id path int true "轮播图ID"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"删除成功"}"
// @Router /banner/delete/{id} [delete]
func DeleteBanner(c *gin.Context) {
	id := c.Param("id")
	if err := service.DeleteBanner(id); err != nil {
		response.FailWithMessage("删除失败", c)
	} else {
		response.OkWithMessage("删除成功", c)
	}
}

// GetBanner 获取轮播图详情
// @Tags Banner
// @Summary 获取轮播图详情
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param id path int true "轮播图ID"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"获取成功"}"
// @Router /banner/detail/{id} [get]
func GetBanner(c *gin.Context) {
	id := c.Param("id")
	banner, err := service.GetBanner(id)
	if err != nil {
		response.FailWithMessage("获取失败", c)
	} else {
		response.OkWithData(banner, c)
	}
}

// GetBannerList 获取启用的轮播图列表(供小程序使用)
// @Tags Banner
// @Summary 获取启用的轮播图列表
// @accept application/json
// @Produce application/json
// @Success 200 {string} string "{"success":true,"data":{},"msg":"获取成功"}"
// @Router /banner/list [get]
func GetBannerList(c *gin.Context) {
	banners, err := service.GetBannerList()
	if err != nil {
		response.FailWithMessage("获取失败", c)
	} else {
		response.OkWithData(banners, c)
	}
}

// GetAllBanner 获取所有轮播图(包括禁用的，供管理后台使用)
// @Tags Banner
// @Summary 获取所有轮播图
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data query model.BannerReq true "分页参数"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"获取成功"}"
// @Router /banner/all [get]
func GetAllBanner(c *gin.Context) {
	var bannerReq model.BannerReq
	if err := c.ShouldBindQuery(&bannerReq); err != nil {
		response.FailWithMessage("参数错误", c)
		return
	}

	total, list, err := service.GetAllBanner(bannerReq)
	if err != nil {
		response.FailWithMessage("获取失败", c)
	} else {
		response.OkWithDetailed(response.PageResult{
			List:  list,
			Total: total,
			Page:  bannerReq.Page,
			Size:  bannerReq.Size,
		}, "获取成功", c)
	}
}
