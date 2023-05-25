package proxy

import (
	"bytes"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/nahid/gohttp"
	"net/http"
	"net/http/httputil"
	"net/url"
	"opsw/database"
	"opsw/utils"
	"text/template"
)

const loginSuccessPageTemplate = `
<!DOCTYPE html>
<html>
<head>
    <meta charset="UTF-8"/>
    <title>Loading...</title>
    <script>
        (function () {
            let globalState;
            try {
                globalState = JSON.parse(window.localStorage.getItem("GlobalState"));
            } catch (e) {
                globalState = {};
            }
			if (typeof globalState.themeConfig !== "object") {
				globalState.themeConfig = {};
			}
            globalState.themeConfig.panelName = "管理面板 [{{ .ServerIp }}]";
            globalState.themeConfig.theme = "dark";
            globalState.isLogin = true;
            window.localStorage.setItem("GlobalState", JSON.stringify(globalState));
            //
            window.location.href = "{{ .RedirectUrl }}"
        })()
    </script>
</head>
<body>
</body>
</html>
`

func ManagePanel(c *gin.Context, userInfo *database.User) {
	ip := c.Query("ip")
	if ip == "" {
		ip, _ = c.Cookie("manage_panel_ip")
	} else {
		c.SetCookie("manage_panel_ip", ip, 604800, "", "", false, false)
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
		PanelPort:     4004,
		PanelUsername: "admin",
		PanelPassword: "admin123",
	}*/
	// 代理地址
	target := fmt.Sprintf("http://%s:%d", server.Ip, server.PanelPort)
	proxyUrl, _ := url.Parse(target)

	// 判断是否登录路由
	if c.Request.URL.Path == "/manage/panel/login" {
		res, err := gohttp.NewRequest().
			JSON(map[string]any{
				"name":          server.PanelUsername,
				"password":      server.PanelPassword,
				"ignoreCaptcha": true,
			}).
			Post(target + "/manage/panel/api/v1/auth/login")
		if err != nil {
			utils.GinResult(c, http.StatusBadRequest, "验证失败失败", gin.H{"error": err.Error()})
			return
		}
		var psession string
		for i := range res.GetResp().Cookies() {
			cc := res.GetResp().Cookies()[i]
			if cc.Name == "psession" {
				psession = cc.Value
			}
		}
		if psession == "" {
			utils.GinResult(c, http.StatusBadRequest, "验证失败失败", gin.H{"error": "psession 为空"})
			return
		}
		c.SetCookie("psession", psession, 604800, "", "", false, false)

		tmpl, err := template.New("text").Parse(loginSuccessPageTemplate)
		defer func() {
			if r := recover(); r != nil {
				utils.GinResult(c, http.StatusBadRequest, "验证模板分析失败", gin.H{"error": fmt.Sprintf("%v", r)})
				return
			}
		}()
		if err != nil {
			utils.GinResult(c, http.StatusBadRequest, "验证模板分析失败", gin.H{"error": err.Error()})
			return
		}
		var buffer bytes.Buffer
		_ = tmpl.Execute(&buffer, gin.H{
			"ServerIp":    server.Ip,
			"RedirectUrl": "/manage/panel/",
		})
		c.Writer.Header().Set("Content-Type", "text/html")
		_, _ = c.Writer.Write(buffer.Bytes())
		return
	}

	// 身份丢失跳转路径
	ifj, err := c.Cookie("identity_failure_jump")
	if err != nil || ifj == "" {
		c.SetCookie("identity_failure_jump", "/manage/panel/login", 604800, "", "", false, false)
	}

	// 代理请求
	c.Request.Header.Set("X-Real-Ip", c.ClientIP())
	c.Request.Header.Set("X-Forwarded-For", c.ClientIP())
	proxy := httputil.NewSingleHostReverseProxy(proxyUrl)
	proxy.ServeHTTP(c.Writer, c.Request)
}
