package vars

type ConfStruct struct {
	Mode    string
	Host    string
	Port    string
	StartAt string
	DB      string
}

type UserModel struct {
	ID        uint32 `json:"id"`
	Email     string `json:"email"`
	Name      string `json:"name"`
	Encrypt   string `json:"encrypt"`
	Password  string `json:"password"`
	Token     string `json:"token"`
	CreatedAt uint32 `json:"created_at"`
	UpdatedAt uint32 `json:"updated_at"`
}

var (
	Config ConfStruct
)
