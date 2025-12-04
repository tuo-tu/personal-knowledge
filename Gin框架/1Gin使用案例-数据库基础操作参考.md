# 基本使用

## 启动MySQL服务器等待使用

采用任意方式都可以，保证有可用的MySQL服务器即可。

本课程采用docker-compose方式管理MySQL服务器。

在案例目录创建文件：docker-compose.yml

键入如下内容：

```yml
services:
  db:
    container_name: gormExampleMySQL
    image: mysql
    environment:
      MYSQL_ROOT_PASSWORD: secret
      MYSQL_DATABASE: gormExample
    ports:
      - "3306:3306"
    volumes:
      - ./mysql/data:/var/lib/mysql
```

启动Docker.

docker compose

```shell
# 启动
docker-compose up -d
# 关闭
docker-compose down

> docker-compose up -d
[+] Running 2/2
 - Network gormexample_default  Created                                                          0.7s
 - Container gormExampleMySQL   Started
```

测试：

```shell
# 输入密码
docker exec -it gormExampleMySQL mysql -p

# 行内密码
docker exec -it gormExampleMySQL mysql -psecret
```

## 使用流程

1. 安装GORM和驱动
2. 连接数据库服务器
3. 定义模型
4. 使用迁移管理表
5. 完成操作
6. 调试、监控

## 抽象层和数据库驱动

![image.png](https://fynotefile.oss-cn-zhangjiakou.aliyuncs.com/fynote/fyfile/13080/1680502356013/db45129fd85e41feab26038078873c4d.png)

抽象层提供操作接口，具体的数据库操作实现由数据库对应的驱动实现。

```go
import (
  // 驱动
  "gorm.io/driver/mysql"
  // 抽象层
  "gorm.io/gorm"
)
```

gorm 为了方便使用，整合了以下驱动：

```go
gorm.io/driver/sqlite
gorm.io/driver/mysql
gorm.io/driver/postgres
gorm.io/driver/sqlserver
```

因此使用GROM时，需要安装gorm和对应的驱动。

## 安装GORM和MySQL驱动

我们创建一个示例module，初始化mod后，安装gorm和驱动：

```shell
go get -u gorm.io/gorm
go get -u gorm.io/driver/mysql
```

## 连接数据库服务器

示例代码：

```go
import (
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"log"
)

func BasicUsage() {
	dsn := "root:secret@tcp(127.0.0.1:3306)/gormExample?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("db:", db)
}
```

GORM使用连接池技术管理连接。

方法：

```go
func Open(dialector Dialector, opts ...Option) (db *DB, err error)
```

用于初始化数据库回话，基于拨号器 Dialector 和选项。返回 `*gorm.DB` 对象和错误。

Dialector，是通过驱动的Open方法创建的，以MySQL为例：

```go
func Open(dsn string) gorm.Dialector
```

需要提供DSN参数。

DSN，Data Source Name，数据源名称，用于描述在哪里找到数据。MySQL的DSN信息：

[MySQL DSN 说明](https://github.com/go-sql-driver/mysql#dsn-data-source-name)

## 定义模型

模型，就是一个struct类型，示例代码为：

```go
type Article struct {
	gorm.Model

	Subject     string
	Likes       uint
	Published   bool
	PublishTime time.Time
	AuthorID    uint
}
```

通常情况下，要嵌入 gorm.Model，用于保有核心字段。

gorm.Model

```go
type Model struct {
	ID        uint `gorm:"primarykey"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt DeletedAt `gorm:"index"`
}
```

## 迁移表结构

Migrate，迁移指的是通过模型来确定表结构。通常使用 AutoMigrate()完成迁移：

```go
// 迁移
db.AutoMigrate(&Article{})
```

执行以上代码后，对应的表结构就自动创建出来了。 AutoMigrate 会创建表、缺失的外键、约束、列和索引。 如果大小、精度、是否为空可以更改，则 AutoMigrate 会改变列的类型。 出于保护您数据的目的，它 **不会** 删除未使用的列。

使用 mysql 客户端，查看表结构的变化：

```
mysql> show create table articles;
```

## 基本CRUD操作

基于模型对象，完成CRUD，模型对象，也就是Article类型的数据示例。Article{}

### 初始化DB对象

```go
var DB *gorm.DB

func init() {
	// 定义DSN
	const dsn = "root:secret@tcp(127.0.0.1:3306)/gormExample?charset=utf8mb4&parseTime=True&loc=Local"

	// 连接服务器（池）
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal(err)
	}

	DB = db
}
```

### Create，创建

步骤：

- 创建模型对象
- db.Create() 完成insert操作

示例：

```go
func Create() {
	// 构建Article类型数据
	article := &Article{
		Subject:     "GORM 的 CRUD 基础操作",
		Likes:       0,
		Published:   true,
		PublishTime: time.Now(),
		AuthorID:    42,
	}

	// DB.Create 完成数据库的insert
	if err := DB.Create(article).Error; err != nil {
		log.Fatal(err)
	}

	// print
	fmt.Println(article)
}

```

测试执行

```shell
func TestCreate(t *testing.T) {
	Create()
}

go test -run Create  
&{{1 2023-04-03 19:22:13.46 +0800 CST 2023-04-03 19:22:13.46 +0800 CST {0001-01-01 00:00:00 +0000 UTC false}}
 GORM 的使用 0 true 2023-04-03 19:22:13.4586223 +0800 CST m=+0.015389201 42}
PASS
ok      gormExample     0.056s

```

在数据库中查看：

```mysql
SELECT * FROM articles\G
```

注意，gorm.Model 嵌入的字段：

- ID，auto_increment
- created_at和updated_at为当前时间
- deleted_at为null

### Retrieve，获取

步骤：

- 给定查询条件，例如PK
- 选择查询单个还是多个
  - Find() 多个
  - First() 单个

示例：

```go
func Retrieve(id uint) {
	// 初始化Article模型，零值
	article := &Article{}

	// DB.First()
	if err := DB.First(article, id).Error; err != nil {
		log.Fatal(err)
	}

	// print
	fmt.Println(article)
}
```

测试：

```shell
func TestRetrieve(t *testing.T) {
	Retrieve(1)
}

go test -run Retrieve
&{{1 2023-04-03 19:22:13.46 +0800 CST 2023-04-03 19:22:13.46 +0800 CST {0001-01-01 00:00:00 +0000 UTC false}}
 GORM 的使用 0 true 2023-04-03 19:22:13.459 +0800 CST 42}
PASS
ok      gormExample     0.052s
```

此处仅仅展示根据主键查询。

使用错误的id，测试查询出错的情况。

```shell
func TestRetrieve(t *testing.T) {
	Retrieve(3)
}

go test -run Retrieve

2023/04/03 19:34:46 D:/apps/courses/gormExample/basic.go:62 record not found
[3.809ms] [rows:0] SELECT * FROM `articles` WHERE `articles`.`id` = 3 AND `articles`.`deleted_at` IS NULL ORD
ER BY `articles`.`id` LIMIT 1
2023/04/03 19:34:46 record not found
exit status 1   
FAIL    gormExample     0.057s
```

### Update, 更新

步骤：

- 先确定更新的对象
- 设置对象属性字段
- 将对象存储

示例：

```go
func Update() {
	// 获取需要更新的对象
	article := &Article{}
	if err := DB.First(article, 1).Error; err != nil {
		log.Fatal(err)
	}

	// 更新对象字段
	article.AuthorID = 23
	article.Likes = 101
	article.Subject = "新的文章标题"

	// 存储，DB.Save()
	if err := DB.Save(article).Error; err != nil {
		log.Fatal(err)
	}

	// print
	fmt.Println(article)
}
```

测试：

```shell
func TestUpdate(t *testing.T) {
	Update()
}

go test -run Update
&{{1 2023-04-03 19:22:13.46 +0800 CST 2023-04-03 19:38:02.379 +0800 CST {0001-01-01 00:00:00 +0000 UTC false}
} 新的文字标题 100 true 2023-04-03 19:22:13.459 +0800 CST 42}
PASS
ok      gormExample     0.066s
```

本例子中，体现的典型的ORM更新。语法上更新对象。gorm同样支持基于条件的更新。

### Delete, 删除

步骤：

- 确定删除的模型对象
- 删除

示例：

```go
func Delete() {
	// 获取模型对象
	article := &Article{}
	if err := DB.First(article, 1).Error; err != nil {
		log.Fatal(err)
	}

	// DB.Delete() 删除
	if err := DB.Delete(article).Error; err != nil {
		log.Fatal(err)
	}

	// print
	fmt.Println("article was deleted")
}


// 当然也可以
DB.Delete(&Article{}, 1)
```

测试：

```shell
func TestDelete(t *testing.T) {
	Delete()
}

go test -run Delete  
article was deleted
PASS
ok      gormExample     0.069s
```

也可以通过主键ID删除，但本例子主要体现ORM的概念。同时在实际业务逻辑中，删除前，往往要对数据做额外的处理，通常也会先查询到的。

查看数据表，会发现记录的deleted_at字段设置了时间，表示该记录被删除。

```shell
mysql> select * from articles\G
*************************** 1. row ***************************
          id: 1
  created_at: 2023-04-03 19:22:13.460
  updated_at: 2023-04-03 19:38:02.379
  deleted_at: 2023-04-03 19:42:39.757
     subject: 新的文字标题
       likes: 100
   published: 1
publish_time: 2023-04-03 19:22:13.459
   author_id: 42
1 row in set (0.00 sec)
```

## Debug和日志

### db.Debug

`db.Debug()`方法用于将当前操作的log级别调整为 info 级别，就是可以获取当前执行的SQL：

```go
func Debug() {
	article := &Article{
		Subject:     "Article Subject",
		PublishTime: time.Now(),
	}
	if err := DB.Debug().Create(article).Error; err != nil {
		log.Fatal(err)
	}

	if err := DB.Debug().First(article, 1).Error; err != nil {
		log.Fatal(err)
	}
}
```

测试：

```powershell
func TestDebug(t *testing.T) {
	Debug()
}

> go test -run Debug

2023/04/04 13:16:59 D:/apps/courses/gormExample/basic.go:100
[8.999ms] [rows:1] INSERT INTO `articles` (`created_at`,`updated_at`,`deleted_at`,`subject`,`likes`,`publishe
d`,`publish_time`,`author_id`) VALUES ('2023-04-04 13:16:59.102','2023-04-04 13:16:59.102',NULL,'Article Subj
ect',0,false,'2023-04-04 13:16:59.1',0)

2023/04/04 13:16:59 D:/apps/courses/gormExample/basic.go:104
[3.000ms] [rows:1] SELECT * FROM `articles` WHERE `articles`.`id` = 4 AND `articles`.`deleted_at` IS NULL AND
 `articles`.`id` = 4 ORDER BY `articles`.`id` LIMIT 1
PASS
ok      gormExample     0.055s

```

### 全局配置日志级别

可以在 `gorm.Open` 时设置全局日志级别，或者通过 `db.Config.Logger`完成配置：

```go
// gorm.Open
db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
    Logger: logger.Default.LogMode(logger.Info),
})

// DB.Config.Logger
DB.Config.Logger = logger.Default.LogMode(logger.Info)
```

配置完成后，后续的操作都会使用Info级别的Log。

示例，更新 init() 方法：

```go
func init() {
	// 定义DSN
	const dsn = "root:secret@tcp(127.0.0.1:3306)/gormExample?charset=utf8mb4&parseTime=True&loc=Local"

	// 连接服务器（池）
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		// 设置Info级别的默认日志
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		log.Fatal(err)
	}

	DB = db
}
```

之前之前的CRUD操作，都会存在SQL输出了。

### 日志级别

GORM定义了4个级别：

- Info，logger.Info，全部消息
- Warn，logger.Warn，默认，警告
- Error, logger.Error，错误
- Silent, logger.Silent，静默

### 配置日志选项

GORM有默认的Logger实现，如下：

```go
Default = New(log.New(os.Stdout, "\r\n", log.LstdFlags), Config{
		SlowThreshold:             200 * time.Millisecond,
		LogLevel:                  Warn,
		IgnoreRecordNotFoundError: false,
		Colorful:                  true,
	})
```

**我们通过自定义Logger，即可控制选项：**

```go
// 自定义日志
var logWriter io.Writer

func init() {
	// 定义DSN
	const dsn = "root:secret@tcp(127.0.0.1:3306)/gormExample?charset=utf8mb4&parseTime=True&loc=Local"
	// 初始化logWriter
	logWriter, _ = os.OpenFile("./sql.log", os.O_CREATE|os.O_APPEND, 0644)
	customLogger := logger.New(log.New(logWriter, "\n", log.LstdFlags),
		logger.Config{
			// 慢查询阈值 200ms
			SlowThreshold: 200 * time.Millisecond,
			// 日志级别
			LogLevel: logger.Info,
			// 是否忽略记录不存在的错误
			IgnoreRecordNotFoundError: false,
			// 不彩色化
			Colorful: false,
		})
	// 连接服务器（池）
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		// 设置为自定义的日志
		Logger: customLogger,
	})
	if err != nil {
		log.Fatal(err)
	}

	DB = db
}
```

后续执行操作时，日志被记录在./sql.log文件中:

2023/04/04 15:11:24 D:/apps/mashibing/gormExample/basic.go:105
[12.441ms] [rows:1] INSERT INTO `articles` (`created_at`,`updated_at`,`deleted_at`,`subject`,`likes`,`published`,`publish_time`,`author_id`) VALUES ('2023-04-04 15:11:24.438','2023-04-04 15:11:24.438',NULL,'GORM 的 CRUD 基础操作',0,true,'2023-04-04 15:11:24.435',42)

2023/04/04 15:12:05 D:/apps/mashibing/gormExample/basic.go:105
[11.999ms] [rows:1] INSERT INTO `articles` (`created_at`,`updated_at`,`deleted_at`,`subject`,`likes`,`published`,`publish_time`,`author_id`) VALUES ('2023-04-04 15:12:05.032','2023-04-04 15:12:05.032',NULL,'GORM 的 CRUD 基础操作',0,true,'2023-04-04 15:12:05.029',42)

2023/04/04 15:12:27 D:/apps/mashibing/gormExample/basic.go:105
[17.603ms] [rows:1] INSERT INTO `articles` (`created_at`,`updated_at`,`deleted_at`,`subject`,`likes`,`published`,`publish_time`,`author_id`) VALUES ('2023-04-04 15:12:27.809','2023-04-04 15:12:27.809',NULL,'GORM 的 CRUD 基础操作',0,true,'2023-04-04 15:12:27.806',42)


# 模型定义

模型，model，在ORM中与关系表映射的应用程序对象，不同语言使用不同数据类型实现，通常为对象类型，在Go语言中是struct结构体类型。具体的struct类型实例就是具体的某个表的模型对象，ORM中与表中记录映射。

```go
type Article struct {
    gorm.Model
	Subject     string
}
```

## 表名定义

### 表名约定

在默认情况下，GORM有约定，使用小写+下划线（蛇形命名）的复数形式作为表名，例如：

| Model        | Table           |
| ------------ | --------------- |
| Post         | posts           |
| Category     | categories      |
| Box          | boxes           |
| PostCategory | post_categories |

示例：

```go
type Post struct {
	gorm.Model
}
type Category struct {
	gorm.Model
}
type PostCategory struct {
	gorm.Model
}
type Box struct {
	gorm.Model
}

func Migrate() {
	if err := DB.Debug().AutoMigrate(&Post{}, &Category{}, &PostCategory{}, &Box{}); err != nil {
		log.Fatal(err)
	}
}

// =============================
// *_test.go
func TestMigrate(t *testing.T) {
	Migrate()
}
> go test -run Migrate

```

查看数据表：

```mysql
mysql> show tables;
+-----------------------+
| Tables_in_gormexample |
+-----------------------+
| boxes                 |
| categories            |
| post_categories       |
| posts                 |
+-----------------------+
4 rows in set (0.01 sec)
```

### 命名策略自定义

GORM的约定，就是使用内置的命名策略来实现的。命名策略是实现了Namer接口的类型：

```go
// gorm.io/gorm/schema

// Namer namer interface
type Namer interface {
	TableName(table string) string
	SchemaName(table string) string
	ColumnName(table, column string) string
	JoinTableName(joinTable string) string
	RelationshipFKName(Relationship) string
	CheckerName(table, column string) string
	IndexName(table, column string) string
}
```

默认的策略（约定）的TableName实现如下：

```go
// gorm.io/gorm/schema

// NamingStrategy tables, columns naming strategy
type NamingStrategy struct {
    // 表名前缀
	TablePrefix   string
    // 是否单数表名
	SingularTable bool
    // 替换器，用于替换特定字符串
	NameReplacer  Replacer
    // 是否为sname_casing形式
	NoLowerCase   bool
}

// 将模型名转为表名
func (ns NamingStrategy) TableName(str string) string {
	// 单数
    if ns.SingularTable {
		return ns.TablePrefix + ns.toDBName(str)
	}
    // 复数
	return ns.TablePrefix + inflection.Plural(ns.toDBName(str))
}

// ColumnName convert string to column name
func (ns NamingStrategy) ColumnName(table, column string) string {
	return ns.toDBName(column)
}
```

### 表名前缀

使用默认的命名策略，来增加表名前缀。

通过修改gorm.Open时的配置实现：

```go
db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
    NamingStrategy: schema.NamingStrategy{
        TablePrefix:   "msb_",
        SingularTable: true,
        NameReplacer:  nil,
        NoLowerCase:   true,
    },
})
```

测试生成的数据表：

```mysql
mysql> show tables;
+-----------------------+
| Tables_in_gormexample |
+-----------------------+
| boxes                 |
| categories            |
| post_categories       |
| posts                 |
| msb_Category          |
| msb_Post              |
| msb_PostCategory      |
| msb_Box               |
+-----------------------+
```

实际项目中，比较喜欢使用小写表名的。所以 NoLowerCase: false比较常见。

> 旧版v1的修改表名前缀的方案，是：
>
> ```go
> gorm.DefaultTableNameHanlder
> ```

### 表名自定义

若需要使用自定义规则表名，模型需要实现 Tabler 接口，Tabler接口：

```go
// gorm.io/gorm/schema
type Tabler interface {
	TableName() string
}
```

示例：

```go
type Box struct {
	gorm.Model
}

func (Box) TableName() string {
	return "my_box"
}
```

得到的表名就是 my_box。

### 临时指定表名

方法：

```go
func (db *DB) Table(name string, args ...interface{}) (tx *DB)DB.Table
```

用于在一个执行周期内，指定临时表名。若配合Migrate使用，可以设置所迁移的表的名字。例如：

```go
DB.Table("temp_articles").AutoMigrate(&Article{})
```

Table方法常用于SQL的执行中。

## 字段类型映射

模型的字段类型可为：

- Go基本数据类型，典型的：`bool, int, uint, float32/64, string, time.Time, []byte`
- Go基本数据类型的指针类型，典型的：`*int, *uint, *float32/64, *string, *time.Time`
- 实现了Scanner和Valuer接口的自定义类型，典型的database/sql包定义的sql.NullType系列类型

示例：

```go
type TypeMap struct {
	gorm.Model

	FInt       int
	FUInt      uint
	FFloat32   float32
	FFloat64   float64
	FString    string
	FTime      time.Time
	FByteSlice []byte

	FIntP     *int
	FUIntP    *uint
	FFloat32P *float32
	FFloat64P *float64
	FStringP  *string
	FTimeP    *time.Time
}
```

基于以上模型，迁移生成的表结构：

```mysql
mysql> desc type_maps;
+--------------+-----------------+------+-----+---------+----------------+
| Field        | Type            | Null | Key | Default | Extra          |
+--------------+-----------------+------+-----+---------+----------------+
| id           | bigint unsigned | NO   | PRI | NULL    | auto_increment |
| created_at   | datetime(3)     | YES  |     | NULL    |                |
| updated_at   | datetime(3)     | YES  |     | NULL    |                |
| deleted_at   | datetime(3)     | YES  | MUL | NULL    |                |
| f_int        | bigint          | YES  |     | NULL    |                |
| f_uint       | bigint unsigned | YES  |     | NULL    |                |
| f_float32    | float           | YES  |     | NULL    |                |
| f_float64    | double          | YES  |     | NULL    |                |
| f_string     | longtext        | YES  |     | NULL    |                |
| f_time       | datetime(3)     | YES  |     | NULL    |                |
| f_byte_slice | longblob        | YES  |     | NULL    |                |
| f_int_p      | bigint          | YES  |     | NULL    |                |
| f_uint_p     | bigint unsigned | YES  |     | NULL    |                |
| f_float32_p  | float           | YES  |     | NULL    |                |
| f_float64_p  | double          | YES  |     | NULL    |                |
| f_string_p   | longtext        | YES  |     | NULL    |                |
| f_time_p     | datetime(3)     | YES  |     | NULL    |                |
+--------------+-----------------+------+-----+---------+----------------+
```

以上就是MySQL数据对应的类型。注意不同数据库有不同的实现，但以上类型是大多数数据库支持的通用类型。

## 指针类型和非指针类型的区别

上面的案例中，*T和T对应的MySQL类型是一致的，这是因为在MySQL的数据字段类型层面，没有指针的概念。

但在Go的类型中，存在区别：

- T，不能表示NULL，如果数据库中字段的值为NULL，那么映射到模型对象的字段为类型零值。同理，无法设置记录字段的值为NULL。
- *T，可以使用nil表示NULL，意味着如果数据库中字段的值为NULL，那么映射到模型对象的字段为nil。同理，可以通过将字段设置为nil，来设置字段为NULL。

插入一条ID:1的记录作为测试：

```mysql
mysql> INSERT INTO `msb_type_map` (`id`) VALUES (1);
```

示例，使用TypeMap模型，完成查询：

```go
func PointerDiff() {
	tm := &TypeMap{}
	fmt.Printf("%+v\n", tm)

	DB.First(tm, 1)
	fmt.Printf("%+v\n", tm)
}
```

比较模型零值时，不同字段的查询。

比较表中记录字段为Null时，模型字段的差异。

```go
&{Model:{ID:0 CreatedAt:0001-01-01 00:00:00 +0000 UTC UpdatedAt:0001-01-01 00:00:00 +0000 UTC DeletedAt:{Time
:0001-01-01 00:00:00 +0000 UTC Valid:false}} FInt:0 FUInt:0 FFloat32:0 FFloat64:0 FString: FTime:0001-01-01 0
0:00:00 +0000 UTC FByteSlice:[] FIntP:<nil> FUIntP:<nil> FFloat32P:<nil> FFloat64P:<nil> FStringP:<nil> FTime
P:<nil>}
&{Model:{ID:1 CreatedAt:0001-01-01 00:00:00 +0000 UTC UpdatedAt:0001-01-01 00:00:00 +0000 UTC DeletedAt:{Time
:0001-01-01 00:00:00 +0000 UTC Valid:false}} FInt:0 FUInt:0 FFloat32:0 FFloat64:0 FString: FTime:0001-01-01 0
0:00:00 +0000 UTC FByteSlice:[] FIntP:<nil> FUIntP:<nil> FFloat32P:<nil> FFloat64P:<nil> FStringP:<nil> FTime
P:<nil>}

```

**结论：若表中字段可以为NULL，那么应该使用*T指针类型，映射字段，nil表示NULL。**

## 自定义字段类型

除标准Go类型外，还可以使用实现database/sql包中Scanner和database/sql/driver包中Valuer接口的自定义类型，以便让 GORM 知道如何将该类型接收、保存到数据库。其中：

Scanner接口：

```go
// database/sql/sql.go

// Scanner is an interface used by Scan.
type Scanner interface {
	// Scan 从数据库中分配一个值，用于查询时设置字段值
	Scan(src any) error
}
```

Valuer接口：

```go
// database/sql/driver/types.go

type Valuer interface {
	// Value 用于获取模型字段的值
	Value() (Value, error)
}
```

以 sql.NullTime 为例：

```go
// database/sql/sql.go

// NullTime represents a time.Time that may be null.
// NullTime implements the Scanner interface so
// it can be used as a scan destination, similar to NullString.
type NullTime struct {
	Time  time.Time
	Valid bool // Valid is true if Time is not NULL
}

// Scan implements the Scanner interface.
func (n *NullTime) Scan(value any) error {
	if value == nil {
		n.Time, n.Valid = time.Time{}, false
		return nil
	}
	n.Valid = true
	return convertAssign(&n.Time, value)
}

// Value implements the driver Valuer interface.
func (n NullTime) Value() (driver.Value, error) {
	if !n.Valid {
		return nil, nil
	}
	return n.Time, nil
}
```

**典型的sql.NullType**

```go
sql.NullTime
sql.NullByte
sql.NullBool
sql.NullFloat64
sql.NullInt16
sql.NullInt64
sql.NullString
sql.NullInt32
```

结构体字段 Valid:false 表示，数据表值是NULL。使用模型时，可以根据对应的Valid字段，判定数据表中数据是否为NULL。

**第三方自定义类型**

```go
github.com/google/uuid
```

install：

```shell
go get github.com/google/uuid
```

**示例**：

```go
type CustomTypeModel struct {
	gorm.Model

	FTime     time.Time
	FNullTime sql.NullTime

	FString     string
	FNullString sql.NullString

	FUUID     uuid.UUID
	FNullUUID uuid.NullUUID
}

func CustomType() {
	//id := uuid.UUID{}
	//id.Scan()  // Scanner
	//id.Value() // Valuer

	// 初始化模型
	ctm := &CustomTypeModel{}
	// 迁移数据表
	DB.AutoMigrate(ctm)

	// 创建
	ctm.FTime = time.Now()         // 当前时间
	ctm.FNullTime = sql.NullTime{} // 零值，Valid默认为false
	ctm.FString = ""
	ctm.FNullString = sql.NullString{}

	ctm.FUUID = uuid.New()
	ctm.FNullUUID = uuid.NullUUID{}

	DB.Create(ctm)

	// 查询
	DB.First(ctm, ctm.ID)

	// 判定字段是否为NULL
	if ctm.FString == "" {
		fmt.Println("FString is NULL")
	} else {
		fmt.Println("FString is NOT NULL")
	}

	if ctm.FNullString.Valid == false {
		fmt.Println("FNullString is NULL")
	} else {
		fmt.Println("FNullString is NOT NULL")
	}
}
```

测试结果：

```shell
#// 测试函数
func TestCustomType(t *testing.T) {
	CustomType()
}

# 运行测试
> go test -run UUIDCreate
&{Model:{ID:4 CreatedAt:2023-04-06 11:45:28.752 +0800 CST UpdatedAt:2023-04-06 11:45:28.752 +0800 CST Deleted
At:{Time:0001-01-01 00:00:00 +0000 UTC Valid:false}} FTime:2023-04-06 11:45:28.751 +0800 CST FString: FNullTi
me:{Time:0001-01-01 00:00:00 +0000 UTC Valid:false} FNullString:{String: Valid:false} FUUID:04c608f2-6de3-4f0
5-a86b-c699b3560a6f}
字段为NULL
PASS  
ok      gormExample     0.096s


# mysql 结果
mysql> select * from msb_custom_type_model\G
*************************** 1. row ***************************
           id: 1
   created_at: 2023-04-06 11:37:26.584
   updated_at: 2023-04-06 11:37:26.584
   deleted_at: NULL
       f_time: 2023-04-06 11:37:26.583
     f_string:
  f_null_time: NULL
f_null_string: NULL
       f_uuid: af5a0505-0f07-495e-ac03-250cd5ccc8bf
1 row in set (0.00 sec)
```

## 字段标签设置字段属性

结构体字段的标签 gorm，可用来对GROM的行为进行设置。

标签语法

```go
type ModelType struct {
	Field Type `gorm:"key:value;key;"`
}
```

常用的字段标签如下表：

| key                    | 功能                                                                                      | 说明                                                                                                       |
| ---------------------- | ----------------------------------------------------------------------------------------- | ---------------------------------------------------------------------------------------------------------- |
| column                 | 列名                                                                                      |                                                                                                            |
| type                   | 类型                                                                                      | 根据是需求设置字段类型，推荐兼容性好的类型。所有数据库都支持 bool、int、uint、float、string、time、bytes。 |
| not null               | 设置为NOT NULL，默认为NULL，不需要指定                                                    |                                                                                                            |
| default                | 默认值                                                                                    | 根据需求设置                                                                                               |
| autoCreateTime         | 创建时追踪当前时间，支持Time和int类型，int类型表示时间戳，通常使用gorm.Model中的定义      | autoCreateTime;autoUpdateTime:milli;autoCreateTime:nano;分别表示秒级，毫秒级，纳秒级                       |
| autoUpdateTime         | 创建/更新时追踪当前时间，支持Time和int类型，int类型表示时间戳，通常使用gorm.Model中的定义 | autoCreateTime;autoUpdateTime:milli;autoCreateTime:nano;分别表示秒级，毫秒级，纳秒级                       |
| autoIncrement          | autoUpdateTime自动增长                                                                    |                                                                                                            |
| autoIncrementIncrement | 自动增长的步长                                                                            |                                                                                                            |
| comment                | 注释                                                                                      |                                                                                                            |
|                        |                                                                                           |                                                                                                            |
| size                   | 列长度，通常在type中指定，例如 varchar(255)                                               |                                                                                                            |
| precision              | 精度，通常在type中指定，例如 decimal(10, 2) 中的10                                        |                                                                                                            |
| scale                  | 小数位数，通常在type中指定，例如 decimal(10, 2)中的2                                      |                                                                                                            |

示例：

```go
type FieldTag struct {
	gorm.Model

	FTypeChar    string `gorm:"type:char(127)"`
	FTypeVarchar string `gorm:"type:varchar(255)"`
	FTypeText    string `gorm:"type:text"`
	FTypeBlob    []byte `gorm:"type:blob"`
	FTypeEnum    string `gorm:"type:enum('Go', 'GORM', 'MySQL')"`
	FTypeSet     string `gorm:"type:set('Go', 'GORM', 'MySQL')"`

    FColName    string `gorm:"column:custom_field"`
	FTypeNotNull string `gorm:"type:varchar(255);not null"`
	FTypeDefault string `gorm:"type:varchar(255);not null;default:hello gorm!"`
	FTypeComment string `gorm:"type:varchar(255);not null;default:hello gorm!;comment:some comment"`
}
```

通过migrate测试创建的表的字段属性。

```mysql
mysql> set names utf8;
Query OK, 0 rows affected, 1 warning (0.01 sec)

mysql> show create table msb_field_tag\G
*************************** 1. row ***************************
       Table: msb_field_tag
Create Table: CREATE TABLE `msb_field_tag` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `created_at` datetime(3) DEFAULT NULL,
  `updated_at` datetime(3) DEFAULT NULL,
  `deleted_at` datetime(3) DEFAULT NULL,
  `f_string_default` longtext,
  `f_type_char` char(32) DEFAULT NULL,
  `f_type_varchar` varchar(255) DEFAULT NULL,
  `f_type_text` text,
  `f_type_blob` blob,
  `f_type_enum` enum('Go','GORM','MySQL') DEFAULT NULL,
  `f_type_set` set('Go','GORM','MySQL') DEFAULT NULL,
  `custom_column_name` longtext,
  `f_col_not_null` varchar(255) NOT NULL,
  `f_col_default` varchar(255) NOT NULL DEFAULT 'gorm middle ware',
  `f_col_comment` varchar(255) DEFAULT NULL COMMENT '带有注释的字段',
  PRIMARY KEY (`id`),
  KEY `idx_msb_field_tag_deleted_at` (`deleted_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
1 row in set (0.00 sec)
```

注意，类型与结构之间的合理性，下面的标签也可以设置成功，但操作时会有很大局限性：

```go
 FTypeUn1 string `gorm:"type:int"`
 FTypeUn2 int    `gorm:"type:tinyint"`
```

## 索引和约束

- 索引，Index，快速检索数据
  - Index
  - UniqueIndex
- 约束，Constraint，表中数据的限制条件，大部分约束都是基于索引实现的。
  - 索引约束, Index
  - 主键约束，PK
  - 外键约束，FK
  - 检查约束，check
  - NOT NULL
  - DEFAULT

MySQL中使用 Key 实现了约束和索引的功能：

索引和约束是通过字段标签来定义的，常用的标签如下：

| key         | 功能     | 说明 |
| ----------- | -------- | ---- |
| primaryKey  | 主键     |      |
| unique      | 唯一键   |      |
| index       | 索引     |      |
| uniqueIndex | 唯一索引 |      |
| check       | 检查约束 |      |

支持创建复合索引，通过名字识别。复合索引中，默认基于模型字段顺序确定索引字段优先级，支持使用priority选项定义。

支持索引选项，index:索引名字,key1:value1,key2:value2 的方式指定选项：

- sort:desc 降序关键字，默认升序
- length:N 前缀N作为关键字
- comment 索引注释
- type:btree 索引类型
- where:CONDITION 过滤条件

示例：

```go
type IAndC struct {
	ID          uint    `gorm:"primaryKey"`
	Email       string  `gorm:"type:varchar(255);unique"`
	Age         int8    `gorm:"index;check:age >= 18 AND email is not null"`
	FirstName   string  `gorm:"index:name"`
	LastName    string  `gorm:"index:name"`
	FirstName1  string  `gorm:"index:name1,priority:2"`
	LastName1   string  `gorm:"index:name1,priority:1"`
	Height      float32 `gorm:"index:,sort:desc"`
	AddressHash string  `gorm:"type:varchar(42);index:,length:12"`
	Telephone   string  `gorm:"type:varchar(16);uniqueIndex:,comment:电话号码唯一索引"`
}
```

利用Migrate测试表结构：

```mysql
mysql> show create table msb_i_and_c\G
*************************** 1. row ***************************
       Table: msb_i_and_c
Create Table: CREATE TABLE `msb_i_and_c` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `email` varchar(255) DEFAULT NULL,
  `age` tinyint DEFAULT NULL,
  `first_name` varchar(191) DEFAULT NULL,
  `last_name` varchar(191) DEFAULT NULL,
  `first_name1` varchar(191) DEFAULT NULL,
  `last_name1` varchar(191) DEFAULT NULL,
  `address_hash` varchar(42) DEFAULT NULL,
  `height` float DEFAULT NULL,
  `telephone` varchar(16) DEFAULT NULL,

  PRIMARY KEY (`id`),
  UNIQUE KEY `email` (`email`),
  UNIQUE KEY `telephone` (`telephone`),
  UNIQUE KEY `idx_msb_i_and_c_telephone` (`telephone`) COMMENT '电话号码唯一索引',
  KEY `idx_msb_i_and_c_address_hash` (`address_hash`(12)),
  KEY `idx_msb_i_and_c_age` (`age`),
  KEY `name` (`first_name`,`last_name`),
  KEY `name1` (`last_name1`,`first_name1`),
  KEY `idx_msb_i_and_c_height` (`height` DESC),
  CONSTRAINT `name_checker` CHECK ((`name` <> _utf8mb4'jinzhu'))
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
1 row in set (0.00 sec)
```

写数据时，就要满足以上的约束了。

测试不满足约束的插入：

```go
func IAndCCreate() {
	iac := &IAndC{}
	//iac.Age = 18
	if err := DB.Create(iac).Error; err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%+v", iac)
}
```

测试：

```shell
func TestIAndCCreate(t *testing.T) {
	IAndCCreate()
}

> go test -run IAndCCreate
2023/04/06 16:17:22 Error 3819 (HY000): Check constraint 'chk_msb_i_and_c_age' is violated.
exit status 1   
FAIL    gormExample     0.063s
```

## 字段操作控制

### 忽略字段

`-` 标签标示忽略该字段的数据库操作。

例如，网络传递的字段是 url，但数据表中存储的字段是：schema, host, path, query_string。就是将URL拆解成多个部分进行存储，但网络传输时使用的是url字段。那么该url字段，就应该忽略数据库的读写（包括迁移）操作，示例：

```go
type Service struct {
	Schema      string
	Host        string
	Path        string
	QueryString string
	Url         string `gorm:"-"`
}
```

以上模型在迁移生成表时，会忽略Url字段，同样在CRUD时，也会被忽略。

迁移测试：

```mysql
mysql> desc msb_service;
+--------------+----------+------+-----+---------+-------+
| Field        | Type     | Null | Key | Default | Extra |
+--------------+----------+------+-----+---------+-------+
| schema       | longtext | YES  |     | NULL    |       |
| host         | longtext | YES  |     | NULL    |       |
| path         | longtext | YES  |     | NULL    |       |
| query_string | longtext | YES  |     | NULL    |       |
+--------------+----------+------+-----+---------+-------+
4 rows in set (0.01 sec)
```

忽略字段可以做更细致的忽略限定：

```
-:migration // 忽略迁移，不忽略CRUD
```

测试，手动增加url字段：

```mysql
mysql> alter table msb_service add column url longtext;
```

测试读写操作：

```go
func ServiceCRUD() {
	s := &Service{}
	s.Schema = "https"
	s.Url = "https://www.mashibing.com/study?key=value"
	DB.Create(s)
	fmt.Printf("%+v\n", s)
}
```

基于tag的不同：

```
// - 不能修改
// -:migration 可以修改
```

### 权限控制

使用tag，可以控制某个字段的读写权限：

```
// 写权限控制
<-:false 无写入
<-:create 仅创建
<-:update 仅更新

// 读权限控制
->:false 无读取
```

示例：

```go
type Service struct {
	Schema      string
	Host        string
	Path        string `gorm:"->:false"`
	QueryString string `gorm:"<-:update"`
	Url string `gorm:"-"`
}
```

测试CRUD，可以看到对应的操作会被控制。

## 字段自动编解码

当我们需要处理数据库不能直接处理的数据时，通常要自定义处理过程。典型的自定义方案有：

- 使用序列化器
- 使用实现Scanner和Valuer接口的自定义类型

### 序列化器 Serializer

![image.png](https://fynotefile.oss-cn-zhangjiakou.aliyuncs.com/fynote/fyfile/13080/1680606521071/e89858df32a54e7f94e77a123066e58e.png)

优先推荐使用GORM提供的序列化器，完成字段的序列化与反序列化。

GORM 提供了一些默认的序列化器：json、gob、unixtime。

使用 serializer 标签进行设置。

示例，处理为JSON编码：

```go
type Paper struct {
	gorm.Model

	Subject string
	//Tags    []string
	// 使用 json 序列化器进行处理
	Tags []string `gorm:"serializer:json"`
}

func PaperCrud() {
	if err := DB.AutoMigrate(&Paper{}); err != nil {
		log.Fatal(err)
	}

	// 常规操作
	paper := &Paper{}
	paper.Subject = "使用Serializer操作Tags字段"
	paper.Tags = []string{"Go", "Serializer", "Gorm", "MySQL"}
	// create 会执行序列化工作,serialize
	if err := DB.Create(paper).Error; err != nil {
		log.Fatal(err)
	}

	// 查询
	newPaper := &Paper{}
	// First会执行反序列化工作，unserialize
	DB.First(newPaper, 5)
	fmt.Printf("%+v\n", newPaper)
}

```

测试，发现可以处理[]string类型的字段，对应的数据库内容是JSON编码内容：

```mysql
*************************** 1. row ***************************
        id: 1
created_at: 2023-04-07 16:09:12.993
updated_at: 2023-04-07 16:09:12.993
deleted_at: NULL
   subject: 使用Serializer操作Tags字段
categories: NULL
      tags: ["Go","Serializer","Gorm","MySQL"]
5 rows in set (0.00 sec)
```

而模型字段是[]类型：

```go
Tags:[Go Serializer Gorm MySQL]
```

### 自定义编解码器

编解码器需要实现下面两个接口：

```go
// import "gorm.io/gorm/schema"

type SerializerInterface interface {
    Scan(ctx context.Context, field *schema.Field, dst reflect.Value, dbValue interface{}) error
    SerializerValuerInterface
}

type SerializerValuerInterface interface {
    Value(ctx context.Context, field *schema.Field, dst reflect.Value, fieldValue interface{}) (interface{}, error)
}
```

以上接口的方法，与 Valuer和Scanner接口几乎一致：

```go
// database/sql/sql.go
// Scanner is an interface used by Scan.
type Scanner interface {
	// Scan 从数据库中分配一个值，用于查询时设置字段值
	Scan(src any) error
}

// database/sql/driver/types.go
type Valuer interface {
	// Value 用于获取模型字段的值
	Value() (Value, error)
}
```

可见，思路一致的。

自定义编码器步骤：

1. 定义实现编码器接口的类型
2. 注册编码器
3. 在模型tag中使用

示例，实现自定义编码器，同样处理[]string类型为CSV格式（Comma-Separated Values，逗号分隔值）：

```go
// 1.定义实现了序列化器接口的类型
type CSVSerializer struct{}

// 实现Scan，unserialize时执行
// ctx Context对象
// field 模型的字段对应的类型
// dst 目标值（最终结果赋值到dst）
// dbValue 从数据库读取的值
// 错误
func (CSVSerializer) Scan(ctx context.Context, field *schema.Field, dst reflect.Value, dbValue interface{}) error {
	// 初始化一个用来存储字段值的变量
	var fieldValue []string
	// 一:解析读取到的数据表的数据
	if dbValue != nil { // 不是 NULL
		// 支持解析的只有string和[]byte
		// 使用类型检测进行判定
		var str string
		switch v := dbValue.(type) {
		case string:
			str = v
		case []byte:
			str = string(v)
		default:
			return fmt.Errorf("failed to unmarshal CSV value: %#v", dbValue)
		}
		// 二：核心：将数据表中的字段使用逗号分割，形成 []string
		fieldValue = strings.Split(str, ",")
	}

	// 三，将处理好的数据，设置到dst上
	field.ReflectValueOf(ctx, dst).Set(reflect.ValueOf(fieldValue))
	return nil
}

// 实现Value, serialize时执行
// fieldValue 模型的的字段值
func (CSVSerializer) Value(ctx context.Context, field *schema.Field, dst reflect.Value, fieldValue interface{}) (interface{}, error) {
	// 将字段值转换为可存储的CSV结构
	return strings.Join(fieldValue.([]string), ","), nil
}
```

测试，csv编码器的使用：

```go
type Paper struct {
	gorm.Model

	Subject string
	//Tags    []string
	// 使用 json 序列化器进行处理
	Tags []string `gorm:"serializer:json"`
	// 使用自定义的编码器
	Categories []string `gorm:"serializer:csv"`
}

// 2.注册到GORM中
// 3.测试
func CustomSerializer() {
	// 注册序列化器
	schema.RegisterSerializer("csv", CSVSerializer{})

	// 测试
	if err := DB.AutoMigrate(&Paper{}); err != nil {
		log.Fatal(err)
	}

	// 常规操作
	paper := &Paper{}
	paper.Subject = "使用自定义的Serializer操作Categories字段"
	paper.Tags = []string{"Go", "Serializer", "Gorm", "MySQL"}
	paper.Categories = []string{"Go", "Serializer", "Gorm", "MySQL"}
	// create 会执行序列化工作,serialize
	if err := DB.Create(paper).Error; err != nil {
		log.Fatal(err)
	}

	// 查询
	newPaper := &Paper{}
	// First会执行反序列化工作，unserialize
	DB.First(newPaper, paper.ID)
	fmt.Printf("%+v\n", newPaper)
}
```

unit test:

```shell
func TestCustomSerializer(t *testing.T) {
	CustomSerializer()
}

> go test -run CustomSerializer
&{Model:{ID:16 CreatedAt:2023-04-07 19:01:53.513 +0800 CST UpdatedAt:2023-04-07 19:01:53.513 +0800 CST DeletedAt:{Time:0001-01-01 00:00:00 +0000 UTC Valid:false}} Subject:使用自定义的Serializer操作Categories字段 Tags:[Go Serializer Gorm MySQL] Categories:[Go Serializer Gorm MySQL]}
PASS
ok      gormExample     0.082s


# mysql
*************************** 16. row ***************************
        id: 16
created_at: 2023-04-07 19:01:53.513
updated_at: 2023-04-07 19:01:53.513
deleted_at: NULL
   subject: 使用自定义的Serializer操作Categories字段
categories: Go,Serializer,Gorm,MySQL
      tags: ["Go","Serializer","Gorm","MySQL"]
16 rows in set (0.00 sec)
```

## 嵌入结构体和gorm.Model

### gorm.Model

GORM 定义一个 `gorm.Model` 结构体，其包括字段 `ID`、`CreatedAt`、`UpdatedAt`、`DeletedAt`：

```go
type Model struct {
    // Primary Key
	ID        uint `gorm:"primarykey"`
    // 创建时间
	CreatedAt time.Time
    // 创建或更新时间
	UpdatedAt time.Time
    // 删除时间
	DeletedAt DeletedAt `gorm:"index"`
}
```

其中，DeleteAt类型的定义：

```go
type DeletedAt sql.NullTime
```

其中 sql.NullTime 类型指的是可以为Null的Time类型，定义如下：

```go
type NullTime struct {
	Time  time.Time
	Valid bool // Valid is true if Time is not NULL
}
```

若需要使用自定义名字的创建和更新时间，可以使用tag：autoCreateTime 和 autoUpdateTime。

强烈推荐，实体模型都要嵌入 gorm.Model。

### 嵌入结构体

除了gorm.Model, 其他结构体也可以嵌入。

结构体字段标签 embeddedPrefix 用于设置嵌入结构体字段映射到DB中的前缀。

示例：

```go
type Blog struct {
	gorm.Model

	BlogBasic `gorm:""`
	Author    `gorm:"embeddedPrefix:author_"`
}

type BlogBasic struct {
	Subject string
	Summary string
	Content string
}

type Author struct {
	Name  string
	Email string
}

```

迁移生成的表：

```mysql
mysql> desc msb_blog;
+---------------+-----------------+------+-----+---------+----------------+
| Field         | Type            | Null | Key | Default | Extra          |
+---------------+-----------------+------+-----+---------+----------------+
| id            | bigint unsigned | NO   | PRI | NULL    | auto_increment |
| created_at    | datetime(3)     | YES  |     | NULL    |                |
| updated_at    | datetime(3)     | YES  |     | NULL    |                |
| deleted_at    | datetime(3)     | YES  | MUL | NULL    |                |
| subject       | longtext        | YES  |     | NULL    |                |
| summary       | longtext        | YES  |     | NULL    |                |
| content       | longtext        | YES  |     | NULL    |                |
| author_name   | longtext        | YES  |     | NULL    |                |
| author_email  | longtext        | YES  |     | NULL    |                |
+---------------+-----------------+------+-----+---------+----------------+
```

### 实际开发时模型结构体通常有三个部分

```go
type Blog struct {
	// 一：基础结构
	gorm.Model

	// 二：实体字段
	BlogBasic `gorm:""`
	Author    `gorm:"embeddedPrefix:author_"`

	// 三：关联关系
	User user.User
}
```

## 小结

核心内容：

- 表名管理
- 类型管理
  - 默认的映射类型
  - 指针和非指针类型
  - 自定义字段类型
  - 字段的编码解码
- 字段属性管理
  - NOT NULL
  - DEFAULT
  - COMMENT
- 索引和约束的管理
  - 主键约束
  - 唯一键，唯一索引
  - 普通索引
  - 复合索引
  - check约束
- 嵌入结构体的使用


# 具体操作

## 终结方法和链式方法

GORM 中有两种类型的方法：

- 终结方法（Finisher Method），用于生成并执行当前语句的方法

  - ```go
    Create, First, Take, Find, Save, Update, Delete, Scan, Row, Rows
    ```
- 链式方法（Chain Method），用于将子句（Clauses）加入当前语句的方法

  - ```go
    Select, Table, Where, Joins, Group, Having, Order, Limit, Offset, Debug
    ```

示例：

```go
type User struct {
	gorm.Model

	Username string
	Name     string
	Email    string
	Birthday *time.Time
}

func OperatorType() {
	DB.AutoMigrate(&User{})

	var users []User
  
  
    // 一步操作
	//DB.Where("birthday IS NOT NULL").Where("email like ?", "@163.com%").Order("name DESC").Find(&users)

    // 分步操作
	query := DB.Where("birthday IS NOT NULL")
	query.Where("email like ?", "@163.com%")
	query.Order("name DESC")
	query.Find(&users)
}
```

终结方法的核心操作是执行SQL，同时处理结果，处理错误，以DB.First()为例：

```go
// First finds the first record ordered by primary key, matching given conditions conds
func (db *DB) First(dest interface{}, conds ...interface{}) (tx *DB) {
	tx = db.Limit(1).Order(clause.OrderByColumn{
		Column: clause.Column{Table: clause.CurrentTable, Name: clause.PrimaryKey},
	})
	if len(conds) > 0 {
		if exprs := tx.Statement.BuildCondition(conds[0], conds[1:]...); len(exprs) > 0 {
			tx.Statement.AddClause(clause.Where{Exprs: exprs})
		}
	}
	tx.Statement.RaiseErrorOnNotFound = true
	tx.Statement.Dest = dest
    // 执行
	return tx.callbacks.Query().Execute(tx)
}
```

链式方法的核心操作是将特定的子句记录在语句中，以DB.Where为例：

```go
// [docs]: https://gorm.io/docs/query.html#Conditions
func (db *DB) Where(query interface{}, args ...interface{}) (tx *DB) {
	tx = db.getInstance()
	if conds := tx.Statement.BuildCondition(query, args...); len(conds) > 0 {
        // 将where作为子句加入语句
		tx.Statement.AddClause(clause.Where{Exprs: conds})
	}
	return
}
```

使用：

- 终结方法，用于最终的执行，执行完毕后，通常要接收结果，或处理错误
- 链式方法，用于设置子句，执行完毕后，需要配合终结方法才会产生最终的结果
- 终结方法和链式方法，通常都会返回*gorm.DB对象，但终结方法通常会设置了DB对象的错误，因此重复执行可能会出问题

## 错误处理

### DB.Error

在终结方法执行完毕后，会将执行的错误记录在Db对象的Error字段上。因此：

**在终结方法执行后检测错误，是强烈推荐的操作**。

示例：

```go
if result := DB.Create(article); result.Error != nil {
    log.Fatal(result.Error)
}

if err := DB.First(article, id).Error; err != nil {
    log.Fatal(err)
}
```

### ErrRecordNotFound错误

当 `First`、`Last`、`Take` 方法（这几个方法都是查找一条的方法）找不到记录时，GORM 会返回 `ErrRecordNotFound` 错误。这个是GORM的行为，数据库层面在没有记录是不会响应错误。

注意，当 `First`、`Last`、`Take` 方法存在错误时，不一定是ErrRecordNotFound类型，也有可能是其他类型。若有需要判定是否为ErrRecordNotFound类型错误，可以通过 errors.Is() 方法进行判断。

示例：

```
func ErrorHandle() {
	user := User{}
	if err := DB.First(user, 42).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			log.Println("Record Not Found")
		} else {
			log.Fatal(err)
		}
	}
}
```

### Gorm定义的错误

```go
var (
	// ErrRecordNotFound record not found error
	ErrRecordNotFound = logger.ErrRecordNotFound
	// ErrInvalidTransaction invalid transaction when you are trying to `Commit` or `Rollback`
	ErrInvalidTransaction = errors.New("invalid transaction")
	// ErrNotImplemented not implemented
	ErrNotImplemented = errors.New("not implemented")
	// ErrMissingWhereClause missing where clause
	ErrMissingWhereClause = errors.New("WHERE conditions required")
	// ErrUnsupportedRelation unsupported relations
	ErrUnsupportedRelation = errors.New("unsupported relations")
	// ErrPrimaryKeyRequired primary keys required
	ErrPrimaryKeyRequired = errors.New("primary key required")
	// ErrModelValueRequired model value required
	ErrModelValueRequired = errors.New("model value required")
	// ErrModelAccessibleFieldsRequired model accessible fields required
	ErrModelAccessibleFieldsRequired = errors.New("model accessible fields required")
	// ErrSubQueryRequired sub query required
	ErrSubQueryRequired = errors.New("sub query required")
	// ErrInvalidData unsupported data
	ErrInvalidData = errors.New("unsupported data")
	// ErrUnsupportedDriver unsupported driver
	ErrUnsupportedDriver = errors.New("unsupported driver")
	// ErrRegistered registered
	ErrRegistered = errors.New("registered")
	// ErrInvalidField invalid field
	ErrInvalidField = errors.New("invalid field")
	// ErrEmptySlice empty slice found
	ErrEmptySlice = errors.New("empty slice found")
	// ErrDryRunModeUnsupported dry run mode unsupported
	ErrDryRunModeUnsupported = errors.New("dry run mode unsupported")
	// ErrInvalidDB invalid db
	ErrInvalidDB = errors.New("invalid db")
	// ErrInvalidValue invalid value
	ErrInvalidValue = errors.New("invalid value, should be pointer to struct or slice")
	// ErrInvalidValueOfLength invalid values do not match length
	ErrInvalidValueOfLength = errors.New("invalid association values, length doesn't match")
	// ErrPreloadNotAllowed preload is not allowed when count is used
	ErrPreloadNotAllowed = errors.New("preload is not allowed when count is used")
)
```

## 创建Create

### 示例模型

```go
type Content struct {
	gorm.Model

	Subject     string
	Likes       uint
	PublishTime *time.Time
}
```

### 插入及结果

创建记录使用DB.Create()方法实现。

插入成功后，最新的AutoIncrement的ID，可从模型上直接获取。DB.RowsAffected可以获取影响的记录数。

典型的ORM操作，示例：

典型的ORM操作，示例：

```go
func CreateBasic() {
	DB.AutoMigrate(&Content{})

	c1 := Content{}
	c1.Subject = "GORM的使用"

	result1 := DB.Create(&c1)
	if result1.Error != nil {
		log.Fatal(result1.Error)
	}
	fmt.Println(c1.ID, result1.RowsAffected)
}
```

### 支持 Map 设置字段值

插入时，除了可以在模型上配置字段完成数据指定外，还可以使用Map类型，完成字段数据的指定：

使用map类型的逻辑意义就是纯粹数据层面的操作，主动放弃了类似创建时间、修改时间的自动更新功能。

map的类型为：`map[string]any`

需要通过Model()指定对应的模型：

```go
func CreateBasic() {
	DB.AutoMigrate(&Content{})

	// 模型映射记录，操作模型字段，就是操作记录的列
	c1 := Content{}
	c1.Subject = "GORM的使用"

	// 执行创建（insert）
	result1 := DB.Create(&c1)
	// 处理错误
	if result1.Error != nil {
		log.Fatal(result1.Error)
	}
	// 最新的ID，和影响的记录数
	fmt.Println(c1.ID, result1.RowsAffected)

	// map 指定数据
	//设置map 的values
	values := map[string]any{
		"Subject":     "Map指定值",
		"PublishTime": time.Now(),
	}
	// create
	result2 := DB.Model(&Content{}).Create(values)
	if result2.Error != nil {
		log.Fatal(result2.Error)
	}
	// 测试输出
	fmt.Println(result2.RowsAffected)
}
```

### 批量插入

当Create()的参数是模型切片或者Map切片时，Create()支持一次性全部插入，例如：

```go
func CreateMulti() {
	DB.AutoMigrate(&Content{})

	// model
	cs := []Content{
		{Subject: "标题1"},
		{Subject: "标题2"},
		{Subject: "标题3"},
	}
	result1 := DB.Create(&cs)
	if result1.Error != nil {
		log.Fatal(result1.Error)
	}
	fmt.Println(result1.RowsAffected)
	for _, c := range cs {
		fmt.Println(c.ID)
	}

	// map
	vs := []map[string]any{
		{"Subject": "标题4"},
		{"Subject": "标题5"},
		{"Subject": "标题6"},
	}
	result2 := DB.Model(&Content{}).Create(vs)
	if result2.Error != nil {
		log.Fatal(result2.Error)
	}
	fmt.Println(result2.RowsAffected)
}
```

create形成的操作为一条SQL插入全部记录。

```mysql
GORM:2023/04/10 21:14:55 D:/apps/courses/gormExample/operator.go:76
[6.871ms] [rows:3] INSERT INTO `msb_content` (`created_at`,`updated_at`,`deleted_at`,`subject`,`likes`,`publish_time`) VALUES ('2023-04-10 21:14:55.156','2023-04-10 21:14:55.156',NULL,'标题1',0,NULL),('2023-04-10 21:14:55.156','2023-04-10 21:14:55.156',NULL,'标题2',0,NULL),('2023-04-10 21:14:55.156','2023-04-10 21:14:55.156',NULL,'标题3',0,NULL)
GORM:2023/04/10 21:14:55 D:/apps/courses/gormExample/operator.go:91
[6.949ms] [rows:3] INSERT INTO `msb_content` (`subject`) VALUES ('标题4'),('标题5'),('标题6')

```

#### 分批插入

当需要执行大量的记录插入时，推荐使用分批插入。也就是不是一次性全部插入，而是每次插入固定的条数，直到全部插入完毕。分批插入的优势可以避免单条SQL过长问题，也会避免当某条记录插入失败而需要全部重新插入的问题。

方法：

```go
func (db *DB) CreateInBatches(value interface{}, batchSize int) (tx *DB)
```

第二个参数表示批次大小。

示例：

```go
func CreateBatch() {
	DB.AutoMigrate(&Content{})

	// model
	cs := []Content{
		{Subject: "标题1"},
		{Subject: "标题2"},
		{Subject: "标题3"},
		{Subject: "标题4"},
		{Subject: "标题5"},
	}
	result1 := DB.CreateInBatches(&cs, 2)
	if result1.Error != nil {
		log.Fatal(result1.Error)
	}
	fmt.Println(result1.RowsAffected)
	for _, c := range cs {
		fmt.Println(c.ID)
	}

	// map
	vs := []map[string]any{
		{"Subject": "标题6"},
		{"Subject": "标题7"},
		{"Subject": "标题8"},
		{"Subject": "标题9"},
		{"Subject": "标题0"},
	}
	result2 := DB.Model(&Content{}).CreateInBatches(vs, 2)
	if result2.Error != nil {
		log.Fatal(result2.Error)
	}
	fmt.Println(result2.RowsAffected)
}
```

分批插入，会形成多条SQL，每条插入batchSize条记录。参考SQL日志。

```sql
[2.673ms] [rows:2] INSERT INTO `msb_content` (`created_at`,`updated_at`,`deleted_at`,`subject`,`likes`,`publish_time`) VALUES ('2023-04-10 21:15:07.809','2023-04-10 21:15:07.809',NULL,'标题1',0,NULL),('2023-04-10 21:15:07.809','2023-04-10 21:15:07.809',NULL,'标题2',0,NULL)
GORM:2023/04/10 21:15:07 D:/apps/courses/gormExample/operator.go:111
[2.132ms] [rows:2] INSERT INTO `msb_content` (`created_at`,`updated_at`,`deleted_at`,`subject`,`likes`,`publish_time`) VALUES ('2023-04-10 21:15:07.811','2023-04-10 21:15:07.811',NULL,'标题3',0,NULL),('2023-04-10 21:15:07.811','2023-04-10 21:15:07.811',NULL,'标题4',0,NULL)
GORM:2023/04/10 21:15:07 D:/apps/courses/gormExample/operator.go:111
[1.606ms] [rows:1] INSERT INTO `msb_content` (`created_at`,`updated_at`,`deleted_at`,`subject`,`likes`,`publish_time`) VALUES ('2023-04-10 21:15:07.813','2023-04-10 21:15:07.813',NULL,'标题5',0,NULL)

GORM:2023/04/10 21:15:07 D:/apps/courses/gormExample/operator.go:128
[3.414ms] [rows:2] INSERT INTO `msb_content` (`subject`) VALUES ('标题6'),('标题7')
GORM:2023/04/10 21:15:07 D:/apps/courses/gormExample/operator.go:128
[3.326ms] [rows:2] INSERT INTO `msb_content` (`subject`) VALUES ('标题8'),('标题9')
GORM:2023/04/10 21:15:07 D:/apps/courses/gormExample/operator.go:128
[2.717ms] [rows:1] INSERT INTO `msb_content` (`subject`) VALUES ('标题0')
```

#### 配置项 CreateBatchSize

初始化Gorm.DB时，可以使用选项CreateBatchSize选项，控制全部的Create操作都按照分批插入的模型运行。

```go
db, err := gorm.Open(Dialector, &gorm.Config{
  CreateBatchSize: 2,
})
```

使用以上配置，重新测试 `CreateMulti()`，也是分批操作。

使用该配置缺点就是不能选择一次性全部插入了，自由度不足。

### UpSert

UpSert, Update Insert的缩写。逻辑是当插入冲突时（主键或唯一键已经存在），执行更新操作。判定是否冲突的要素通常是主键。

GROM通过 clause.OnConflict{} 类型实现控制冲突后的行为：

```go
type OnConflict struct {
    // 冲突列
	Columns      []Column
	// 冲突条件
    Where        Where
    // 冲突目标条件
	TargetWhere  Where
    // 冲突约束
	OnConstraint string
    // 什么都不做
	DoNothing    bool
    // 做的更新操作，通常是指定更新的具体字段
	DoUpdates    Set
    // 更新全部字段
	UpdateAll    bool
}
```

使用DB.Clauses()子句，传入以上类型实例来配置Create()冲突后的操作。

示例，冲突的插入， 冲突后更新全部字段， ，冲突后更新部分字段：

```go
func UpSert() {
	c := Content{}
	c.Subject = "原始标题"
	c.Likes = 12
	DB.Create(&c)
	fmt.Println(c)

	c2 := Content{}
	c2.ID = c.ID
	c2.Subject = "新标题"
	c2.Likes = 20
	if err := DB.Create(&c2).Error; err != nil {
		log.Fatal(err)
		//Error 1062 (23000): Duplicate entry '13' for key 'msb_content.PRIMARY'
	}

	//c3 := Content{}
	//c3.ID = c.ID
	//c3.Subject = "新标题"
	//c3.Likes = 20
	//if err := DB.
	//	Clauses(clause.OnConflict{UpdateAll: true}).
	//	Create(&c3).Error; err != nil {
	//	log.Fatal(err)
	//}

	//c4 := Content{}
	//c4.ID = c.ID
	//c4.Subject = "新标题"
	//c4.Likes = 20
	//if err := DB.
	//	Clauses(clause.OnConflict{DoUpdates: clause.AssignmentColumns([]string{"likes"})}).
	//	Create(&c4).Error; err != nil {
	//	log.Fatal(err)
	//}
}
```

在MySQL中，使用的是On duplicate key 实现的UpSert。

### 默认值处理

GORM支持使用 default 标签，设置字段的默认值：

```go
Likes       uint  `gorm:"default:100"`
Views       *uint `gorm:"default:100"`
```

**当字段为类型零值时，触发使用默认值**。

问题：由于是类型零值触发默认值，那意味着类似：0, false, "",等都不会保存到数据表中。

方案：可以使用指针类型，或 sql.NullT类型来避免问题。

示例：

```go
type Content struct {
	gorm.Model

	Subject     string
	Likes       uint  `gorm:"default:100"`
	Views       *uint `gorm:"default:100"`
	PublishTime *time.Time
}

func DefaultValue() {
	DB.AutoMigrate(&Content{})

	c := Content{}
	c.Subject = "原始内容"
	likes, views := uint(0), uint(0)
	c.Likes = likes
	c.Views = &views
	DB.Create(&c)
	fmt.Println(c.Likes, *c.Views)
}
```

实操中通常使用模型的创建方法来初始化默认值，不通过定义default标签的方案：

```go
const (
	defaultViews = 99
	defaultLikes = 99
)

func NewContent() Content {
	return Content{
		Likes: defaultLikes,
		Views: defaultViews,
	}
}

func DefaultValueOften() {
	DB.AutoMigrate(&Content{})

	c := NewContent()
	c.Subject = "原始内容"
	DB.Create(&c)
	fmt.Println(c.Likes, c.Views)
}
```

### 插入特定字段

当仅需要操作部分字段时，可以使用方法：

```go
// 选择需要操作的字段
func (db *DB) Select(query interface{}, args ...interface{}) (tx *DB)

// 选择不需要操作的字段
func (db *DB) Omit(columns ...string) (tx *DB)
```

示例：

```go
func SelectCol() {
	DB.AutoMigrate(&Content{})

	c := Content{}
	c.Views = 99
	c.Likes = 7
	c.Subject = "标题"
	now := time.Now()
	c.PublishTime = &now

    // 选择字段
	DB.Select("Subject", "Views", "CreatedAt").Create(&c)
	// INSERT INTO `msb_content` (`created_at`,`updated_at`,`subject`,`views`) VALUES ('2023-04-11 17:51:39.895','2023-04-11 17:51:39.895','标题',99)

    // 忽略字段
	DB.Omit("Subject", "Views", "CreatedAt").Create(&c)
	// INSERT INTO `msb_content` (`updated_at`,`deleted_at`,`likes`,`publish_time`) VALUES ('2023-04-11 17:52:29.034',NULL,7,'2023-04-11 17:52:29.032')
}
```

### 钩子

钩子，Hook，在特定执行时间点执行的方法。类似于事件驱动。CRUD都有对应的钩子方法。

如果我们定义了钩子方法，那么在操作时就会执行对应的钩子方法。

创建时，可用的钩子方法有：

```go
// 开始事务
BeforeSave()
BeforeCreate()
Create()
AfterCreate()
AfterSave()
// 提交或回滚事务

```

依据上面的顺序执行。其中BeforeSave和AfterSave是创建和更新操作通用的。

钩子方法是具体的某个模型的方法，其签名为：

```go
func(*gorm.DB) error
```

在钩子方法中，典型的功能：

- 业务逻辑代码
- 通用配置代码

两类。

示例：

```go
func (c *Content) BeforeCreate(db *gorm.DB) error {
	// 业务代码
	if c.PublishTime == nil {
		now := time.Now()
		c.PublishTime = &now
	}

	// 配置代码
	db.Statement.AddClause(clause.OnConflict{UpdateAll: true})

	return nil
}
```

测试创建：

```go
func CreateHook() {
	DB.AutoMigrate(&Content{})
	c1 := Content{}

	err := DB.Create(&c1).Error
	if err != nil {
		log.Fatal(err)
	}
	//INSERT INTO `msb_content` (`created_at`,`updated_at`,`deleted_at`,`subject`,`likes`,`views`,`publish_time`) VALUES ('2023-04-11 18:44:56.62','2023-04-11 18:44:56.62',NULL,'',0,0,'2023-04-11 18:44:56.62') ON DUPLICATE KEY UPDATE `updated_at`='2023-04-11 18:44:56.62',`deleted_at`=VALUES(`deleted_at`),`subject`=VALUES(`subject`),`likes`=VALUES(`likes`),`views`=VALUES(`views`),`publish_time`=VALUES(`publish_time`)

}
```

不指定PublishTime，会存储当前的时间。同时增加了On Duplicate key update 子句。

注意，包括钩子和Create在内的全部方法处在一个事务中，若钩子方法返回的错误，会导致事务回滚，不会执行后续的操作，即使是AfterX的钩子也是如此。

示例：

```go
func (c *Content) AfterCreate(db *gorm.DB) error {
	return errors.New("custom error")
}
```

此时，模型Content上的Create操作是不会影响DB的。Create方法会得到对应的错误。

## 查询操作

查询操作主要使用方法：

```go
// 查询单条
db.First()
db.Last()
db.Take()

// 查询多条
db.Find()
```

### 主键查询

支持基于主键查询1条或多条：

```go
// 查询一条
db.First(&model, PK)
// 模型的主键字段存在值，自动构建基于主键的查询
model := Model{ID:10}
db.First(&model)

// 查询多条
db.Find(&[]model, []PK{PK1, PK2, ...})
```

若主键为string类型，需要使用条件表达式：

```go
// 查询一条
db.First(&model, "pk = ?", "stringPK")

// 查询多条
db.Find(&[]model, "pk IN ?", []PK{"stringPK1", "stringPK2", ...})
```

示例代码：

```go
type ContentStrPK struct {
	ID          string `gorm:"primaryKey"`
	Subject     string
	Likes       uint
	Views       uint
	PublishTime *time.Time
}

func GetByPk() {
	DB.AutoMigrate(&Content{}, &ContentStrPK{})

	c := Content{}
	if err := DB.First(&c, 10).Error; err != nil {
		log.Println(err)
	}

	cStrPk := ContentStrPK{}
	if err := DB.First(&cStrPk, "id=?", "some id").Error; err != nil {
		log.Println(err)
	}

	var cs []Content
	if err := DB.Find(&cs, []uint{10, 11, 12}).Error; err != nil {
		log.Println(err)
	}

	var cStrPks []ContentStrPK
	if err := DB.Find(&cStrPks, "id IN ?", []string{"some", "id"}).Error; err != nil {
		log.Println(err)
	}
}
```

测试并查看SQL日志。通过SQL了解：

```mysql
SELECT * FROM `msb_content` WHERE `msb_content`.`id` = 10 AND `msb_content`.`deleted_at` IS NULL ORDER BY `msb_content`.`id` LIMIT 1

SELECT * FROM `msb_content_str_pk` WHERE id='some id' ORDER BY `msb_content_str_pk`.`id` LIMIT 1

SELECT * FROM `msb_content` WHERE `msb_content`.`id` IN (10,11,12) AND `msb_content`.`deleted_at` IS NULL

SELECT * FROM `msb_content_str_pk` WHERE id IN ('some','id')
```

测试字符串类型不使用条件表达式的语法：

```go
cStrPk := ContentStrPK{}
if err := DB.First(&cStrPk, "some id").Error; err != nil {
    log.Println(err)
}
```

该查询会触发错误，GORM将"some id"当做查询条件处理了，SQL：

```mysql
GORM:2023/04/14 15:33:09 D:/apps/courses/gormExample/retrieve.go:14 Error 1064 (42000): You have an error in your SQL syntax; check the manual that corresponds to your MySQL server version for the right syntax to use near 'id ORDER BY `msb_content_str_pk`.`id` LIMIT 1' at line 1
[0.536ms] [rows:0] SELECT * FROM `msb_content_str_pk` WHERE some id ORDER BY `msb_content_str_pk`.`id` LIMIT 1
```

### 查询单条

```go
// 查询单条
db.First()
db.Last()
db.Take()

// 带有Limit的Find
db.Limit(1).Find(&model)
```

查询单条可以使用以上三个方法，区别为：

- db.First，主键升序排序的第一条
- db.Last，主键降序排序的第一条
- db.Take，不拼凑排序子句的第一条，数据库的默认返回顺序
- 带有Limit的Find，若Find的结果传递为单模型的引用，也可以查询单条记录。但一定要配合Limit使用，否者会出现查询到多条，然后从中选择一条的严重查询性能Bug。

可见，区别在于使用非主键查询时，有差异。若使用主键查询，结果永远是确定的那条记录。

示例测试：

```go
func GetOne() {
	c := Content{}
	if err := DB.First(&c, "id > ?", 42).Error; err != nil {
		log.Println(err)
	}

	o := Content{}
	if err := DB.Last(&o, "id > ?", 42).Error; err != nil {
		log.Println(err)
	}

	n := Content{}
	if err := DB.Take(&n, "id > ?", 42).Error; err != nil {
		log.Println(err)
	}
  
    f := Content{}
	if err := DB.Limit(1).Find(&f, "id > ?", 42).Error; err != nil {
		log.Println(err)
	}

	fs := Content{}
	if err := DB.Find(&fs, "id > ?", 42).Error; err != nil {
		log.Println(err)
	}
}
```

生成的SQL：

```mysql
[4.794ms] [rows:0] SELECT * FROM `msb_content` WHERE id > 42 AND `msb_content`.`deleted_at` IS NULL ORDER BY `msb_content`.`id` LIMIT 1

[3.784ms] [rows:0] SELECT * FROM `msb_content` WHERE id > 42 AND `msb_content`.`deleted_at` IS NULL ORDER BY `msb_content`.`id` DESC LIMIT 1

[2.700ms] [rows:0] SELECT * FROM `msb_content` WHERE id > 42 AND `msb_content`.`deleted_at` IS NULL LIMIT 1

[2.689ms] [rows:0] SELECT * FROM `msb_content` WHERE id > 42 AND `msb_content`.`deleted_at` IS NULL LIMIT 1

[2.819ms] [rows:0] SELECT * FROM `msb_content` WHERE id > 42 AND `msb_content`.`deleted_at` IS NULL
```

### 扫描查询结果至Map映射表

Gorm除了允许将结果扫描至model或[]model外，还支持将结果扫描至map:

```go
// 单条
map[string]any
// 多条
[]map[string]any
```

由于没有模型了，因此需要通过.Model()方法指定查询的模型。

示例：

```go
func GetToMap() {
	// 单条
	c := map[string]any{} //map[string]interface{}{}
	if err := DB.Model(&Content{}).First(&c, 13).Error; err != nil {
		log.Println(err)
	}
	//fmt.Println(c["id"], c["id"].(uint) == 13)
	// 需要接口类型断言，才能继续处理
	if c["id"].(uint) == 13 {
		fmt.Println("id bingo")
	}
	// time类型的处理
	fmt.Println(c["created_at"])
	t, err := time.Parse("2006-01-02 15:04:05.000 -0700 CST", "2023-04-10 22:00:11.582 +0800 CST")
	if err != nil {
		log.Println(err)
	}
	if c["created_at"].(time.Time) == t {
		fmt.Println("created_at bingo")
	}

	// 多条
	var cs []map[string]any
	if err := DB.Model(&Content{}).Find(&cs, []uint{13, 14, 15}).Error; err != nil {
		log.Println(err)
	}
	for _, c := range cs {
		fmt.Println(c["id"].(uint), c["subject"].(string), c["created_at"].(time.Time))
	}
}
```

- key为字段
- value为字段值

value为any类型，使用时需要类型测试或断言。

### 查询单列Pluck()

除了查询记录，还可以查询单列，使用方法DB.Pluck()实现。

需要使用Model确定映射的表名。

查询的结果是切片类型。

示例：

```go
func GetPluck() {
	var subjects []sql.NullString
	if err := DB.Model(&Content{}).Where("id > ?", 30).Pluck("subject", &subjects).Error; err != nil {
		log.Println(err)
	}

	//if err := DB.Model(&Content{}).Where("id > ?", 30).Pluck("concat(coalesce(subject, 'NULL'), '-', likes)", &subjects).Error; err != nil {
	//	log.Println(err)
	//}
	for _, subject := range subjects {
		if subject.Valid {
			fmt.Println(subject.String)
		} else {
			fmt.Println("NULL")
		}
	}
}
```

使用的是，注意数据表中NULL的处理，因此，若字段允许为NULL，那么尽量使用sql.NullType系列类型。

可以使用数据库函数，构造单列的值，例如连接concat等。

### select字段选择子句

查询时，使用Select()方法指定需要从数据库查询的字段，默认为*全部字段：

```go
func GetSelect() {
	var c Content
	if err := DB.Select("subject", "likes").First(&c, 13).Error; err != nil {
		log.Fatalln(err)
	}
	fmt.Printf("%+v\n", c)
}
```

形成的SQL：

```mysql
[3.001ms] [rows:1] SELECT `subject`,`likes` FROM `msb_content` WHERE `msb_content`.`id` = 13 AND `msb_content`.`deleted_at` IS NULL ORDER BY `msb_content`.`id` LIMIT 1

```

此时，模型的字段值，除了subject和likes之外，都是go类型零值。

```go
{Model:{ID:0 CreatedAt:0001-01-01 00:00:00 +0000 UTC UpdatedAt:0001-01-01 00:00:00 +0000 UTC DeletedA
t:{Time:0001-01-01 00:00:00 +0000 UTC Valid:false}} Subject:原始内容 Likes:12 Views:0 PublishTime:<ni
l>}
```

同样，字段还可以使用表达式代替，示例：

```go
func GetSelect() {
	DB.AutoMigrate(&Content{})

	var c Content
	if err := DB.Select("subject", "concat(subject, views) as sv").First(&c, 13).Error; err != nil {
		log.Fatalln(err)
	}
	fmt.Printf("%+v\n", c)
}
```

形成的SQL：

```mysql
[6.003ms] [rows:1] SELECT `subject`,concat(subject, views) as sv FROM `msb_content` WHERE `msb_content`.`id` = 13 AND `msb_content`.`deleted_at` IS NULL ORDER BY `msb_content`.`id` LIMIT 1

```

如果此时，使用模型接收扫描的结果，那么模型上要具备对应的字段，本例子中就是Sv。由于Sv字段是非结构字段，因此配合权限控制达到目的：

```go
type Content struct {
	gorm.Model
	// 其他字段略
    // 无写权限，无迁移权限
	Sv string `gorm:"<-:false;-:migration"`
}
```

除了在模型字段上处理外，这种自定义的字段还可以使用map结构接收，或者通过执行SQL的方式接收，例如.Row()方法。

### distinct去重子句

查询时去掉重复的行，主要配合Find()使用：

```go
func GetDistinct() {
	var cs []Content

	if err := DB.Distinct("*").Find(&cs).Error; err != nil {
		log.Fatalln(err)
	}
}
```

Distinct()方法需要提供字段列表作为参数，用于表示Select Distinct 后的字段。

形成的SQL：

```mysql
[8.667ms] [rows:42] SELECT DISTINCT * FROM `msb_content` WHERE `msb_content`.`deleted_at` IS NULL

```

### where条件子句

#### 条件设置方法

SQL的where子句用于控制查询条件，使用如下方法来控制where子句。

- DB.Find(), DB.First()，内联条件，将条件放在查询方法中。典型是明确查询条件时使用
- DB.Where()，最典型的条件写法，当动态拼凑条件使用
- DB.Or(), OR条件逻辑或运算
- DB.Not(), Not条件逻辑非运算

示例：

```go
func WhereMethod() {
	var cs []Content

	// SELECT * FROM `msb_content` WHERE likes > 100 AND subject like 'gorm%' AND `msb_content`.`deleted_at` IS NULL
	//if err := DB.Find(&cs, "likes > ? and subject like ?", 100, "gorm%").Error; err != nil {
	//	log.Fatalln(err)
	//}

	// SELECT * FROM `msb_content` WHERE likes > 100 AND subject like 'gorm%' AND `msb_content`.`deleted_at` IS NULL
	//query := DB.Where("likes > ?", 100)
	//query.Where("subject like ?", "gorm%")
	//if err := query.Find(&cs).Error; err != nil {
	//	log.Fatalln(err)
	//}

	// SELECT * FROM `msb_content` WHERE (likes > 100 OR subject like 'gorm%') AND `msb_content`.`deleted_at` IS NULL
	//query := DB.Where("likes > ?", 100)
	//query.Or("subject like ?", "gorm%")
	//if err := query.Find(&cs).Error; err != nil {
	//	log.Fatalln(err)
	//}

	// SELECT * FROM `msb_content` WHERE NOT likes > 100 AND subject like 'gorm%' AND `msb_content`.`deleted_at` IS NULL
	query := DB.Not("likes > ?", 100)
	query.Where("subject like ?", "gorm%")
	if err := query.Find(&cs).Error; err != nil {
		log.Fatalln(err)
	}
}
```

#### 条件表述语法

使用如下的类型来表述方式条件：

- **字符串，配合占位符（匿名和具名占位符）构建条件，最典型的结构，推荐。**
- *gorm.DB类型，分组条件，用于构建复杂的逻辑运算。应该从初始的DB对象进行构建。
- number，主键匹配
- slice, In条件的值，未指定字段，则使用主键
- map，key为字段，value为字段值，通常是=，in运算
- Struct，field为字段，value为字段值，为=运算，零值不视作条件

示例：

```go
func WhereType() {
	var cs []Content
	//
	//query := DB.Where("likes > ? and subject like ?", 100, "gorm%")
	//if err := query.Find(&cs).Error; err != nil {
	//	log.Fatalln(err)
	//}

	//// (1 or 2) and (3 and (4 or 5))
	//cond1 := DB.Where("likes > ?", 100).Or("likes < ?", 1000)
	//cond2 := DB.Where("views > ?", 2000).Or("views < ?", 20000)
	//cond3 := DB.Where("subject like ?", "gorm%").Where(cond2)
	//query := DB.Where(cond1).Where(cond3)
	//if err := query.Find(&cs).Error; err != nil {
	//	log.Fatalln(err)
	//}

	//query := DB.Where("likes = ? AND views IN ?", 100, []uint{1, 2, 3})
	//if err := query.Find(&cs).Error; err != nil {
	//	log.Fatalln(err)
	//}

	//query := DB.Where(map[string]any{
	//	"likes": 100,
	//	"views": []uint{1, 2, 3},
	//})
	//if err := query.Find(&cs).Error; err != nil {
	//	log.Fatalln(err)
	//}

	query := DB.Where(Content{
		Likes: 100,
		Views: 1000,
	})
	if err := query.Find(&cs).Error; err != nil {
		log.Fatalln(err)
	}
}
```

#### 占位符

- ? 匿名占位符，通过索引顺序与数据绑定
- @someName 具名占位符，通过名字与数据绑定

示例：

```go
func PlaceHolder() {
	var cs []Content

	// sql.Named
	query := DB.Where("likes = @like AND views IN @view", sql.Named("view", 1000), sql.Named("like", 100))

	// map[string]any
	//query := DB.Where("likes = @like AND views IN @view", map[string]any{
	//	"view": []uint{1, 2, 3},
	//	"name": 100,
	//})
	if err := query.Find(&cs).Error; err != nil {
		log.Fatalln(err)
	}
}
```

### order排序子句

使用Db.Order()完成排序。

常用的参数为字符串类型，设置order by子句。

支持连续调用，设置多个排序字段。（多个排序字段拼凑成一个字符串也可以）

```go
// 官网示例
db.Order("age desc, name").Find(&users)
// SELECT * FROM users ORDER BY age desc, name;

db.Order("age desc").Order("name").Find(&users)
// SELECT * FROM users ORDER BY age desc, name;
```

还支持子句类型gorm.Clause.OrderBy{}类型，用于构建带有表达式的排序子句。

示例，按照某个值list进行排序：

```go
func OrderBy() {
	var cs []Content

	query := DB.Clauses(clause.OrderBy{
		Expression: clause.Expr{
			SQL:                "field(id, ?)",
			Vars:               []any{[]uint{2, 1, 3}},
			WithoutParentheses: true,
		},
	})
	if err := query.Find(&cs).Error; err != nil {
		log.Fatalln(err)
	}
}

```

拼凑的ORDER By field(id, 2, 1, 3)。

### limit结果限制子句

Gorm使用：

- Limit(n) 限定查询的结果数量，n若为-1，表示不使用限制子句。
- Offset(n) 控制偏移，n若为-1，表示不使用偏移子句。

示例（官网例子）：

```go
db.Limit(3).Find(&users)
// SELECT * FROM users LIMIT 3;

db.Offset(3).Find(&users)
// SELECT * FROM users OFFSET 3;

db.Limit(3).Offset(3).Find(&users)
// SELECT * FROM users LIMIT 3 OFFSET 3;
```

典型的应用场景为分页：

基于用户请求的页码数和每页需要的记录数量（这个也可以后端控制），来确定limit和offset的参数。

示例：

```go
// 定义分页必要数据结构
type Pager struct {
	Page, PageSize int
}

// 默认的值
const (
	DefaultPage     = 1
	DefaultPageSize = 12
)

// 翻页程序
func Pagination(pager Pager) {
	// 确定page, offset 和 pagesize

	page := DefaultPage
	if pager.Page != 0 {
		page = pager.Page
	}

	pagesize := DefaultPageSize
	if pager.PageSize != 0 {
		pagesize = pager.PageSize
	}

	// 计算offset
	// page, pagesize, offset
	// 1, 10, 0
	// 2, 10, 10
	// 3, 10, 20
	offset := pagesize * (page - 1)

	var cs []Content
	// SELECT * FROM `msb_content` WHERE `msb_content`.`deleted_at` IS NULL LIMIT 15 OFFSET 30
	if err := DB.Offset(offset).Limit(pagesize).Find(&cs).Error; err != nil {
		log.Fatalln(err)
	}
}
```

测试 SQL。

#### 使用Scope重用翻页代码

由于分页是非常典型的查询列表操作。因此可以将翻页逻辑重用。

使用GORM提供的Scopes，可以重用。

作用域允许你复用通用的逻辑，这种共享逻辑需要定义为类型 `func(*gorm.DB) *gorm.DB`。

需要通过DB.Scope(func(*gorm.DB) *gorm.DB)来复用。

示例：

```go
// 用于得到func(db *gorm.DB) *gorm.DB类型函数
// 为什么不直接定义函数，因为需要func(db *gorm.DB) *gorm.DB与分页信息产生联系。
func Paginate(pager Pager) func(db *gorm.DB) *gorm.DB {
	// 计算page
	page := DefaultPage
	if pager.Page != 0 {
		page = pager.Page
	}

	// 计算pagesize
	pagesize := DefaultPageSize
	if pager.PageSize != 0 {
		pagesize = pager.PageSize
	}

	// 计算offset
	// page, pagesize, offset
	// 1, 10, 0
	// 2, 10, 10
	// 3, 10, 20
	offset := pagesize * (page - 1)

	return func(db *gorm.DB) *gorm.DB {
		// 使用闭包的变量，实现翻页的业务逻辑
		return db.Offset(offset).Limit(pagesize)
	}
}

// 测试重用的分页查询
func PaginationScope(pager Pager) {

	var cs []Content
	// SELECT * FROM `msb_content` WHERE `msb_content`.`deleted_at` IS NULL LIMIT 15 OFFSET 30
	if err := DB.Scopes(Paginate(pager)).Find(&cs).Error; err != nil {
		log.Fatalln(err)
	}

	var ps []Post
	// SELECT * FROM `msb_post` WHERE `msb_content`.`deleted_at` IS NULL LIMIT 15 OFFSET 30
	if err := DB.Scopes(Paginate(pager)).Find(&ps).Error; err != nil {
		log.Fatalln(err)
	}
}
```

测试生成的SQL.

```go

func TestPaginationScope(t *testing.T) {
	request := Pager{3, 15}
	PaginationScope(request)
}
```

### Count查询

Count用于统计满足条件的记录数量。

GORM提供了独立的Count终结方法完成记录数合计操作。

Count的使用，通常配合翻页使用，用于获取总记录数，以便于统计总页数。

```go
func Count(pager Pager) {

	// 集中的条件，用于统计数量和获取某页记录
	query := DB.Model(&Content{}).
		Where("likes > ?", 99)

	// total rows count
	var count int64
	if err := query.Count(&count).Error; err != nil {
		log.Fatalln(err)
	}
	// SELECT count(*) FROM `msb_content` WHERE likes > 99 AND `msb_content`.`deleted_at` IS NULL
	// 计算总页数 ceil( count / pagesize)

	// rows per page
	var cs []Content
	if err := query.Scopes(Paginate(pager)).Find(&cs).Error; err != nil {
		log.Fatalln(err)
	}
	// SELECT * FROM `msb_content` WHERE likes > 99 AND `msb_content`.`deleted_at` IS NULL LIMIT 15 OFFSET 30
}
```

测试，与PaginationScope参数一致。查看SQL。注意统计记录数和查询记录使用的条件保持一致。

### group和having分组和过滤子句

- group 子句，将结果进行分组
- Having子句，基于分组的结果进行过滤。通常分组后会执行合计操作，例如，count，max，avg等，基于这些操作的结果进行过滤。

分组合计过滤后，得到的数据通常就不是典型的模型或模型集合了，而是自定义的结构体类型或map类型。因此在Find时，通常给的都是自定义的结构体切片。

示例：

```go
type Content struct {
	gorm.Model
	// 其他字段略
    // 增加一个内容作者ID，用该字段分组
	AuthorID uint
}

func GroupHaving() {
	DB.AutoMigrate(&Content{})

	type Result struct {
		UserID     uint
		TotalLikes int
		TotalViews int
		AvgViews   int
	}

	var rs []Result
	if err := DB.Select("author_id", "SUM(likes) as total_likes", "SUM(views) as total_views", "AVG(views) as avg_views").
		Group("author_id").Having("total_views > ?", 99).
		Find(&rs).Error; err != nil {
		log.Fatalln(err)
	}
	// SELECT `author_id`,SUM(likes) as total_likes,SUM(views) as total_views,AVG(views) as avg_views FROM `msb_content` WHERE `msb_content`.`deleted_at` IS NULL GROUP BY `author_id` HAVING total_views > 99

}
```

### 迭代查询

迭代查询，用来减少一次性处理大量数据的压力。也称为流式查询。

使用以下方法：

- DB.Rows()，查询结果
- DB.ScanRows()，将结果扫描至结构体

示例：

```go
func Iterator() {
	// 利用DB.Rows() 获取Rows对象
	rows, err := DB.Model(&Content{}).Rows()
	if err != nil {
		log.Fatalln(err)
	}
	// [rows:-] SELECT * FROM `msb_content` WHERE `msb_content`.`deleted_at` IS NULL

	// 注意：保证使用过后关闭rows结果集
	defer func() {
		_ = rows.Close()
	}()
	fmt.Println(rows)

	// 迭代的从Rows中扫描记录到模型
	for rows.Next() {
		// 还有记录存在与结果集中
		var c Content
		if err := DB.ScanRows(rows, &c); err != nil {
			log.Fatalln(err)
		}
		fmt.Println(c.Subject)
	}
}

```

配合 for 循环，也可以完成结果集中全部记录的遍历。此时应用程序中，每次仅处理一条记录。

**注意：保证使用过后关闭rows结果集。**

### 锁子句

GORM支持在查询时加锁，使用子句 clause.Locking实现：

Locking结构如下：

```go
type Locking struct {
    // 锁强度（类型），典型的SHARE，UPDATE
	Strength string
    // 对应的表
	Table    Table
    // 选项，例如NOWAIT非阻塞
	Options  string
}
```

典型的使用Strength控制共享锁或独占锁。

示例：

```go
func Locking() {
	var cs []Content
	if err := DB.Clauses(clause.Locking{Strength: "UPDATE"}).Find(&cs, "likes > ?", 10).Error; err != nil {
		log.Fatalln(err)
	}
	//[4.904ms] [rows:19] SELECT * FROM `msb_content` WHERE likes > 10 AND `msb_content`.`deleted_at` IS NULL FOR UPDATE

	if err := DB.Clauses(clause.Locking{Strength: "SHARE"}).Find(&cs, "likes > ?", 10).Error; err != nil {
		log.Fatalln(err)
	}
	// [2.663ms] [rows:19] SELECT * FROM `msb_content` WHERE likes > 10 AND `msb_content`.`deleted_at` IS NULL FOR SHARE
}
```

### 子查询

子查询，subquery，嵌入在其他语句中的查询，称为子查询。例如：

```mysql
# 条件子查询
select * from content where author_id in (select id from author where status=0);

# from 子查询
select * from (select subject, likes from content where publish_time is null) as temp where likes > 10;
```

gorm，支持直接使用gorm.DB类型作为参数来构建子查询。

示例：

```go
// Author模型
type Author struct {
	gorm.Model
	Status int

	Name  string
	Email string
}
```

```go
func SubQuery() {
	DB.AutoMigrate(&Author{}, &Content{})

	authorIDs := DB.Model(&Author{}).Select("id").Where("status=?", 0)
	var cs []Content
	if err := DB.
		Where("author_id IN (?)", authorIDs).
		Find(&cs).Error; err != nil {
		log.Fatalln(err)
	}
	// [2.128ms] [rows:0] SELECT * FROM `msb_content` WHERE author_id IN (SELECT `id` FROM `msb_author` WHERE status=0 AND `msb_author`.`deleted_at` IS NULL) AND `msb_content`.`deleted_at` IS NULL

	type Result struct {
		Subject string
		Likes   int
	}
	var rs []Result
	fromQuery := DB.Model(&Content{}).Select("subject", "likes").Where("publish_time is null")
	if err := DB.Table("(?) as temp", fromQuery).
		Where("likes > ?", 10).
		Find(&rs).Error; err != nil {
		log.Fatalln(err)
	}
	// [3.278ms] [rows:17] SELECT * FROM (SELECT `subject`,`likes` FROM `msb_content` WHERE publish_time is null AND `msb_content`.`deleted_at` IS NULL) as temp WHERE likes > 10
}
```

注意，子查询需要使用括号包裹！

### 查询钩子

查询操作支持一个钩子方法：

```go
 AfterFind(tx *gorm.DB) (err error)
```

在查询后执行。先Find，在AfterFind，最后返回。

通常可以用来对数据补充处理。

示例：

```go
func (c *Content) AfterFind(db *gorm.DB) error {
	// 业务代码
	if c.AuthorID == 0 {
		c.AuthorID = 1 // 1 is default author
	}

	return nil
}
```

本例中假定 id==1的author为默认author，例如：

```go
Author{
	ID: 1,
	Name: "默认作者"
}
```

查询测试：

```go
func FindHook() {
	var c Content
	if err := DB.First(&c, 13).Error; err != nil {
		log.Fatalln(err)
	}

	fmt.Println(c.AuthorID)
}
```

c.AuthorID 就是1。

## 更新操作

### 主键更新

模型更新的典型方法 `Save()`，用来存储模型的字段。会基于模型**主键**是否存在有效值（非零值）决定执行Insert或Update操作。

示例：

```go
func UpdatePK() {
	c := Content{}
	if err := DB.Save(&c).Error; err != nil {
		log.Fatalln(err)
	}
	fmt.Println(c)

    if err := DB.Save(&c).Error; err != nil {
		log.Fatalln(err)
	}
	fmt.Println(c)
}
```

因此，更新操作，一定是查询完毕，再执行Save。

更新时，会自动维护 UpdatedAt字段为当前时间。

### 条件更新及更新行数

可以使用struct或map[string]any表示字段和值的对应，来更新满足条件的记录。通常使用：

- Updates() ，执行update操作
- Where(), 设置where条件子句
- Model(), 设置from表名子句

完成更新。

示例：

```go
func UpdateWhere() {
	values := map[string]any{
		"subject": "new subject",
		"likes":   10001,
	}
	result := DB.
		Model(&Content{}).
		Where("likes = ?", 0).
		Updates(values)
	if result.Error != nil {
		log.Fatalln(result.Error)
	}
	log.Println("updated rows num: ", result.RowsAffected)
}
```

条件更新时，通常需要知道更新的记录数量，通过 result.RowsAffected 来获取。

推荐使用map结构表示数据。

若使用struct结构，零值的字段不会被更新。语法示例：

```go
// struct 结构，
values := Content{
    Subject: "new subject",
    Likes:   10001,
}
```

### 阻止无条件的更新

若条件更新时未指定条件，那么GORM不会更新记录，同时会返回 `ErrMissingWhereClause`错误。

此行为的目的，是为了保护失误情况下的全局更新。

示例：

```go
func UpdateNoWhere() {
	// map结构
	values := map[string]any{
		"subject": "new subject",
		"likes":   10001,
	}
	result := DB.
		Model(&Content{}).
		Updates(values)
	if result.Error != nil {
		log.Fatalln(result.Error)
		// WHERE conditions required
	}
	log.Println("updated rows num: ", result.RowsAffected)
}
```

若确实需要全局更新，则设置一个永远为真（1）的条件即可：

```go
Where("1=1")
```

### 表达式值更新

若带更新字段的值为表达式，例如+10。需要使用clause.Expr{}类型进行表示。推荐使用gorm.Expr()方法来构建该类型：

示例：

```go
func UpdateExpr() {
	// 更新的字段值数据
	// map推荐
	values := map[string]any{
		"subject": "Where Update Row",
		// 值为表达式计算的结果时，使用Expr类型
		"likes": gorm.Expr("likes + ?", 10),
		//"likes": "likes + 10",
		// Incorrect integer value: 'likes + 10' for column 'likes' at row 1
	}

	// 执行带有条件的更新
	result := DB.Model(&Content{}).
		Where("likes > ?", 100).
		Updates(values)
	// [17.011ms] [rows:51] UPDATE `msb_content` SET `likes`=likes + 10,`subject`='Where Update Row',`updated_at`='2023-04-21 17:28:45.498' WHERE likes > 100 AND `msb_content`.`deleted_at` IS NULL

	if result.Error != nil {
		log.Fatalln(result.Error)
	}

	// 获取更新结果，更新的记录数量（受影响的记录数）
	// 指的是修改的记录数，而不是满足条件的记录数
	log.Println("updated rows num: ", result.RowsAffected)
}
```

由于likes字段为整型，直接使用 "likes + 10" 是类型不匹配的。

### 更新Hook

类似创建，更新支持4个Hook：

```go
// 执行顺序
// 启动事务
BeforeSave, 
BeforeUpdate, 
Updates()
AfterUpdate
AfterSave,
// 提交事务
```

钩子方法的函数签名：

```go
func(*gorm.DB) error
```

更新操作的Hook中，两个特殊的操作比较常用：

- db.Statement.SetColumn，修改某个特定字段的值，用于before钩子方法。通常可以使用模型直接操作。
- tx.Statement.Changed, 判定某些字段是否发生变化，用于before钩子方法，在使用update或updates方法时起作用。通过与model的字段值比较，判定是否变化。

## 删除操作

### 主键删除

将具有主键的模型作为参数传递给 `DB.Delete()` 方法，会删除该模型对应的记录。

参考，基础操作删除部分。默认删除，是通过将DeleteAt字段设置为删除时间来实现的。若不存在DeleteAt字段，会执行Delete操作完成删除。

示例：

```go
func Delete() {
	// 获取模型对象
	article := &Article{}
	if err := DB.First(article, 1).Error; err != nil {
		log.Fatal(err)
	}

	// DB.Delete() 删除
	if err := DB.Delete(article).Error; err != nil {
		log.Fatal(err)
	}

	// print
	fmt.Println("article was deleted")

}

// 当然也可以
DB.Delete(&Article{}, 1)
```

### 条件删除

通过Where子句，或Delete的内联条件，可以删除满足条件的记录。

示例：

```go
func DeleteWhere() {
	if err := DB.Delete(&Content{}, "likes < ?", 100).Error; err != nil {
		log.Fatalln(err)
	}
	// [7.893ms] [rows:0] UPDATE `msb_content` SET `deleted_at`='2023-04-21 18:57:13.338' WHERE likes < 100 AND `msb_content`.`deleted_at` IS NULL
	if err := DB.Where("likes < ?", 100).Delete(&Content{}).Error; err != nil {
		log.Fatalln(err)
	}
}

```

### 阻止无条件的删除

同样，若没有指定删除条件，既没有ID，也没有条件。GORM将不会运行并返回 `ErrMissingWhereClause` 错误。也是处于数据安全的角度考虑。

要删除全部，通过执行永远为真的条件即可。

```go
.Where("1=1")
```

### 删除的行数

通过 result.RowsAffected，返回删除的行数。

```go
fmt.Println(result.RowsAffected)
```

### 删除的行数

通过 result.RowsAffected，返回删除的行数。

```go
fmt.Println(result.RowsAffected)
```

### 逻辑删除

也叫软删除。

如果模型中包含了 gorm.DeletedAt 字段，模型自动获得软删除能力。

可以通过嵌入gorm.Model实现，也可以通过定义gorm.DeletedAt类型来实现：

```go
type M struct {
	ID uint
	DeletedAt gorm.DeletedAt
}
```

当调用 `Delete`时，GORM并不会从数据库中删除该记录，而是将该记录的 `DeleteAt`设置为当前时间，而后的一般查询方法将无法查找到此条记录。也就是会自动增加And deleted_at is null 的where条件。

示例：

```go
// user's ID is `111`
db.Delete(&user)
// UPDATE users SET deleted_at="2013-10-29 10:23" WHERE id = 111;

// Batch Delete
db.Where("age = ?", 20).Delete(&User{})
// UPDATE users SET deleted_at="2013-10-29 10:23" WHERE age = 20;

// Soft deleted records will be ignored when querying
db.Where("age = 20").Find(&user)
// SELECT * FROM users WHERE age = 20 AND deleted_at IS NULL;

```

#### 查询被删除记录

使用 `DB.Unscoped`能发来查询到被软删除的记录：

```go
func FindDeleted() {
	var c Content
	DB.Delete(&c, 13)

	if err := DB.First(&c, 13).Error; err != nil {
		log.Println(err)
	}
	//[4.604ms] [rows:0] SELECT * FROM `msb_content` WHERE `msb_content`.`id` = 13 AND `msb_content`.`deleted_at` IS NULL ORDER BY `msb_content`.`id` LIMIT 1

	if err := DB.Unscoped().First(&c, 13).Error; err != nil {
		log.Println(err)
	}
	// [3.320ms] [rows:1] SELECT * FROM `msb_content` WHERE `msb_content`.`id` = 13 ORDER BY `msb_content`.`id` LIMIT 1
	fmt.Printf("%+v\n", c)
}
```

#### 物理删除

你可以使用 `DB.Unscoped`方法来永久删除匹配的记录

```go
func DeleteHard() {
	var c Content
	if err := DB.Unscoped().Delete(&c, 14).Error; err != nil {
		log.Fatalln(err)
	}
	//	[8.135ms] [rows:0] DELETE FROM `msb_content` WHERE `msb_content`.`id` = 14
}
```

## 原生SQL

原生SQL，指的是我们使用 标准SQL语句完成DB操作。

当需要执行原生的SQL的时，将SQL中的参数使用占位符占位，之后提供变量，拼凑构造成完整的SQL进行执行。

查询时，使用各种子句方法，其实就是再构造SQL的各个部分。

执行SQL，分为两种：

- 查询类型，存在返回数据的，典型就是select, show, desc等。利用DB.Raw()和DB.Scan()完成，通常需要定义响应结果结构。
- 非查询类，没有返回数据的， insert， update，delete，DDL等。利用DB.Exec()完成。

示例：

```go
// 原生查询测试
func RawSelect() {
	// 结果类型
	type Result struct {
		ID           uint
		Subject      string
		Likes, Views int
	}
	var rs []Result

	// SQL
	sql := "SELECT `id`, `subject`, `likes`, `views` FROM `msb_content` WHERE `likes` > ? ORDER BY `likes` DESC LIMIT ?"

	// 执行SQL，并扫描结果
	if err := DB.Raw(sql, 99, 12).Scan(&rs).Error; err != nil {
		log.Fatalln(err)
	}
	// [8.298ms] [rows:12] SELECT `id`, `subject`, `likes`, `views` FROM `msb_content` WHERE `likes` > 99 ORDER BY `likes` DESC LIMIT 12

	log.Println(rs)
}

// 执行类的SQL原生
func RawExec() {
	// SQL
	sql := "UPDATE `msb_content` SET `subject` = CONCAT(`subject`, '-new postfix') WHERE `id` BETWEEN ? AND ?"

	// 执行，获取结果
	result := DB.Exec(sql, 30, 40)
	if result.Error != nil {
		log.Fatalln(result.Error)
	}
	// [13.369ms] [rows:10] UPDATE `msb_content` SET `subject` = CONCAT(`subject`, '-new postfix') WHERE `id` BETWEEN 30 AND 40

	log.Println(result.RowsAffected)

}
```

### sql.Row和sql.Rows

若需要在原生SQL查询时，使用标准结构sql.Row和sql.Rows时，调用DB.Row()和DB.Rows()方法：

```go
func (db *DB) Rows() (*sql.Rows, error) 
func (db *DB) Row() *sql.Row
```

获取Row或Rows后，需要扫描到结果变量来使用：

- 将Row扫描到单独变量，row.Scan
- 将Row扫描到整体结构体类型，DB.ScanRow()
- 判断Rows中是否存在记录，rows.Next()

示例代码：

```go
// 原生查询测试
func RawSelect() {
	// 结果类型
	type Result struct {
		ID           uint
		Subject      string
		Likes, Views int
	}
	var rs []Result

	// SQL
	sql := "SELECT `id`, `subject`, `likes`, `views` FROM `msb_content` WHERE `likes` > ? ORDER BY `likes` DESC LIMIT ?"

	// 执行SQL，并扫描结果
	if err := DB.Raw(sql, 99, 12).Scan(&rs).Error; err != nil {
		log.Fatalln(err)
	}
	// [8.298ms] [rows:12] SELECT `id`, `subject`, `likes`, `views` FROM `msb_content` WHERE `likes` > 99 ORDER BY `likes` DESC LIMIT 12

	log.Println(rs)
}
```

取决于项目对数据类型的需要，选择模型方案或者行方案。

## 会话模式

### 会话模式的基本使用

GORM支持链式操作，意味着前边的操作会对后影响后边的调用，示例：

```go
func SessionIssue() {
	//
	db := DB.Model(&Content{}).Where("views > ?", 100)
	db.Where("likes > ?", 99)
	var cs []Content
	db.Find(&cs)
}
```

这在很多时候很好用。

但当我们需要连续执行多次查询时，就可能出问题，导致子句的重叠，示例：

```go
func SessionIssue() {
	//
	//db := DB.Model(&Content{}).Where("views > ?", 100)
	//db.Where("likes > ?", 99)
	//var cs []Content
	//db.Find(&cs)
	//[3.259ms] [rows:0] SELECT * FROM `msb_content` WHERE likes > 99 AND likes < 199 AND `msb_content`.`deleted_at` IS NULL

	db := DB.Model(&Content{}).Where("views > ?", 100)
	var cs1 []Content
	db.Where("likes > ?", 99).Find(&cs1)
	// [3.777ms] [rows:0] SELECT * FROM `msb_content` WHERE views > 100 AND likes > 99 AND `msb_content`.`deleted_at` IS NULL

	var cs2 []Content
	db.Where("likes > ?", 199).Find(&cs2)
	// [2.638ms] [rows:0] SELECT * FROM `msb_content` WHERE views > 100 AND likes > 99 AND `msb_content`.`deleted_at` IS NULL AND likes > 199
}
```

注意 cs2 的查询，有两个 likes > ? 条件，这在逻辑上产生了冲突。

要解决以上的条件（或其他子句）重叠的问题，通常有两个方案：

- 每个查询都从DB开始构建，DB调用的第一个方法，会重新初始化新的DB对象、reinitialize
- 使用Session会话，可以让某些子句重用

从DB开始示例：

```go
func SessionDB() {
	// 连续执行查询
	// 1
	// Where("views > ?", 100).Where("likes > ?", 9)
	db1 := DB.Model(&Content{}).Where("views > ?", 100)
	db1.Where("likes > ?", 9)
	var cs1 []Content
	db1.Find(&cs1)
	// [10.683ms] [rows:0] SELECT * FROM `msb_content` WHERE views > 100 AND likes > 9 AND `msb_content`.`deleted_at` IS NULL

	// 2,找到likes<5
	// Where("views > ?", 100).Where("likes < ?", 5)
	db2 := DB.Model(&Content{}).Where("views > ?", 100)
	db2.Where("likes < ?", 5)
	var cs2 []Content
	db2.Find(&cs2)
	// [4.139ms] [rows:0] SELECT * FROM `msb_content` WHERE views > 100 AND likes < 5 AND `msb_content`.`deleted_at` IS NULL
}
```

Session示例：

```go
func SessionNew() {

	// 需要重复使用的部分
	// 将Session方法前的配置，记录到了当前的会话中
	// 后边再次调用db的方法直到终结方法，会保持会话中的子句选项
	// 执行完终结方法后，再次调用db的方法到终结方法，可以重用会话中的子句选项。
	db := DB.Model(&Content{}).Where("views > ?", 100).Session(&gorm.Session{})

	// 连续执行查询
	// 1
	// Where("views > ?", 100).Where("likes > ?", 9)
	var cs1 []Content
	db.Where("likes > ?", 9).Find(&cs1)
	// [4.633ms] [rows:0] SELECT * FROM `msb_content` WHERE views > 100 AND likes > 9 AND `msb_content`.`deleted_at` IS NULL

	// 2,找到likes<5
	// Where("views > ?", 100).Where("likes < ?", 5)
	var cs2 []Content
	db.Where("likes < ?", 5).Find(&cs2)
	// [3.846ms] [rows:0] SELECT * FROM `msb_content` WHERE views > 100 AND likes < 5 AND `msb_content`.`deleted_at` IS NULL
}
```

上面代码使用Session方法启动了新的Session，意味着，db对象可以保持会话开启前的状态。继续使用db对象时，到执行终结方法前，都是以该会话状态为初始化状态的。这就可以保证，会话中的子句可以重用。

### Session会话的常用选项

```go
// Session session config when create session with Session() method
type Session struct {
	DryRun                   bool
	PrepareStmt              bool
	NewDB                    bool
	Initialized              bool
	SkipHooks                bool
	SkipDefaultTransaction   bool
	DisableNestedTransaction bool
	AllowGlobalUpdate        bool
	FullSaveAssociations     bool
	QueryFields              bool
	Context                  context.Context
	Logger                   logger.Interface
	NowFunc                  func() time.Time
	CreateBatchSize          int
}
```

#### 禁用Hook

官方例子：

```go
DB.Session(&gorm.Session{SkipHooks: true}).Create(&user)
DB.Session(&gorm.Session{SkipHooks: true}).Delete(&user)
DB.Session(&gorm.Session{SkipHooks: true}).Find(&user)
DB.Session(&gorm.Session{SkipHooks: true}).Model(User{}).Where("age > ?", 18).Updates(&user)
DB.Session(&gorm.Session{SkipHooks: true}).Save(&user)

```

示例：

```go
func SessionOption() {
	db := DB.Model(&Content{}).Session(&gorm.Session{
		SkipHooks: true,
	})
	db.Save(&Content{Subject: "no hooks"})
}

func (c *Content) BeforeCreate(db *gorm.DB) error {
	log.Println("content before create hook")
	return nil
}
```

通过修改 SkipHooks: true， 可以看到是否有输出content before create hook。

#### DryRun模式

生成 `SQL` 但不执行。 它可以用于准备或测试生成的 SQL：

```go
&Session{DryRun: true}
```

示例：

```
func SessionOption() {
	// DryRun
	db := DB.Model(&Content{}).Session(&gorm.Session{
		DryRun: true,
	})
	stmt := db.Save(&Content{Subject: "no hooks"}).Statement
	log.Println(stmt.SQL.String())
	log.Println(stmt.Vars)
}
```

#### 预编译模式

在执行 SQL 时都会创建一个 prepared statement 并将其缓存，以提高后续的效率：

```go
// Session配置
&Session{
    PrepareStmt: true,
}

// gorm.Open配置
&gorm.Config{
 	PrepareStmt: true,
}

```

当连续执行结构一致，但数据不同的SQL时，可以利用预编译的SQL缓存，提升效率。

示例：

```go
func SessionOption() {
	// prepare
	db := DB.Model(&Content{}).Session(&gorm.Session{
		PrepareStmt: true,
	})

	stmtManger, ok := db.ConnPool.(*gorm.PreparedStmtDB)
	if !ok {
		log.Fatalln("*gorm.PreparedStmtDB assert failed")
	}
	log.Println(stmtManger.PreparedSQL)

	var c1 Content
	db.First(&c1, 13)
	log.Println(stmtManger.PreparedSQL)
	var c2 Content
	db.First(&c2, 13)
	var c3 Content
	db.First(&c3, 13)
}
```

#### 允许全局Update/Delete

MissingWhereClause

```go
&gorm.Session{
  AllowGlobalUpdate: true,
}
```

**不要这么做！**

#### Debug()

Debug 利用将日志级别更改为logger.Info来实现。

```go
func (db *DB) Debug() (tx *DB) {
  return db.Session(&Session{
    Logger:         db.Logger.LogMode(logger.Info),
  })
}
```

示例：

```go
DB.Debug().First(&c, 13)
```

#### 初始化

得到一个新初始化的DB对象，官网例子。

目的是取消之前的全部链式方法。

```go
tx := db.Session(&gorm.Session{Initialized: true})

```

## Context支持

GORM支持Context：

使用 DB.WithContext() 或 &Session{Context: ctx} 字段进行配置。

示例，控制执行时间的Context：

```go
func ContextTOCancel() {
	ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond)
	defer cancel()

	var cs []Content
	if err := DB.WithContext(ctx).Limit(10).Find(&cs).Error; err != nil {
		log.Fatalln(err)
	}
	fmt.Println(cs)
}

```

在预设的时间没有执行完毕的话，DB会返回错误。

```go
> go test -run ContextTOCancel
2023/04/24 19:16:32 context deadline exceeded
exit status 1       
FAIL    gormExample     0.052s

```

可以在测试方法中定义context对象，传递到功能方法中：

```go
func ContextTimeoutCancel(ctx context.Context) {
	// 传递Context执行
	var cs []Content
	if err := DB.WithContext(ctx).Limit(5).Find(&cs).Error; err != nil {
		log.Fatalln(err)
	}
	fmt.Println(cs)
}

func TestContextTimeoutCancel(t *testing.T) {
	// 设置一个定时Cancel的Context
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Millisecond)
	defer cancel()

	ContextTimeoutCancel(ctx)
}
```

通常在整体请求周期内，设置一个Deadline，保证不会一直持久执行。

DB.Statement.Context 可以用来访问Context对象，完成自定义操作。例如在Hook中。
