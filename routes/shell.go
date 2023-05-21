package routes

import (
	"fmt"
	"net/http"
	"opsw/utils"
)

// NoAuthApiShellStartSh 启动脚本（添加服务器）
func (app *AppStruct) NoAuthApiShellStartSh() {
	token := utils.GinInput(app.Context, "token")
	app.Context.String(http.StatusOK, utils.Shell("/start.sh", map[string]any{
		"URL":   fmt.Sprintf("%s/ws", utils.GinHomeUrl(app.Context)),
		"TOKEN": token,
	}))
}

// NoAuthApiShellEndSh 结束脚本（删除服务器）
func (app *AppStruct) NoAuthApiShellEndSh() {
	app.Context.String(http.StatusOK, utils.Shell("/end.sh", map[string]any{}))
}
