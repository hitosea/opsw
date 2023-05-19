package command

import (
	"fmt"
	"github.com/nahid/gohttp"
	"github.com/spf13/cobra"
	"opsw/utils"
	"opsw/utils/logger"
	"opsw/vars"
	"os"
	"strconv"
	"strings"
	"time"
)

var cfg = &vars.ExecStruct{}

// execCmd represents the exec command
var execCmd = &cobra.Command{
	Use:   "exec",
	Short: "远程执行命令",
	PreRun: func(cmd *cobra.Command, args []string) {
		if !utils.CheckOs() {
			utils.PrintError("暂不支持的操作系统")
			os.Exit(1)
		}
		if len(cfg.Host) == 0 || cfg.Cmd == "" {
			utils.PrintError("必须填写：host、cmd")
			os.Exit(0)
		}
		ip := cfg.Host
		port := "22"
		if ipport := strings.Split(cfg.Host, ":"); len(ipport) == 2 {
			ip = ipport[0]
			port = ipport[1]
		}
		if utils.StringToIP(ip) == nil {
			utils.PrintError(fmt.Sprintf("ip[%s]无效", ip))
			os.Exit(1)
		}
		cfg.Host = fmt.Sprintf("%s:%s", ip, port)
		if cfg.SSHConfig.User == "" {
			cfg.SSHConfig.User = "root"
		}
		if cfg.SSHConfig.Password != "" {
			cfg.SSHConfig.Password = utils.Base64Decode(cfg.SSHConfig.Password)
		}
		if cfg.SSHConfig.PkFile != "" {
			cfg.SSHConfig.PkPassword = cfg.SSHConfig.Password
		}
		if len(cfg.Cmd) > 0 {
			cfg.Cmd = utils.Base64Decode(cfg.Cmd)
		}
		if len(cfg.Param) > 0 {
			cfg.Param = utils.Base64Decode(cfg.Param)
		}
	},
	Run: func(cmd *cobra.Command, args []string) {
		execStart()
	},
}

func execStart() {
	if len(cfg.LogFile) > 0 {
		_ = logger.SetLogger(fmt.Sprintf(`{"File":{"filename":"%s","level":"TRAC","daily":true,"maxlines":100000,"maxsize":10,"maxdays":3,"append":true,"permit":"0660"}}`, cfg.LogFile))
	}

	logger.Info("---------- exec start ----------")

	key := utils.GenerateString(32)
	cmdFile := utils.CacheDir("/tmp/.exec_%s", key)
	resFile := utils.CacheDir("/tmp/.exec_%s_result", key)
	err := cfg.SSHConfig.SaveFileAndChmodX(cfg.Host, cmdFile, utils.Shell("/exec.sh", map[string]any{
		"CMD":      cfg.Cmd,
		"END_TAG":  "success",
		"END_PATH": resFile,
	}))
	if err != nil {
		response(err)
		return
	}

	err = cfg.SSHConfig.CmdAsync(cfg.Host, fmt.Sprintf("sudo %s %s", cmdFile, cfg.Param))
	if err != nil {
		response(err)
		return
	}

	result := cfg.SSHConfig.CmdToStringNoLog(cfg.Host, fmt.Sprintf("cat %s", resFile), "")
	if result != "success" {
		response(fmt.Errorf("结果错误"))
		return
	}

	response(nil)

	_ = cfg.SSHConfig.CmdAsync(cfg.Host, fmt.Sprintf("rm -f %s", cmdFile))
	_ = cfg.SSHConfig.CmdAsync(cfg.Host, fmt.Sprintf("rm -f %s", resFile))
}

func response(err error) {
	status := "success"
	errorMsg := ""
	if err != nil {
		status = "error"
		errorMsg = err.Error()
	}

	if strings.HasPrefix(cfg.Url, "http://") || strings.HasPrefix(cfg.Url, "https://") {
		logger.Info("---------- callback start ----------")
		ip, _ := utils.GetIpAndPort(cfg.Host)
		timestamp := strconv.FormatInt(time.Now().Unix(), 10)
		_, err := gohttp.NewRequest().
			FormData(map[string]string{
				"ip":        ip,
				"status":    status,
				"error":     errorMsg,
				"timestamp": timestamp,
			}).
			Post(cfg.Url)
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
	execCmd.Flags().StringVar(&cfg.Host, "host", "", "服务器IP: 192.168.0.5 or 192.168.0.5:22")
	execCmd.Flags().StringVar(&cfg.SSHConfig.User, "user", "root", "用户名，默认: root")
	execCmd.Flags().StringVar(&cfg.SSHConfig.Password, "password", "", "密码, 必须经过base64编码（如果设置pkfile，则为pkfile的密码）")
	execCmd.Flags().StringVar(&cfg.SSHConfig.PkFile, "pkfile", "", "密钥路径，如果设置，则使用密钥登录")
	execCmd.Flags().StringVar(&cfg.Cmd, "cmd", "", "执行的命令, 必须经过base64编码")
	execCmd.Flags().StringVar(&cfg.Param, "param", "", "参数, 必须经过base64编码")
	execCmd.Flags().StringVar(&cfg.Url, "url", "", "回调地址, 必须以 \"http://\" or \"https://\" 前缀")
	execCmd.Flags().StringVar(&cfg.LogFile, "log", "", "日志文件路径")
}
