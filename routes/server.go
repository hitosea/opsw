package routes

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"net/http"
	"opsw/database"
	"opsw/utils"
	"opsw/vars"
)

// AuthApiServerCreate 添加服务器
func (app *AppStruct) AuthApiServerCreate() {
	var (
		ip       = utils.GinInput(app.Context, "ip")
		username = utils.GinInput(app.Context, "username")
		password = utils.GinInput(app.Context, "password")
		port     = utils.GinInput(app.Context, "port")
		remark   = utils.GinInput(app.Context, "remark")
	)
	if utils.StringToIP(ip) == nil {
		utils.GinResult(app.Context, http.StatusBadRequest, "IP格式不正确")
		return
	}
	if username == "" {
		username = "root"
	}
	if port == "" {
		port = "22"
	}
	//
	db, err := database.InDB(vars.Config.DB)
	if err != nil {
		utils.GinResult(app.Context, http.StatusBadRequest, "数据库连接失败", gin.H{"error": err.Error()})
	}
	defer database.CloseDB(db)
	//
	var server = &database.Server{}
	db.Where(map[string]any{
		"ip": ip,
	}).Last(&server)
	if server.Id > 0 {
		utils.GinResult(app.Context, http.StatusBadRequest, "服务器已存在")
		return
	}
	server = &database.Server{
		Ip:       ip,
		Username: username,
		Password: password,
		Port:     port,
		Remark:   remark,
		State:    "Creating",
	}
	err = db.Transaction(func(tx *gorm.DB) error {
		err = tx.Create(server).Error
		if err != nil {
			return err
		}
		err = tx.Create(&database.ServerUser{
			ServerId: server.Id,
			UserId:   app.UserInfo.Id,
			OwnerId:  app.UserInfo.Id,
		}).Error
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		utils.GinResult(app.Context, http.StatusBadRequest, "添加服务器失败", gin.H{"error": err.Error()})
	} else {
		utils.GinResult(app.Context, http.StatusOK, "添加服务器成功", server)
	}
}
