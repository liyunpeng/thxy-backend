package model

type Config struct {
	Model
	BaseUrl               string `json:"base_url"`
	ServiceCurrentVersion string `json:"service_current_version"`
	Mp3SourceRouter       string `json:"mp3_source_router"`
}

func (Config) TableName() string {
	return "config"
}

func init() {
	autoMigrateModels = append(autoMigrateModels, &Config{})
}

func FindConfig() (config *Config, err error) {
	config = new(Config)

	err = db.Debug().Model(&Config{}).Select("*").First(config).Error
	return
}