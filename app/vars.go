package app

type RunModel struct {
	Mode string
	Host string
	Port string
}

var (
	RunConf RunModel
)
