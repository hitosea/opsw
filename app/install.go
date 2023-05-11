package app

import (
	"fmt"
	"github.com/spf13/cobra"
	"opsw/utils"
	"os"
)

var cmdFile string

var installApp = &cobra.Command{
	Use:   "install",
	Short: "安装服务",
	PreRun: func(cmd *cobra.Command, args []string) {
		cmdFile = utils.RunDir("/.opsw/install/cmd")
		//
		if !utils.CheckOs() {
			utils.PrintError("暂不支持的操作系统")
			os.Exit(1)
		}
		err := utils.WriteFile(cmdFile, utils.AssetsContent("install.sh", map[string]any{}))
		if err != nil {
			utils.PrintError(fmt.Sprintf("保存安装文件失败：%s", err.Error()))
			os.Exit(1)
		}
	},
	Run: func(cmd *cobra.Command, args []string) {
		utils.CmdFile(cmdFile)
	},
}

func init() {
	rootApp.AddCommand(installApp)
}
