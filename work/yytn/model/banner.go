package model

import (
	"time"

	"github.com/gogf/gf/v2/net/ghttp"
	"github.com/golang-module/carbon/v2"
	"gorm.io/gorm"
)

// Banner 轮播图模型
type Banner struct {
	ID        uint           `gorm:"primarykey" json:"id"`
	Title     string         `json:"title" gorm:"comment:轮播图标题"`       // 轮播图标题
	ImageURL  string         `json:"image_url" gorm:"comment:图片URL"`    // 图片URL
	LinkURL   string         `json:"link_url" gorm:"comment:跳转链接"`     // 跳转链接
	Sort      int            `json:"sort" gorm:"comment:排序值，越大越靠前"`  // 排序值
	Status    int            `json:"status" gorm:"comment:状态 0-禁用 1-启用"` // 状态
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at"`
}

// TableName 轮播图表名
func (Banner) TableName() string {
	return "banner"
}

// BannerReq 轮播图请求参数
type BannerReq struct {
	ID       uint   `json:"id" form:"id"`
	Title    string `json:"title" form:"title"`
	ImageURL string `json:"image_url" form:"image_url"`
	LinkURL  string `json:"link_url" form:"link_url"`
	Sort     int    `json:"sort" form:"sort"`
	Status   int    `json:"status" form:"status"`
	Page     int    `json:"page" form:"page"`
	Size     int    `json:"size" form:"size"`
}

// BannerRes 轮播图响应参数
type BannerRes struct {
	ID        uint           `json:"id"`
	Title     string         `json:"title"`
	ImageURL  string         `json:"image_url"`
	LinkURL   string         `json:"link_url"`
	Sort      int            `json:"sort"`
	Status    int            `json:"status"`
	CreatedAt carbon.DateTime `json:"created_at"`
	UpdatedAt carbon.DateTime `json:"updated_at"`
}

// BindRequest 绑定请求参数
func (r *BannerReq) BindRequest(ctx *ghttp.Request) error {
	return ctx.Parse(r)
}
