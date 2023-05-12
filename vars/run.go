package vars

import "time"

type RunStruct struct {
	Mode    string
	Host    string
	Port    string
	StartAt string
}

type UserStruct struct {
	Email     string    `json:"email"`
	Name      string    `json:"name"`
	Token     string    `json:"token"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

var (
	RunConf RunStruct
)
