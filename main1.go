package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/github"
)

// 配置 GitHub OAuth 客户端信息
var githubOAuthConfig = &oauth2.Config{
	ClientID:     "Ov23liW6umAHoi7wbATz",                     // GitHub 应用的 Client ID
	ClientSecret: "735521fcf0c9d9a8ec1c351f3bf755c8a5a2167b", // GitHub 应用的 Client Secret
	RedirectURL:  "http://localhost:8080/oauth/redirect",     // 授权回调地址
	Scopes:       []string{"read:user", "user:email"},        // 需要的权限
	Endpoint:     github.Endpoint,
}

var state = "random-generated-state" // 防止 CSRF 攻击的随机字符串

func main() {
	// 初始化 Gin 引擎
	r := gin.Default()

	// 定义路由
	r.GET("/login", loginHandler)             // 登录路由
	r.GET("/oauth/redirect", callbackHandler) // GitHub 授权回调路由

	// 启动服务器
	log.Println("Server is running at http://localhost:8080...")
	if err := r.Run("0.0.0.0:8080"); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}

// Step 1: 用户请求登录，跳转到 GitHub 授权页面
func loginHandler(c *gin.Context) {
	fmt.Println(111)
	// 生成 GitHub OAuth 授权 URL
	url := githubOAuthConfig.AuthCodeURL(state)
	fmt.Println(22)
	log.Printf("Generated GitHub OAuth URL: %s", url) // 打印 URL 以便调试

	// 重定向到 GitHub 授权页面
	c.Redirect(http.StatusFound, url)
}

// Step 2: 授权回调处理
func callbackHandler(c *gin.Context) {
	// 验证 state 参数
	if c.Query("state") != state {
		c.JSON(http.StatusBadRequest, gin.H{"error": "State is invalid"})
		return
	}

	// 获取授权码
	code := c.Query("code")
	if code == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Code not found"})
		return
	}

	// 用授权码换取访问令牌
	token, err := githubOAuthConfig.Exchange(c, code)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to exchange token", "details": err.Error()})
		return
	}

	// 使用访问令牌获取用户信息
	client := githubOAuthConfig.Client(c, token)
	resp, err := client.Get("https://api.github.com/user")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get user info", "details": err.Error()})
		return
	}
	defer resp.Body.Close()

	// 解析用户信息
	var userInfo map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&userInfo); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to decode user info", "details": err.Error()})
		return
	}

	// 返回用户信息
	c.JSON(http.StatusOK, userInfo)
}
