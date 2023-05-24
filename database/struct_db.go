package database

// 数据库模型

type User struct {
	Base
	Email    string `json:"email"`
	Name     string `json:"name"`
	Encrypt  string `json:"encrypt"`
	Password string `json:"password"`
	Token    string `json:"token"`
	Avatar   string `json:"avatar"`
}

type Server struct {
	Base
	Ip       string `json:"ip"`
	Username string `json:"username"`
	Password string `json:"password"`
	Port     string `json:"port"`
	Remark   string `json:"remark"`
	State    string `json:"state"`
	Token    string `json:"token"`

	PanelPort     int32  `json:"panel_port"`
	PanelUsername string `json:"panel_username"`
	PanelPassword string `json:"panel_password"`
}

type ServerUser struct {
	Base
	ServerId int32 `json:"server_id"`
	UserId   int32 `json:"user_id"`
	OwnerId  int32 `json:"owner_id"`
}

type ServerInfo struct {
	Base
	ServerId             int32  `json:"server_id"`
	Hostname             string `json:"hostname"`
	OS                   string `json:"os"`
	Platform             string `json:"platform"`
	PlatformFamily       string `json:"platform_family"`
	PlatformVersion      string `json:"platform_version"`
	KernelArch           string `json:"kernel_arch"`
	KernelVersion        string `json:"kernel_version"`
	VirtualizationSystem string `json:"virtualization_system"`

	CPUCores        int    `json:"cpu_cores"`
	CPULogicalCores int    `json:"cpu_logical_cores"`
	CPUModelName    string `json:"cpu_model_name"`

	CurrentInfo string `json:"current_info"`

	Version   string `json:"version"`
	CommitSHA string `json:"commit_sha"`
}
