package routes

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"opsw/database"
	"opsw/utils"
	"opsw/vars"
	"strconv"
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
	utils.GinSetCookie(app.Context, "token", user.Token, 30*24*86400)
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
	utils.GinSetCookie(app.Context, "token", user.Token, 30*24*86400)
	utils.GinResult(app.Context, http.StatusOK, "注册成功", user)
}

// NoAuthApiUserLogout 退出
func (app *AppStruct) NoAuthApiUserLogout() {
	utils.GinRemoveCookie(app.Context, "token")
	utils.GinResult(app.Context, http.StatusOK, "退出成功")
}

// AuthApiUserInfo 用户信息
func (app *AppStruct) AuthApiUserInfo() {
	app.UserInfo.Encrypt = ""
	app.UserInfo.Password = ""
	utils.GinResult(app.Context, http.StatusOK, "获取成功", app.UserInfo)
}

// AuthApiUserShareServer 用户分享服务器
func (app *AppStruct) AuthApiUserShareServer() {
	var (
		serverId = utils.GinInput(app.Context, "server_id")
		userIds  = app.Context.PostFormArray("user_ids[]")
	)
	if len(userIds) < 1 {
		utils.GinResult(app.Context, http.StatusBadRequest, "用户参数错误")
		return
	}
	sid, err := strconv.Atoi(serverId)
	if err != nil {
		utils.GinResult(app.Context, http.StatusBadRequest, "服务器参数错误", gin.H{"error": err.Error()})
		return
	}

	db, err := database.InDB(vars.Config.DB)
	if err != nil {
		utils.GinResult(app.Context, http.StatusBadRequest, "数据库连接失败", gin.H{"error": err.Error()})
		return
	}
	defer database.CloseDB(db)

	var serverUser database.ServerUser
	db.Model(&database.ServerUser{}).Where("owner_id = ? and server_id = ?", app.UserInfo.Id, sid).First(&serverUser)
	if serverUser.Id == 0 {
		utils.GinResult(app.Context, http.StatusBadRequest, "服务器不存在或已被删除")
		return
	}

	for _, userId := range userIds {
		var user database.User
		db.Model(&database.User{}).Where("id = ?", userId).First(&user)
		if user.Id == 0 {
			utils.GinResult(app.Context, http.StatusBadRequest, "用户不存在或已被删除")
			return
		}
		if user.Id == serverUser.UserId {
			continue
		}
		db.Model(&database.ServerUser{}).Where("user_id = ? and server_id = ?", user.Id, serverUser.ServerId).FirstOrCreate(&database.ServerUser{
			ServerId: serverUser.ServerId,
			UserId:   user.Id,
			OwnerId:  app.UserInfo.Id,
		})
	}

	utils.GinResult(app.Context, http.StatusOK, "分享成功")
}

// AuthApiUserShareOptions 用户分享选项
func (app *AppStruct) AuthApiUserShareOptions() {
	type option struct {
		Label    string `json:"label"`
		Value    int32  `json:"value"`
		Disabled bool   `json:"disabled"`
	}
	var options []option
	db, err := database.InDB(vars.Config.DB)
	if err != nil {
		utils.GinResult(app.Context, http.StatusBadRequest, "数据库连接失败", gin.H{"error": err.Error()})
		return
	}
	defer database.CloseDB(db)
	var users []database.User
	db.Model(&database.User{}).Where("id <> ?", app.UserInfo.Id).Find(&users)
	for _, user := range users {
		options = append(options, option{
			Label:    user.Email,
			Value:    user.Id,
			Disabled: false,
		})
	}
	utils.GinResult(app.Context, http.StatusOK, "分享用户选项", options)
}
