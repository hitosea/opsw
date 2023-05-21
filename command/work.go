package command

import (
	"encoding/json"
	"fmt"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/togettoyou/wsc"
	"opsw/utils"
	"opsw/utils/logger"
	"opsw/vars"
	"os"
	"strings"
	"time"
)

var (
	workLogDir = utils.CacheDir("/logs/work")
	workConf   = &vars.WorkStruct{}
	workRid    = ""
)

var workCmd = &cobra.Command{
	Use:   "work",
	Short: "启动工作模式",
	PreRun: func(cmd *cobra.Command, args []string) {
		if workConf.Conf != "" && utils.IsFile(workConf.Conf) {
			viper.SetConfigFile(workConf.Conf)
			err := viper.ReadInConfig()
			if err == nil {
				workConf.Url = viper.GetString("url")
				workConf.Mode = viper.GetString("mode")
				workConf.Token = viper.GetString("token")
			}
		}
		if workConf.Url == "" {
			utils.PrintError("请填写服务端url")
			os.Exit(0)
		}
		if workConf.Mode == "" {
			utils.PrintError("请填写客户端类型")
			os.Exit(0)
		}
		if workConf.Token == "" {
			utils.PrintError("请填写客户端token")
			os.Exit(0)
		}
		if !strings.HasPrefix(workConf.Url, "ws://") &&
			!strings.HasPrefix(workConf.Url, "wss://") &&
			!strings.HasPrefix(workConf.Url, "http://") &&
			!strings.HasPrefix(workConf.Url, "https://") {
			utils.PrintError("服务端url必须以ws://或wss://开头")
			os.Exit(0)
		}
	},
	Run: func(cmd *cobra.Command, args []string) {
		workStart()
	},
}

func workStart() {
	workConf.Url = strings.Replace(workConf.Url, "http://", "ws://", 1)
	workConf.Url = strings.Replace(workConf.Url, "https://", "wss://", 1)
	nodeName, _ := os.Hostname()
	if strings.Contains(workConf.Url, "?") {
		workConf.Url = fmt.Sprintf("%s&mode=%s&token=%s&name=%s", workConf.Url, workConf.Mode, workConf.Token, nodeName)
	} else {
		workConf.Url = fmt.Sprintf("%s?mode=%s&token=%s&name=%s", workConf.Url, workConf.Mode, workConf.Token, nodeName)
	}
	//
	err := utils.Mkdir(workLogDir, 0755)
	if err != nil {
		logger.Error("Failed to create log dir: %s\n", err.Error())
		os.Exit(1)
	}
	_ = logger.SetLogger(fmt.Sprintf(`{"File":{"filename":"%s/work.log","level":"TRAC","daily":true,"maxlines":100000,"maxsize":10,"maxdays":3,"append":true,"permit":"0660"}}`, workLogDir))
	//
	done := make(chan bool)
	ws := wsc.New(workConf.Url)
	// 自定义配置
	ws.SetConfig(&wsc.Config{
		WriteWait:         10 * time.Second,
		MaxMessageSize:    512 * 1024, // 512KB
		MinRecTime:        2 * time.Second,
		MaxRecTime:        30 * time.Second,
		RecFactor:         1.5,
		MessageBufferSize: 1024,
	})
	// 设置回调处理
	ws.OnConnected(func() {
		logger.Debug("OnConnected: ", ws.WebSocket.Url)
		logger.SetWebsocket(ws)
	})
	ws.OnConnectError(func(err error) {
		logger.Debug("OnConnectError: ", err.Error())
	})
	ws.OnDisconnected(func(err error) {
		logger.Debug("OnDisconnected: ", err.Error())
	})
	ws.OnClose(func(code int, text string) {
		logger.Debug("OnClose: ", code, text)
		done <- true
	})
	ws.OnTextMessageSent(func(message string) {
		logger.Debug("OnTextMessageSent: ", message)
	})
	ws.OnBinaryMessageSent(func(data []byte) {
		logger.Debug("OnBinaryMessageSent: ", string(data))
	})
	ws.OnSentError(func(err error) {
		logger.Debug("OnSentError: ", err.Error())
	})
	ws.OnPingReceived(func(appData string) {
		logger.Debug("OnPingReceived: ", appData)
	})
	ws.OnPongReceived(func(appData string) {
		logger.Debug("OnPongReceived: ", appData)
	})
	ws.OnTextMessageReceived(func(message string) {
		logger.Debug("OnTextMessageReceived: ", message)
		handleMessageReceived(ws, message)
	})
	ws.OnBinaryMessageReceived(func(data []byte) {
		logger.Debug("OnBinaryMessageReceived: ", string(data))
	})
	// 开始连接
	go ws.Connect()
	for {
		select {
		case <-done:
			return
		}
	}
}

// 处理消息
func handleMessageReceived(ws *wsc.Wsc, message string) {
	var msg map[string]any
	if ok := json.Unmarshal([]byte(message), &msg); ok == nil {
		msgType, _ := msg["type"].(float64)
		msgData, _ := msg["data"].(any)
		if msgType == vars.WsOnline {
			if dataMap, _ := msgData.(map[string]any); dataMap != nil {
				if own, _ := dataMap["own"].(float64); own == 1 {
					workRid, _ = dataMap["rid"].(string)
				}
			}
		}
	}
}

func init() {
	rootCommand.AddCommand(workCmd)
	workCmd.Flags().StringVar(&workConf.Url, "url", "", "服务端url")
	workCmd.Flags().StringVar(&workConf.Mode, "mode", "", "客户端类型")
	workCmd.Flags().StringVar(&workConf.Token, "token", "", "客户端token")
	workCmd.Flags().StringVar(&workConf.Conf, "conf", "", "配置文件路径")
}
