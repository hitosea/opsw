package routes

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"opsw/database"
	"opsw/utils"
)

// NoAuthApiUserLogin 登录
func (app *AppStruct) NoAuthApiUserLogin() {
	var (
		email    = utils.GinInput(app.Context, "email")
		password = utils.GinInput(app.Context, "password")
	)
	if !utils.IsEmail(email) {
		utils.GinResult(app.Context, http.StatusBadRequest, "邮箱格式不正确")
		return
	}
	user, err := database.UserCheck(email, password)
	if err != nil {
		utils.GinResult(app.Context, http.StatusBadRequest, "登录失败", gin.H{"error": err.Error()})
		return
	}
	utils.GinSetCookie(app.Context, "user_token", user.Token, 30*24*86400)
	utils.GinResult(app.Context, http.StatusOK, "登录成功", user)
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
	user, err := database.UserCreate(email, "", password)
	if err != nil {
		utils.GinResult(app.Context, http.StatusBadRequest, "注册失败", gin.H{"error": err.Error()})
		return
	}
	utils.GinSetCookie(app.Context, "user_token", user.Token, 30*24*86400)
	utils.GinResult(app.Context, http.StatusOK, "注册成功", user)
}

// NoAuthApiUserLogout 退出
func (app *AppStruct) NoAuthApiUserLogout() {
	utils.GinRemoveCookie(app.Context, "user_token")
	utils.GinResult(app.Context, http.StatusOK, "退出成功")
}

// AuthApiUserInfo 用户信息
func (app *AppStruct) AuthApiUserInfo() {
	app.UserInfo.Encrypt = ""
	app.UserInfo.Password = ""
	utils.GinResult(app.Context, http.StatusOK, "获取成功", app.UserInfo)
}
