package database

// 扩展模型

type Base struct {
	Id        int32 `json:"id"`
	CreatedAt int32 `json:"created_at"`
	UpdatedAt int32 `json:"updated_at"`
}

type ServerItem struct {
	Server

	ServerId int32 `json:"server_id"`
	UserId   int32 `json:"user_id"`
	OwnerId  int32 `json:"owner_id"`

	Hostname        string `json:"hostname"`
	Platform        string `json:"platform"`
	PlatformVersion string `json:"platform_version"`
	Version         string `json:"version"`
	CurrentInfo     string `json:"current_info"`

	Upgrade string `json:"upgrade"`
}

type ServerInfoCurrent struct {
	Uptime          uint64 `json:"uptime"`
	TimeSinceUptime string `json:"time_since_uptime"`

	Procs uint64 `json:"procs"`

	Load1            float64 `json:"load1"`
	Load5            float64 `json:"load5"`
	Load15           float64 `json:"load15"`
	LoadUsagePercent float64 `json:"load_usage_percent"`

	CPUPercent     []float64 `json:"cpu_percent"`
	CPUUsedPercent float64   `json:"cpu_used_percent"`
	CPUUsed        float64   `json:"cpu_used"`
	CPUTotal       int       `json:"cpu_total"`

	MemoryTotal       uint64  `json:"memory_total"`
	MemoryAvailable   uint64  `json:"memory_available"`
	MemoryUsed        uint64  `json:"memory_used"`
	MemoryUsedPercent float64 `json:"memory_used_percent"`

	IOReadBytes  uint64 `json:"io_read_bytes"`
	IOWriteBytes uint64 `json:"io_write_bytes"`
	IOCount      uint64 `json:"io_count"`
	IOReadTime   uint64 `json:"io_read_time"`
	IOWriteTime  uint64 `json:"io_write_time"`

	DiskData []ServerInfoDiskInfo `json:"disk_data"`

	NetBytesSent uint64 `json:"net_bytes_sent"`
	NetBytesRecv uint64 `json:"net_bytes_recv"`

	ShotTime string `json:"shot_time"`
}

type ServerInfoDiskInfo struct {
	Path        string  `json:"path"`
	Type        string  `json:"type"`
	Device      string  `json:"device"`
	Total       uint64  `json:"total"`
	Free        uint64  `json:"free"`
	Used        uint64  `json:"used"`
	UsedPercent float64 `json:"used_percent"`

	InodesTotal       uint64  `json:"inodes_total"`
	InodesUsed        uint64  `json:"inodes_used"`
	InodesFree        uint64  `json:"inodes_free"`
	InodesUsedPercent float64 `json:"inodes_used_percent"`
}

type ServerInfoDiskMount struct {
	Type   string
	Mount  string
	Device string
}
