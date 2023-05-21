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
		command := fmt.Sprintf("content://%s/api/shell/start.sh?token=%s", utils.GinHomeUrl(app.Context), server.Token)
		url := fmt.Sprintf("%s/api/server/create/notify?token=%s", utils.GinHomeUrl(app.Context), server.Token)
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
		fmt.Println(cmd)
		_ = utils.WriteFile(logf, fmt.Sprintf("开始添加服务器 %s\n", utils.FormatYmdHis(time.Now())))
		_, _ = utils.Cmd("-c", cmd)
		utils.GinResult(app.Context, http.StatusOK, "添加服务器成功", server)
	}
}

// NoAuthApiServerCreateNotify 添加服务器时通知回调
func (app *AppStruct) NoAuthApiServerCreateNotify() {
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
	if server.State == "Installing" {
		_ = utils.AppendToFile(logf, fmt.Sprintf("添加服务器结束，时间：%s\n", time_))
		if state == "error" {
			server.State = "Error"
			_ = utils.AppendToFile(logf, fmt.Sprintf("添加服务器失败，原因：%s\n", error_))
		} else {
			server.State = "Installed"
			_ = utils.AppendToFile(logf, fmt.Sprintf("添加服务器成功\n"))
		}
		err = database.ServerUpdate(server)
		if err != nil {
			utils.GinResult(app.Context, http.StatusBadRequest, "更新服务器状态失败", gin.H{"error": err.Error()})
			return
		}
	}
	utils.GinResult(app.Context, http.StatusOK, "添加服务器成功", gin.H{
		"ip":    ip,
		"state": state,
		"error": error_,
	})
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
		for i, s := range *list {
			(*list)[i].Password = "******"
			(*list)[i].Token = "******"
			if s.State == "Installed" {
				(*list)[i].State = "Offline"
				for _, v := range vars.WsClients {
					if v.Type == "server" && v.Uid == s.Id {
						(*list)[i].State = "Online"
					}
				}
			}
		}
		utils.GinResult(app.Context, http.StatusOK, "服务器列表", gin.H{"list": list})
	}
}

// AuthApiServerLog 查看服务器日志
func (app *AppStruct) AuthApiServerLog() {
	ip := app.Context.Query("ip")
	tail, _ := strconv.Atoi(app.Context.Query("tail"))
	if tail <= 0 {
		tail = 200
	}
	if tail > 10000 {
		tail = 10000
	}
	if _, err := database.ServerGet(map[string]any{
		"ip": ip,
	}, app.UserInfo.Id, false); err != nil {
		utils.GinResult(app.Context, http.StatusBadRequest, "读取失败", gin.H{"error": err.Error()})
		return
	}
	logf := utils.CacheDir("/logs/server/%s/serve.log", ip)
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
	ip := app.Context.Query("ip")
	operation := app.Context.Query("operation")
	//
	if operation == "delete" {
		runFile, err := exec.LookPath(os.Args[0])
		if err != nil {
			utils.GinResult(app.Context, http.StatusBadRequest, "获取当前路径失败", gin.H{"error": err.Error()})
			return
		}
		server, err := database.ServerDelete(map[string]any{
			"ip": ip,
		}, app.UserInfo.Id)
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
		fmt.Println(cmd)
		_ = utils.AppendToFile(logf, fmt.Sprintf("------------------------\n开始删除服务器 %s\n", utils.FormatYmdHis(time.Now())))
		_, _ = utils.Cmd("-c", cmd)
		utils.GinResult(app.Context, http.StatusOK, "操作成功")
	} else {
		utils.GinResult(app.Context, http.StatusBadRequest, "操作类型错误")
	}
}
