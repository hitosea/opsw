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
	Context  *gin.Context
	UserInfo *database.User
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
	app.UserInfo.Token = utils.GinInput(app.Context, "token")
	if app.UserInfo.Token == "" {
		app.UserInfo.Token = utils.GinGetCookie(app.Context, "token")
	}
	if app.UserInfo.Token != "" {
		userInfo, err := database.UserGet(map[string]any{
			"token": app.UserInfo.Token,
		})
		if err != nil {
			app.UserInfo.Token = ""
		} else {
			app.UserInfo = userInfo
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
