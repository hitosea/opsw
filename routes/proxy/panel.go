package proxy

import (
	"github.com/gin-gonic/gin"
	"net/http/httputil"
	"net/url"
	"opsw/database"
)

type AppStruct struct {
	Context    *gin.Context
	UserInfo   *database.User
	ServerInfo *database.Server
}

func (app *AppStruct) Panel() {
	c := app.Context

	target := "http://192.168.64.44:22636"
	proxyUrl, _ := url.Parse(target)

	c.Request.Header.Set("X-Real-Ip", c.ClientIP())
	c.Request.Header.Set("X-Forwarded-For", c.ClientIP())
	c.Request.URL.Path = "/ee715a9914"

	proxy := httputil.NewSingleHostReverseProxy(proxyUrl)
	proxy.ServeHTTP(c.Writer, c.Request)
}
