package app

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/spf13/cobra"
	"html/template"
	"io"
	"opsw/resources/assets"
	"opsw/utils"
	"os"
	"strings"
	"time"
)

var serviceCmd = &cobra.Command{
	Use:   "service",
	Short: "启动服务",
	PreRun: func(cmd *cobra.Command, args []string) {
		if !utils.CheckOs() {
			utils.PrintError("暂不支持的操作系统")
			os.Exit(1)
		}
		err := utils.WriteFile(utils.WorkDir("run"), utils.FormatYmdHis(time.Now()))
		if err != nil {
			utils.PrintError("无法写入文件")
			os.Exit(1)
		}
	},
	Run: func(cmd *cobra.Command, args []string) {
		if RunConf.Mode == "debug" {
			gin.SetMode(gin.DebugMode)
		} else if RunConf.Mode == "test" {
			gin.SetMode(gin.TestMode)
		} else {
			gin.SetMode(gin.ReleaseMode)
		}
		router := gin.Default()
		templates, err := loadTemplate()
		if err != nil {
			utils.PrintError(err.Error())
			os.Exit(1)
		}
		router.SetHTMLTemplate(templates)
		//
		router.Any("/*path", func(c *gin.Context) {
			// todo
		})
		//
		_ = router.Run(fmt.Sprintf("%s:%s", RunConf.Host, RunConf.Port))
	},
}

func loadTemplate() (*template.Template, error) {
	t := template.New("")
	for name, file := range assets.AssetsWeb.Files {
		if file.IsDir() {
			continue
		}
		if strings.HasPrefix(name, "/resources/web/dist/assets/") {
			h, err := io.ReadAll(file)
			if err != nil {
				return nil, err
			}
			err = utils.WriteByte(utils.WorkDir(name), h)
			if err != nil {
				return nil, err
			}
		}
		if strings.HasSuffix(name, ".html") {
			h, err := io.ReadAll(file)
			if err != nil {
				return nil, err
			}
			t, err = t.New(name).Parse(string(h))
			if err != nil {
				return nil, err
			}
		}
	}
	return t, nil
}

func init() {
	rootApp.AddCommand(serviceCmd)
	serviceCmd.Flags().StringVar(&RunConf.Host, "host", "", "主机名，默认：0.0.0.0")
	serviceCmd.Flags().StringVar(&RunConf.Port, "port", "", "端口号，默认：8080")
	serviceCmd.Flags().StringVar(&RunConf.Mode, "mode", "release", "运行模式，可选：debug|test|release")
}
