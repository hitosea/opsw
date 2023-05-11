package utils

import (
	"bytes"
	"fmt"
	"io"
	"opsw/resources/assets"
	"strings"
	"text/template"
)

// AssetsContent 从模板中获取内容
func AssetsContent(name string, envMap map[string]interface{}) string {
	content := ""
	for key, file := range assets.AssetsShell.Files {
		if file.IsDir() {
			continue
		}
		if strings.HasSuffix(key, fmt.Sprintf("/%s", name)) {
			h, err := io.ReadAll(file)
			if err == nil {
				content = string(h)
				break
			}
		}
	}
	return TemplateContent(content, envMap)
}

// TemplateContent 从模板中获取内容
func TemplateContent(templateContent string, envMap map[string]interface{}) string {
	templateContent = strings.ReplaceAll(templateContent, "\t", "    ")
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
	var buffer bytes.Buffer
	_ = tmpl.Execute(&buffer, envMap)
	return string(buffer.Bytes())
}
