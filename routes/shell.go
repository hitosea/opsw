package routes

import (
	"fmt"
	"net/http"
	"opsw/utils"
)

// NoAuthApiShellStartSh 启动脚本（添加服务器）
func (app *AppStruct) NoAuthApiShellStartSh() {
	var (
		action        = utils.GinInput(app.Context, "action")
		token         = utils.GinInput(app.Context, "token")
		panelPort     = utils.GinInput(app.Context, "panel_port")
		panelUsername = utils.GinInput(app.Context, "panel_username")
		panelPassword = utils.GinInput(app.Context, "panel_password")
	)
	app.Context.String(http.StatusOK, utils.Shell("/start.sh", map[string]any{
		"URL":    fmt.Sprintf("%s/ws", utils.GinHomeUrl(app.Context)),
		"TOKEN":  token,
		"ACTION": action,

		"PANEL_PORT":     panelPort,
		"PANEL_USERNAME": panelUsername,
		"PANEL_PASSWORD": panelPassword,
	}))
}

// NoAuthApiShellEndSh 结束脚本（删除服务器）
func (app *AppStruct) NoAuthApiShellEndSh() {
	app.Context.String(http.StatusOK, utils.Shell("/end.sh", map[string]any{}))
}
