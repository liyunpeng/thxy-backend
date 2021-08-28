package model

import "github.com/jinzhu/gorm"

type User struct {
	Model
	Code    string `json:"code" gorm:"size:40;comment:'用户编码'"`
	Account string `json:"account" gorm:"size:40;comment:'用户账号'"`
	Phone   string `json:"phone" gorm:"size:11;comment:'手机号'"`
	ListenedCourseIds string `json:"listened_course_ids" gorm:"size:256;comment:''"`
}

func (User) TableName() string {
	return "user"
}

func init() {
	autoMigrateModels = append(autoMigrateModels, &User{})
}

func InsertUser( tx *gorm.DB, u *User) (err error) {
	err = tx.Debug().Model(&User{}).Create(u).Error
	return
}