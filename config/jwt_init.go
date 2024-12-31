package config

import (
	"awesomeProject/model"
	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/gin-gonic/gin"
	"time"
)

// Login 是用于接受登录的用户名与密码
// ” 表示接受form数据时应从username中获取，响应json时应放到username字段，在绑定form和json时这个字段是必须的
type login struct {
	Username string `form:"username" json:"username" binding:"required"`
	Password string `form:"password" json:"password" binding:"required"`
}

func InitJWTMiddleware() (*jwt.GinJWTMiddleware, error) {
	return jwt.New(initParams())
}

func initParams() *jwt.GinJWTMiddleware {
	return &jwt.GinJWTMiddleware{
		Realm:            "test zone",          //标识
		SigningAlgorithm: "HS256",              //加密算法
		Key:              []byte("secret key"), //密钥
		Timeout:          7 * 24 * time.Hour,   //7天过期
		MaxRefresh:       7 * 24 * time.Hour,   //刷新最大延长时间
		IdentityKey:      model.IdentityKey,    //指定cookie的id
		PayloadFunc: func(data interface{}) jwt.MapClaims { //负载，这里可以定义返回jwt中的payload数据
			if v, ok := data.(*model.User); ok {
				return jwt.MapClaims{
					model.IdentityKey: v.UserName,
				}
			}
			return jwt.MapClaims{}
		},
		IdentityHandler: func(c *gin.Context) interface{} {
			claims := jwt.ExtractClaims(c)
			return &model.User{
				UserName: claims[model.IdentityKey].(string),
			}
		},
		Authenticator: Authenticator, //在这里可以写我们的登录验证逻辑
		Authorizator: func(data interface{}, c *gin.Context) bool { //当用户通过token请求受限接口时，会经过这段逻辑
			if v, ok := data.(*model.User); ok && v.UserName == "admin" {
				return true
			}

			return false
		},
		Unauthorized: func(c *gin.Context, code int, message string) { //错误时响应
			c.JSON(code, gin.H{
				"code":    code,
				"message": message,
			})
		},
		// 指定从哪里获取token 其格式为："<source>:<name>" 如有多个，用逗号隔开
		TokenLookup:   "header: Authorization, query: token, cookie: jwt",
		TokenHeadName: "Bearer",
		TimeFunc:      time.Now,
	}
}

func Authenticator(c *gin.Context) (interface{}, error) {
	// 从context中绑定loginVals，并验证用户输入的用户名和密码与服务端存储的是否一致。
	// 一致则返回用户信息，否则返回错误信息
	var loginVals login
	if err := c.ShouldBind(&loginVals); err != nil {
		return "", jwt.ErrMissingLoginValues
	}
	userID := loginVals.Username
	password := loginVals.Password

	if (userID == "admin" && password == "admin") || (userID == "test" && password == "test") {
		return &model.User{
			UserName:  userID,
			LastName:  "Wang",
			FirstName: "Ruijie",
		}, nil
	}

	return nil, jwt.ErrFailedAuthentication
}
