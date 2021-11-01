package model

type Log struct {
	Model

	FileName      string `json:"file_name"`      // log文件名
	AppVersion    string `json:"app_version"`    //  软件版本
	SystemVersion string `json:"system_version"` // 手机系统版本号
	Brand         string `json:"brand"`          // 手机厂商
	ModelVersion  string `json:"model_version"`  // 手机型号
}

func (Log) TableName() string {
	return "log"
}

func init() {
	autoMigrateModels = append(autoMigrateModels, &Log{})
}

func GetLogList() (log []Log, err error) {
	err = db.Model(&Log{}).Select("*").Find(&log).Error
	return
}

func InsertLog(c *Log) (err error) {
	err = db.Debug().Create(c).Error
	return
}
