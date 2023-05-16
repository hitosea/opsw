package command

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/spf13/cobra"
	"html/template"
	"io"
	"opsw/assets"
	"opsw/database"
	"opsw/routes"
	"opsw/utils"
	"opsw/vars"
	"os"
	"strings"
)

var runCommand = &cobra.Command{
	Use:   "run",
	Short: "启动服务",
	PreRun: func(cmd *cobra.Command, args []string) {
		if !utils.CheckOs() {
			utils.PrintError("暂不支持的操作系统")
			os.Exit(1)
		}
		if vars.Config.Host == "" {
			vars.Config.Host = "0.0.0.0"
		}
		if vars.Config.Port == "" {
			vars.Config.Port = "8080"
		}
		if vars.Config.DB == "" {
			vars.Config.DB = fmt.Sprintf("sqlite3://%s", utils.CacheDir("/run.db"))
		}
		err := utils.WriteFile(utils.CacheDir("/run.json"), utils.StructToJson(vars.Config))
		if err != nil {
			utils.PrintError(fmt.Sprintf("配置文件写入失败: %s", err.Error()))
			os.Exit(1)
		}
		_, err = database.Init()
		if err != nil {
			utils.PrintError(fmt.Sprintf("数据库初始化失败: %s", err.Error()))
			os.Exit(1)
		}
	},
	Run: func(cmd *cobra.Command, args []string) {
		if vars.Config.Mode == "debug" {
			gin.SetMode(gin.DebugMode)
		} else if vars.Config.Mode == "test" {
			gin.SetMode(gin.TestMode)
		} else {
			gin.SetMode(gin.ReleaseMode)
		}
		router := gin.Default()
		templates, err := runTemplate()
		if err != nil {
			utils.PrintError(fmt.Sprintf("模板初始化失败: %s", err.Error()))
			os.Exit(1)
		}
		router.SetHTMLTemplate(templates)
		//
		router.Any("/*path", func(context *gin.Context) {
			app := routes.AppStruct{
				Context:  context,
				UserInfo: &vars.UserModel{},
			}
			app.Entry()
		})
		//
		_ = router.Run(fmt.Sprintf("%s:%s", vars.Config.Host, vars.Config.Port))
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
	runCommand.Flags().StringVar(&vars.Config.Host, "host", "", "主机名，默认：0.0.0.0")
	runCommand.Flags().StringVar(&vars.Config.Port, "port", "", "端口号，默认：8080")
	runCommand.Flags().StringVar(&vars.Config.Mode, "mode", "release", "运行模式，可选：debug|test|release")
	runCommand.Flags().StringVar(&vars.Config.DB, "db", "", "数据库连接地址，如：sqlite://run.db")
}
