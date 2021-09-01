package model

import (
	"fmt"
	"github.com/jinzhu/gorm"
	"thxy/logger"
	"thxy/types"
	"thxy/utils"
	"time"
)

type User struct {
	Model
	Code              string           `json:"code" gorm:"size:40;comment:'用户编码'"`
	Account           string           `json:"account" gorm:"size:40;comment:'用户账号'"`
	Phone             string           `json:"phone" gorm:"size:11;comment:'手机号'"`
	Pwd               string           `json:"pwd" gorm:"size:255;comment: '密码'"`
	LogTime           types.NormalTime `json:"log_time" gorm:"comment: '最近登录时间'"`
	ListenedCourseIds string           `json:"listened_course_ids" gorm:"size:256;comment:''"`
	RegisterTime      types.NormalTime `json:"register_time" gorm:"comment: '注册时间'"`
}

type Loginuser struct {
	SessionId            string `json:"session_id"`
	UserId               int    `json:"user_id"`
	UserCode             string `json:"user_code"`
	Headurl              string `json:"headurl"`
	Phone                string `json:"phone"`
	Email                string `json:"email"`
	Country              string `json:"country"`
	CountryCode          string `json:"country_code"`
	LastTime             int64  `json:"last_time"`
	LogTime              string `json:"log_time"`
	Authentication       int    `json:"authentication"`
	AuthenticationRemark string `json:"authentication_remark"`
}
type Vcode_t struct {
	Code     string
	Lasttime int64 //发送时间
}

func (User) TableName() string {
	return "user"
}

func init() {
	autoMigrateModels = append(autoMigrateModels, &User{})
}

func InsertUser(tx *gorm.DB, u *User) (err error) {
	err = tx.Debug().Model(&User{}).Create(u).Error
	return
}

func GetUserbyCode(code string) (v *User, err error) {
	v = new(User)
	err = db.Model(v).Where("code = ?", code).First(v).Error
	return
}

func UpdateUserLogTime(user *User) (err error) {
	if user.Id <= 0 {
		err = fmt.Errorf("try to update user login time without user id")
		return
	}
	err = db.Model(user).Update("log_time", types.NewNormalTime(time.Now())).Error
	return err
}

func RegisterUser(params *types.RegisterParams, byType string) (u *User, err error) {

	// 生成user_code 确保其唯一性
	userCode := utils.GenUserCode()
	for {
		// 查询 user 表中是否有该userCode
		var c int
		err = db.Table("user").Where("code = ?", userCode).Count(&c).Error
		if err != nil {
			return
		}
		// 如果没有 跳出循环 使用上面生成的 userCode
		if c == 0 {
			break
		}
		userCode = utils.GenUserCode()
	}

	pwd, err := utils.EncryptPassword(params.Pwd)
	if err != nil {
		logger.Error.Println("设置[%s]密码出错: %v", params.Key, err)
		return nil, err
	}
	//pwd := utils.Password(params.Pwd)
	now := types.NewNormalTime(time.Now())
	user := &User{
		Code:         userCode,
		Pwd:          pwd,
		LogTime:      now,
		RegisterTime: now,
	}

	//if byType == "email" {
	//	user.Email = params.Key
	//	user.Verification = 2
	//} else {
	//	user.Phone = params.Key
	//	user.Verification = 1
	//}
	tx := db.Begin()
	// 插入用户记录
	err = tx.Create(user).Error
	if err != nil {
		tx.Rollback()
		return
	}
	u = user
	// 为用创建balance记录 coin_type 1 coin_type 2 fil
	//b1 := &Balance{
	//	Code:     userCode,
	//	CoinType: 1,
	//}
	//err = tx.Create(b1).Error
	//if err != nil {
	//	tx.Rollback()
	//	return
	//}
	//b2 := &Balance{
	//	Code:     userCode,
	//	CoinType: 2,
	//}
	//err = tx.Create(b2).Error
	//if err != nil {
	//	tx.Rollback()
	//	return
	//}
	err = tx.Commit().Error
	return
}
