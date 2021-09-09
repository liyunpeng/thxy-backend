package model

import "github.com/jinzhu/gorm"

type UserListenedCourseFile struct {
	Model
	Code                     string `json:"code" gorm:"index:idx_code; size:40;comment:'用户编码'"`
	CourseId                 int    `json:"course_id" gorm:"index:idx_course_id; size:40;comment:'课程id'"`
	LastListenedCourseFileId int    `json:"last_listened_course_file_id" gorm:"size:40;comment:'上次听过的课程文件id'"`
	ListenedFiles            string `json:"listened_files" gorm:"size:1024;comment:'已听文件的序列化字符串'"`
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

func FindUserListenedCourseByUserCodeAndCourseId(code string, courseId int) (a []*UserListenedCourseFile, err error) {
	err = db.Debug().Model(&UserListenedCourseFile{}).Select("*").Where(" code = ? and course_id = ? ", code, courseId).Find(&a).Error
	return
}

func UpdateUserListenedCourseByUserCodeAndCourseId(listenedFiles string, code string, courseId, lsstListenedFileId int) (err error) {
	err = db.Debug().Exec(" update user_listened_course set listened_files = ?, last_listened_course_file_id = ?  where code = ? and course_id = ? ",
		listenedFiles, lsstListenedFileId,  code, courseId).Error
	return
}
