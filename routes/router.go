package routes

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
	"net/http"
	"opsw/database"
	"opsw/utils"
	"reflect"
	"strings"
)

type AppStruct struct {
	Context    *gin.Context
	UserInfo   *database.User
	ServerInfo *database.Server
}

func (app *AppStruct) Entry() {
	urlPath := app.Context.Request.URL.Path
	methodName := urlToName(urlPath)
	// 静态资源
	if strings.HasPrefix(urlPath, "/assets") {
		app.Context.File(utils.CacheDir("/web/dist%s", urlPath))
		return
	}
	if strings.HasSuffix(urlPath, "/favicon.ico") {
		app.Context.Status(http.StatusOK)
		return
	}
	// 读取身份
	token := utils.GinToken(app.Context)
	if strings.HasPrefix(utils.Base64Decode(token), "u:") {
		if info, err := database.UserGet(map[string]any{
			"token": token,
		}); err == nil {
			app.UserInfo = info
		}
	} else if strings.HasPrefix(utils.Base64Decode(token), "s:") {
		if info, err := database.ServerGet(map[string]any{
			"token": token,
		}, -1, false); err == nil {
			app.ServerInfo = info
		}
	}
	// 动态路由（不需要登录）
	if callByName(false, methodName, app) {
		return
	}
	// 登录验证
	if app.UserInfo.Token == "" {
		utils.GinResult(app.Context, http.StatusUnauthorized, "请先登录")
		return
	}
	// 动态路由（需要登录）
	if callByName(true, methodName, app) {
		return
	}
	// 页面输出
	app.Context.HTML(http.StatusOK, "/web/dist/index.html", gin.H{
		"CODE": "",
		"MSG":  "",
	})
}

func urlToName(urlPath string) string {
	if strings.Contains(urlPath, "/") || strings.Contains(urlPath, "_") || strings.Contains(urlPath, "-") {
		caser := cases.Title(language.Und)
		urlPath = strings.ReplaceAll(urlPath, "/", " ")
		urlPath = strings.ReplaceAll(urlPath, "_", " ")
		urlPath = strings.ReplaceAll(urlPath, "-", " ")
		urlPath = strings.ReplaceAll(urlPath, ".", " ")
		urlPath = strings.ReplaceAll(caser.String(urlPath), " ", "")
	}
	if urlPath == "Entry" {
		return ""
	}
	return urlPath
}

func callByName(requireAuth bool, methodName string, object interface{}) bool {
	if methodName == "" {
		return false
	}
	if requireAuth {
		methodName = fmt.Sprintf("Auth%s", methodName)
	} else {
		methodName = fmt.Sprintf("NoAuth%s", methodName)
	}
	method := reflect.ValueOf(object).MethodByName(methodName)
	if method.IsValid() {
		method.Call(nil)
		return true
	} else {
		return false
	}
}
