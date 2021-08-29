package model

import "github.com/jinzhu/gorm"

type CourseFile struct {
	Model
	CourseId        int     `json:"course_id"`
	Number          int     `json:"number"` // 第几节课
	ImgSrc          string  `json:"img_src"`
	Mp3FileName     string  `json:"mp3_file_name"`
	Introduce       string  `json:"introduce"`
	Provider        string  `json:"provider"`
	GroupId         int     `json:"group_id"`
	WordFilePath    string  `json:"word_file_path"`
	ListenedPercent int `json:"listened_percent" gorm:"-"`
}

func (CourseFile) TableName() string {
	return "course_file"
}

func init() {
	autoMigrateModels = append(autoMigrateModels, &CourseFile{})
}

func FindCourseFileById(id int) (a *CourseFile, err error) {
	a = new(CourseFile)
	err = db.Debug().Model(&CourseFile{}).Select("*").Where("id = ? ", id).Find(a).Error
	return
}

func FindCourseFileByCourseId(courseId int) (a []*CourseFile, err error) {
	err = db.Debug().Model(&CourseFile{}).Select("*").Where("course_id = ? ", courseId).Find(&a).Error
	return
}

func FindCourseFileListLatest(limit int) (a []*CourseFile, err error) {
	err = db.Debug().Model(&CourseFile{}).Select("*").Order(" id desc ").Find(&a).Limit(limit).Error
	return
}

func InsertCourseFile(tx *gorm.DB, c *CourseFile) (err error) {
	err = tx.Debug().Model(&CourseFile{}).Create(c).Error
	return
}
