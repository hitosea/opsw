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

type WsClientStruct struct {
	Conn   *websocket.Conn `json:"conn"`
	UserId int32           `json:"user_id"`
	RandId string          `json:"rand_id"`
}

type WsMsgStruct struct {
	State int `json:"state"`
	Data  any `json:"data"`
}
