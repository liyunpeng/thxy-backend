package model

import (
	"github.com/jinzhu/gorm"
	"thxy/types"
)

type CourseFile struct {
	Model
	CourseId         int    `json:"course_id" gorm:"default:0"`                     // 所在课程id
	Number           int    `json:"number" gorm:"index:idx_number; comment:'课程编号'"` // 第几节课
	ImgFileName      string `json:"img_file_name"`
	Mp3FileName      string `json:"mp3_file_name" gorm:"size:128"`
	Introduce        string `json:"introduce"`
	Provider         string `json:"provider"`
	GroupId          int    `json:"group_id"`
	Duration         string `json:"duration" gorm:"size:16"`
	WordFilePath     string `json:"word_file_path"`
	ListenedPercent  int    `json:"listened_percent" gorm:"-"`
	ListenedPosition int    `json:"listened_position" gorm:"-"`
}

func (CourseFile) TableName() string {
	return "course_file"
}

func init() {
	autoMigrateModels = append(autoMigrateModels, &CourseFile{})
}

func FindCourseFileById(id int) (courseFile *CourseFile, err error) {
	courseFile = new(CourseFile)
	err = db.Debug().Model(&CourseFile{}).Select("*").Where("id = ? ", id).Find(courseFile).Error
	return
}

func FindCourseFileByCourseId(courseId int) (courseFiles []*CourseFile, err error) {
	err = db.Debug().Model(&CourseFile{}).Select("*").Where("course_id = ? ", courseId).Find(&courseFiles).Error
	return
}

func FindCourseFileByCourseIdAndCourseFileId(courseId, courseFileId int) (a []*CourseFile, err error) {
	err = db.Debug().Model(&CourseFile{}).Select("*").Where("course_id = ? and id >= ? ", courseId, courseFileId).Find(&a).Error
	return
}

func FindCourseFileCountByCourseId(courseId int) (count int, err error) {
	c := new(types.CountType)
	err = db.Debug().Table("course_file").Select("count(*) as count").Where("course_id = ? ", courseId).First(c).Error
	count = c.Count
	return
}

func FindCourseFileListLatest(limit int) (courseFiles []*CourseFile, err error) {
	err = db.Debug().Model(&CourseFile{}).Select("*").Order(" id desc ").Limit(limit).Find(&courseFiles).Error
	return
}

func InsertCourseFile(tx *gorm.DB, courseFile *CourseFile) (err error) {
	err = tx.Debug().Model(&CourseFile{}).Create(courseFile).Error
	return
}
