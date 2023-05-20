package command

import (
	"github.com/spf13/cobra"
	"opsw/utils"
	"opsw/vars"
	"os"
	"strings"
)

var workConf = &vars.WorkStruct{}

// workCmd represents the websocket command
var workCmd = &cobra.Command{
	Use:   "work",
	Short: "启动工作模式",
	PreRun: func(cmd *cobra.Command, args []string) {
		workConf.Url = utils.Base64Decode(workConf.Url)
		if workConf.Url == "" {
			utils.PrintError("请填写服务端url")
			os.Exit(0)
		}
		if !strings.HasPrefix(workConf.Url, "ws://") && !strings.HasPrefix(workConf.Url, "wss://") {
			utils.PrintError("服务端url必须以ws://或wss://开头")
			os.Exit(0)
		}
		if workConf.Mode == "" {
			utils.PrintError("请填写客户端类型")
			os.Exit(0)
		}
	},
	Run: func(cmd *cobra.Command, args []string) {
		workStart()
	},
}

func workStart() {

}

func init() {
	rootCommand.AddCommand(workCmd)
	workCmd.Flags().StringVar(&workConf.Url, "url", "", "服务端url，必须经过base64编码")
	workCmd.Flags().StringVar(&workConf.Mode, "mode", "", "客户端类型")
}
