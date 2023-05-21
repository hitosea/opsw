package routes

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"net/http"
	"opsw/database"
	"opsw/utils"
	"opsw/vars"
	"os"
	"os/exec"
	"time"
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
	runFile, err := exec.LookPath(os.Args[0])
	if err != nil {
		utils.GinResult(app.Context, http.StatusBadRequest, "获取当前路径失败", gin.H{"error": err.Error()})
		return
	}
	//
	db, err := database.InDB(vars.Config.DB)
	if err != nil {
		utils.GinResult(app.Context, http.StatusBadRequest, "数据库连接失败", gin.H{"error": err.Error()})
	}
	defer database.CloseDB(db)
	//
	var server = &database.Server{}
	var serverUser = &database.ServerUser{}
	db.Where(map[string]any{
		"ip": ip,
	}).Last(&server)
	if server.Id > 0 {
		utils.GinResult(app.Context, http.StatusBadRequest, fmt.Sprintf("服务器 [%s] 已存在", ip))
		return
	}
	server = &database.Server{
		Ip:       ip,
		Username: username,
		Password: password,
		Port:     port,
		Remark:   remark,
		State:    "Installing",
		Token:    utils.Base64Encode("s:%s", utils.GenerateString(22)),
	}
	serverUser = &database.ServerUser{
		UserId:  app.UserInfo.Id,
		OwnerId: app.UserInfo.Id,
	}
	err = db.Transaction(func(tx *gorm.DB) error {
		err = tx.Create(server).Error
		if err != nil {
			return err
		}
		serverUser.ServerId = server.Id
		err = tx.Create(serverUser).Error
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		utils.GinResult(app.Context, http.StatusBadRequest, "添加服务器失败", gin.H{"error": err.Error()})
	} else {
		start := fmt.Sprintf("%s/api/shell/start.sh?token=%s", utils.GinHomeUrl(app.Context), server.Token)
		curl := fmt.Sprintf("curl -sSL '%s' | bash", start)
		logf := utils.CacheDir("/logs/server/%s/deploy.log", server.Ip)
		cmd := fmt.Sprintf("%s exec --host %s:%s --user %s --password %s --cmd %s --log %s >/dev/null 2>&1 &",
			runFile,
			server.Ip,
			server.Port,
			server.Username,
			utils.Base64Encode(server.Password),
			utils.Base64Encode(curl),
			logf)
		fmt.Println(cmd)
		_ = utils.WriteFile(logf, fmt.Sprintf("开始部署服务器 %s\n", utils.FormatYmdHis(time.Now())))
		_, _ = utils.Cmd("-c", cmd)
		utils.GinResult(app.Context, http.StatusOK, "添加服务器成功", server)
	}
}

// AuthApiServerList 服务器列表
func (app *AppStruct) AuthApiServerList() {
	db, err := database.InDB(vars.Config.DB)
	if err != nil {
		utils.GinResult(app.Context, http.StatusBadRequest, "数据库连接失败", gin.H{"error": err.Error()})
	}
	defer database.CloseDB(db)
	//
	list := &[]database.ServerList{}
	err = db.
		Model(&database.Server{}).
		Select("servers.*, server_users.user_id, server_users.server_id, server_users.owner_id").
		Joins("left join server_users on server_users.server_id = servers.id").
		Limit(100).
		Scan(&list).Error
	if err != nil {
		utils.GinResult(app.Context, http.StatusBadRequest, "获取服务器列表失败", gin.H{"error": err.Error()})
	} else {
		for i := range *list {
			(*list)[i].Password = "******"
			(*list)[i].Token = "******"
		}
		utils.GinResult(app.Context, http.StatusOK, "服务器列表", gin.H{"list": list})
	}
}
