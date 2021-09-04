package types

import "time"

const (
	ResErrorCode    = 5
	ResSuccessCode  = 3
	ResExpireCode   = 6
	ResMaintainCode = 7
)

type TomlSession struct {
	Type string
	Path string
	Life int
}

type AppConfig struct {
	//App TomlConfig
	WxAppid   string
	WxSecret  string
	WxGateway string
	Test      TomlEnv
	Session   TomlSession
}
type TomlServer struct {
	Host         string
	Webhook      string
	SignSecret   string
	CronSpec     string
	FileDownload string
}
type TomlDatabase struct {
	Type        string
	User        string
	Password    string
	Host        string
	Name        string
	RonglaiName string
}

// AdminLoginParams - /forcepool/adminLogin
type AdminLoginParams struct {
	Account string `json:"account"`
	Pwd     string `json:"pwd"`
}

type TomlRedis struct {
	Type        string
	User        string
	Password    string
	Host        string
	Name        string
	RonglaiName string
}

type TomlFileStore struct {
	FileStorePath string
}

type TomlEnv struct {
	FilStore TomlFileStore
	Server   TomlServer
	Database TomlDatabase
	Session  TomlSession
	Redis    TomlRedis
}

type KeyPassParams struct {
	Key   string `json:"key"`
	Pwd   string `json:"pwd"`
	Vcode string `json:"vcode"`
}

type RegisterParams struct {
	KeyPassParams
	ConfirmPwd string `json:"confirm_pwd"`
}

type NormalTime struct {
	time.Time
}

func NewNormalTime(t time.Time) NormalTime {
	nt := NormalTime{}
	nt.Time = t
	return nt
}
