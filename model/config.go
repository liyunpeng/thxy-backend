package model

type Config struct {
	Model
	BaseUrl string `json:"base_url"`
	Mp3SourceRouter string `json:"mp3_source_router"`
}


func (Config) TableName() string {
	return "config"
}

func init() {
	autoMigrateModels = append(autoMigrateModels, &Config{})
}

func FindConfig() (a *Config, err error) {
	a = new(Config)

	err = db.Debug().Model(&Config{}).Select("*").First(a).Error
	return
}