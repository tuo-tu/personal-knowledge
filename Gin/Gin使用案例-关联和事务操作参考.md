# 关联操作

## 标准模型结构

关系型数据库中，二维表间的关系：

- 一对一
- 多对一，一对多
- 多对多

如图所示：

![image.png](https://fynotefile.oss-cn-zhangjiakou.aliyuncs.com/fynote/fyfile/13080/1682393879066/ad09c9ccbcc744e2acd05de7d372d468.png)

其中：

- Author 和 Author间，在Author的角度是一对多，在Author的角度是多对一
- Author和Tag间，是多对多
- Author和AuthorMate间，既可以是一对多，也可以做一对一，看业务逻辑，本例中我们采用一对一

在GORM中，可以在模型中定义关联的方式，实现以上的对应的关系：

- 使用模型类型，表示对应一个的关系
- 使用模型切片类型，表示对应多个的关系
- 使用tag，many2many表示多对多关系，需要制定关联表名
- 需要使用外键字段确保关联。默认的关联字段是模型+ID的形式。
  - 例如Author一对多关联Author，那么Author中就应该有AuthorID作为关联字段
  - 允许自定义

示例代码：

```go
// Author模型
type Author struct {
	gorm.Model
	Status int
	Name   string
	Email  string

	// 拥有多个论文内容
	Essays []Essay
}

// 论文内容
type Essay struct {
	gorm.Model
	Subject string
	Content string

	// 外键字段
	AuthorID *uint

	// 属于某个作者
	Author Author

	// 拥有一个论文元信息
	EssayMate EssayMate

	// 拥有多个Tag
	Tags []Tag `gorm:"many2many:essay_tag"`
}

// 论文元信息
type EssayMate struct {
	gorm.Model
	Keyword     string
	Description string

	// 外键字段
	EssayID *uint

	// 属于一个论文内容，比较少用
	//Essay *Essay
}

type Tag struct {
	gorm.Model
	Title string

	// 拥有多个Essay
	Essays []Essay `gorm:"many2many:essay_tag"`
}
```

使用Migrate创建表。以上模型会创建5张表，会自动创建多对多关联表essay_tag。

创建表，及对应的SQL查看外键索引和约束：

```go
func StdAssocModel() {
	// 利用migrate创建表
	// 以及多对多的关联表
	// 以及外键约束
	if err := DB.AutoMigrate(&Author{}, &Essay{}, &Tag{}, &EssayMate{}); err != nil {
		log.Fatalln(err)
	}
	// CREATE TABLE `msb_author` (
	//  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
	//  `created_at` datetime(3) DEFAULT NULL,
	//  `updated_at` datetime(3) DEFAULT NULL,
	//  `deleted_at` datetime(3) DEFAULT NULL,
	//  `status` bigint DEFAULT NULL,
	//  `name` longtext,
	//  `email` longtext,
	//  PRIMARY KEY (`id`),
	//  KEY `idx_msb_author_deleted_at` (`deleted_at`)
	//) ENGINE=InnoDB AUTO_INCREMENT=11 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci

	// CREATE TABLE `msb_essay` (
	//  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
	//  `created_at` datetime(3) DEFAULT NULL,
	//  `updated_at` datetime(3) DEFAULT NULL,
	//  `deleted_at` datetime(3) DEFAULT NULL,
	//  `subject` longtext,
	//  `content` longtext,
	//  `author_id` bigint unsigned DEFAULT NULL,
	//  PRIMARY KEY (`id`),
	//  KEY `idx_msb_essay_deleted_at` (`deleted_at`),
	//  KEY `fk_msb_author_essays` (`author_id`),
	//  CONSTRAINT `fk_msb_author_essays` FOREIGN KEY (`author_id`) REFERENCES `msb_author` (`id`)
	//) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci

	// CREATE TABLE `msb_essay_mate` (
	//  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
	//  `created_at` datetime(3) DEFAULT NULL,
	//  `updated_at` datetime(3) DEFAULT NULL,
	//  `deleted_at` datetime(3) DEFAULT NULL,
	//  `keyword` longtext,
	//  `description` longtext,
	//  `essay_id` bigint unsigned DEFAULT NULL,
	//  PRIMARY KEY (`id`),
	//  KEY `idx_msb_essay_mate_deleted_at` (`deleted_at`),
	//  KEY `fk_msb_essay_essay_mate` (`essay_id`),
	//  CONSTRAINT `fk_msb_essay_essay_mate` FOREIGN KEY (`essay_id`) REFERENCES `msb_essay` (`id`)
	//) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci

	// CREATE TABLE `msb_tag` (
	//  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
	//  `created_at` datetime(3) DEFAULT NULL,
	//  `updated_at` datetime(3) DEFAULT NULL,
	//  `deleted_at` datetime(3) DEFAULT NULL,
	//  `title` longtext,
	//  PRIMARY KEY (`id`),
	//  KEY `idx_msb_tag_deleted_at` (`deleted_at`)
	//) ENGINE=InnoDB AUTO_INCREMENT=9 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci

	// CREATE TABLE `msb_essay_tag` (
	//  `tag_id` bigint unsigned NOT NULL,
	//  `essay_id` bigint unsigned NOT NULL,
	//  PRIMARY KEY (`tag_id`,`essay_id`),
	//  KEY `fk_msb_essay_tag_essay` (`essay_id`),
	//  CONSTRAINT `fk_msb_essay_tag_essay` FOREIGN KEY (`essay_id`) REFERENCES `msb_essay` (`id`),
	//  CONSTRAINT `fk_msb_essay_tag_tag` FOREIGN KEY (`tag_id`) REFERENCES `msb_tag` (`id`)
	//) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
	log.Println("migrate successful")
}
```

在GORM中，将模型关联分为四种：

- Has Many，一对多关系中，一端有多个多端，外键定义在多端
- Belongs To，一对多关系中，多端属于一端，外键定义在多端
- Many to Many，多对多，外键定义在关联表中
- Has One，一对多关系中，一端有一个多端，外键定义在多端

注意Author和AuthorMate的关系定义：

- 当前结构中，可以表示一对多，也可以表示一对一
- 本例中，选择了一对一
- 若需要一对多，那么增加Author中的关联定义 `AuthorMates []AuthorMate`

> 其实本质上就一种关系，就是外键关系。一重外键关系，就是一对多，多重外键关系就是多对多。

## 建立关联

操作关联时，使用方法

```go
db.Model(&model).Association("Association")
```

完成关联的建立。参数是模型中定义的关联字段，具体的关联类型取决于模型的定义。

要求model的主键不能为空。

关联建立后，即可完成关联的管理。

## 添加关联

.Append() 方法添加关联。

参数为需要关联的模型，或模型切片。取决于是对一还是对多。

其中：

- `many to many`、`has many` 添加新的关联
- `has one`, `belongs to` 替换当前的关联

示例：

```go
// 添加关联
func AssocAppend() {
	// A：一对多的关系, Author 1:n Essay
	// 创建测试数据
	var a Author
	a.Name = "一位作者"
	if err := DB.Create(&a).Error; err != nil {
		log.Println(err)
	}
	log.Println("a:", a.ID)
	var e1, e2 Essay
	e1.Subject = "一篇内容"
	//e1.AuthorID = a.ID
	e2.Subject = "另一篇内容"
	if err := DB.Create([]*Essay{&e1, &e2}).Error; err != nil {
		log.Println(err)
	}
	log.Println("e1, e2: ", e1.ID, e2.ID)

	// 添加关联
	if err := DB.Model(&a).Association("Essays").Append([]Essay{e1}); err != nil {
		log.Println(err)
	}
	fmt.Println(len(a.Essays))
	// 基于当前的基础上，添加关联
	if err := DB.Model(&a).Association("Essays").Append([]Essay{e2}); err != nil {
		log.Println(err)
	}
	fmt.Println(len(a.Essays))
	// 添加后，a模型对象的Essays字段，自动包含了关联的Essay模型
	//fmt.Println(a.Essays)

	// B: Essay M:N TAg
	var t1, t2, t3 Tag
	t1.Title = "Go"
	t2.Title = "GORM"
	t3.Title = "Ma"
	if err := DB.Create([]*Tag{&t1, &t2, &t3}).Error; err != nil {
		log.Println(err)
	}
	log.Println("t1, t2, t3: ", t1.ID, t2.ID, t3.ID)

	// e1 t1, t3
	// e2 t1, t2, t3
	if err := DB.Model(&e1).Association("Tags").Append([]Tag{t1, t3}); err != nil {
		log.Println(err)
	}

	if err := DB.Model(&e2).Association("Tags").Append([]Tag{t1, t2, t3}); err != nil {
		log.Println(err)
	}

	// 关联表查看
	// mysql> select * from msb_essay_tag;
	//+--------+----------+
	//| tag_id | essay_id |
	//+--------+----------+
	//|      1 |       12 |
	//|      3 |       12 |
	//|      1 |       13 |
	//|      2 |       13 |
	//|      3 |       13 |
	//+--------+----------+

	// C, Belongs To. Essay N:1 Author
	var e3 Essay
	e3.Subject = "第三篇内容"
	if err := DB.Create([]*Essay{&e3}).Error; err != nil {
		log.Println(err)
	}
	log.Println("e3: ", e3.ID)

	log.Println(e3.Author)
	// 关联
	if err := DB.Model(&e3).Association("Author").Append(&a); err != nil {
		log.Println(err)
	}
	log.Println(e3.Author.ID)

	// 对一的关联，会导致关联被更新
	var a2 Author
	a2.Name = "另一位作者"
	if err := DB.Create(&a2).Error; err != nil {
		log.Println(err)
	}
	log.Println("a2:", a2.ID)
	if err := DB.Model(&e3).Association("Author").Append(&a2); err != nil {
		log.Println(err)
	}
	log.Println(e3.Author.ID)

}
```

查看数据表，注意关联外键字段，是否记录了关联关系。

## 替换关联

使用新的关联关系，替换旧的关系。使用方法：.Replace() 完成

主要用在对多的关系上。

示例：

```go
func AssocReplace() {
	// A. 替换
	// 创建测试数据
	var a Author
	a.Name = "一位作者"
	if err := DB.Create(&a).Error; err != nil {
		log.Println(err)
	}
	log.Println("a:", a.ID)

	var e1, e2, e3 Essay
	e1.Subject = "一篇内容"
	e2.Subject = "另一篇内容"
	e3.Subject = "第三篇内容"
	if err := DB.Create([]*Essay{&e1, &e2, &e3}).Error; err != nil {
		log.Println(err)
	}
	log.Println("e1, e2, e3: ", e1.ID, e2.ID, e3.ID)

	// 添加关联
	if err := DB.Model(&a).Association("Essays").Replace([]Essay{e1, e3}); err != nil {
		log.Println(err)
	}
	fmt.Println(len(a.Essays))
	// 基于当前的基础上，添加关联
	if err := DB.Model(&a).Association("Essays").Replace([]Essay{e2, e3}); err != nil {
		log.Println(err)
	}
	fmt.Println(len(a.Essays))

}
```

## 删除关联

删除与某模型间的关联关系：使用方法：.Delete() 完成

- 多对一、一对多，删除关联字段
- 多对多，删除关联记录
- 对应的实体记录不会删除

示例：

```go
// 参考清空关联
```

## 清空关联

删除全部关联。：使用方法：.Clear() 完成

```go
func AssocDelete() {
	// B. 删除，外键的
	// 创建测试数据
	var a Author
	a.Name = "一位作者"
	if err := DB.Create(&a).Error; err != nil {
		log.Println(err)
	}
	log.Println("a:", a.ID)

	var e1, e2, e3 Essay
	e1.Subject = "一篇内容"
	e2.Subject = "另一篇内容"
	e3.Subject = "第三篇内容"
	if err := DB.Create([]*Essay{&e1, &e2, &e3}).Error; err != nil {
		log.Println(err)
	}
	log.Println("e1, e2, e3: ", e1.ID, e2.ID, e3.ID)

	// 添加关联
	if err := DB.Model(&a).Association("Essays").Replace([]Essay{e1, e2, e3}); err != nil {
		log.Println(err)
	}
	fmt.Println(len(a.Essays))

	if err := DB.Model(&a).Association("Essays").Delete([]Essay{e1, e3}); err != nil {
		log.Println(err)
	}
	fmt.Println(len(a.Essays))
	fmt.Println("------------------------")

	// B. 删除，多对多，关联表
	var t1, t2, t3 Tag
	t1.Title = "Go"
	t2.Title = "GORM"
	t3.Title = "Ma"
	if err := DB.Create([]*Tag{&t1, &t2, &t3}).Error; err != nil {
		log.Println(err)
	}
	log.Println("t1, t2, t3: ", t1.ID, t2.ID, t3.ID)
	// e1 t1, t3
	// e2 t1, t2, t3
	if err := DB.Model(&e1).Association("Tags").Append([]Tag{t1, t2, t3}); err != nil {
		log.Println(err)
	}
	fmt.Println(len(e1.Tags))

	if err := DB.Model(&e1).Association("Tags").Delete([]Tag{t1, t3}); err != nil {
		log.Println(err)
	}
	fmt.Println(len(e1.Tags))

	// C. 清空关联
	if err := DB.Model(&e1).Association("Tags").Clear(); err != nil {
		log.Println(err)
	}
	fmt.Println(len(e1.Tags))
}
```

## 关联查询

使用Find方法，可以查找关联。查找的结果通常是关联的模型或模型切片，支持子句过滤，例如条件，排序，Limit等：

示例：

```go
func AssocFind() {
	//
	e := Essay{}
	DB.First(&e, 18)

	// 查询关联的tags
	//var ts []Tag
	if err := DB.Model(&e).Association("Tags").Find(&e.Tags); err != nil {
		log.Println(err)
	}
	log.Println(e.Tags)

	// 子句，要写在Association()方法前面
	if err := DB.Model(&e).
		Where("tag_id > ?", 7).
		Order("tag_id DESC").
		Association("Tags").Find(&e.Tags); err != nil {
		log.Println(err)
	}
	log.Println(e.Tags)

	// 查询关联的模型的数量
	count := DB.Model(&e).Association("Tags").Count()
	log.Println("count:", count)

}
```

## 关联统计

.Count()方法可以返回关联的数量，不用查询到全部的关联实体。

```go
func AssocFind() {
	//
	e := Essay{}
	DB.First(&e, 18)

	// 查询关联的模型的数量
	count := DB.Model(&e).Association("Tags").Count()
	log.Println("count:", count)

}
```

## 自动存储关联

在创建或更新模型时，如果关联模型存在，GORM会自动存储关联：

示例：

```go
var t1 Tag
DB.First(&t1, 10)

e := Essay{
    Subject: "一个组合的Save",
    Author:  Author{Name: "马士兵"},
    Tags: []Tag{
        t1,
        {Title: "Go"},
        {Title: "GORM"},
    },
}
```

上面的数据，可以一次性在数据表中更新完成。

执行Save：

```go
func AssocSave() {
	var t1 Tag
	DB.First(&t1, 10)

	e := Essay{
		Subject: "一个组合的Save",
		Author:  Author{Name: "马士兵"},
		Tags: []Tag{
			t1,
			{Title: "Ma"},
			{Title: "GORM"},
		},
	}

	if err := DB.Save(&e).Error; err != nil {
		log.Println(err)
	}

	log.Printf("%+v\n", e)

}
```

通过执行多条SQL完成，保证数据的完整性。

## 预加载

预加载，指的是在查询当前模型时，自动查询关联的模型，使用方法：

```go
.Preload("Association", conds)
```

来实现。

支持多次链式调用，来预加载多个关联。

支持指定关联查询条件。

示例：

```go
// Preload
func AssocPreload() {
	// A.直接一步查询Author对应的Essays
	a := Author{}
	if err := DB.
		Preload("Essays").
		First(&a, 1).Error; err != nil {
		log.Fatalln(err)
	}
	// [3.840ms] [rows:2] SELECT * FROM `msb_essay` WHERE `msb_essay`.`author_id` = 1 AND `msb_essay`.`deleted_at` IS NULL
	// [13.014ms] [rows:1] SELECT * FROM `msb_author` WHERE `msb_author`.`id` = 1 AND `msb_author`.`deleted_at` IS NULL ORDER BY `msb_author`.`id` LIMIT 1
	log.Println(a.Essays)
	log.Println("--------------------")

	// B.支持条件过滤
	if err := DB.
		Preload("Essays", "id IN ?", []uint{2, 3, 4}).
		First(&a, 1).Error; err != nil {
		log.Fatalln(err)
	}
	// [3.217ms] [rows:1] SELECT * FROM `msb_essay` WHERE `msb_essay`.`author_id` = 1 AND id IN (2,3,4) AND `msb_essay`.`deleted_at` IS NULL
	log.Println(a.Essays)
	log.Println("-----------------------")

	// C. 支持多次链式调用，同时预加载多个关联
	e := Essay{}
	if err := DB.
		Preload("Author").
		Preload("EssayMate").
		Preload("Tags").
		First(&e, 1).Error; err != nil {
		log.Fatalln(err)
	}
	log.Println(e)
	// [2.776ms] [rows:1] SELECT * FROM `msb_author` WHERE `msb_author`.`id` = 1 AND `msb_author`.`deleted_at` IS NULL
	// [10.398ms] [rows:0] SELECT * FROM `msb_essay_mate` WHERE `msb_essay_mate`.`essay_id` = 1 AND `msb_essay_mate`.`deleted_at` IS NULL
	// [3.260ms] [rows:2] SELECT * FROM `msb_essay_tag` WHERE `msb_essay_tag`.`essay_id` = 1
	// [3.264ms] [rows:2] SELECT * FROM `msb_tag` WHERE `msb_tag`.`id` IN (1,3) AND `msb_tag`.`deleted_at` IS NULL
	// [28.067ms] [rows:1] SELECT * FROM `msb_essay` WHERE `msb_essay`.`id` = 1 AND `msb_essay`.`deleted_at` IS NULL ORDER BY `msb_essay`.`id` LIMIT 1
}
```

### 多级预加载

.Preload()的参数支持层级语法：

```go
.Preload("Association1.Association2.Assocaition3", conds)
```

默认情况下，GORM仅仅会加载一级的关联。使用多级语法，可以预加载多级关联数据。

示例：

```go
// 多级
func AssocLevelPreload() {
	a := Author{}
	if err := DB.
		//Preload("Essays").
		// 多级关联
		Preload("Essays.Tags").
		First(&a, 1).Error; err != nil {
		log.Fatalln(err)
	}
	// [3.843ms] [rows:5] SELECT * FROM `msb_essay_tag` WHERE `msb_essay_tag`.`essay_id` IN (1,2)
	// [3.284ms] [rows:3] SELECT * FROM `msb_tag` WHERE `msb_tag`.`id` IN (1,3,2) AND `msb_tag`.`deleted_at` IS NULL
	// [10.396ms] [rows:2] SELECT * FROM `msb_essay` WHERE `msb_essay`.`author_id` = 1 AND `msb_essay`.`deleted_at` IS NULL
	// [17.609ms] [rows:1] SELECT * FROM `msb_author` WHERE `msb_author`.`id` = 1 AND `msb_author`.`deleted_at` IS NULL ORDER BY `msb_author`.`id` LIMIT 1

	log.Println(a.Essays[0].Tags)
	log.Println(a.Essays[1].Tags)
}

```

### 预加载全部

若需要全部的关联都预加载，除了链式调用全部的关联之外，还可以使用子句：

```go
.Preload(clause.Associations）
```

`clause.Associations`不会预加载层级的关联，可以配合多级预加载一起使用。

示例：

```go
e := Essay{}
	if err := DB.
		//Preload("Author").
		//Preload("EssayMate").
		//Preload("Tags").
		Preload(clause.Associations).
		First(&e, 1).Error; err != nil {
		log.Fatalln(err)
	}
	log.Println(e)
```

## 自定义外键关联属性

当未使用标准的字段进行关联时，需要对关联属性进行设置。

推荐尽量采用标准的模型定义。

典型的需要自定义的情况：

- 复合主键
- 数据库结构已定
- 多个关联，例如，Essay和Author关联了多次，有第一作者，校订作者，通讯作者等。

### 外键字段

使用gorm标签：foreignKey来自定义外键字段，要求与关联字段类型一致。

### 引用字段

使用gorm标签：references来自定义引用字段，要求与外键字段类型一致。

### 约束操作

使用gorm标签：constraint来自定义约束操作：

- OnUpdate
  - CASCADE，级联更新
  - SET NULL，外键设置NULL
  - RESTRICT，限制更新
- OnDelete
  - CASCADE，级联删除
  - SET NULL，外键设为NULL
  - RESTRICT，限制删除，默认

示例定义外键字段：

```go
// 作者模型
// Author模型
type Author struct {
	gorm.Model
	Status int
	Name   string
	Email  string

	// 拥有多个论文内容
	// has many

	// 默认关联
	Essays []Essay

	// 第一作者论文列表
	FirstEssays []Essay `gorm:"foreignKey:FirstAuthorID;references:;"`
	// 第二作者论文列表
	SecondEssays []Essay `gorm:"foreignKey:SecondAuthorID;references:;"`
}

// CREATE TABLE `msb_essay` (
//  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
//  `first_author_id` bigint unsigned DEFAULT NULL,
//  `second_author_id` bigint unsigned DEFAULT NULL,
//  `author_id` bigint unsigned DEFAULT NULL,
//  PRIMARY KEY (`id`),
//  KEY `idx_msb_essay_deleted_at` (`deleted_at`),
//  KEY `fk_msb_author_essays` (`author_id`),
//  KEY `fk_msb_author_first_essays` (`first_author_id`),
//  KEY `fk_msb_author_second_essays` (`second_author_id`),
//  CONSTRAINT `fk_msb_author_essays` FOREIGN KEY (`author_id`) REFERENCES `msb_author` (`id`),
//  CONSTRAINT `fk_msb_author_first_essays` FOREIGN KEY (`first_author_id`) REFERENCES `msb_author` (`id`),
//  CONSTRAINT `fk_msb_author_second_essays` FOREIGN KEY (`second_author_id`) REFERENCES `msb_author` (`id`)
//) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci

// 论文内容
type Essay struct {
	gorm.Model
	Subject string
	Content string

	// 自定义关联字段
	FirstAuthorID  uint
	SecondAuthorID uint

	FirstAuthor  Author `gorm:"foreignKey:FirstAuthorID;references:;"`
	SecondAuthor Author `gorm:"foreignKey:SecondAuthorID;references:;"`

	// 外键字段
	AuthorID *uint
	// 属于某个作者
	// belongs to
	Author Author

	// 拥有一个论文元信息
	// has one
	EssayMate EssayMate

	// 拥有多个Tag
	// many to many
	Tags []Tag `gorm:"many2many:essay_tag"`
}
```

注意，表结构，注意不同字段，映射的不同外键约束。

限制约束操作的示例，删除时，外键设为NULL：

```go
// 作者模型
// Author模型
type Author struct {
	gorm.Model
	Status int
	Name   string
	Email  string

	// 拥有多个论文内容
	// has many

	// 默认关联
	Essays []Essay `gorm:"constraint:OnDelete:SET NULL;"`
}
```


# 事务操作

事务，Transaction，指的是一组数据库操作组成的执行单元，要不全部的操作都成功对数据库产生影响，要不全部的操作都不会对数据库产生影响。在数据库系统中用于保证数据的完整性和一致性。其典型特征为：

- 原子性（**A**tomicity）: 事务 `要么全部完成，要么全部取消`。 如果事务崩溃，状态回到事务之前（事务回滚）。
- 隔离性（**I**solation）: 如果2个事务 T1 和 T2 同时运行，事务 T1 和 T2 最终的结果是相同的，不管 T1和T2谁先结束。
- 持久性（**D**urability）: 一旦事务提交，不管发生什么（崩溃或者出错），数据要保存在数据库中。
- 一致性（**C**onsistency）: 只有合法的数据（依照关系约束和函数约束）才能写入数据库。

## 事务方法

Gorm支持如下方法关联事务：

```go
// 开始事务
tx := DB.Begin()
// 回滚事务
tx.Rollback()
// 提交事务
tx.Commit()
```

注意，DB.Begin()方法返回开始事务的数据库对象，后续的本事务操作应该基于该对象完成，包括数据的CRUD等。

示例：

```go
type Author struct {
	gorm.Model
	Name   string
	// 积分
	Points int
}
```

```go
func TXDemo() {
	// 初始化测试数据
	if err := DB.AutoMigrate(&Author{}); err != nil {
		log.Fatalln(err)
	}
	var a1, a2 Author
	a1.Name = "库里"
	a2.Name = "莫兰特"
	a1.Points = 1600
	a2.Points = 200
	if err := DB.Create([]*Author{&a1, &a2}).Error; err != nil {
		log.Fatalln(err)
	}

	// 事务操作
	// a1 赠送 a2 2000 积分
	p := 2000
	// 开始事务
	tx := DB.Begin()
	// 有时需要考虑数据库是否支持事务的情景
	if tx.Error != nil {
		log.Fatalln(tx.Error)
	}

	// 执行赠送操作
	a1.Points -= p
	a2.Points += p

	// 1执行SQL，可能导致的错误
	if err := tx.Save(&a1).Error; err != nil {
		tx.Rollback()
		return
	}

	if err := tx.Save(&a2).Error; err != nil {
		// 回滚事务
		tx.Rollback()
		return
	}

	// 2业务逻辑可能导致的错误
	// 要求author的积分不能为负数
	if a1.Points < 0 || a2.Points < 0 {
		log.Println("a1.Points < 0 || a2.Points < 0")
		// 回滚事务
		if err := tx.Rollback().Error; err != nil {
			log.Fatalln(err)
		}
	}

	// 提交事务
	if err := tx.Commit().Error; err != nil {
		log.Fatalln(err)
	}

	// 决定回滚还是提交，集中的处理错误风格
	//if err1 != nil || err2 != nil {
	//	tx.Rollback()
	//} else {
	//	tx.Commit()
	//}
}

```

测试，本例中，会导致回滚。

```shell
> go test -run TXDemo
2023/04/27 11:10:10 a1.Points < 0 || a2.Points < 0
2023/04/27 11:10:10 sql: transaction has already been committed or rolled back
exit status 1
FAIL    gormExample     0.147s
```

## 回调函数方式操作事务

除了手动调用 .Begin() tx.Commit() tx.Rollback() 外，GORM提供了一个回调函数的方案来执行事务，自动的完成开启事务和提交或回滚事务：

```go
func (db *DB) Transaction(fc func(tx *DB) error, opts ...*sql.TxOptions) (err error)
```

将业务逻辑代码直接由 func(tx *DB) error 函数实现即可。

同时 Transaction 也会得到具体错误。

推荐优先使用该方法完成事务。原因有2：

1. 编码重点在业务逻辑，不用关心事务的处理
2. 支持嵌套事务

示例：

```go
func TXCallback() {
	// 初始化测试数据
	if err := DB.AutoMigrate(&Author{}); err != nil {
		log.Fatalln(err)
	}
	var a1, a2 Author
	a1.Name = "库里"
	a2.Name = "莫兰特"
	a1.Points = 1600
	a2.Points = 200
	if err := DB.Create([]*Author{&a1, &a2}).Error; err != nil {
		log.Fatalln(err)
	}
	log.Println(a1.ID, a2.ID)

	// 实现事务
	if err := DB.Transaction(func(tx *gorm.DB) error {
		// a1 赠送 a2 2000 积分
		p := 200
		// 执行赠送操作
		a1.Points -= p
		a2.Points += p

		// 1执行SQL，可能导致的错误
		if err := tx.Save(&a1).Error; err != nil {
			return err
		}

		if err := tx.Save(&a2).Error; err != nil {
			return err
		}

		// 2业务逻辑可能导致的错误
		// 要求author的积分不能为负数
		if a1.Points < 0 || a2.Points < 0 {
			return errors.New("a1.Points < 0 || a2.Points < 0")
		}

		// nil 的返回，会导致事务提交
		return nil
	}); err != nil {
		// 返回错误，为了后续的业务逻辑处理
		// 为了通知我们，事务成功还是失败
		// 返回错误，不影响事务的提交和回滚
		log.Println(err)
	}
}
```

Transaction方法还支持嵌套调用，用于支持嵌套事务，示例：

嵌套事务，在实操中，主要用于实现，或的逻辑，例如：a1 转给 a2 2000 积分，若失败，可以a3转给a2。此时，仅需将a1转的事务回滚即可。

示例：

```go
func TXNested() {
	// 初始化测试数据
	if err := DB.AutoMigrate(&Author{}); err != nil {
		log.Fatalln(err)
	}
	var a1, a2, a3 Author
	a1.Name = "库里"
	a2.Name = "莫兰特"
	a3.Name = "欧文"
	a1.Points = 1600
	a2.Points = 200
	a3.Points = 4000
	if err := DB.Create([]*Author{&a1, &a2, &a3}).Error; err != nil {
		log.Fatalln(err)
	}
	log.Println(a1.ID, a2.ID, a3.ID)

	// 实现事务
	if err := DB.Transaction(func(tx *gorm.DB) error {
		// a1 赠送 a2 2000 积分
		p := 20000

		// 执行赠送操作

		// a2 多了积分
		a2.Points += p
		if err := tx.Save(&a2).Error; err != nil {
			return err
		}

		// a1 赠送，使用嵌套事务完成
		errA1 := tx.Transaction(func(tx *gorm.DB) error {
			a1.Points -= p
			// 1执行SQL，可能导致的错误
			if err := tx.Save(&a1).Error; err != nil {
				return err
			}
			if a1.Points < 0 {
				return errors.New("a1.Points < 0")
			}
			// 没有错误成功
			return nil
		})

		// a1 发送失败，才需要a3
		if errA1 != nil {
			// a3 赠送，使用嵌套事务完成
			errA3 := DB.Transaction(func(tx *gorm.DB) error {
				a3.Points -= p
				if err := tx.Save(&a3).Error; err != nil {
					return err
				}
				if a3.Points < 0 {
					return errors.New("a3.Points < 0")
				}
				return nil
			})
			// a3 同样失败
			if errA3 != nil {
				return errors.New("a1 and a3 all send points failed")
			}
		}

		// nil 的返回，会导致事务提交
		return nil
	}); err != nil {
		// 返回错误，为了后续的业务逻辑处理
		// 为了通知我们，事务成功还是失败
		// 返回错误，不影响事务的提交和回滚
		log.Println(err)
	}
}
```

## SavePoint

GORM也提供了对事务逻辑存储点，及回到逻辑存储点的支持：

- SavePoint，定义SavePoint
- Rollbackto，回到SavePoint

示例，实现相同的a1 和 a3 给Points到a2的逻辑：

```go
func TXSavePoint() {
	// 初始化测试数据
	if err := DB.AutoMigrate(&Author{}); err != nil {
		log.Fatalln(err)
	}
	var a1, a2, a3 Author
	a1.Name = "库里"
	a2.Name = "莫兰特"
	a3.Name = "欧文"
	a1.Points = 1600
	a2.Points = 200
	a3.Points = 4000
	if err := DB.Create([]*Author{&a1, &a2, &a3}).Error; err != nil {
		log.Fatalln(err)
	}
	log.Println(a1.ID, a2.ID, a3.ID)

	// 事务操作
	// a1 赠送 a2 2000 积分
	p := 20000
	// 开始事务
	tx := DB.Begin()
	// 有时需要考虑数据库是否支持事务的情景
	if tx.Error != nil {
		log.Fatalln(tx.Error)
	}

	// 执行赠送操作
	// a2 得到积分
	a2.Points += p
	// 1执行SQL，可能导致的错误
	if err := tx.Save(&a2).Error; err != nil {
		tx.Rollback()
		return
	}

	// 逻辑记录发送points是否成功
	var flagSend bool

	// a1 先给 a2 send
	// 设置一个 savepoint
	tx.SavePoint("beforeA1")
	a1.Points -= p
	if err := tx.Save(&a1).Error; err != nil || a1.Points < 0 {
		// 回滚到 beforeA1
		tx.RollbackTo("beforeA1")

		// a3 to a2
		tx.SavePoint("beforeA3")
		a3.Points -= p
		if err := tx.Save(&a3).Error; err != nil || a3.Points < 0 {
			// 回滚到 beforeA3
			tx.RollbackTo("beforeA3")
		} else {
			flagSend = true
		}
	} else {
		flagSend = true
	}

	// 判定发送是否成功
	if flagSend {
		// 提交事务
		if err := tx.Commit().Error; err != nil {
			log.Fatalln(err)
		}
	} else {
		// 回滚事务
		tx.Rollback()
	}
}
```

## 禁用默认事务

了确保数据一致性，GORM 会在事务里执行写入操作（创建、更新、删除）。

如果没有这方面的要求，您可以在初始化时禁用它，这将获得大约 30%+ 性能提升。

官网例子：

```go
// 全局禁用
db, err := gorm.Open(sqlite.Open("gorm.db"), &gorm.Config{
  SkipDefaultTransaction: true,
})

// 持续会话模式
tx := db.Session(&Session{SkipDefaultTransaction: true})
tx.First(&user, 1)
tx.Find(&users)
tx.Model(&user).Update("Age", 18)
```
