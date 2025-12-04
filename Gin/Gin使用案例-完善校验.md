# 数据校验概述

数据校验的三个层面：

1. 前端（web、app），提升用户的体验，不需要等到服务器端返回结果，前端即可了解数据是否合理。（前端开发做的）
2. **应用层面，当后端接口接收到前端数据时，对数据的合理性（合法性）进行校验。（我们后端开发要做的）**
3. 数据存储层面，当数据入库时，必须要保证数据约束的满足。例如唯一、类型、null、外键关联约束。（数据库维护者）

框架中，通常继承数据验证的工具。

gin中，在binding时完成验证。Gin使用 go-playground/validator/v10 进行验证。

验证的操作：

1. 设置验证规则
2. 自定义验证规则
3. 定义验证消息

Gin的验证的流程：

1. 解析数据，将请求数据解析到结构体字段中
2. 规则的校验

意味着，如果请求数据不能够正确的解析到结构体字段中，后边的验证也无从说起。

![image.png](https://fynotefile.oss-cn-zhangjiakou.aliyuncs.com/fynote/fyfile/13080/1722342738052/6f0128cac7e8466992ec3851f3fc97c2.png)

# 字段不重复的校验

自定义验证器。

validator 提供的验证规则，不能满足需求。

## 定义验证器

验证其函数的签名：

```go
func (validator.FieldLevel) bool
```

定义在具体handler的message中。

handlers/role/message.go

```go
func roleTitleUnique(fieldLevel validator.FieldLevel) bool {
	// title的值
	value := fieldLevel.Field().Interface().(string)

	// 校验是否重复
	row := models.Role{}
	utils.DB().Where("`title` = ?", value).Unscoped().First(&row)

	// 判断是否查询到了
	return row.ID == 0
}
```

## 注册该验证器

在Bind操作前完成注册。

选择在handlers/role包的init函数中完成：

```go
func init() {
	// 注册本包逻辑的验证器
	registerValidator()
}
func registerValidator() {
	if validate, ok := binding.Validator.Engine().(*validator.Validate); ok {
		// 注册
		_ = validate.RegisterValidation("roleTitleUnique", roleTitleUnique)
	}
}

```

## 使用验证器

在struct的tag中完成使用：

handlers/role/message.go

以添加为例：

```go
// 添加请求消息
type AddReq struct {
	models.Role
	// 需要额外校验的字段
	Title string `json:"title" binding:"required,roleTitleUnique"`
	Key   string `json:"key" binding:"required"`
}
```

## 测试

![image.png](https://fynotefile.oss-cn-zhangjiakou.aliyuncs.com/fynote/fyfile/13080/1722342738052/ad1c460fec1b4b9da36296f9a0a235f0.png)

通过！

# 添加和更新共用验证器

角色更新时的title验证， 也是需要不能重复。

handlers/role/message.go

```go
// EditBodyReq 更新主体参数
type EditBodyReq struct {
	Title   *string `json:"title" field:"title" binding:"omitempty,roleTitleUnique"`
	Key     *string `json:"key" field:"key"`
	Enabled *bool   `json:"enabled" field:"enabled"`
	Weight  *int    `json:"weight" field:"weight"`
	Comment *string `json:"comment" field:"comment"`
}
```

问题在于，编辑时，验证不重复的条件，应该是和除了自己外的其他记录比较。

## 数据上提供当前ID

因此需要指导当前编辑的是哪条记录。

增加请求body的ID字段。

handlers/role/message.go

```go
// EditBodyReq 更新主体参数
type EditBodyReq struct {
	ID      uint
	Title   *string `json:"title" field:"title" binding:"omitempty,roleTitleUnique"`
	Key     *string `json:"key" field:"key"`
	Enabled *bool   `json:"enabled" field:"enabled"`
	Weight  *int    `json:"weight" field:"weight"`
	Comment *string `json:"comment" field:"comment"`
}
```

初始化EditBodyReq时，提供ID

handlers/role/handlers.go Edit()

```go
// 2. 解析Body请求数据
	body := EditBodyReq{
		ID: uri.ID,
	}
```

## 验证器加入ID条件

条件就变成了，title == ? && id != ?

handlers/role/message.go

```go
func roleTitleUnique(fieldLevel validator.FieldLevel) bool {
	// title的值
	value := fieldLevel.Field().Interface().(string)
	// id的值
	id := fieldLevel.Parent().FieldByName("ID").Interface().(uint)

	// 校验是否重复
	row := models.Role{}
	utils.DB().Where("`title` = ? && `id` != ?", value, id).Unscoped().First(&row)

	// 判断是否查询到了
	return row.ID == 0
}
```

## 需要的位置使用校验

handlers/role/message.go

```go
// EditBodyReq 更新主体参数
type EditBodyReq struct {
	Title   *string `json:"title" field:"title" binding:"omitempty,roleTitleUnique"`
	Key     *string `json:"key" field:"key"`
	Enabled *bool   `json:"enabled" field:"enabled"`
	Weight  *int    `json:"weight" field:"weight"`
	Comment *string `json:"comment" field:"comment"`
}
```

## 测试

![image.png](https://fynotefile.oss-cn-zhangjiakou.aliyuncs.com/fynote/fyfile/13080/1722342738052/1ef02b7d690b4c5e99dac9e4162833ed.png)

通过！

# 验证消息翻译

* 内置规则的消息
* 自定义规则的消息

## 内置规则消息翻译

需要：

1. 翻译器
2. 翻译内容，需要哪种语言？

安装翻译器：

```bash
go get github.com/go-playground/universal-translator
go get github.com/go-playground/validator/v10/translations
```

## 翻译内置消息

通用的操作，定义在common

handlers/common/translation.go

```go
var translator ut.Translator

func translateMessage() {
	// 通用的翻译器
	universalTranslator := ut.New(zh.New())
	// 具体验证引擎
	validate := binding.Validator.Engine().(*validator.Validate)
	// 具体的翻译器
	translator, _ = universalTranslator.GetTranslator("zh")
	// 注册为默认的翻译器
	if err := zhTranslations.RegisterDefaultTranslations(validate, translator); err != nil {
		utils.Logger().Warn(err.Error())
	}
}

func init() {
	// 翻译消息
	translateMessage()
}
```

## 定义翻译的方法

在消息响应时，具体翻译某个消息。

handlers/common/translation.go

```go
func Translate(err error) gin.H {
	// 仅翻译验证消息
	errs, ok := err.(validator.ValidationErrors)
	if !ok {
		return nil
	}

	// 翻译
	msg := gin.H{}
	for _, err := range errs {
		msg[err.Field()] = err.Translate(translator)
	}
	return msg
}
```

## 效果

在响应时使用：

handler.go

```go
func Add(ctx *gin.Context) {
	// 1. 解析请求数据
	req := AddReq{}
	if err := ctx.ShouldBind(&req); err != nil {
		// 记录日志
		utils.Logger().Error(err.Error())
		// 直接响应
		ctx.JSON(http.StatusOK, gin.H{
			"code":         100,
			"message":      err.Error(),
			"transMessage": common.Translate(err),
		})
		return
	}
}
```

![image.png](https://fynotefile.oss-cn-zhangjiakou.aliyuncs.com/fynote/fyfile/13080/1722342738052/c7bc61bf0fc64f93867e3263e8a05a41.png)

## 自定义TagName解析

目的：将错误消息中的字段名于json中的名字完全匹配。

为验证器，注册一个自定义方法，用来完成名字的解析：

handlers/common/translation.go

```go
func translateMessage() {
	// 通用的翻译器
	universalTranslator := ut.New(zh.New())
	// 具体验证引擎
	validate := binding.Validator.Engine().(*validator.Validate)
	// 具体的翻译器
	translator, _ = universalTranslator.GetTranslator("zh")
	// 注册为默认的翻译器
	if err := zhTranslations.RegisterDefaultTranslations(validate, translator); err != nil {
		utils.Logger().Warn(err.Error())
	}

	// 注册TagName的自定义函数
	// 从json这个tag中获取
	validate.RegisterTagNameFunc(func(field reflect.StructField) string {
		return field.Tag.Get("json")
	})
}
```

## 测试

![image.png](https://fynotefile.oss-cn-zhangjiakou.aliyuncs.com/fynote/fyfile/13080/1722342738052/d81b72f355014046a9f208628da4f662.png)

通过！

# 自定义消息翻译

## 定义消息映射

handlers/common/translation.go

```go
// 自定义错误消息
var customMsg = map[string]string{
	"roleTitleUnique": "{0}对应的角色已经存在",
}
```

## 注册错误消息

handlers/common/translation.go

```go
func translateMessage() {
	// 略

	// 注册消息
	translateFn := func(ut ut.Translator, fe validator.FieldError) string {
		msg, err := ut.T(fe.Tag(), fe.Field())
		if err != nil {
			utils.Logger().Warn(err.Error())
			return ""
		}
		return msg
	}
	// 遍历全部的错误消息
	for tag, text := range customMsg {
		// 注册
		// 注册函数，和翻译函数
		if err := validate.RegisterTranslation(tag, translator, func(ut ut.Translator) error {
			if err := ut.Add(tag, text, false); err != nil {
				return err
			}
			return nil
		}, translateFn); err != nil {
			utils.Logger().Warn(err.Error())
		}
	}
}
```

## 测试

![image.png](https://fynotefile.oss-cn-zhangjiakou.aliyuncs.com/fynote/fyfile/13080/1722342738052/118d26d99c724473a7ad0e7b86a7883f.png)

通过！
