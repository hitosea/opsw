package routes

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"opsw/utils"
	"os"
)

// UserLogout 退出登录
func UserLogout(c *gin.Context) {
	userToken := utils.GinGetCookie(c, "user_token")
	if userToken != "" {
		apiFile := utils.CacheDir(fmt.Sprintf("/users/%s", userToken))
		if utils.IsFile(apiFile) {
			_ = os.Remove(apiFile)
		}
	}
	utils.GinRemoveCookie(c, "user_token")
	utils.GinResult(c, http.StatusOK, "退出成功")
}
