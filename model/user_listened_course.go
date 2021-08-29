package model

import "github.com/jinzhu/gorm"

type UserListenedCourseFile struct {
	Model
	//Openid    string           `json:"openid" gorm:"unique_index:idx_openid;size:255;comment:''"`
	Code string `json:"code" gorm:"index:idx_code; size:40;comment:'用户编码'"`
	//Account string `json:"account" gorm:"size:40;comment:'用户账号'"`
	//Phone   string `json:"phone" gorm:"size:11;comment:'手机号'"`
	CourseId      int    `json:"course_id gorm:"index:idx_course_id; size:40;comment:'用户编码'"`
	ListenedFiles string `json:"listened_files" gorm:"size:1024;comment:'用户账号'"`
}

func (UserListenedCourseFile) TableName() string {
	return "user_listened_course"
}

func init() {
	autoMigrateModels = append(autoMigrateModels, &UserListenedCourseFile{})
}

func InsertUserListenedCourse(tx *gorm.DB, u *UserListenedCourseFile) (err error) {
	err = tx.Debug().Model(&UserListenedCourseFile{}).Create(u).Error
	return
}

func FindUserListenedCourseByUserCodeAndCourseId(code string, courseId int)  (  a []*UserListenedCourseFile, err error) {
	err = db.Model(&UserListenedCourseFile{}).Select("*").Where(" code = ? and course_id = ? ", code, courseId).Find(&a).Error
	return
}

func UpdateUserListenedCourseByUserCodeAndCourseId( listenedFiles string , code string, courseId int)  ( err error) {
	err = db.Debug().Exec(" update user_listened_course set listened_files = ?  where code = ? and course_id = ? ",
		listenedFiles, code, courseId).Error
	return
}
