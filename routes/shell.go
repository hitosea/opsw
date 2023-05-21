package routes

import (
	"fmt"
	"net/http"
	"opsw/utils"
)

// NoAuthApiShellStartSh 登录
func (app *AppStruct) NoAuthApiShellStartSh() {
	token := utils.GinInput(app.Context, "token")
	app.Context.String(http.StatusOK, utils.Shell("/start.sh", map[string]any{
		"URL":   fmt.Sprintf("%s/ws", utils.GinHomeUrl(app.Context)),
		"TOKEN": token,
	}))
}
