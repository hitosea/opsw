package vars

type RunStruct struct {
	Mode    string
	Host    string
	Port    string
	StartAt string
}

var (
	RunConf RunStruct
)
