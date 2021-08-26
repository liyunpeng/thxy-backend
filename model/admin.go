package model

import "github.com/jinzhu/gorm"

type Admin struct {
	Model
	Account string `json:"account" gorm:"size:20;comment:'账户'"`
	Phone   string `json:"phone" gorm:"size:20;comment:'手机号'"`
	Pwd     string `json:"pwd" gorm:"size:255;comment:'密码'"`
	Marks   string `json:"marks" gorm:"size:255;comment:'备注'"`
	//LogTime           types.NormalTime `json:"log_time" gorm:"comment:'最近登录时间'"`
	//LastChangePwdTime types.NormalTime `json:"last_change_pwd_time" gorm:"type:datetime;comment:'上一次修改password的时间'"`
	Name      string   `json:"name" gorm:"size:255;comment:'管理员姓名'"`
	Roles     []string `json:"roles" gorm:"-"`
	IsDeleted bool     `json:"is_deleted" gorm:"default:false"`
}

func (Admin) TableName() string {
	return "admin"
}

func init() {
	autoMigrateModels = append(autoMigrateModels, &Admin{})
}

func InsertAdmin(tx *gorm.DB, u *Admin) (err error) {
	err = tx.Debug().Model(&Admin{}).Create(u).Error
	return
}
