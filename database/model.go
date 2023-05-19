package database

type Base struct {
	Id        int32 `json:"id"`
	CreatedAt int32 `json:"created_at"`
	UpdatedAt int32 `json:"updated_at"`
}

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
	Systems  string `json:"systems"`
}

type ServerUser struct {
	Base
	ServerId int32 `json:"server_id"`
	UserId   int32 `json:"user_id"`
	OwnerId  int32 `json:"owner_id"`
}

type ServerList struct {
	Server
	ServerId int32 `json:"server_id"`
	UserId   int32 `json:"user_id"`
	OwnerId  int32 `json:"owner_id"`
}
