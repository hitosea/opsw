package routes

import (
	"fmt"
	"net/http"
	"opsw/utils"
	"os"
)

// NoAuthApiUserLogin 登录
func (app *AppStruct) NoAuthApiUserLogin() {
	//email := app.Context.Query("email")
	//password := app.Context.Query("password")
	utils.GinResult(app.Context, http.StatusOK, "登录成功")
}

// NoAuthApiUserReg 注册
func (app *AppStruct) NoAuthApiUserReg() {
	//email := app.Context.Query("email")
	//password := app.Context.Query("password")
	//password2 := app.Context.Query("password2")
}

// NoAuthApiUserLogout 退出
func (app *AppStruct) NoAuthApiUserLogout() {
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
