package main

import (
	"awesomeProject/config"
	"awesomeProject/model"
	"log"
	"net/http"

	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/gin-gonic/gin"
)

func main() {

	// 定义一个Gin的中间件
	authMiddleware, err := config.InitJWTMiddleware()
	if err != nil {
		log.Fatal("JWT Error:" + err.Error())
	}

	// r 创建一个 gin 的路由实例
	r := gin.Default()

	// 登录接口（注册单个路由），并指定处理函数
	r.POST("/login", authMiddleware.LoginHandler)

	// 注册 /auth 路由组
	auth := r.Group("/auth")
	// 退出登录（路由组中的单个路由），并指定处理函数
	auth.POST("/logout", authMiddleware.LogoutHandler)
	// 刷新token，延长token的有效期
	auth.POST("/refresh_token", authMiddleware.RefreshHandler)
	// MiddlewareFunc 是 JWT 中间件的核心函数，它会在请求进入路由处理函数之前进行 JWT 验证。
	auth.GET("/hello", authMiddleware.MiddlewareFunc(), helloHandler)

	if err := http.ListenAndServe(":8080", r); err != nil {
		log.Fatal(err)
	}
}

// 处理/hellow路由的控制器
func helloHandler(c *gin.Context) {
	claims := jwt.ExtractClaims(c)
	user, _ := c.Get(model.IdentityKey)
	c.JSON(200, gin.H{
		"userID":   claims[model.IdentityKey],
		"userName": user.(*model.User).UserName,
		"text":     "Hello World.",
	})
}

//package main
//
//import (
//	"fmt"
//	"github.com/gin-gonic/gin"
//	"github.com/golang-jwt/jwt"
//	"net/http"
//)
//
//type User struct {
//	Name     string `json:"name1"`
//	Email    string `json:"email"`
//	Password string `json:"password"`
//}
//
//func main() {
//	r := gin.Default()
//	r.GET("/query/:id", func(c *gin.Context) {
//		id := c.Param("id")
//		name := c.PostForm("name")
//		a := c.Query("file")
//		fmt.Println(a)
//		c.JSON(http.StatusOK, gin.H{
//			"message": "Hello " + name + a + id,
//		})
//	})
//	r.POST("/login", func(c *gin.Context) {
//		a := User{}
//		fmt.Println(a)
//		if err := c.BindJSON(&a); err != nil {
//			c.JSON(http.StatusBadRequest, gin.H{"status": http.StatusBadRequest, "error": err.Error()})
//			return
//		}
//		fmt.Println(a)
//	})
//	r.Run(":8080")
//}
//
//// 解析结构体
