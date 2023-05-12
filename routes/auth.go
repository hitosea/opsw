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

func Entry(c *gin.Context) {
	urlPath := c.Request.URL.Path
	// 静态资源
	if strings.HasPrefix(urlPath, "/assets") {
		c.File(utils.CacheDir(fmt.Sprintf("/web/dist%s", urlPath)))
		return
	}
	// 动态路由
	if callFuncByName(urlPath, c) {
		return
	}
	// 页面输出
	c.HTML(http.StatusOK, "/web/dist/index.html", gin.H{
		"CODE": "",
		"MSG":  "",
	})
}

func callFuncByName(funcName string, args ...interface{}) bool {
	if strings.Contains(funcName, "/") || strings.Contains(funcName, "_") || strings.Contains(funcName, "-") {
		caser := cases.Title(language.Und)
		funcName = strings.ReplaceAll(funcName, "/", " ")
		funcName = strings.ReplaceAll(funcName, "_", " ")
		funcName = strings.ReplaceAll(funcName, "-", " ")
		funcName = strings.ReplaceAll(caser.String(funcName), " ", "")
	}
	function := reflect.ValueOf(funcName)
	if function.Kind() != reflect.Func {
		return false
	}
	if function.IsValid() && function.Type().NumIn() == len(args) {
		inputs := make([]reflect.Value, len(args))
		for i, arg := range args {
			inputs[i] = reflect.ValueOf(arg)
		}
		function.Call(inputs)
		return true
	} else {
		return false
	}
}
