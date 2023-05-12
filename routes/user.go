package routes

import (
	"fmt"
	"net/http"
	"opsw/utils"
	"os"
)

// UserLogout 退出登录
func (app *AppStruct) UserLogout() {
	userToken := utils.GinGetCookie(app.Context, "user_token")
	if userToken != "" {
		apiFile := utils.CacheDir(fmt.Sprintf("/users/%s", userToken))
		if utils.IsFile(apiFile) {
			_ = os.Remove(apiFile)
		}
	}
	utils.GinRemoveCookie(app.Context, "user_token")
	utils.GinResult(app.Context, http.StatusOK, "退出成功")
}
