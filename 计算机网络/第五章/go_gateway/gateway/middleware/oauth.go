package middleware

import (
	"encoding/base64"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"go_gateway/bussiness/mvc/dao"
	"go_gateway/bussiness/util"
	"go_gateway/common"
	"strings"
	"time"
)

type OAuthController struct{}

func OAuthRegister(group *gin.RouterGroup) {
	oauth := &OAuthController{}
	group.POST("/tokens", oauth.Tokens)
}

// Tokens godoc
// @Summary 获取TOKEN
// @Description 获取TOKEN
// @Tags OAUTH
// @ID /oauth/tokens
// @Accept  json
// @Produce  json
// @Param body body dto.TokensInput true "body"
// @Success 200 {object} Response{data=dto.TokensOutput} "success"
// @Router /oauth/tokens [post]
func (oauth *OAuthController) Tokens(c *gin.Context) {
	params := &TokensInput{}
	if err := params.BindValidParam(c); err != nil {
		ResponseError(c, 2000, err)
		return
	}

	splits := strings.Split(c.GetHeader("Authorization"), " ")
	if len(splits) != 2 {
		ResponseError(c, 2001, errors.New("用户名或密码格式错误"))
		return
	}

	appSecret, err := base64.StdEncoding.DecodeString(splits[1])
	if err != nil {
		ResponseError(c, 2002, err)
		return
	}
	//fmt.Println("appSecret", string(appSecret))

	//  取出 app_id secret
	//  生成 app_list
	//  匹配 app_id
	//  基于 jwt生成token
	//  生成 output
	parts := strings.Split(string(appSecret), ":")
	if len(parts) != 2 {
		ResponseError(c, 2003, errors.New("用户名或密码格式错误"))
		return
	}

	appList := dao.AppManagerHandler.GetAppList()
	for _, appInfo := range appList {
		if appInfo.AppID == parts[0] && appInfo.Secret == parts[1] {
			claims := jwt.StandardClaims{
				Issuer:    appInfo.AppID,
				ExpiresAt: time.Now().Add(common.JwtExpires * time.Second).In(common.TimeLocation).Unix(),
			}
			token, err := common.JwtEncode(claims)
			if err != nil {
				ResponseError(c, 2004, err)
				return
			}
			output := &TokensOutput{
				ExpiresIn:   common.JwtExpires,
				TokenType:   "Bearer",
				AccessToken: token,
				Scope:       "read_write",
			}
			ResponseSuccess(c, output)
			return
		}
	}
	ResponseError(c, 2005, errors.New("未匹配正确APP信息"))
}

// AdminLogin godoc
// @Summary 管理员退出
// @Description 管理员退出
// @Tags 管理员接口
// @ID /admin_login/logout
// @Accept  json
// @Produce  json
// @Success 200 {object} Response{data=string} "success"
// @Router /admin_login/logout [get]
func (adminlogin *OAuthController) AdminLoginOut(c *gin.Context) {
	sess := sessions.Default(c)
	sess.Delete(common.AdminSessionInfoKey)
	sess.Save()
	ResponseSuccess(c, "")
}

type TokensInput struct {
	GrantType string `json:"grant_type" form:"grant_type" comment:"授权类型" example:"client_credentials" validate:"required"` //授权类型
	Scope     string `json:"scope" form:"scope" comment:"权限范围" example:"read_write" validate:"required"`                   //权限范围
}

func (param *TokensInput) BindValidParam(c *gin.Context) error {
	return util.DefaultGetValidParams(c, param)
}

type TokensOutput struct {
	AccessToken string `json:"access_token" form:"access_token"` //access_token
	ExpiresIn   int    `json:"expires_in" form:"expires_in"`     //expires_in
	TokenType   string `json:"token_type" form:"token_type"`     //token_type
	Scope       string `json:"scope" form:"scope"`               //scope
}
