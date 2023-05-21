package command

import (
	"fmt"
	"github.com/nahid/gohttp"
	"github.com/spf13/cobra"
	"opsw/utils"
	"opsw/utils/logger"
	"opsw/vars"
	"os"
	"strings"
	"time"
)

var execConf = &vars.ExecStruct{}

var execCmd = &cobra.Command{
	Use:   "exec",
	Short: "远程执行命令",
	PreRun: func(cmd *cobra.Command, args []string) {
		if !utils.CheckOs() {
			utils.PrintError("暂不支持的操作系统")
			os.Exit(1)
		}
		if len(execConf.Host) == 0 || execConf.Cmd == "" {
			utils.PrintError("必须填写：host、cmd")
			os.Exit(0)
		}
		ip := execConf.Host
		port := "22"
		if ipport := strings.Split(execConf.Host, ":"); len(ipport) == 2 {
			ip = ipport[0]
			port = ipport[1]
		}
		if utils.StringToIP(ip) == nil {
			utils.PrintError(fmt.Sprintf("ip[%s]无效", ip))
			os.Exit(1)
		}
		execConf.Host = fmt.Sprintf("%s:%s", ip, port)
		if execConf.SSHConfig.User == "" {
			execConf.SSHConfig.User = "root"
		}
		if execConf.SSHConfig.Password != "" {
			execConf.SSHConfig.Password = utils.Base64Decode(execConf.SSHConfig.Password)
		}
		if execConf.SSHConfig.PkFile != "" {
			execConf.SSHConfig.PkPassword = execConf.SSHConfig.Password
		}
		if len(execConf.Cmd) > 0 {
			execConf.Cmd = utils.Base64Decode(execConf.Cmd)
		}
		if len(execConf.Param) > 0 {
			execConf.Param = utils.Base64Decode(execConf.Param)
		}
	},
	Run: func(cmd *cobra.Command, args []string) {
		execStart()
	},
}

func execStart() {
	if len(execConf.LogFile) > 0 {
		_ = logger.SetLogger(fmt.Sprintf(`{"File":{"filename":"%s","level":"TRAC","daily":true,"maxlines":100000,"maxsize":10,"maxdays":3,"append":true,"permit":"0660"}}`, execConf.LogFile))
	}

	logger.Info("---------- exec start ----------")

	key := utils.GenerateString(32)
	cmdFile := fmt.Sprintf("/tmp/.exec_%s", key)
	resFile := fmt.Sprintf("/tmp/.exec_%s_result", key)
	conFile := fmt.Sprintf("/tmp/.exec_%s_content", key)
	if strings.HasPrefix(execConf.Cmd, "content://") {
		execConf.Cmd = execConf.Cmd[10:]
		err := execConf.SSHConfig.CmdAsync(execConf.Host, fmt.Sprintf("curl -o %s -sSL '%s'", conFile, execConf.Cmd))
		if err != nil {
			response(err)
			return
		}
		execConf.Cmd = execConf.SSHConfig.CmdToStringNoLog(execConf.Host, fmt.Sprintf("cat %s", conFile), "\n")
	}
	err := execConf.SSHConfig.SaveFileAndChmodX(execConf.Host, cmdFile, utils.Shell("/exec.sh", map[string]any{
		"CMD":      execConf.Cmd,
		"END_TAG":  "success",
		"END_PATH": resFile,
	}))
	if err != nil {
		response(err)
		return
	}

	err = execConf.SSHConfig.CmdAsync(execConf.Host, fmt.Sprintf("sudo %s %s", cmdFile, execConf.Param))
	if err != nil {
		response(err)
		return
	}

	result := execConf.SSHConfig.CmdToStringNoLog(execConf.Host, fmt.Sprintf("cat %s", resFile), "")
	if result != "success" {
		response(fmt.Errorf("结果错误"))
		return
	}

	response(nil)

	_ = execConf.SSHConfig.CmdAsync(execConf.Host, fmt.Sprintf("rm -f %s", cmdFile))
	_ = execConf.SSHConfig.CmdAsync(execConf.Host, fmt.Sprintf("rm -f %s", resFile))
	_ = execConf.SSHConfig.CmdAsync(execConf.Host, fmt.Sprintf("rm -f %s", conFile))
}

func response(err error) {
	state := "success"
	errorMsg := ""
	if err != nil {
		state = "error"
		errorMsg = err.Error()
	}

	if strings.HasPrefix(execConf.Url, "http://") || strings.HasPrefix(execConf.Url, "https://") {
		logger.Info("---------- callback start ----------")
		ip, _ := utils.GetIpAndPort(execConf.Host)
		_, err := gohttp.NewRequest().
			FormData(map[string]string{
				"ip":    ip,
				"state": state,
				"error": errorMsg,
				"time":  utils.FormatYmdHis(time.Now()),
			}).
			Post(execConf.Url)
		if err != nil {
			logger.Info("---------- callback error ----------")
		} else {
			logger.Info("---------- callback end ----------")
		}
	}

	logger.Info("---------- exec end ----------")
}

func init() {
	rootCommand.AddCommand(execCmd)
	execCmd.Flags().StringVar(&execConf.Host, "host", "", "服务器IP: 192.168.0.5 or 192.168.0.5:22")
	execCmd.Flags().StringVar(&execConf.SSHConfig.User, "user", "root", "用户名，默认: root")
	execCmd.Flags().StringVar(&execConf.SSHConfig.Password, "password", "", "密码, 必须经过base64编码（如果设置pkfile，则为pkfile的密码）")
	execCmd.Flags().StringVar(&execConf.SSHConfig.PkFile, "pkfile", "", "密钥路径，如果设置，则使用密钥登录")
	execCmd.Flags().StringVar(&execConf.Cmd, "cmd", "", "执行的命令，如果想执行url内容请以\"content://\"开头, 必须经过base64编码")
	execCmd.Flags().StringVar(&execConf.Param, "param", "", "执行命令参数, 必须经过base64编码")
	execCmd.Flags().StringVar(&execConf.Url, "url", "", "回调地址, 必须以 \"http://\" or \"https://\" 前缀")
	execCmd.Flags().StringVar(&execConf.LogFile, "log", "", "日志文件路径")
}
