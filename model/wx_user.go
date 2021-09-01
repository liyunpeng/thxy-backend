package model

import (
	"fmt"
	"thxy/types"
)

type WXUser struct {
	Model
	Openid    string           `json:"openid" gorm:"unique_index:idx_openid;size:255;comment:''"`
	Unionid   string           `json:"unionid" gorm:"size:255;comment:''"`
	AvatarUrl string           `json:"avatar_url" gorm:"size:255; comment:'头像地址'"`
	Country   string           `json:"country" gorm:"size:255;comment:'国家'"`
	Province  string           `json:"province" gorm:"size:255;comment:'省'"`
	City      string           `json:"city" gorm:"size:255;comment:'市'"`
	Language  string           `json:"language" gorm:"size:255;comment:'语言'"`
	Nickname  string           `json:"nickname" gorm:"size:255;comment:'国家'"`
	UserCode  string           `json:"user_code" gorm:"size:40;comment:'用户编码'"`
	Gender    int              `json:"gender" gorm:"comment:'性别'"`
	BindTime  types.NormalTime `json:"bind_time" gorm:"comment: '绑定时间'"`
}

func init() {
	autoMigrateModels = append(autoMigrateModels, &WXUser{})
}

func (WXUser) TableName() string {
	return "wx_user"
}

func GetWXUserByOpenid(openid string) (v *WXUser, err error) {
	v = new(WXUser)
	err = db.Model(v).Where("openid = ?", openid).First(v).Error
	return
}

func CreateWXUser(usr *WXUser) (err error) {
	if usr == nil {
		return fmt.Errorf("create wx user: missing usr info")
	}
	err = db.Create(usr).Error
	return
}

func UpdateWXUser(usr *WXUser) (err error) {
	if usr == nil || usr.Id == 0 {
		return fmt.Errorf("update wx user: missing usr info")
	}

	err = db.Model(usr).Updates(usr).Error
	return
}
