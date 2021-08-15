package types

type AppConfig struct {
	//App TomlConfig
	Test TomlEnv
}
type TomlServer struct {
	Host       string
	Webhook    string
	SignSecret string
	CronSpec    string
}
type TomlDatabase struct {
	Type        string
	User        string
	Password    string
	Host        string
	Name        string
	RonglaiName string
}

type TomlSession struct {
}

type TomlRedis struct {
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
