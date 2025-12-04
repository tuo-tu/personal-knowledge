package utils

import (
	"github.com/golang-module/carbon/v2"
	"gorm.io/gorm"
	"time"
)

// DB 全局数据库连接，实际项目中应使用配置初始化
var DB *gorm.DB

// ConvertTimeToCarbon 将time.Time转换为carbon.DateTime
func ConvertTimeToCarbon(t time.Time) carbon.DateTime {
	return carbon.DateTime{
		DateTime: t,
	}
}
