package database

type User struct {
	Id        int32  `json:"id"`
	Email     string `json:"email"`
	Name      string `json:"name"`
	Encrypt   string `json:"encrypt"`
	Password  string `json:"password"`
	Token     string `json:"token"`
	Avatar    string `json:"avatar"`
	CreatedAt int32  `json:"created_at"`
	UpdatedAt int32  `json:"updated_at"`
}

type Server struct {
	Id        int32  `json:"id"`
	Ip        string `json:"ip"`
	Username  string `json:"username"`
	Password  string `json:"password"`
	Port      string `json:"port"`
	Remark    string `json:"remark"`
	State     string `json:"state"`
	Systems   string `json:"systems"`
	CreatedAt int32  `json:"created_at"`
	UpdatedAt int32  `json:"updated_at"`
}

type ServerUser struct {
	Id        int32 `json:"id"`
	ServerId  int32 `json:"server_id"`
	UserId    int32 `json:"user_id"`
	OwnerId   int32 `json:"owner_id"`
	CreatedAt int32 `json:"created_at"`
	UpdatedAt int32 `json:"updated_at"`
}
