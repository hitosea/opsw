package vars

import (
	"github.com/gorilla/websocket"
	"opsw/utils/sshcmd/sshutil"
)

type ConfStruct struct {
	Mode    string
	Host    string
	Port    string
	StartAt string
	DB      string
}

type ExecStruct struct {
	Host      string
	Cmd       string
	Param     string
	Url       string
	LogFile   string
	SSHConfig sshutil.SSH
}

type WorkStruct struct {
	Url   string
	Mode  string
	Token string
	Conf  string
}

type PageStruct struct {
	Page      int   `json:"page"`
	PageSize  int   `json:"page_size"`
	PrevPage  int   `json:"prev_page"`
	NextPage  int   `json:"next_page"`
	PageCount int   `json:"page_count"`
	Data      any   `json:"data"`
	Total     int64 `json:"total"`
}

type WsClientStruct struct {
	Conn *websocket.Conn `json:"conn"`

	Type string `json:"type"` // 客户端类型：user、server
	Cid  int32  `json:"cid"`  // 客户端ID：用户ID、服务器ID
	Rid  string `json:"rid"`  // 客户端随机ID
}

type WsMsgStruct struct {
	Action int `json:"action"` // 消息类型：1、上线；2、下线；3、消息
	Data   any `json:"data"`   // 消息内容

	Type string `json:"type"` // 客户端类型：user、server
	Cid  int32  `json:"cid"`  // 客户端ID：用户ID、服务器ID
	Rid  string `json:"rid"`  // 客户端随机ID
}
