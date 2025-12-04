# 业务需求

使用用户角色表，来实现第一组CRUD操作。

涉及的操作：

* 增，添加
* 删，逻辑删除，恢复和永久删除
* 改，更新一个或多个字段
* 查，查询单条或多条，查询多条的翻页，排序和条件过滤

RESTful风格接口。

# models/表模型目录

利用gorm的migrate功能管理表结构。

models/目录，项目的模型目录，全部的模型都在该目录下。

这里指的模型，是table model，表模型。与表结构息息相关。也可以改为 tables/， dao/目录。

# 自定义基础Model

针对json编码的字段名做了控制。

models/base.go

```go
type Model struct {
	ID        uint           `gorm:"primarykey" json:"id"`
	CreatedAt time.Time      `json:"createdAt"`
	UpdatedAt time.Time      `json:"updatedAt"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

```

# 角色的模型定义

models/role.go

```go
// Role 角色模型
type Role struct {
	Model
	Title   string `gorm:"type:varchar(255);uniqueIndex" json:"title"`
	Key     string `gorm:"type:varchar(255);uniqueIndex" json:"key"`
	Enabled bool   `gorm:"" json:"enabled"`
	Weight  int    `gorm:"index;" json:"weight"`
	Comment string `gorm:"type:text" json:"comment"`
}
```

# 表迁移

gorm提供的，基于模型的定义，完成表的创建。

提供模型的初始化功能，（类似于handler），完成迁移。

models/init.go

```go
// Init 初始化模型
func Init() {
	// migrate
	migrate()

	// seed
}

// 表结构迁移
func migrate() {
	// 自动迁移
	if err := utils.DB().AutoMigrate(
		&Role{},
	); err != nil {
		log.Fatalln(err)
	}
}
```

在程序初始阶段，完成模型的初始化：

main.go

```go
func main() {
	// 解析配置
	utils.ParseConfig()
	// 设置应用模式
	utils.SetMode()
	// 设置日志
	utils.SetLogger()

	// 初始化数据库连接
	utils.InitDB()
	// 初始化模型
	models.Init()

	// 初始化路由引擎
	r := handlers.InitEngine()
	// 使用logger输出应用日志
	utils.Logger().Info("service is listening", "addr", viper.GetString("app.addr"))
	r.Run(viper.GetString("app.addr")) // 监听并在 0.0.0.0:8080 上启动服务
}
```

运行时，debug模式下，会gorm使用info级别的日志，将全部的表操作SQL打印：

```
$ go run .

2024/07/22 17:42:51 D:/apps/mashibing/ginCms/models/init.go:19
[1.194ms] [rows:-] SELECT DATABASE()

2024/07/22 17:42:51 D:/apps/mashibing/ginCms/models/init.go:19
[7.550ms] [rows:1] SELECT SCHEMA_NAME from Information_schema.SCHEMATA where SCHEMA_NAME LIKE 'gincms%' ORDER BY SCHEMA_NAME='gincms' DESC,SCHEMA_NAME limit 1

2024/07/22 17:42:51 D:/apps/mashibing/ginCms/models/init.go:19
[48.148ms] [rows:-] SELECT count(*) FROM information_schema.tables WHERE table_schema = 'gincms' AND table_name = 'role' AND table_type = 'BASE TABLE'

2024/07/22 17:42:51 D:/apps/mashibing/ginCms/models/init.go:19
[120.247ms] [rows:0] CREATE TABLE `role` (`id` bigint unsigned AUTO_INCREMENT,`created_at` datetime(3) NULL,`updated_at` datetime(3) NULL,`deleted_at` datetime(3) NULL,`title` varchar(255),`key` varchar(255),`enabled` boolean,`weight` bigint,`comment` text,PRIMARY KEY (`id`),INDEX `idx_role_deleted_at` (`deleted_at`),UNIQUE INDEX `idx_role_title` (`title`),UNIQUE INDEX `idx_role_key` (`key`),INDEX `idx_role_weight` (`weight`))
[GIN-debug] [WARNING] Creating an Engine instance with the Logger and Recovery middleware already attached.

[GIN-debug] [WARNING] Running in "debug" mode. Switch to "release" mode in production.
 - using env:   export GIN_MODE=release
 - using code:  gin.SetMode(gin.ReleaseMode)

[GIN-debug] GET    /ping                     --> ginCms/handlers/system.Ping (3 handlers)
{"time":"2024-07-22T17:42:51.4197693+08:00","level":"INFO","msg":"service is listening","addr":":8084"}
[GIN-debug] [WARNING] You trusted all proxies, this is NOT safe. We recommend you to set a value.
Please check https://pkg.go.dev/github.com/gin-gonic/gin#readme-don-t-trust-all-proxies for details.
[GIN-debug] Listening and serving HTTP on :8084

```

通过容器内的mysql服务来查看：cmd。

```
>docker exec -it gincms-mysql-1 mysql -p
Enter password:
Welcome to the MySQL monitor.  Commands end with ; or \g.
Your MySQL connection id is 11
Server version: 8.2.0 MySQL Community Server - GPL

Copyright (c) 2000, 2023, Oracle and/or its affiliates.

Oracle is a registered trademark of Oracle Corporation and/or its
affiliates. Other names may be trademarks of their respective
owners.

Type 'help;' or '\h' for help. Type '\c' to clear the current input statement.

mysql>
```

```
> desc gincms.role;
+------------+-----------------+------+-----+---------+----------------+
| Field      | Type            | Null | Key | Default | Extra          |
+------------+-----------------+------+-----+---------+----------------+
| id         | bigint unsigned | NO   | PRI | NULL    | auto_increment |
| created_at | datetime(3)     | YES  |     | NULL    |                |
| updated_at | datetime(3)     | YES  |     | NULL    |                |
| deleted_at | datetime(3)     | YES  | MUL | NULL    |                |
| title      | varchar(255)    | YES  | UNI | NULL    |                |
| key        | varchar(255)    | YES  | UNI | NULL    |                |
| enabled    | tinyint(1)      | YES  |     | NULL    |                |
| weight     | bigint          | YES  | MUL | NULL    |                |
| comment    | text            | YES  |     | NULL    |                |
+------------+-----------------+------+-----+---------+----------------+
9 rows in set (0.00 sec)
```

# 数据填充

当存在预制数据时，使用Seed，数据填充功能。

角色：默认的几个角色。

* 管理员
* 普通用户

更新 models/init.go

```go
// Init 初始化模型
func Init() {
	// migrate
	migrate()

	// seed
	seed()
}

// 数据填充
func seed() {
	roleSeed()
}


```

由具体的模型，提供具体的填充数据，需要填充的模型表，提供自己的seed方法，完成填充。在init中的seed负责调用具体的模型seed方法。

models/role.go

```go
// 填充数据
func roleSeed() {
	// 构建数据
	rows := []Role{
		{
			Title:   "管理员",
			Key:     "administrator",
			Enabled: true,
			Model:   Model{ID: 1},
		},
		{
			Title:   "常规用户",
			Key:     "regular",
			Enabled: true,
			Model:   Model{ID: 2},
		},
	}

	// 插入
	for _, row := range rows {
		if err := utils.DB().FirstOrCreate(&row, row.ID).Error; err != nil {
			utils.Logger().With(err.Error())
		}
	}
}
```

运行测试即可。

cmd mysql 查看：

```
mysql> set names utf8mb4;
Query OK, 0 rows affected (0.00 sec)

mysql> select * from gincms.role\G
*************************** 1. row ***************************
        id: 1
created_at: 2024-07-22 18:00:30.390
updated_at: 2024-07-22 18:00:30.390
deleted_at: NULL
     title: 管理员
       key: administrator
   enabled: 1
    weight: 0
   comment:
*************************** 2. row ***************************
        id: 2
created_at: 2024-07-22 18:00:30.490
updated_at: 2024-07-22 18:00:30.490
deleted_at: NULL
     title: 常规用户
       key: regular
   enabled: 1
    weight: 0
   comment:
2 rows in set (0.01 sec)
```

# 角色查询单条

流程：

1. 接收主键ID参数
2. 根据主键ID查询
3. 响应查询结果

## 定义路由

handlers/role/router.go

```go
package role

import "github.com/gin-gonic/gin"

func Router(r *gin.Engine) {
	r.GET("role", GetRow) // GET /role?id=21
}

```

需要在handlers/init.go中，InitEngine方法调用路由定义的方法：

```go
// 初始化路由引擎
func InitEngine() *gin.Engine {
	// 1. 初始化路由引擎
	r := gin.Default()

	// 2. 注册不同模块的路由
	system.Router(r)
	role.Router(r)

	return r
}
```

提供Handler

handlers/role/handler.go

```go
package role

import "github.com/gin-gonic/gin"

func GetRow(ctx *gin.Context) {

}

```

## 解析请求参数

规范的做法：每个API，都有对应的请求参数和响应参数。

### 定义请求消息类型

handlers/role/message.go

```go
package role

// GetRowReq GetRow接口的请求消息类型
type GetRowReq struct {
	ID uint `form:"id" binding:"required,gt=0"`
}
```

### Handler中解析请求消息

handlers/role/handler.go

```go
func GetRow(ctx *gin.Context) {
	// 1. 解析请求数据（消息）
	req := GetRowReq{}
	if err := ctx.ShouldBindQuery(&req); err != nil {
		// 记录日志
		utils.Logger().Error(err.Error())
		// 直接响应
		ctx.JSON(http.StatusOK, gin.H{
			"code":    100,
			"message": err.Error(),
		})
	}

	log.Println(req)
}
```

### API测试

利用api请求工具完成即可：

![image.png](https://fynotefile.oss-cn-zhangjiakou.aliyuncs.com/fynote/fyfile/13080/1721639416053/f33de2fe9b514b74973077c0e24779ce.png)

![image.png](https://fynotefile.oss-cn-zhangjiakou.aliyuncs.com/fynote/fyfile/13080/1721639416053/da58239ecd4e47d8a5af6c4c297c0693.png)

启动service后，发出请求：

![image.png](https://fynotefile.oss-cn-zhangjiakou.aliyuncs.com/fynote/fyfile/13080/1721639416053/708cd10d14b3437ab43f969f43f03d14.png)

![image.png](https://fynotefile.oss-cn-zhangjiakou.aliyuncs.com/fynote/fyfile/13080/1721639416053/48ff03a0f4d14dbfa5f20058170d1580.png)

## 数据查询及响应

过程：

* 模型定义查询方法
* handler调用查询方法

### 模型定义查询方法

models/role.go

```go
// 根据条件查询单条
// assoc 是否查询管理数据
// where, args 查询条件
func RoleFetchRow(assoc bool, where any, args ...any) (*Role, error) {
	// 查询本条
	row := &Role{}
	if err := utils.DB().Where(where, args...).First(&row).Error; err != nil {
		return nil, err
	}

	// 关联查询
	if assoc {
	}

	return row, nil
}
```

### Handler调用查询并响应

handlers/role/handler.go

GetRow()

```go
func GetRow(ctx *gin.Context) {
	// 1. 解析请求数据（消息）
	req := GetRowReq{}
	if err := ctx.ShouldBindQuery(&req); err != nil {
		// 记录日志
		utils.Logger().Error(err.Error())
		// 直接响应
		ctx.JSON(http.StatusOK, gin.H{
			"code":    100,
			"message": err.Error(),
		})
		return
	}

	// 2. 利用模型完成查询
	row, err := models.RoleFetchRow(false, "`id` = ?", req.ID)
	if err != nil {
		// 记录日志
		utils.Logger().Error(err.Error())
		// 直接响应
		ctx.JSON(http.StatusOK, gin.H{
			"code":    100,
			"message": err.Error(),
		})
		return
	}

	// 3. 响应
	ctx.JSON(http.StatusOK, gin.H{
		"code": 0,
		"data": row,
	})
}
```

以上代码，同时完成了响应部分！

### 测试

![image.png](https://fynotefile.oss-cn-zhangjiakou.aliyuncs.com/fynote/fyfile/13080/1721639416053/e80633d158834089a447fcfe499e6817.png)

# 角色查询列表

流程：

1. 定义路由
2. 解析请求参数
   * 过滤参数
   * 排序参数
   * 翻页参数
3. 根据参数查询数据
4. 响应

## 定义路由

handlers/role/router.go

```go
func Router(r *gin.Engine) {
	g := r.Group("role")
	g.GET("", GetRow)      // GET /role?id=21
	g.GET("list", GetList) // GET /role/list?
}
```

使用了路由分组，来集中设置当前模块的路径前缀。

创建Handler

handlers/role/handler.go

```go
func GetList(ctx *gin.Context) {

}
```

## 解析请求参数

### 定义请求类型

列表的请求参数结构类似，都需要包含：

* 过滤
* 排序
* 翻页

几乎任何的列表查询都需要。

将公共的结构，单独定义，再嵌入到具体的查询请求参数中。

定义公共的列表查询类型：

创建handler公共的目录，handlers/common/message.go

```go
// 通用的查询列表过滤类型
type Filter struct {
	// 指针类型表示该字段可以不填
	// omitempty, 非零值才校验
	// gt, 字符串长度>0
	Keyword *string `form:"keyword" binding:"omitempty,gt=0"`
}

// 通用的查询列表排序类型
type Sorter struct {
	// 排序字段
	SortField *string `form:"sortField" binding:"omitempty,gt=0"`
	// 排序方式 asc,desc
	// oneof，多个选项之一
	SortMethod *string `form:"sortMethod" binding:"omitempty,oneof=asc desc"`
}

// 通用的查询列表翻页类型
type Pager struct {
	// 页码索引
	PageNum *int `form:"pageNum" binding:"omitempty,gt=1"`
	// 每页记录数
	PageSize *int `form:"pageSize" binding:"omitempty,gt=0"`
}

```

定义请求参数类型：

hanlers/roles/message.go

```go
// GetListReq GetList请求参数类型
type GetListReq struct {
	// 过滤
	common.Filter
	// 排序
	common.Sorter
	// 翻页
	common.Pager
}
```

### 解析请求参数

handler中解析

handlers/role/handler.go GetList()

```go
func GetList(ctx *gin.Context) {
	// 1. 解析请求消息
	req := GetListReq{}
	if err := ctx.ShouldBindQuery(&req); err != nil {
		// 记录日志
		utils.Logger().Error(err.Error())
		// 直接响应
		ctx.JSON(http.StatusOK, gin.H{
			"code":    100,
			"message": err.Error(),
		})
		return
	}

	log.Println(req)
}
```

逐步测试阶段通过！

### 整理参数

为请求参数初始化默认值！

为每个独立的部分，定义独立的方法进行整理。

handlers/common/message.go

```go
// 一组通用的常量
const (
	PageNumDefault  = 1
	PageSizeDefault = 10
	PageSizeMax     = 100

	SortFieldDefault  = "id"
	SortMethodDefault = "DESC"
)

// Clean 整理Filter
func (f *Filter) Clean() {
	if f.Keyword == nil {
		temp := ""
		f.Keyword = &temp
	}
}

// Clean 整理Sorter
func (s *Sorter) Clean() {
	if s.SortField == nil {
		temp := SortFieldDefault
		s.SortField = &temp
	}
	if s.SortMethod == nil {
		temp := SortMethodDefault
		s.SortMethod = &temp
	}
}

// Clean 整理Pager
func (p *Pager) Clean() {
	if p.PageNum == nil {
		temp := PageNumDefault
		p.PageNum = &temp
	}
	if p.PageSize == nil {
		temp := PageSizeDefault
		p.PageSize = &temp
	}
	if *p.PageSize > PageSizeMax {
		temp := PageSizeMax
		p.PageSize = &temp
	}
}
```

整理全部的请求消息，再每个具体请求类型处定义：

handler/role/message.go

```go
// Clean 查询列表参数清理
func (req *GetListReq) Clean() {
	req.Filter.Clean()
	req.Sorter.Clean()
	req.Pager.Clean()
}
```

handler中，解析请求参数之后，完成参数的清理：

handlers/role/handler.go

```go
func GetList(ctx *gin.Context) {
	// 1. 解析请求消息
	req := GetListReq{}
	if err := ctx.ShouldBindQuery(&req); err != nil {
		// 记录日志
		utils.Logger().Error(err.Error())
		// 直接响应
		ctx.JSON(http.StatusOK, gin.H{
			"code":    100,
			"message": err.Error(),
		})
		return
	}

	// 2. 整理请求参数
	req.Clean()

	log.Println(*req.Keyword, *req.SortField, *req.SortMethod, *req.PageNum, *req.PageSize)
}
```

阶段测试通过！

## 查询参数结构调整

更新结构上的循环问题。

代码架构来说**：handler调用model**，不应该出现model调用handler的情况。包括models和handlers包中的资源。

目前：

Filter，Sorter，Pager都是在Handler中定义。

而模型查询时，需要以上几个数据，那么就出现了反向调用的情况。

思路：将功能由model提供，如果handler需要，嵌入即可。

models/base.go

```go
package models

import (
	"gorm.io/gorm"
	"time"
)

type Model struct {
	ID        uint           `gorm:"primarykey" json:"id"`
	CreatedAt time.Time      `json:"createdAt"`
	UpdatedAt time.Time      `json:"updatedAt"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

// 通用的查询列表过滤类型
type Filter struct {
	// 指针类型表示该字段可以不填
	// omitempty, 非零值才校验
	// gt, 字符串长度>0
	Keyword *string `form:"keyword" binding:"omitempty,gt=0"`
}

// 通用的查询列表排序类型
type Sorter struct {
	// 排序字段
	SortField *string `form:"sortField" binding:"omitempty,gt=0"`
	// 排序方式 asc,desc
	// oneof，多个选项之一
	SortMethod *string `form:"sortMethod" binding:"omitempty,oneof=asc desc"`
}

// 通用的查询列表翻页类型
type Pager struct {
	// 页码索引
	PageNum *int `form:"pageNum" binding:"omitempty,gt=0"`
	// 每页记录数
	PageSize *int `form:"pageSize" binding:"omitempty,gt=0"`
}

const (
	PageNumDefault  = 1
	PageSizeDefault = 10
	PageSizeMax     = 100

	SortFieldDefault  = "id"
	SortMethodDefault = "DESC"
)

// Clean 整理Filter
func (f *Filter) Clean() {
	if f.Keyword == nil {
		temp := ""
		f.Keyword = &temp
	}
}

// Clean 整理Sorter
func (s *Sorter) Clean() {
	if s.SortField == nil {
		temp := SortFieldDefault
		s.SortField = &temp
	}
	if s.SortMethod == nil {
		temp := SortMethodDefault
		s.SortMethod = &temp
	}
}

// Clean 整理Pager
func (p *Pager) Clean() {
	if p.PageNum == nil {
		temp := PageNumDefault
		p.PageNum = &temp
	}
	if p.PageSize == nil {
		temp := PageSizeDefault
		p.PageSize = &temp
	}
	if *p.PageSize > PageSizeMax {
		temp := PageSizeMax
		p.PageSize = &temp
	}
}

```

handler层还需要这些类型，通过嵌入的方案进行重用：

handles/role/message.go 由models类型组合而来：

```
// GetListReq GetList请求参数类型
type GetListReq struct {
	// 过滤
	models.Filter
	// 排序
	models.Sorter
	// 翻页
	models.Pager
}

```

handlers/common/message.go 不需要了。

将Filter整理到具体的model中，Filter与模型直接相关，而Sorter和Pager通用的。

Filter的代码定义，由models/base.go 转移到 models/role.go，并更新为和Role相关的Filter：

models/role.go

```go
// 通用的查询列表过滤类型
type RoleFilter struct {
	// 指针类型表示该字段可以不填
	// omitempty, 非零值才校验
	// gt, 字符串长度>0
	Keyword *string `form:"keyword" binding:"omitempty,gt=0"`
}

// Clean 整理Filter
func (f *RoleFilter) Clean() {
	if f.Keyword == nil {
		temp := ""
		f.Keyword = &temp
	}
}
```

handlers/role/message.go

```go
// GetListReq GetList请求参数类型
type GetListReq struct {
	// 过滤
	models.RoleFilter
	// 排序
	models.Sorter
	// 翻页
	models.Pager
}

// Clean 查询列表参数清理
func (req *GetListReq) Clean() {
	req.RoleFilter.Clean()
	req.Sorter.Clean()
	req.Pager.Clean()
}
```

测试通过！功能未变，结构改变。满足handler调用model的理念！

## 查询数据并响应

### 模型实现查询方法

models/role.go

查询，需要考虑，过滤，排序，和翻页。

```go
// RoleFetchList 查询列表
// @param assoc bool 是否查询关联
// @param filter RoleFilter 过滤参数
// @param sorter Sorter 排序参数
// @param pager Pager 翻页参数
// @return []*Role Role列表
// @return error
func RoleFetchList(assoc bool, filter RoleFilter, sorter Sorter, pager Pager) ([]*Role, error) {
	// 初始化query
	query := utils.DB().Model(&Role{})

	// 1. 过滤
	if *filter.Keyword != "" {
		query.Where("`title` LIKE ?", "%"+*filter.Keyword+"%")
	}
	// 其他字段过滤

	// 2. 排序
	query.Order(fmt.Sprintf("`%s` %s", *sorter.SortField, strings.ToUpper(*sorter.SortMethod)))

	// 3. 翻页 offset limit
	// 在pagesize>0时，才进行翻页
	if *pager.PageSize > 0 {
		// 偏移 ==（当前页码 - 1） 乘以 每页记录数
		offset := (*pager.PageNum - 1) * *pager.PageSize
		query.Offset(offset).Limit(*pager.PageSize)
	}

	// 4. 查询
	var rows []*Role
	if err := query.Find(&rows).Error; err != nil {
		return nil, err
	}

	// 5. 关联查询
	if assoc {
	}

	// 返回
	return rows, nil
}
```

### handler调用并响应

handlers/role/handler.go

```go
func GetList(ctx *gin.Context) {
	// 1. 解析请求消息
	req := GetListReq{}
	if err := ctx.ShouldBindQuery(&req); err != nil {
		// 记录日志
		utils.Logger().Error(err.Error())
		// 直接响应
		ctx.JSON(http.StatusOK, gin.H{
			"code":    100,
			"message": err.Error(),
		})
		return
	}

	// 2. 整理请求参数
	req.Clean()

	log.Println(*req.Keyword, *req.SortField, *req.SortMethod, *req.PageNum, *req.PageSize)

	// 3. 基于model查询
	rows, err := models.RoleFetchList(false, req.RoleFilter, req.Sorter, req.Pager)
	if err != nil {
		// 记录日志
		utils.Logger().Error(err.Error())
		// 直接响应
		ctx.JSON(http.StatusOK, gin.H{
			"code":    100,
			"message": "查询错误",
		})
		return
	}

	// 4. 响应
	ctx.JSON(http.StatusOK, gin.H{
		"code": 0,
		"data": rows,
	})

}
```

## 测试

![image.png](https://fynotefile.oss-cn-zhangjiakou.aliyuncs.com/fynote/fyfile/13080/1721639416053/2a88ea5578664d84a67597a0e0580072.png)

![image.png](https://fynotefile.oss-cn-zhangjiakou.aliyuncs.com/fynote/fyfile/13080/1721639416053/ccb012c34def4e4a9d0c056f78a2d86c.png)

接口测试通过！

# 角色添加

步骤：

1. 开放路由。（API一切从路由开始）
2. 解析添加参数
3. 插入到数据库
4. 响应

## 路由创建

handlers/role/router.go

```go
func Router(r *gin.Engine) {
	g := r.Group("role")
	g.GET("", GetRow)      // GET /role?id=21
	g.GET("list", GetList) // GET /role/list?
	g.POST("", Add)        // GET /role/list?
}
```

资源相同，意味着URI一致。操作不同，通过不同的Method来区分，就是Restful风格的API。

提供handler

hanlers/role/hanler.go

```go
func Add(ctx *gin.Context) {

}
```

## 请求数据参数

### 定义请求数据类型

需要校验的字段，需要额外定义。

请求数据和模型不同的字段，需要额外定义。

handlers/role/message.go

```go
// 添加请求消息
type AddReq struct {
	models.Role
	// 需要额外校验的字段
	Title string `json:"title" binding:"required"`
	Key   string `json:"key" binding:"required"`
}
```

### 解析请求数据

handlers/role/handler.go

```go
func Add(ctx *gin.Context) {
	req := AddReq{}
	if err := ctx.ShouldBind(&req); err != nil {
		// 记录日志
		utils.Logger().Error(err.Error())
		// 直接响应
		ctx.JSON(http.StatusOK, gin.H{
			"code":    100,
			"message": err.Error(),
		})
		return
	}

	log.Println(req)
}
```

阶段测试通过，获取了解析数据。

![image.png](https://fynotefile.oss-cn-zhangjiakou.aliyuncs.com/fynote/fyfile/13080/1721639416053/06fda6e99ecb49ee8486a4dbc585b36c.png)

## 数据入库

步骤：

1. 从 AddReq 获取 Role 模型
2. 模型提供入库方法
3. handler调用

### 从 AddReq 获取 Role 模型

handlers/role/message.go

```go
// AddReq to Role
func (req AddReq) ToRole() *models.Role {
	row := req.Role
	row.Title = req.Title
	row.Key = req.Key
	return &row
}
```

### 模型提供入库方法

models/role.go

```go
func RoleInsert(row *Role) error {
	// 将insert操作在事务里完成，插入时，有时会涉及到关联数据的处理。
	// 数据及关联数据的插入，放在一个事务中
	return utils.DB().Transaction(func(tx *gorm.DB) error {
		// 完成插入
		if err := tx.Create(&row).Error; err != nil {
			return err
		}
		return nil
	})
}
```

### Handler中调用实现数据插入

handlers/role/handler.go

```go
func Add(ctx *gin.Context) {
	// 1. 解析请求数据
	req := AddReq{}
	if err := ctx.ShouldBind(&req); err != nil {
		// 记录日志
		utils.Logger().Error(err.Error())
		// 直接响应
		ctx.JSON(http.StatusOK, gin.H{
			"code":    100,
			"message": err.Error(),
		})
		return
	}

	// 2. 利用模型完成插入
	role := req.ToRole()
	if err := models.RoleInsert(role); err != nil {
		// 记录日志
		utils.Logger().Error(err.Error())
		// 直接响应
		ctx.JSON(http.StatusOK, gin.H{
			"code":    100,
			"message": "数据插入错误",
		})
		return
	}

	// 3. 响应
	// 往往需要重新查询一边，获取最新的role信息
	row, err := models.RoleFetchRow(false, "`id` = ?", role.ID)
	if err != nil {
		// 记录日志
		utils.Logger().Error(err.Error())
		// 直接响应
		ctx.JSON(http.StatusOK, gin.H{
			"code":    100,
			"message": "查询错误",
		})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"code": 100,
		"data": row,
	})
}
```

上面代码，响应前，查询了一下最新的数据，保证处理好了查询的关联、最新更新等。

### 提供一个只传递ID的查询单条方法

整理下代码。

models/role.go

```go
func RoleFetch(id uint, assoc bool) (*Role, error) {
	return RoleFetchRow(assoc, "`id` = ?", id)
}
```

handlers/role/handler.go中直接调用：

```go
	// 3. 响应
	// 往往需要重新查询一边，获取最新的role信息
	row, err := models.RoleFetch(role.ID, false)
	if err != nil {
		// 记录日志
		utils.Logger().Error(err.Error())
		// 直接响应
		ctx.JSON(http.StatusOK, gin.H{
			"code":    100,
			"message": "查询错误",
		})
		return
	}
```

## 测试

![image.png](https://fynotefile.oss-cn-zhangjiakou.aliyuncs.com/fynote/fyfile/13080/1721639416053/660880f5d5d84e1fbcae27a8c330681a.png)

![image.png](https://fynotefile.oss-cn-zhangjiakou.aliyuncs.com/fynote/fyfile/13080/1721639416053/2fe96628b4334675abe51ecfab3fe3d8.png)

# 角色删除

涉及三个操作：

1. 删除
2. 还原
3. 永久删除

还提供回收站功能！回收站本身是查询功能，查询的条件是已删除的记录。

支持回收站的条件，数据记录不是delete删除，而是逻辑删除。

逻辑删除：通过在记录上，增加标识，判断记录是否被删除。

使用 gorm，使用了 deleted_at字段，记录被删除的时间。该字段为null，标识未删除。

```go
type Model struct {
	ID        uint           `gorm:"primarykey" json:"id"`
	CreatedAt time.Time      `json:"createdAt"`
	UpdatedAt time.Time      `json:"updatedAt"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}
```

## 路由开始

handlers/role/router.go

```go
func Router(r *gin.Engine) {
	g := r.Group("role")
	g.GET("", GetRow)      // GET /role?id=21
	g.GET("list", GetList) // GET /role/list?
	g.POST("", Add)        // GET /role/list?

	g.DELETE("", Delete) // DELETE /role?id=22&id=33&id=44
}
```

提供Delete的处理器：

handlers/role/handler.go

```go
func Delete(ctx *gin.Context) {

}
```

## 解析删除请求参数

支持删除多条，提供几个ID，就删除几个记录。

响应给用户删除的记录数。

### 定义参数类型

应该可以解析多个id参数。

handlers/role/message.go

```go
// DeleteReq 删除的请求消息
type DeleteReq struct {
	IDList []uint `form:"id" binding:"gt=0"`
}
```

至少要有1个id。

### 解析

在Delete的处理器中解析：

handler/role/handler.go

Delete()

```go
func Delete(ctx *gin.Context) {
	// 1. 解析请求数据
	req := DeleteReq{}
	if err := ctx.ShouldBind(&req); err != nil {
		// 记录日志
		utils.Logger().Error(err.Error())
		// 直接响应
		ctx.JSON(http.StatusOK, gin.H{
			"code":    100,
			"message": err.Error(),
		})
		return
	}
	log.Println(req)
}
```

阶段测试通过：

![image.png](https://fynotefile.oss-cn-zhangjiakou.aliyuncs.com/fynote/fyfile/13080/1721639416053/2b24687928934c37b90a5a8e7a51b378.png)

## 删除数据及响应

定义模型的方法，完成删除，再处理器中调用。

计划删除后，响应给用户删除的记录数。

### 定义模型删除方法

models/role.go

```go
// RoleDelete 角色删除
// @return 删除的记录数，error
func RoleDelete(idList []uint) (int64, error) {
	// 将delete操作在事务里完成，删除时，有时会涉及到关联数据的处理
	rowsNum := int64(0)
	err := utils.DB().Transaction(func(tx *gorm.DB) error {
		result := tx.Delete(&Role{}, idList)
		if result.Error != nil {
			return result.Error
		} else {
			// 删除成功
			rowsNum = result.RowsAffected
		}
		return nil
	})

	return rowsNum, err
}
```

### 控制器调用

handlers/role/handler.go

```go
func Delete(ctx *gin.Context) {
	// 1. 解析请求数据
	req := DeleteReq{}
	if err := ctx.ShouldBindQuery(&req); err != nil {
		// 记录日志
		utils.Logger().Error(err.Error())
		// 直接响应
		ctx.JSON(http.StatusOK, gin.H{
			"code":    100,
			"message": err.Error(),
		})
		return
	}
	//log.Println(req)

	// 2. 删除数据
	rowsNum, err := models.RoleDelete(req.IDList)
	if err != nil {
		// 记录日志
		utils.Logger().Error(err.Error())
		// 直接响应
		ctx.JSON(http.StatusOK, gin.H{
			"code":    100,
			"message": "数据删除错误",
		})
		return
	}

	// 3. 响应
	ctx.JSON(http.StatusOK, gin.H{
		"code": 0,
		"data": rowsNum,
	})
}
```

## 测试

通过！

![image.png](https://fynotefile.oss-cn-zhangjiakou.aliyuncs.com/fynote/fyfile/13080/1721639416053/830460ce06744651880f98653711d704.png)

数据表中，被删除的，记录的deleted_at为删除时间:

```
mysql> select id, deleted_at from role;
+----+-------------------------+
| id | deleted_at              |
+----+-------------------------+
|  5 | NULL                    |
|  1 | 2024-07-26 14:41:39.408 |
|  2 | 2024-07-26 14:41:39.408 |
|  3 | 2024-07-26 14:41:39.408 |
+----+-------------------------+
4 rows in set (0.00 sec)
```

删除的记录，不会被直接查询到。

测试查询列表，生成的SQL带有deleted_at的过滤条件：

```sql
[3.664ms] [rows:1] SELECT * FROM `role` WHERE `role`.`deleted_at` IS NULL ORDER BY `id` DESC LIMIT 10
[GIN] 2024/07/26 - 14:45:44 | 200 |      3.9293ms |       127.0.0.1 | GET      "/role/list"

```

# 回收站

## 删除数据的查询

查到那些被删除的记录。

更新查询列表的操作。

回收站和查询列表的区别，是判断是否查询deleted_at is null的记录。

gorm中，通过模型的方法：

```
*DB.Unscoped()
```

来查询全部记录。语法去掉 deleted_at IS NULL 的条件。

相对来说：

```go
// 常规列表，不用处理

// 回收站
*DB.Unscoped().Where("deleted_at IS NOT NULL") 条件进行删除
```

### 定义标识查询范围的常量

models/base.go

```go
// 查询范围的常量
const (
	SCOPE_ALL = iota
	SCOPE_UNDELETED
	SCOPE_DELETED
)
```

### 模型查询列表方法增加范围参数

models/role.go

```go
// RoleFetchList 查询列表
// @param filter RoleFilter 过滤参数
// @param sorter Sorter 排序参数
// @param pager Pager 翻页参数
// @param scope uint8 范围参数
// @param assoc bool 是否查询关联
// @return []*Role Role列表
// @return error
func RoleFetchList(filter RoleFilter, sorter Sorter, pager Pager, scope uint8, assoc bool) ([]*Role, error) {
	// 初始化query
	query := utils.DB().Model(&Role{})

	// 1. 过滤
	// 查询范围
	switch scope {
	case SCOPE_ALL:
		query.Unscoped()
	case SCOPE_DELETED:
		query.Unscoped().Where("`deleted_at` IS NOT NULL")
	case SCOPE_UNDELETED:
		fallthrough
	default:
		// do nothing. default case
	}
	// 条件过滤
	if *filter.Keyword != "" {
		query.Where("`title` LIKE ?", "%"+*filter.Keyword+"%")
	}
	// 其他字段过滤


}
```

上面代码，只修改了参数和范围条件部分。

handler中，调用该方法，是根据需要传递参数。例如

查询列表，handlers/role/handler.go GetList()

```go
// 3. 基于model查询
	rows, err := models.RoleFetchList(req.RoleFilter, req.Sorter, req.Pager, models.SCOPE_UNDELETED, false)
```

## 回收站查询

### 路由

handlers/role/router.go

```go
func Router(r *gin.Engine) {
	g := r.Group("role")
	g.GET("", GetRow)      // GET /role?id=21
	g.GET("list", GetList) // GET /role/list?
	g.POST("", Add)        // GET /role/list?

	g.DELETE("", Delete)      // DELETE /role?id=22&id=33&id=44
	g.GET("recycle", Recycle) // DELETE /role?id=22&id=33&id=44
}
```

### 实现回收站处理器

思路是与 GetList() 这个handler重用代码。不过是获取的范围参数不同。

定义内部方法，实现GetList的全部功能，将变化的部分，设计为参数。在Get List和Recycle中进行调用传参。

handlers/role/handler.go

```go
func list(ctx *gin.Context, scope uint8, assoc bool) {
	// 1. 解析请求消息
	req := GetListReq{}
	if err := ctx.ShouldBindQuery(&req); err != nil {
		// 记录日志
		utils.Logger().Error(err.Error())
		// 直接响应
		ctx.JSON(http.StatusOK, gin.H{
			"code":    100,
			"message": err.Error(),
		})
		return
	}

	// 2. 整理请求参数
	req.Clean()

	// 3. 基于model查询
	rows, err := models.RoleFetchList(req.RoleFilter, req.Sorter, req.Pager, scope, assoc)
	if err != nil {
		// 记录日志
		utils.Logger().Error(err.Error())
		// 直接响应
		ctx.JSON(http.StatusOK, gin.H{
			"code":    100,
			"message": "查询错误",
		})
		return
	}

	// 4. 响应
	ctx.JSON(http.StatusOK, gin.H{
		"code": 0,
		"data": rows,
	})

}
```

在Get List和Recycle中进行调用传参：

handlers/role/handler.go

```go
func Recycle(ctx *gin.Context) {
	list(ctx, models.SCOPE_DELETED, false)
}
func GetList(ctx *gin.Context) {
	list(ctx, models.SCOPE_UNDELETED, true)
}
```

### 测试

通过！

![image.png](https://fynotefile.oss-cn-zhangjiakou.aliyuncs.com/fynote/fyfile/13080/1721639416053/32bdb3fc7bdb486f9554d256a6cb3ffd.png)

## 还原数据

核心思路：将deleted_at字段设置为null。

### 路由

handlers/role/router.go

```go
func Router(r *gin.Engine) {
	g := r.Group("role")
	// 查询一条
	g.GET("", GetRow) // GET /role?id=21
	// 查询多条
	g.GET("list", GetList) // GET /role/list?
	// 添加
	g.POST("", Add) // POST /role
	// 删除
	g.DELETE("", Delete) // DELETE /role?id=22&id=33&id=44
	// 查询回收站
	g.GET("recycle", Recycle) // GET /role/recycle?
	// 回收站还原
	g.PUT("restore", Restore) // PUT /role?id=22&id=33&id=44
}
```

处理器

handlers/role/handler.go

```
func Restore(ctx *gin.Context) {

}
```

### 解析请求参数

定义：handlers/role/message.go

```go
// RestoreReq 还原的请求消息
type RestoreReq struct {
	IDList []uint `form:"id" binding:"gt=0"`
}
```

handler中解析: handlers/role/handler.go

```go
func Restore(ctx *gin.Context) {
	// 1. 解析请求数据
	req := RestoreReq{}
	if err := ctx.ShouldBindQuery(&req); err != nil {
		// 记录日志
		utils.Logger().Error(err.Error())
		// 直接响应
		ctx.JSON(http.StatusOK, gin.H{
			"code":    100,
			"message": err.Error(),
		})
		return
	}
}
```

### 模型还原方法

models/role.go

```go
// RoleRestore 还原
func RoleRestore(idList []uint) (int64, error) {
	// 还原的记录数
	rowsNum := int64(0)
	err := utils.DB().Transaction(func(tx *gorm.DB) error {
		result := tx.Model(&Role{}).Unscoped().Where("`id` IN ?", idList).Update("deleted_at", nil)
		if result.Error != nil {
			return result.Error
		} else {
			// 更新成功
			rowsNum = result.RowsAffected
		}
		return nil
	})

	return rowsNum, err
}
```

### 控制器调用及响应

handlers/role/handler.go

```go
func Restore(ctx *gin.Context) {
	// 1. 解析请求数据
	req := RestoreReq{}
	if err := ctx.ShouldBindQuery(&req); err != nil {
		// 记录日志
		utils.Logger().Error(err.Error())
		// 直接响应
		ctx.JSON(http.StatusOK, gin.H{
			"code":    100,
			"message": err.Error(),
		})
		return
	}

	// 2. 还原数据
	rowsNum, err := models.RoleRestore(req.IDList)
	if err != nil {
		// 记录日志
		utils.Logger().Error(err.Error())
		// 直接响应
		ctx.JSON(http.StatusOK, gin.H{
			"code":    100,
			"message": "数据还原错误",
		})
		return
	}

	// 3. 响应
	ctx.JSON(http.StatusOK, gin.H{
		"code": 0,
		"data": rowsNum,
	})
}
```

### 测试

通过。

![image.png](https://fynotefile.oss-cn-zhangjiakou.aliyuncs.com/fynote/fyfile/13080/1721639416053/f574b9767b8841e5b93263a992636b07.png)

数据表中：

```mysql
mysql> select id, deleted_at from role;
+----+-------------------------+
| id | deleted_at              |
+----+-------------------------+
|  1 | NULL                    |
|  2 | NULL                    |
|  5 | NULL                    |
|  3 | 2024-07-26 14:41:39.408 |
+----+-------------------------+
4 rows in set (0.00 sec)
```

## 永久删除

核心：.Unscope().Delete()

思路：更新删除的操作，增加一个请求参数，表示是否强制（永久）删除。

### 更新删除的请求参数

handlers/role/message.go

```go
// DeleteReq 删除的请求消息
type DeleteReq struct {
	IDList []uint `form:"id" binding:"gt=0"`
	Force  bool   `form:"force" binding:""`
}
```

force == true, 表示永久删除。

### 更新删除的模型方法

提供force参数的支持。

models/role.go

```go
// RoleDelete 角色删除
// @param force bool 是否强制删除
// @return 删除的记录数，error
func RoleDelete(idList []uint, force bool) (int64, error) {
	// 将delete操作在事务里完成，删除时，有时会涉及到关联数据的处理
	rowsNum := int64(0)
	err := utils.DB().Transaction(func(tx *gorm.DB) error {
		query := tx.Model(&Role{})
		// 强制
		if force {
			query.Unscoped()
		}
		result := query.Delete(&Role{}, idList)
		if result.Error != nil {
			return result.Error
		} else {
			// 删除成功
			rowsNum = result.RowsAffected
		}
		return nil
	})

	return rowsNum, err
}
```

### 更新删除控制器调用时的传参

handlers/role/handler.go Delete()

```go
// 2. 删除数据
	rowsNum, err := models.RoleDelete(req.IDList, req.Force)
```

### 测试

通过。

![image.png](https://fynotefile.oss-cn-zhangjiakou.aliyuncs.com/fynote/fyfile/13080/1721639416053/a606deb425a64b61b680ef52fa6102a7.png)

数据表：

```mysql
mysql> select id, deleted_at from role;
+----+-------------------------+
| id | deleted_at              |
+----+-------------------------+
|  1 | NULL                    |
|  2 | NULL                    |
|  3 | 2024-07-26 14:41:39.408 |
+----+-------------------------+
3 rows in set (0.00 sec)
```

# 角色单条更新

更新具体某个角色的多个属性字段。

角色多条更新

更新多个角色的某个属性，将id为11，22，33的enabled设为false。

每个多条的更新，更新内容都是特定的，会存在多个角色多条更新的API。

## 路由设置

handlers/role/router.go

```go
// 更新单条
g.PUT(":id", Edit) // PUT /role/33
```

称为Uri参数，也叫路由参数。

handler

handlers/role/handle.go

```
func Edit(ctx *gin.Context) {}
```

## 解析请求参数

有2个参数：

1. URI上的id参数
2. 更新的内容，请求body参数

### URI参数定义

handlers/role/message.go

```go
// EditUriReq URI上的id参数
type EditUriReq struct {
	ID uint `uri:"id" binding:"required,gt=0"` // 可以考虑加一个id存在的校验
}
```

### URI参数解析

handlers/role/handler.go

```go
func Edit(ctx *gin.Context) {
	// 1. 解析URI请求数据
	uri := EditUriReq{}
	if err := ctx.ShouldBindUri(&uri); err != nil {
		// 记录日志
		utils.Logger().Error(err.Error())
		// 直接响应
		ctx.JSON(http.StatusOK, gin.H{
			"code":    100,
			"message": err.Error(),
		})
		return
	}
}
```

### 主体参数的定义

主体参数，表示更新的属性信息。

设计上，允许用户仅仅传递部分字段。需要在请求参数中，识别出来用户传递了哪些字段。进而才能做到去更新传递的字段属性信息。

因此，设计主体参数类型时，以Model为标准，将需要可能更新的字段，设置为指针类型。通过指针类型是否为nil，来判断用户是否传递了该属性。

例如 string 非指针类型，如果值为空字符串"", 如何判断用户是没有传递还是需要将字段设置为空呢？判断不了。

使用 *string 类型，这时：

* nil 用户未传递
* "" 或其他值，表示用户需要设置该值

handlers/role/message.go

```go
// EditBodyReq 更新主体参数
type EditBodyReq struct {
	Title   *string `json:"title"`
	Key     *string `json:"key"`
	Enabled *bool   `json:"enabled"`
	Weight  *int    `json:"weight"`
	Comment *string `json:"comment"`
}
```

### 主体参数的解析

handlers/role/handler.go

```go
func Edit(ctx *gin.Context) {
	// 1. 解析URI请求数据
	uri := EditUriReq{}
	if err := ctx.ShouldBindUri(&uri); err != nil {
		// 记录日志
		utils.Logger().Error(err.Error())
		// 直接响应
		ctx.JSON(http.StatusOK, gin.H{
			"code":    100,
			"message": err.Error(),
		})
		return
	}

	// 2. 解析Body请求数据
	body := EditBodyReq{}
	if err := ctx.ShouldBind(&body); err != nil {
		// 记录日志
		utils.Logger().Error(err.Error())
		// 直接响应
		ctx.JSON(http.StatusOK, gin.H{
			"code":    100,
			"message": err.Error(),
		})
		return
	}
}
```

阶段测试，通过：

![image.png](https://fynotefile.oss-cn-zhangjiakou.aliyuncs.com/fynote/fyfile/13080/1721639416053/aa2404b829214f5896c5c7b7dcf1195b.png)

![image.png](https://fynotefile.oss-cn-zhangjiakou.aliyuncs.com/fynote/fyfile/13080/1721639416053/e4370d1afa1742d5879eb49c1331565b.png)

## 更新数据映射

可以从参数中识别，哪些用户传递了，哪些用户未传递。

通过将请求数据转换为map操作，将用户传递的数据摘出来。

后边在做数据库操作时，map结构可以用于更新部分字段。

### 定义字段map结构

models/base.go

```go
// key: 字段名
// value: 字段值
type FieldMap = map[string]any
```

### 定义转换为FieldMap的方法

涉及遍历结构体字段的操作。

go 没有提供该功能，需要配合reflect反射实现。

handlers/role/message.go

```go
// EditBodyReq 更新主体参数
type EditBodyReq struct {
	Title   *string `json:"title" field:"title"`
	Key     *string `json:"key" field:"key"`
	Enabled *bool   `json:"enabled" field:"enabled"`
	Weight  *int    `json:"weight" field:"weight"`
	Comment *string `json:"comment" field:"comment"`
}

func (req EditBodyReq) ToFieldMap() models.FieldMap {
	// 1. 初始化map
	m := models.FieldMap{}

	// 2. 利用反射来遍历req结构的全部字段
	reqType := reflect.TypeOf(req)
	reqValue := reflect.ValueOf(req)
	// 通过字段数量，进行遍历
	for i, nums := 0, reqType.NumField(); i < nums; i++ {
		// 获取 field tag
		fieldTag := reqType.Field(i).Tag.Get("field")
		// 存在 field tag才自动处理
		if fieldTag == "" {
			continue
		}
		// 判断字段是否为nil，这个值的判断
		if !reqValue.Field(i).IsNil() {
			if fieldTag == "some_field" {
				// 考虑特殊字段的处理情况
			} else {
				// 放入map
				m[fieldTag] = reqValue.Field(i).Elem().Interface()
			}
		}
	}

	return m
}
```

### handler中调用

handlers/role/handler.go

```go
// 3. req to map
	fieldMap := body.ToFieldMap()
	log.Println(fieldMap)
```

阶段测试欧克！

![image.png](https://fynotefile.oss-cn-zhangjiakou.aliyuncs.com/fynote/fyfile/13080/1721639416053/12d14bd8be11418ca7bb8fa9bbea8196.png)![image.png](https://fynotefile.oss-cn-zhangjiakou.aliyuncs.com/fynote/fyfile/13080/1721639416053/d4b278e7d6574f9992e43cd2f7a2357b.png)

## 数据更新和响应

### 模型中增加更新方法

models/role.go

```go
func RoleUpdates(fieldMap FieldMap, id uint) error {
	return utils.DB().Transaction(func(tx *gorm.DB) error {
		// 完成插入
		if err := tx.Model(&Role{}).Where("`id` = ?", id).Updates(fieldMap).Error; err != nil {
			return err
		}
		return nil
	})
}
```

### Handler调用

handlers/role/handler.go

```go
func Edit(ctx *gin.Context) {
	// 1. 解析URI请求数据

	// 2. 解析Body请求数据

	// 3. req to map
	fieldMap := body.ToFieldMap()
	//log.Println(fieldMap)
	if err := models.RoleUpdates(fieldMap, uri.ID); err != nil {
		// 记录日志
		utils.Logger().Error(err.Error())
		// 直接响应
		ctx.JSON(http.StatusOK, gin.H{
			"code":    100,
			"message": "数据更新错误",
		})
		return
	}

	// 3. 响应
	// 往往需要重新查询一边，获取最新的role信息
	row, err := models.RoleFetch(uri.ID, false)
	if err != nil {
		// 记录日志
		utils.Logger().Error(err.Error())
		// 直接响应
		ctx.JSON(http.StatusOK, gin.H{
			"code":    100,
			"message": "查询错误",
		})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"code": 100,
		"data": row,
	})
}
```

## 测试

![image.png](https://fynotefile.oss-cn-zhangjiakou.aliyuncs.com/fynote/fyfile/13080/1721639416053/c8cf3b67924d4c3aa2e17449bcda1ade.png)

通过！

# 角色修改多条的特定属性

典型的例子：同时修改多条记录的enabled设置。

请求格式：

* 角色表示：多个ID
* 更新字段：结构化好的1个或多个字段

## 定义路由

handlers/role/router.go

```go
func Router(r *gin.Engine) {
	g := r.Group("role")

	// 更新多条的enabled
	g.PUT("/enabled", EditEnabled) // PUT /role?id=11&id=22
}
```

提供Handler

```go
func EditEnabled(ctx *gin.Context) {

}
```

## 解析请求参数

handlers/role/message.go

```go
// EditEnabledQueryReq 更新enabled的请求消息
// 将全部的Enabled字段，设置为相同的值
type EditEnabledQueryReq struct {
	IDList []uint `form:"id" binding:"gt=0"`
}

type EditEnabledBodyReq struct {
	Enabled bool `json:"enabled"`
}
```

handlers/role/handler.go

```go
func EditEnabled(ctx *gin.Context) {
	// 1. 解析Query请求数据
	query := EditEnabledQueryReq{}
	if err := ctx.ShouldBindQuery(&query); err != nil {
		// 记录日志
		utils.Logger().Error(err.Error())
		// 直接响应
		ctx.JSON(http.StatusOK, gin.H{
			"code":    100,
			"message": err.Error(),
		})
		return
	}
	//log.Println(query)

	// 2. 解析Body请求数据
	body := EditEnabledBodyReq{}
	if err := ctx.ShouldBind(&body); err != nil {
		// 记录日志
		utils.Logger().Error(err.Error())
		// 直接响应
		ctx.JSON(http.StatusOK, gin.H{
			"code":    100,
			"message": err.Error(),
		})
		return
	}
	//log.Println(body)
}
```

## 更新数据表

定义模型的更新方法：

models/role.go

```go
// RoleUpdateEnabled 更新多个字段的Enabled值
func RoleUpdateEnabled(idList []uint, enabled bool) (int64, error) {
	rowsNum := int64(0)
	err := utils.DB().Transaction(func(tx *gorm.DB) error {
		result := tx.Model(&Role{}).Where("`id` IN ?", idList).Update("enabled", enabled)
		if result.Error != nil {
			return result.Error
		} else {
			// 更新成功
			rowsNum = result.RowsAffected
		}
		return nil
	})

	return rowsNum, err
}
```

handlers/role/handler.go EditEnabled()

```
func EditEnabled(ctx *gin.Context) {
	// 1. 解析Query请求数据


	// 2. 解析Body请求数据


	// 3. 更新数据
	rowsNum, err := models.RoleUpdateEnabled(query.IDList, body.Enabled)
	if err != nil {
		// 记录日志
		utils.Logger().Error(err.Error())
		// 直接响应
		ctx.JSON(http.StatusOK, gin.H{
			"code":    100,
			"message": "数据更新错误",
		})
		return
	}

	// 4. 响应
	ctx.JSON(http.StatusOK, gin.H{
		"code": 0,
		"data": rowsNum,
	})
}
```

## 测试

![image.png](https://fynotefile.oss-cn-zhangjiakou.aliyuncs.com/fynote/fyfile/13080/1721639416053/94476994fae84bd09fb2266e614385bf.png)

欧克！

# CRUD操作小结

单表的CRUD结束了。

![image.png](https://fynotefile.oss-cn-zhangjiakou.aliyuncs.com/fynote/fyfile/13080/1721639416053/3d015ef3a7b84fb3a55a2fdbb6a8f0cb.png)

Restful路由：

```go
[GIN-debug] GET    /role                     --> ginCms/handlers/role.GetRow (3 handlers)
[GIN-debug] GET    /role/list                --> ginCms/handlers/role.GetList (3 handlers)
[GIN-debug] POST   /role                     --> ginCms/handlers/role.Add (3 handlers)
[GIN-debug] DELETE /role                     --> ginCms/handlers/role.Delete (3 handlers)
[GIN-debug] GET    /role/recycle             --> ginCms/handlers/role.Recycle (3 handlers)
[GIN-debug] PUT    /role/restore             --> ginCms/handlers/role.Restore (3 handlers)
[GIN-debug] PUT    /role/:id                 --> ginCms/handlers/role.Edit (3 handlers)
[GIN-debug] PUT    /role/enabled             --> ginCms/handlers/role.EditEnabled (3 handlers)

```

分层，MVC分层，在前后端分离架构下是没有View。：

* 路由接收请求（在MVC分层中，也成为前端控制器），分发请求到具体的处理器（也称为控制器的动作）
* handler处理器完成业务逻辑处理，接收和处理请求参数，数据（业务）的操作由处理器调用模型层来实现，做出响应
* 业务逻辑的具体实现由模型层完成

处理流程：

* 路由确定操作，操作的资源和动作
* 解析请求参数
  * Query，资源标识
  * Uri，资源标识
  * Body，主体数据
* 校验请求参数
* 请求数据的处理，组合，在handler层实现
* 核心数据的处理，在模型层处理
* 构造JSON响应

# CORS中间件

浏览器在请求接口时，存在同源策略安全机制。默认只允许请求同源的数据。

常规情况，前端和后端接口经常被独立部署，不在同源。此时就发生跨域请求。

后端接口因该允许跨域的特定资源请求，实现原理，通过**控制响应头**来告知浏览器，哪些资源可以跨域共享。

项目中，通过增加中间件的方案来实现：

安装cors中间件：

```
go get github.com/gin-contrib/cors
```

```
$ go get github.com/gin-contrib/cors
go: added github.com/gin-contrib/cors v1.7.2
```

定义中间件：

handlers/common/mwCors.go

```go
func UseCors(engine *gin.Engine) {
	// 1. 设置中间件
	cfg := cors.DefaultConfig()
	cfg.AllowAllOrigins = true
	cfg.AllowCredentials = true
	cfg.AddAllowHeaders("Authorization")
	// 2. 初始化，并使用中间件
	engine.Use(cors.New(cfg))
}
```

初始路由时，添加该中间件即可：

handlers/init.go

```go
// 初始化路由引擎
func InitEngine() *gin.Engine {
	// 1. 初始化路由引擎
	r := gin.Default()

	// 设置中间件
	common.UseCors(r)

	// 2. 注册不同模块的路由
	system.Router(r)
	role.Router(r)

	return r
}
```


```
[GIN-debug] GET    /ping                     --> ginCms/handlers/system.Ping (4 handlers)
[GIN-debug] GET    /role                     --> ginCms/handlers/role.GetRow (4 handlers)
[GIN-debug] GET    /role/list                --> ginCms/handlers/role.GetList (4 handlers)
```

此时有4个handlers:

1. Logger()
2. Recover()
3. Cors()
4. 业务处理器
