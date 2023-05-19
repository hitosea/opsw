package vars

import "opsw/utils/sshcmd/sshutil"

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

type UserModel struct {
	ID        uint32 `json:"id"`
	Email     string `json:"email"`
	Name      string `json:"name"`
	Encrypt   string `json:"encrypt"`
	Password  string `json:"password"`
	Token     string `json:"token"`
	Avatar    string `json:"avatar"`
	CreatedAt uint32 `json:"created_at"`
	UpdatedAt uint32 `json:"updated_at"`
}

var (
	Config ConfStruct
)
