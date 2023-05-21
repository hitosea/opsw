package utils

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"net/url"
	"strings"
)

// GinToken Gin获取Token（Header、Query、Cookie）
func GinToken(c *gin.Context) string {
	token := c.GetHeader("token")
	if token == "" {
		token = GinInput(c, "token")
	}
	if token == "" {
		token = GinGetCookie(c, "token")
	}
	return token
}

// GinInput Gin获取参数（优先POST、取Query）
func GinInput(c *gin.Context, key string) string {
	if c.PostForm(key) != "" {
		return strings.TrimSpace(c.PostForm(key))
	}
	return strings.TrimSpace(c.Query(key))
}

// GinScheme Gin获取Scheme
func GinScheme(c *gin.Context) string {
	scheme := "http://"
	if c.Request.TLS != nil || c.Request.Header.Get("X-Forwarded-Proto") == "https" {
		scheme = "https://"
	}
	return scheme
}

// GinHomeUrl Gin获取HomeUrl
func GinHomeUrl(c *gin.Context) string {
	return fmt.Sprintf("%s%s", GinScheme(c), c.Request.Host)
}

// GinGetCookie Gin获取Cookie
func GinGetCookie(c *gin.Context, name string) string {
	value, _ := c.Cookie(name)
	return value
}

// GinSetCookie Gin设置Cookie
func GinSetCookie(c *gin.Context, name, value string, maxAge int) {
	c.SetCookie(name, value, maxAge, "/", "", false, false)
}

// GinRemoveCookie Gin删除Cookie
func GinRemoveCookie(c *gin.Context, name string) {
	c.SetCookie(name, "", -1, "/", "", false, false)
}

// GinResult 返回结果
func GinResult(c *gin.Context, code int, content string, values ...any) {
	c.Header("Expires", "-1")
	c.Header("Cache-Control", "no-cache")
	c.Header("Pragma", "no-cache")
	var data any
	if len(values) == 1 {
		data = values[0]
	} else if len(values) == 0 {
		data = gin.H{}
	} else {
		data = values
	}
	//
	if strings.Contains(c.GetHeader("Accept"), "application/json") {
		// 接口返回
		c.JSON(http.StatusOK, gin.H{
			"code": code,
			"msg":  content,
			"data": data,
		})
	} else {
		// 页面返回
		if code == http.StatusMovedPermanently {
			c.Redirect(code, content)
		} else {
			c.HTML(http.StatusOK, "/web/dist/index.html", gin.H{
				"CODE": code,
				"MSG":  url.QueryEscape(content),
			})
		}
	}
}
