package routes

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"opsw/database"
	"opsw/utils"
	"opsw/vars"
	"os"
	"os/exec"
	"regexp"
	"strconv"
	"strings"
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
	if server, _ := database.ServerGet(map[string]any{
		"ip": ip,
	}, -1, false); server != nil && server.Id > 0 {
		utils.GinResult(app.Context, http.StatusBadRequest, fmt.Sprintf("服务器 [%s] 已存在", server.Ip))
		return
	}
	//
	server := &database.Server{
		Ip:       ip,
		Username: username,
		Password: password,
		Port:     port,
		Remark:   remark,
		State:    "Installing",
		Token:    utils.Base64Encode("s:%s", utils.GenerateString(22)),
	}
	serverUser := &database.ServerUser{
		UserId:  app.UserInfo.Id,
		OwnerId: app.UserInfo.Id,
	}
	err = database.ServerCreate(server, serverUser)
	if err != nil {
		utils.GinResult(app.Context, http.StatusBadRequest, "添加服务器失败", gin.H{"error": err.Error()})
	} else {
		command := fmt.Sprintf("content://%s/api/shell/start.sh?action=install&token=%s", utils.GinHomeUrl(app.Context), server.Token)
		url := fmt.Sprintf("%s/api/server/notify?token=%s", utils.GinHomeUrl(app.Context), server.Token)
		logf := utils.CacheDir("/logs/server/%s/serve.log", server.Ip)
		cmd := fmt.Sprintf("%s exec --host %s:%s --user %s --password %s --cmd %s --url %s --log %s >/dev/null 2>&1 &",
			runFile,
			server.Ip,
			server.Port,
			server.Username,
			utils.Base64Encode(server.Password),
			utils.Base64Encode(command),
			url,
			logf)
		_ = utils.WriteFile(logf, fmt.Sprintf("开始添加服务器 %s\n", utils.FormatYmdHis(time.Now())))
		_, _ = utils.Cmd("-c", cmd)
		utils.GinResult(app.Context, http.StatusOK, "添加服务器成功", server)
	}
}

// NoAuthApiServerNotify 添加/升级服务器时通知回调
func (app *AppStruct) NoAuthApiServerNotify() {
	var (
		ip     = utils.GinInput(app.Context, "ip")
		token  = utils.GinInput(app.Context, "token")
		state  = utils.GinInput(app.Context, "state")
		error_ = utils.GinInput(app.Context, "error")
		time_  = utils.GinInput(app.Context, "time")
	)
	if ip == "" {
		utils.GinResult(app.Context, http.StatusBadRequest, "IP不能为空")
		return
	}
	if token == "" {
		utils.GinResult(app.Context, http.StatusBadRequest, "Token不能为空")
		return
	}
	server, err := database.ServerGet(map[string]any{
		"ip":    ip,
		"token": token,
	}, -1, false)
	if err != nil {
		utils.GinResult(app.Context, http.StatusBadRequest, "服务器不存在", gin.H{"error": err.Error()})
		return
	}
	//
	logf := utils.CacheDir("/logs/server/%s/serve.log", server.Ip)
	action := "添加"
	if server.State == "Upgrading" {
		action = "升级"
	}
	if server.State == "Installing" || server.State == "Upgrading" {
		_ = utils.AppendToFile(logf, fmt.Sprintf("%s服务器结束，时间：%s\n", action, time_))
		if state == "error" {
			server.State = "Error"
			_ = utils.AppendToFile(logf, fmt.Sprintf("添加服务器失败，原因：%s\n", error_))
		} else {
			server.State = "Installed"
			_ = utils.AppendToFile(logf, fmt.Sprintf("%s服务器成功\n", action))
		}
		err = database.ServerUpdate(server)
		if err != nil {
			utils.GinResult(app.Context, http.StatusBadRequest, "更新服务器状态失败", gin.H{"error": err.Error()})
			return
		}
		app.wsNotifyServerUsers(vars.WsServerNotify, vars.WsClientStruct{
			Conn: nil,
			Type: vars.WsIsServer,
			Cid:  server.Id,
			Rid:  "api",
		})
	}
	utils.GinResult(app.Context, http.StatusOK, fmt.Sprintf("%s服务器成功", action), gin.H{
		"ip":    ip,
		"state": state,
		"error": error_,
	})
}

// AuthApiServerList 服务器列表
func (app *AppStruct) AuthApiServerList() {
	var (
		key         = utils.GinInput(app.Context, "key")
		page, _     = strconv.Atoi(utils.GinInput(app.Context, "page"))
		pageSize, _ = strconv.Atoi(utils.GinInput(app.Context, "page_size"))
	)
	db, err := database.InDB(vars.Config.DB)
	if err != nil {
		utils.GinResult(app.Context, http.StatusBadRequest, "数据库连接失败", gin.H{"error": err.Error()})
	}
	defer database.CloseDB(db)
	//
	query := db.Model(&database.Server{}).
		Select("servers.*, server_users.user_id, server_users.server_id, server_users.owner_id, server_infos.hostname, server_infos.platform, server_infos.platform_version, server_infos.version").
		Joins("left join server_users on server_users.server_id = servers.id").
		Joins("left join server_infos on server_infos.server_id = servers.id")
	if key != "" {
		query = query.Where("servers.ip like ? OR servers.remark like ? OR server_infos.hostname like ?", "%"+key+"%", "%"+key+"%", "%"+key+"%")
	}
	data := &vars.PageStruct{
		Page:     page,
		PageSize: pageSize,
		Data:     &[]*database.ServerItem{},
	}
	_, err = database.Page(query, data)
	if err != nil {
		utils.GinResult(app.Context, http.StatusBadRequest, "获取服务器列表失败", gin.H{"error": err.Error()})
	} else {
		for _, s := range *data.Data.(*[]*database.ServerItem) {
			database.ServerFormat(s)
		}
		utils.GinResult(app.Context, http.StatusOK, "服务器列表", data)
	}
}

// AuthApiServerOne 服务器详情
func (app *AppStruct) AuthApiServerOne() {
	id, _ := strconv.Atoi(app.Context.Query("id"))
	ip := app.Context.Query("ip")
	//
	where := map[string]any{}
	if id > 0 {
		where["servers.id"] = id
	} else {
		where["servers.ip"] = ip
	}
	//
	db, err := database.InDB(vars.Config.DB)
	if err != nil {
		utils.GinResult(app.Context, http.StatusBadRequest, "数据库连接失败", gin.H{"error": err.Error()})
	}
	defer database.CloseDB(db)
	//
	one := &database.ServerItem{}
	err = db.
		Model(&database.Server{}).
		Where(where).
		Select("servers.*, server_users.user_id, server_users.server_id, server_users.owner_id, server_infos.hostname, server_infos.platform, server_infos.platform_version, server_infos.version").
		Joins("left join server_users on server_users.server_id = servers.id").
		Joins("left join server_infos on server_infos.server_id = servers.id").
		Last(&one).Error
	//
	if err != nil {
		utils.GinResult(app.Context, http.StatusBadRequest, "获取服务器详情失败", gin.H{"error": err.Error()})
		return
	}
	utils.GinResult(app.Context, http.StatusOK, "服务器详情", database.ServerFormat(one))
}

// AuthApiServerLog 查看服务器日志
func (app *AppStruct) AuthApiServerLog() {
	id, _ := strconv.Atoi(app.Context.Query("id"))
	ip := app.Context.Query("ip")
	tail, _ := strconv.Atoi(app.Context.Query("tail"))
	//
	where := map[string]any{}
	if id > 0 {
		where["id"] = id
	} else {
		where["ip"] = ip
	}
	if tail <= 0 {
		tail = 200
	}
	if tail > 10000 {
		tail = 10000
	}
	//
	server, err := database.ServerGet(where, app.UserInfo.Id, false)
	if err != nil {
		utils.GinResult(app.Context, http.StatusBadRequest, "读取失败", gin.H{"error": err.Error()})
		return
	}
	logf := utils.CacheDir("/logs/server/%s/serve.log", server.Ip)
	if !utils.IsFile(logf) {
		utils.GinResult(app.Context, http.StatusBadRequest, "日志文件不存在")
		return
	}
	content, _ := utils.Cmd("-c", fmt.Sprintf("tail -%d %s", tail, logf))
	content = strings.TrimSpace(content)
	content = regexp.MustCompile(`\[\d+m`).ReplaceAllString(content, "")
	content = regexp.MustCompile(`\[\d+;\d+m`).ReplaceAllString(content, "")
	content = regexp.MustCompile(`\[([^[]*\.[^[]+):\d+]\x20?`).ReplaceAllString(content, "")
	utils.GinResult(app.Context, http.StatusOK, "读取成功", gin.H{
		"log": content,
	})
}

func (app *AppStruct) AuthApiServerOperation() {
	id, _ := strconv.Atoi(app.Context.Query("id"))
	ip := app.Context.Query("ip")
	operation := app.Context.Query("operation")
	remark := app.Context.Query("remark")
	//
	where := map[string]any{}
	if id > 0 {
		where["id"] = id
	} else {
		where["ip"] = ip
	}
	//
	runFile, err := exec.LookPath(os.Args[0])
	if err != nil {
		utils.GinResult(app.Context, http.StatusBadRequest, "获取当前路径失败", gin.H{"error": err.Error()})
		return
	}
	//
	if operation == "delete" {
		// 删除服务器
		server, err := database.ServerDelete(where, app.UserInfo.Id)
		if err != nil {
			utils.GinResult(app.Context, http.StatusBadRequest, "操作失败", gin.H{"error": err.Error()})
			return
		}
		command := fmt.Sprintf("content://%s/api/shell/end.sh", utils.GinHomeUrl(app.Context))
		logf := utils.CacheDir("/logs/server/%s/serve.log", server.Ip)
		cmd := fmt.Sprintf("%s exec --host %s:%s --user %s --password %s --cmd %s --log %s >/dev/null 2>&1 &",
			runFile,
			server.Ip,
			server.Port,
			server.Username,
			utils.Base64Encode(server.Password),
			utils.Base64Encode(command),
			logf)
		_ = utils.AppendToFile(logf, fmt.Sprintf("--------------------\n开始删除服务器 %s\n", utils.FormatYmdHis(time.Now())))
		_, _ = utils.Cmd("-c", cmd)
		utils.GinResult(app.Context, http.StatusOK, "操作成功")
	} else if operation == "upgrade" {
		// 升级服务器
		server, err := database.ServerGet(where, app.UserInfo.Id, true)
		if err != nil {
			utils.GinResult(app.Context, http.StatusBadRequest, "操作失败", gin.H{"error": err.Error()})
			return
		}
		serverInfo, _ := database.ServerInfoGet(server.Id)
		if serverInfo == nil || strings.Compare(vars.Version, serverInfo.Version) > 0 {
			server.State = "Upgrading"
			err = database.ServerUpdate(server)
			if err != nil {
				utils.GinResult(app.Context, http.StatusBadRequest, "操作失败", gin.H{"error": err.Error()})
				return
			}
			command := fmt.Sprintf("content://%s/api/shell/start.sh?action=upgrade&token=%s", utils.GinHomeUrl(app.Context), server.Token)
			url := fmt.Sprintf("%s/api/server/notify?token=%s", utils.GinHomeUrl(app.Context), server.Token)
			logf := utils.CacheDir("/logs/server/%s/serve.log", server.Ip)
			cmd := fmt.Sprintf("%s exec --host %s:%s --user %s --password %s --cmd %s --url %s --log %s >/dev/null 2>&1 &",
				runFile,
				server.Ip,
				server.Port,
				server.Username,
				utils.Base64Encode(server.Password),
				utils.Base64Encode(command),
				url,
				logf)
			_ = utils.AppendToFile(logf, fmt.Sprintf("--------------------\n开始升级服务器 %s\n", utils.FormatYmdHis(time.Now())))
			_, _ = utils.Cmd("-c", cmd)
			utils.GinResult(app.Context, http.StatusOK, "操作成功")
		} else {
			utils.GinResult(app.Context, http.StatusBadRequest, "服务器已经是最新版本")
		}
	} else if operation == "remark" {
		// 修改备注
		server, err := database.ServerGet(where, app.UserInfo.Id, false)
		if err != nil {
			utils.GinResult(app.Context, http.StatusBadRequest, "操作失败", gin.H{"error": err.Error()})
			return
		}
		server.Remark = remark
		err = database.ServerUpdate(server)
		if err != nil {
			utils.GinResult(app.Context, http.StatusBadRequest, "操作失败", gin.H{"error": err.Error()})
			return
		}
		utils.GinResult(app.Context, http.StatusOK, "修改成功")
	} else {
		utils.GinResult(app.Context, http.StatusBadRequest, "操作类型错误")
	}
}
