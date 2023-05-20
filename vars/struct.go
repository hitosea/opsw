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
}

type WsClientStruct struct {
	Conn *websocket.Conn `json:"conn"`
	Type string          `json:"type"` // 用户类型：user、server
	Uid  int32           `json:"uid"`  // 用户ID、服务器ID
	Rid  string          `json:"rid"`  // 用户随机ID
}

type WsMsgStruct struct {
	Type int `json:"type"`
	Data any `json:"data"`
}
