package command

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/spf13/cobra"
	"html/template"
	"io"
	"opsw/assets"
	"opsw/routes"
	"opsw/utils"
	"opsw/vars"
	"os"
	"strings"
	"time"
)

var runCommand = &cobra.Command{
	Use:   "run",
	Short: "启动服务",
	PreRun: func(cmd *cobra.Command, args []string) {
		if !utils.CheckOs() {
			utils.PrintError("暂不支持的操作系统")
			os.Exit(1)
		}
		err := utils.WriteFile(utils.CacheDir("/run"), utils.FormatYmdHis(time.Now()))
		if err != nil {
			utils.PrintError("无法写入文件")
			os.Exit(1)
		}
		if vars.RunConf.Host == "" {
			vars.RunConf.Host = "0.0.0.0"
		}
		if vars.RunConf.Port == "" {
			vars.RunConf.Port = "8080"
		}
	},
	Run: func(cmd *cobra.Command, args []string) {
		if vars.RunConf.Mode == "debug" {
			gin.SetMode(gin.DebugMode)
		} else if vars.RunConf.Mode == "test" {
			gin.SetMode(gin.TestMode)
		} else {
			gin.SetMode(gin.ReleaseMode)
		}
		router := gin.Default()
		templates, err := runTemplate()
		if err != nil {
			utils.PrintError(err.Error())
			os.Exit(1)
		}
		router.SetHTMLTemplate(templates)
		//
		router.Any("/*path", func(c *gin.Context) {
			routes.Entry(c)
		})
		//
		_ = router.Run(fmt.Sprintf("%s:%s", vars.RunConf.Host, vars.RunConf.Port))
	},
}

func runTemplate() (*template.Template, error) {
	distPath := "/web/dist/"
	if utils.IsDir(utils.CacheDir(distPath)) {
		_ = os.RemoveAll(utils.CacheDir(distPath))
	}
	t := template.New("")
	for name, file := range assets.Web.Files {
		if file.IsDir() {
			continue
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
		} else if strings.HasPrefix(name, distPath) {
			h, err := io.ReadAll(file)
			if err != nil {
				return nil, err
			}
			err = utils.WriteByte(utils.CacheDir(name), h)
			if err != nil {
				return nil, err
			}
		}
	}
	return t, nil
}

func init() {
	rootCommand.AddCommand(runCommand)
	runCommand.Flags().StringVar(&vars.RunConf.Host, "host", "", "主机名，默认：0.0.0.0")
	runCommand.Flags().StringVar(&vars.RunConf.Port, "port", "", "端口号，默认：8080")
	runCommand.Flags().StringVar(&vars.RunConf.Mode, "mode", "release", "运行模式，可选：debug|test|release")
}
