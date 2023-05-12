package routes

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
	"net/http"
	"opsw/utils"
	"reflect"
	"strings"
)

type AppStruct struct {
	Context *gin.Context
}

func (app *AppStruct) Entry() {
	urlPath := app.Context.Request.URL.Path
	// 静态资源
	if strings.HasPrefix(urlPath, "/assets") {
		app.Context.File(utils.CacheDir(fmt.Sprintf("/web/dist%s", urlPath)))
		return
	}
	// 动态路由
	if callAppMethodByName(urlPath, app) {
		return
	}
	// 页面输出
	app.Context.HTML(http.StatusOK, "/web/dist/index.html", gin.H{
		"CODE": "",
		"MSG":  "",
	})
}

func callAppMethodByName(methodName string, object interface{}) bool {
	if strings.Contains(methodName, "/") || strings.Contains(methodName, "_") || strings.Contains(methodName, "-") {
		caser := cases.Title(language.Und)
		methodName = strings.ReplaceAll(methodName, "/", " ")
		methodName = strings.ReplaceAll(methodName, "_", " ")
		methodName = strings.ReplaceAll(methodName, "-", " ")
		methodName = strings.ReplaceAll(caser.String(methodName), " ", "")
	}
	if methodName == "Entry" {
		return false
	}
	method := reflect.ValueOf(object).MethodByName(methodName)
	if method.IsValid() {
		method.Call(nil)
		return true
	} else {
		return false
	}
}
