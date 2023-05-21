package routes

import (
	"fmt"
	"net/http"
	"opsw/utils"
)

// NoAuthShellStartSh 登录
func (app *AppStruct) NoAuthShellStartSh() {
	token := utils.GinInput(app.Context, "token")
	app.Context.String(http.StatusOK, utils.Shell("/start.sh", map[string]any{
		"URL":   fmt.Sprintf("%s%s/ws", utils.GinScheme(app.Context), app.Context.Request.Host),
		"TOKEN": token,
	}))
}
