package routes

import (
	"net/http"
	"opsw/utils"
)

// NoAuthShellStartSh 登录
func (app *AppStruct) NoAuthShellStartSh() {
	app.Context.String(http.StatusOK, utils.Shell("/start.sh", map[string]any{}))
}
