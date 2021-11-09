package model

type Config struct {
	Model
	//BaseUrl                 string `json:"base_url"`
	ServiceCurrentVersion   string `json:"service_current_version" gorm:"default:0"`
	//Mp3SourceRouter         string `json:"mp3_source_router"`
	CourseTypeUpdateVersion string `json:"course_type_update_version" gorm:"default:0"`
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

func UpdateConfigUpdateVersion() (err error) {
	err = db.Debug().Exec(" update config set course_type_update_version=course_type_update_version+1  ").Error
	return
}