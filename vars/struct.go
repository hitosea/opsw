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
