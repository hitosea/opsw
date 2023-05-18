package utils

import (
	"bytes"
	"fmt"
	"io"
	"opsw/assets"
	"strings"
	"text/template"
)

var (
	assetsDict = make(map[string]string)
	sqlDict    = make(map[string]string)
)

// Assets 从模板中获取内容
func Assets(name string, envMap map[string]interface{}) string {
	if content, ok := assetsDict[name]; ok {
		return Template(content, envMap)
	}
	assetsDict[name] = ""
	for key, file := range assets.Shell.Files {
		if file.IsDir() {
			continue
		}
		if strings.HasSuffix(key, name) {
			h, err := io.ReadAll(file)
			if err == nil {
				assetsDict[name] = strings.ReplaceAll(string(h), "\t", "    ")
				break
			}
		}
	}
	return Template(assetsDict[name], envMap)
}

// Sql 从模板中获取内容
func Sql(name, autoIncrement string) []string {
	if _, ok := sqlDict[name]; ok {
		content := Template(sqlDict[name], map[string]any{
			"INCREMENT": autoIncrement,
		})
		return strings.Split(content, ";")
	}
	sqlDict[name] = ""
	for key, file := range assets.Database.Files {
		if file.IsDir() {
			continue
		}
		if strings.HasSuffix(key, name) {
			h, err := io.ReadAll(file)
			if err == nil {
				sqlDict[name] = string(h)
				break
			}
		}
	}
	content := Template(sqlDict[name], map[string]any{
		"INCREMENT": autoIncrement,
	})
	return strings.Split(content, ";")
}

// Template 从模板中获取内容
func Template(templateContent string, envMap map[string]interface{}) string {
	tmpl, err := template.New("text").Parse(templateContent)
	defer func() {
		if r := recover(); r != nil {
			PrintError(fmt.Sprintf("模板分析失败: %s", err))
		}
	}()
	if err != nil {
		panic(1)
	}
	envMap["RUN_PATH"] = RunDir("")
	envMap["WORK_PATH"] = CacheDir("")
	var buffer bytes.Buffer
	_ = tmpl.Execute(&buffer, envMap)
	return string(buffer.Bytes())
}
