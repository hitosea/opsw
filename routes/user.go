package routes

import (
	"fmt"
	"net/http"
	"opsw/database"
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
	var (
		email     = utils.GinInput(app.Context, "email")
		password  = utils.GinInput(app.Context, "password")
		password2 = utils.GinInput(app.Context, "password2")
	)
	if !utils.IsEmail(email) {
		utils.GinResult(app.Context, http.StatusBadRequest, "邮箱格式不正确")
		return
	}
	if password != password2 {
		utils.GinResult(app.Context, http.StatusBadRequest, "两次密码不一致")
		return
	}
	user, err := database.CreateUser(email, "", password)
	if err != nil {
		utils.GinResult(app.Context, http.StatusBadRequest, fmt.Sprintf("注册失败：%s", err.Error()))
		return
	}
	utils.GinResult(app.Context, http.StatusOK, "注册成功", user)
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
