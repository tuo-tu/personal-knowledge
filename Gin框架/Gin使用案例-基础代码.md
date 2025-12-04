# 案例介绍

# 案例项目初始化

初始化的内容：

* 项目创建，module
* 代码管理，版本库，git
* 目录结构

## 项目创建

手动执行go mod init（也可以利用goland或其他编辑器完成工作)。

```
go mod init <module name>
```

```bash
$ mkdir ginCms
$ cd ginCms
$ go mod init ginCms
go: creating new go.mod: module ginCms

```

利用编辑器打开即可：

![image.png](https://fynotefile.oss-cn-zhangjiakou.aliyuncs.com/fynote/fyfile/13080/1721623860015/cdf52a80c9a9404f8ffba07992dac0fe.png)

## 代码版本管理

git：

* 初始化
* 做第一次提交
* 可选的远程版本库
* 配置.gitignore

初始化：

```
$ git init
Initialized empty Git repository in D:/apps/mashibing/ginCms/.git/

```

做第一次提交：

推荐创建一个README.md文件。

```bash
$ git add go.mod README.md
$ git commit -m 'first commit'
[master (root-commit) dd4e029] first commit
 2 files changed, 6 insertions(+)
 create mode 100644 README.md
 create mode 100644 go.mod

```

注意：课程中会手动提交，实操时按照自己的习惯来。

![image.png](https://fynotefile.oss-cn-zhangjiakou.aliyuncs.com/fynote/fyfile/13080/1721623860015/2815ec249593459583e9cb6c8835453e.png)

远程版本库：

在git.mashibing.com上创建版本库，https://git.mashibing.com/msb_59143/ginCms.git

在远程绑定:

```
$ git remote add origin https://git.mashibing.com/msb_59143/ginCms.git
$ git push -u origin master
Enumerating objects: 4, done.
Counting objects: 100% (4/4), done.
Delta compression using up to 8 threads
Compressing objects: 100% (3/3), done.
Writing objects: 100% (4/4), 348 bytes | 348.00 KiB/s, done.
Total 4 (delta 0), reused 0 (delta 0), pack-reused 0
remote: . Processing 1 references
remote: Processed 1 references in total
To https://git.mashibing.com/msb_59143/ginCms.git
 * [new branch]      master -> master
branch 'master' set up to track 'origin/master'.

```

配置.gitignore：

```
.DS_Store

# local env files
.env.local
.env.*.local

# Log files
npm-debug.log*
yarn-debug.log*
yarn-error.log*

# Editor directories and files
.idea
.vscode
*.suo
*.ntvs*
*.njsproj
*.sln
*.sw?

/server/log/
/server/gva
/server/latest_log

*.iml
web/.pnpm-debug.log
web/pnpm-lock.yaml

/web/node_modules
/web/dist

volumes/


```

再次提交。

# 服务启动与Handler管理

## 服务启动

拷贝基本代码到main.go

main.go

```go
package main

import "github.com/gin-gonic/gin"

func main() {
	r := gin.Default()
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})
	r.Run() // 监听并在 0.0.0.0:8080 上启动服务
}
```

运行：

```
go mod tidy
go run .
```

```
$ go mod tidy
go: finding module for package github.com/gin-gonic/gin
go: downloading github.com/gin-gonic/gin v1.10.0
go: found github.com/gin-gonic/gin in github.com/gin-gonic/gin v1.10.0
go: downloading github.com/pelletier/go-toml/v2 v2.2.2
go: downloading golang.org/x/crypto v0.23.0
go: downloading golang.org/x/arch v0.8.0

```

```
$ go run .
[GIN-debug] [WARNING] Creating an Engine instance with the Logger and Recovery middleware already attached.

[GIN-debug] [WARNING] Running in "debug" mode. Switch to "release" mode in production.
 - using env:   export GIN_MODE=release
 - using code:  gin.SetMode(gin.ReleaseMode)

[GIN-debug] GET    /ping                     --> main.main.func1 (3 handlers)
[GIN-debug] [WARNING] You trusted all proxies, this is NOT safe. We recommend you to set a value.
Please check https://pkg.go.dev/github.com/gin-gonic/gin#readme-don-t-trust-all-proxies for details.
[GIN-debug] Environment variable PORT is undefined. Using port :8080 by default
[GIN-debug] Listening and serving HTTP on :8080

```

访问：

```
curl http://localhost:8080/ping
{"message":"pong"}
```

## Handler管理

将Handler代码放在特定的目录中。

放在handlers目录中，同时基于业务逻辑进行分类：

```
/handlers/system
/handlers/user
/handlers/content
```

创建文件：

/handlers/system/hanler.go

```
package system

import "github.com/gin-gonic/gin"

func Ping(ctx *gin.Context) {
	ctx.JSON(200, gin.H{
		"message": "pong",
	})
}

```

main.go

```
func main() {
	r := gin.Default()
	r.GET("/ping", system.Ping)
	r.Run() // 监听并在 0.0.0.0:8080 上启动服务
}
```

# 集中管理路由

路由管理的思路：

* 修改路由不需要更新main.go文件
* 初始化路由引擎的操作单独处理
* 根据业务，将路由分散到具体的业务处理中。考虑handlers/子目录。

## 初始化路由引擎的操作单独处理

增加文件handlers/init.go用于初始化和handler相关的内容，包括路由：

handlers/init.go

```
// 初始化路由引擎
func InitEngine() *gin.Engine {
	r := gin.Default()
	r.GET("/ping", system.Ping)
	return r
}
```

main.go main()中，仅需要调用InitEngine()即可：

```
func main() {
	// 初始化路由引擎
	r := handlers.InitEngine()
	r.Run() // 监听并在 0.0.0.0:8080 上启动服务
}

```

## 根据业务将路由分散到具体的业务处理中

在handlers/system/router.go中，完成system相关的路由的初始化：

```
func Router(r *gin.Engine) {
	r.GET("/ping", Ping)
}
```

在handler/init.go的InitEngine()中调用

```
// 初始化路由引擎
func InitEngine() *gin.Engine {
	// 1. 初始化路由引擎
	r := gin.Default()

	// 2. 注册不同模块的路由
	system.Router(r)

	return r
}
```

# 管理配置

配置的核心操作：

* 存储配置
* 解析配置
* 使用配置

## 存储配置

通常将配置存储在特定格式的配置文件中，例如：

* json
* **yaml**
* xml
* ini

以yaml格式为例，存储配置。

创建配置文件：

configs.yaml

```
app:
  addr: ":8080"
```

configs.yaml 配置文件不是程序源码的一部分。时独立的文件。最终会存在一个程序的执行性程序和配置文件：

* server.exe  server
* configs.yaml

实操的时候，配置内容还可以来自其他服务。

## 解析配置

核心操作：

* 读取文件内容
* 解析内容到特定格式

往往需要附加操作：

* 不同配置格式的支持
* 不同位置的支持
* 不同的存储方案的支持
* 默认值

推荐使用 viper 包实现：https://github.com/spf13/viper。

特性：

* 默认配置
* 从 JSON, TOML, YAML, HCL 和 Java 属性配置文件读取数据
* 实时查看和重新读取配置文件（可选）
* 从环境变量中读取
* 从远程配置系统(etcd 或 Consul)读取数据并监听变化
* 从命令行参数读取
* 从 buffer 中读取
* 设置显式值

### 安装viper

```
go get github.com/spf13/viper
```

```
$ go get github.com/spf13/viper
go: downloading github.com/spf13/viper v1.19.0
go: added github.com/fsnotify/fsnotify v1.7.0
go: added github.com/hashicorp/hcl v1.0.0
go: added github.com/magiconair/properties v1.8.7
go: added github.com/mitchellh/mapstructure v1.5.0
go: added github.com/sagikazarmark/locafero v0.4.0
go: added github.com/sagikazarmark/slog-shim v0.1.0
go: added github.com/sourcegraph/conc v0.3.0
go: added github.com/spf13/afero v1.11.0
go: added github.com/spf13/cast v1.6.0
go: added github.com/spf13/pflag v1.0.5
go: added github.com/spf13/viper v1.19.0
go: added github.com/subosito/gotenv v1.6.0
go: added go.uber.org/atomic v1.9.0
go: added go.uber.org/multierr v1.9.0
go: added golang.org/x/exp v0.0.0-20230905200255-921286631fa9
go: added gopkg.in/ini.v1 v1.67.0

```

### 编写Viper初始配置函数

系统的核心工具，集中管理。utils/目录就是系统工具的操作代码目录

utils/config.go

```
// 默认配置
func defaultConfig() {
	viper.SetDefault("app.addr", ":8080")
}

// ParseConfig 解析配置
func ParseConfig() {
	// 1. 默认配置
	defaultConfig()

	// 2. 配置解析参数
	viper.AddConfigPath(".")       // 从哪些目录搜索配置文件
	viper.SetConfigName("configs") // 配置文件名字
	viper.SetConfigType("yaml")    // 配置类型（格式）

	// 3. 执行解析
	if err := viper.ReadInConfig(); err != nil {
		log.Fatal(err)
	}
}

```

在项目初始化时调用配置解析，在比较靠前的位置：

main.go

```
func main() {
	// 解析配置
	utils.ParseConfig()

	// 初始化路由引擎
	r := handlers.InitEngine()
	r.Run() // 监听并在 0.0.0.0:8080 上启动服务
}
```

## 使用配置

以配置监听端口为例：

main.go

```
func main() {
	// 解析配置
	utils.ParseConfig()

	// 初始化路由引擎
	r := handlers.InitEngine()
	r.Run(viper.GetString("app.addr")) // 监听并在 0.0.0.0:8080 上启动服务
}
```

通过修改配置文件，进行测试：

configs.yaml

```
app:
  addr: ":8084"
```

启动服务：

```
$ go run .
[GIN-debug] [WARNING] Creating an Engine instance with the Logger and Recovery middleware already attached.

[GIN-debug] [WARNING] Running in "debug" mode. Switch to "release" mode in production.
 - using env:   export GIN_MODE=release
 - using code:  gin.SetMode(gin.ReleaseMode)

[GIN-debug] GET    /ping                     --> ginCms/handlers/system.Ping (3 handlers)
[GIN-debug] [WARNING] You trusted all proxies, this is NOT safe. We recommend you to set a value.
Please check https://pkg.go.dev/github.com/gin-gonic/gin#readme-don-t-trust-all-proxies for details.
[GIN-debug] Listening and serving HTTP on :8084

```

# 应用模式配置

模式：

* 发布
* 调试模式
* 开发模式

gin支持模式。

基于配置，来修改gin的模式：

## 增加配置项

configs.yaml

```
app:
  mode: "debug" # release, debug, test
  addr: ":8084"
```

utils/config.go

```
func defaultConfig() {
	viper.SetDefault("app.mode", "debug")
	viper.SetDefault("app.addr", ":8080")
}
```

## 设置模式

工具级别：

utils/mode.go

```
// SetMode 设置应用模式
func SetMode() {
	switch strings.ToLower(viper.GetString("app.mode")) {
	case "release":
		gin.SetMode(gin.ReleaseMode)
	case "test":
		gin.SetMode(gin.TestMode)
	case "debug":
		fallthrough
	default:
		gin.SetMode(gin.DebugMode)
	}

}
```

项目运行初始化时调用：

main.go

```
func main() {
	// 解析配置
	utils.ParseConfig()
	// 设置应用模式
	utils.SetMode()

	// 初始化路由引擎
	r := handlers.InitEngine()
	r.Run(viper.GetString("app.addr")) // 监听并在 0.0.0.0:8080 上启动服务
}
```

# 管理日志

将日志集中管理。

推荐之一：Logrus，过去很流行的一个结构化日志包。

**本课推荐：slog, log/slog，标准库中提供的结构化日志包。https://pkg.go.dev/log/slog**

工作：

* 设置集中的日志Writer
* 配置日志信息，例如格式等

## 设置集中的日志Writer

增加日志工具，用于初始化logger：

utils/log.go

```
package utils

import (
	"github.com/gin-gonic/gin"
	"io"
	"log"
	"log/slog"
	"os"
)

// SetLogger 设置日志
func SetLogger() {
	// 1. 设置集中的日志Writer
	setLoggerWriter()

	// 2. 初始化
	initLogger()
}

// logger
var logger *slog.Logger

func Logger() *slog.Logger {
	return logger
}

// 公共的writer变量
var logWriter io.Writer

func LogWriter() io.Writer {
	return logWriter
}

// 设置writer
func setLoggerWriter() {
	// 根据不同的mode，选择不同的writer
	switch gin.Mode() {
	case gin.ReleaseMode:
		// 打开文件
		logfile := "./logs/app.log"
		if file, err := os.OpenFile(logfile, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0666); err != nil {
			log.Println(err)
			return
		} else {
			logWriter = file
		}
	case gin.TestMode, gin.DebugMode:
		fallthrough
	default:
		logWriter = os.Stdout
	}
}

// 初始化日志
func initLogger() {
	// 使用json模式记录
	logger = slog.New(slog.NewJSONHandler(logWriter, &slog.HandlerOptions{}))
}

```

## 应用初始化时设置日志

main.go

```
func main() {
	// 解析配置
	utils.ParseConfig()
	// 设置应用模式
	utils.SetMode()
	// 设置日志
	utils.SetLogger()

	// 初始化路由引擎
	r := handlers.InitEngine()
	r.Run(viper.GetString("app.addr")) // 监听并在 0.0.0.0:8080 上启动服务
}
```

## 使用应用日志

main.go 中演示：

```go
func main() {
	// 解析配置
	utils.ParseConfig()
	// 设置应用模式
	utils.SetMode()
	// 设置日志
	utils.SetLogger()

	// 初始化路由引擎
	r := handlers.InitEngine()
	// 使用logger输出应用日志
	utils.Logger().Info("service is listening", "addr", viper.GetString("app.addr"))
	r.Run(viper.GetString("app.addr")) // 监听并在 0.0.0.0:8080 上启动服务
}
```

```
$ go run .
[GIN-debug] GET    /ping                     --> ginCms/handlers/system.Ping (3 handlers)
{"time":"2024-07-22T15:32:46.1847233+08:00","level":"INFO","msg":"service is listening","addr":":8084"}


```

上面为debug模式。

测试记录到文件，使用release模式：

记录的结果：

logs/app.log

```
{"time":"2024-07-22T15:35:58.9737951+08:00","level":"INFO","msg":"service is listening","addr":":8084"}

```

## 分割日志文件

分割的逻辑：

* 日期分割
* 尺寸分割
* 其他属性分割

选择基于mouth月创建日志文件：

```
app-202403.log
```

### 基于时间创建日志文件名

utils/log.go setLoggerWriter函数：

```go
// 打开文件
		month := time.Now().Format("200601")
		logfile := fmt.Sprintf("./logs/app-%s.log", month)
```

## 配置日志存储目录

增加配置项：

configs.yaml

```
app:
  mode: "release" # release, debug, test
  addr: ":8084"
  log:
    path: "./logs"
```

设置默认值：

utils/config.go

```
// 默认配置
func defaultConfig() {
	viper.SetDefault("app.mode", "debug")
	viper.SetDefault("app.addr", ":8080")
	viper.SetDefault("app.log.path", "./logs")
}
```

更新创建日志文件的代码。

utils/logs.go setLoggerWriter

```go
// 设置writer
func setLoggerWriter() {
	// 根据不同的mode，选择不同的writer
	switch gin.Mode() {
	case gin.ReleaseMode:
		// 打开文件
		month := time.Now().Format("200601")
		logfile := viper.GetString("app.log.path")
		logfile += fmt.Sprintf("/app-%s.log", month)

		if file, err := os.OpenFile(logfile, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0666); err != nil {
			log.Println(err)
			return
		} else {
			logWriter = file
		}
	case gin.TestMode, gin.DebugMode:
		fallthrough
	default:
		logWriter = os.Stdout
	}
}
```

测试通过！

代码提交！

# 初始化数据库连接

工作内容：

* 启动数据库服务，与代码无关。只需得到dsn即可
* 连接数据库服务

## 启动数据库服务

使用docker-composer的方案管理数据库服务。

```
version: "3"

services:
  mysql:
    image: mysql:8
    command: mysqld --character-set-server=utf8mb4
    restart: always
    environment:
      MYSQL_ROOT_PASSWORD: secret
      MYSQL_DATABASE: gincms
    ports:
      - "3304:3306"
    volumes:
      - ./volumes/mysql/data:/var/lib/mysql
```

```
$ docker-compose up -d
 Network gincms_default  Creating
 Network gincms_default  Created
 Container gincms-mysql-1  Creating
 Container gincms-mysql-1  Created
 Container gincms-mysql-1  Starting
 Container gincms-mysql-1  Started

```

```
$ docker ps
CONTAINER ID   IMAGE     COMMAND                  CREATED          STATUS          PORTS                               NAMES
0caa6ef658b8   mysql:8   "docker-entrypoint.s…"   43 seconds ago   Up 41 seconds   33060/tcp, 0.0.0.0:3304->3306/tcp   gincms-mysql-1

```

启动成功！

## 连接数据库服务

### 使用配置管理DSN

configs.yaml

```
app:
  mode: "release" # release, debug, test
  addr: ":8084"
  log:
    path: "./logs"
db:
  dsn: root:secret@tcp(localhost:3304)/gincms?charset=utf8mb4&parseTime=True&loc=Local
```

增加默认：

utils/config.go

```go

// 默认配置
func defaultConfig() {
	viper.SetDefault("app.mode", "debug")
	viper.SetDefault("app.addr", ":8080")
	viper.SetDefault("app.log.path", "./logs")
	viper.SetDefault("db.dsn", "user:password@tcp(host:port)/dbname?charset=utf8mb4&parseTime=True&loc=Local")
}
```

### 使用GORM操作数据库

安装

```
go get gorm.io/gorm
go get gorm.io/driver/mysql
```

```
$ go get gorm.io/gorm
go get gorm.io/driver/mysq
go: downloading gorm.io/gorm v1.25.11
go: added github.com/jinzhu/inflection v1.0.0
go: added github.com/jinzhu/now v1.1.5
go: added gorm.io/gorm v1.25.11

```

```
$ go get gorm.io/driver/mysql
go: downloading gorm.io/driver/mysql v1.5.7
go: added github.com/go-sql-driver/mysql v1.7.0
go: added gorm.io/driver/mysql v1.5.7

```

### 增加连接数据库方法

uitls/db.go

```go
import (
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	gormLogger "gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
	"log"
	"time"
)

// InitDB 初始化数据库连接
func InitDB() {
	// 1. 配置
	logLevel := gormLogger.Warn
	// 根据应用的mod，控制级别
	switch gin.Mode() {
	case gin.ReleaseMode:
		logLevel = gormLogger.Warn
	case gin.TestMode, gin.DebugMode:
		fallthrough
	default:
		// 最多的日志
		logLevel = gormLogger.Info
	}
	// db 日志
	unionLogger := gormLogger.New(
		log.New(LogWriter(), "\n", log.LstdFlags),
		gormLogger.Config{
			SlowThreshold:             time.Second,
			Colorful:                  false,
			IgnoreRecordNotFoundError: false,
			ParameterizedQueries:      false,
			LogLevel:                  logLevel,
		},
	)
	// gorm 连接配置
	conf := &gorm.Config{
		SkipDefaultTransaction: false,
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true, // 单数表名
		}, // 数据表命名策略
		FullSaveAssociations:                     false,
		Logger:                                   unionLogger,
		NowFunc:                                  nil,
		DryRun:                                   false,
		PrepareStmt:                              false,
		DisableAutomaticPing:                     false,
		DisableForeignKeyConstraintWhenMigrating: true, // 数据表迁移时禁用外键约束
		IgnoreRelationshipsWhenMigrating:         false,
		DisableNestedTransaction:                 false,
		AllowGlobalUpdate:                        false,
		QueryFields:                              false,
		CreateBatchSize:                          0,
		TranslateError:                           false,
		PropagateUnscoped:                        false,
		ClauseBuilders:                           nil,
		ConnPool:                                 nil,
		Dialector:                                nil,
		Plugins:                                  nil,
	}

	// 2. 创建db对象
	dsn := viper.GetString("db.dsn")
	if dbNew, err := gorm.Open(mysql.Open(dsn), conf); err != nil {
		log.Fatalln(err)
	} else {
		db = dbNew
	}
}

// 全局的db对象
var db *gorm.DB

// DB 全局访问db对象的方法
func DB() *gorm.DB {
	return db
}

```

### 应用初始时完成调用

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

	// 初始化路由引擎
	r := handlers.InitEngine()
	// 使用logger输出应用日志
	utils.Logger().Info("service is listening", "addr", viper.GetString("app.addr"))
	r.Run(viper.GetString("app.addr")) // 监听并在 0.0.0.0:8080 上启动服务
}
```

应用程序中需要使用，则下面handler示例代码：

```
func Ping(ctx *gin.Context) {
	utils.DB().Save()
}
```
