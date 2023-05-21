package command

import (
	"encoding/json"
	"fmt"
	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/disk"
	"github.com/shirou/gopsutil/v3/host"
	"github.com/shirou/gopsutil/v3/load"
	"github.com/shirou/gopsutil/v3/mem"
	"github.com/shirou/gopsutil/v3/net"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/togettoyou/wsc"
	"opsw/database"
	"opsw/utils"
	"opsw/utils/logger"
	"opsw/vars"
	"os"
	"strings"
	"time"
)

var (
	ws *wsc.Wsc

	workLogDir = "/var/log/opsw"
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
	ws = wsc.New(workConf.Url)
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
		onConnected()
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

// 连接成功
func onConnected() {
	data, err := loadBaseInfo("all", "all")
	if err != nil {
		logger.Error("获取基础信息失败: ", err.Error())
	} else {
		sendJson(vars.WsMsgStruct{Type: vars.WsServerInfo, Data: data})
	}
}

func sendJson(data any) {
	if data == nil {
		return
	}
	ss, err := json.Marshal(data)
	if err != nil {
		logger.Error("发送JSON消息序列化失败: ", err.Error())
		return
	}
	err = ws.SendTextMessage(string(ss))
	if err != nil {
		logger.Error("发送JSON消息失败: ", err.Error())
	}
}

// 读取服务器基础信息
func loadBaseInfo(ioOption string, netOption string) (*database.ServerInfo, error) {
	var baseInfo database.ServerInfo
	hostInfo, err := host.Info()
	if err != nil {
		return nil, err
	}
	baseInfo.Hostname = hostInfo.Hostname
	baseInfo.OS = hostInfo.OS
	baseInfo.Platform = hostInfo.Platform
	baseInfo.PlatformFamily = hostInfo.PlatformFamily
	baseInfo.PlatformVersion = hostInfo.PlatformVersion
	baseInfo.KernelArch = hostInfo.KernelArch
	baseInfo.KernelVersion = hostInfo.KernelVersion
	ss, _ := json.Marshal(hostInfo)
	baseInfo.VirtualizationSystem = string(ss)

	cpuInfo, err := cpu.Info()
	if err == nil {
		baseInfo.CPUModelName = cpuInfo[0].ModelName
	}

	baseInfo.CPUCores, _ = cpu.Counts(false)
	baseInfo.CPULogicalCores, _ = cpu.Counts(true)

	baseInfo.CurrentInfo = *loadCurrentInfo(ioOption, netOption)

	baseInfo.Version = vars.Version
	baseInfo.CommitSHA = vars.CommitSHA
	return &baseInfo, nil
}

// 读取服务器当前信息
func loadCurrentInfo(ioOption string, netOption string) *database.ServerInfoCurrent {
	var currentInfo database.ServerInfoCurrent
	hostInfo, _ := host.Info()
	currentInfo.Uptime = hostInfo.Uptime
	currentInfo.TimeSinceUptime = time.Now().Add(-time.Duration(hostInfo.Uptime) * time.Second).Format("2006-01-02 15:04:05")
	currentInfo.Procs = hostInfo.Procs

	currentInfo.CPUTotal, _ = cpu.Counts(true)
	totalPercent, _ := cpu.Percent(0, false)
	if len(totalPercent) == 1 {
		currentInfo.CPUUsedPercent = totalPercent[0]
		currentInfo.CPUUsed = currentInfo.CPUUsedPercent * 0.01 * float64(currentInfo.CPUTotal)
	}
	currentInfo.CPUPercent, _ = cpu.Percent(0, true)

	loadInfo, _ := load.Avg()
	currentInfo.Load1 = loadInfo.Load1
	currentInfo.Load5 = loadInfo.Load5
	currentInfo.Load15 = loadInfo.Load15
	currentInfo.LoadUsagePercent = loadInfo.Load1 / (float64(currentInfo.CPUTotal*2) * 0.75) * 100

	memoryInfo, _ := mem.VirtualMemory()
	currentInfo.MemoryTotal = memoryInfo.Total
	currentInfo.MemoryAvailable = memoryInfo.Available
	currentInfo.MemoryUsed = memoryInfo.Used
	currentInfo.MemoryUsedPercent = memoryInfo.UsedPercent

	currentInfo.DiskData = loadDiskInfo()

	if ioOption == "all" {
		diskInfo, _ := disk.IOCounters()
		for _, state := range diskInfo {
			currentInfo.IOReadBytes += state.ReadBytes
			currentInfo.IOWriteBytes += state.WriteBytes
			currentInfo.IOCount += state.ReadCount + state.WriteCount
			currentInfo.IOReadTime += state.ReadTime
			currentInfo.IOWriteTime += state.WriteTime
		}
	} else {
		diskInfo, _ := disk.IOCounters(ioOption)
		for _, state := range diskInfo {
			currentInfo.IOReadBytes += state.ReadBytes
			currentInfo.IOWriteBytes += state.WriteBytes
			currentInfo.IOCount += state.ReadCount + state.WriteCount
			currentInfo.IOReadTime += state.ReadTime
			currentInfo.IOWriteTime += state.WriteTime
		}
	}

	if netOption == "all" {
		netInfo, _ := net.IOCounters(false)
		if len(netInfo) != 0 {
			currentInfo.NetBytesSent = netInfo[0].BytesSent
			currentInfo.NetBytesRecv = netInfo[0].BytesRecv
		}
	} else {
		netInfo, _ := net.IOCounters(true)
		for _, state := range netInfo {
			if state.Name == netOption {
				currentInfo.NetBytesSent = state.BytesSent
				currentInfo.NetBytesRecv = state.BytesRecv
			}
		}
	}

	currentInfo.ShotTime = time.Now()
	return &currentInfo
}

// 读取服务器磁盘信息
func loadDiskInfo() []database.ServerInfoDiskInfo {
	var datas []database.ServerInfoDiskInfo
	stdout, err := utils.Exec("df -hT -P|grep '/'|grep -v tmpfs|grep -v 'snap/core'|grep -v udev")
	if err != nil {
		return datas
	}
	lines := strings.Split(stdout, "\n")

	var mounts []database.ServerInfoDiskMount
	var excludes = []string{"/mnt/cdrom", "/boot", "/boot/efi", "/dev", "/dev/shm", "/run/lock", "/run", "/run/shm", "/run/user"}
	for _, line := range lines {
		fields := strings.Fields(line)
		if len(fields) < 7 {
			continue
		}
		if fields[1] == "tmpfs" {
			continue
		}
		if strings.Contains(fields[2], "M") || strings.Contains(fields[2], "K") {
			continue
		}
		if strings.Contains(fields[6], "docker") {
			continue
		}
		isExclude := false
		for _, exclude := range excludes {
			if exclude == fields[6] {
				isExclude = true
			}
		}
		if isExclude {
			continue
		}
		mounts = append(mounts, database.ServerInfoDiskMount{Type: fields[1], Device: fields[0], Mount: fields[6]})
	}

	for i := 0; i < len(mounts); i++ {
		state, err := disk.Usage(mounts[i].Mount)
		if err != nil {
			continue
		}
		var itemData database.ServerInfoDiskInfo
		itemData.Path = mounts[i].Mount
		itemData.Type = mounts[i].Type
		itemData.Device = mounts[i].Device
		itemData.Total = state.Total
		itemData.Free = state.Free
		itemData.Used = state.Used
		itemData.UsedPercent = state.UsedPercent
		itemData.InodesTotal = state.InodesTotal
		itemData.InodesUsed = state.InodesUsed
		itemData.InodesFree = state.InodesFree
		itemData.InodesUsedPercent = state.InodesUsedPercent
		datas = append(datas, itemData)
	}
	return datas
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
