package route

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"opsw/utils"
	"os"
	"strings"
)

func Auth(c *gin.Context) {
	urlPath := c.Request.URL.Path
	// 静态资源
	if strings.HasPrefix(urlPath, "/assets") {
		c.File(utils.WorkDir(fmt.Sprintf("/resources/web/dist%s", urlPath)))
		return
	}
	// 退出登录
	if strings.HasPrefix(urlPath, "/oauth/logout") {
		userToken := utils.GinGetCookie(c, "user_token")
		if userToken != "" {
			apiFile := utils.WorkDir(fmt.Sprintf("/users/%s", userToken))
			if utils.IsFile(apiFile) {
				_ = os.Remove(apiFile)
			}
		}
		utils.GinRemoveCookie(c, "user_token")
		utils.GinResult(c, http.StatusOK, "退出成功")
		return
	}
	// 页面输出
	c.HTML(http.StatusOK, "/resources/web/dist/index.html", gin.H{
		"CODE": "",
		"MSG":  "",
	})
}
