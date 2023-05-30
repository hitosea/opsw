package proxy

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"net/http/httputil"
	"net/url"
	"opsw/database"
	"opsw/utils"
)

func ManagePanel(c *gin.Context, userInfo *database.User) {
	ip := c.Query("ip")
	if ip == "" {
		ip, _ = c.Cookie("x-panel-ip")
	}
	if ip == "" {
		utils.GinResult(c, http.StatusBadRequest, "缺少参数 ip")
		return
	}
	server, err := database.ServerGet(map[string]any{
		"ip": ip,
	}, userInfo.Id, false)
	if err != nil {
		utils.GinResult(c, http.StatusBadRequest, "无法访问", gin.H{"error": err.Error()})
		return
	}
	/*server := database.Server{
		Ip:            "192.168.100.22",
		PanelPort:     9999,
		PanelUsername: "admin",
		PanelPassword: "admin123",
	}*/

	// 地址
	target := fmt.Sprintf("http://%s:%d", server.Ip, server.PanelPort)
	proxyUrl, _ := url.Parse(target)

	// Cookie
	c.SetCookie("x-panel-ip", ip, 86400*7, "", "", false, false)
	c.SetCookie("x-panel-frame", "opsw", 86400*7, "", "", false, false)

	// Header
	c.Request.Header.Set("X-Real-Ip", c.ClientIP())
	c.Request.Header.Set("X-Forwarded-For", c.ClientIP())
	c.Request.Header.Set("X-Panel-Username", server.PanelUsername)
	c.Request.Header.Set("X-Panel-Password", server.PanelPassword)

	// 代理请求
	proxy := httputil.NewSingleHostReverseProxy(proxyUrl)
	proxy.ServeHTTP(c.Writer, c.Request)
}
