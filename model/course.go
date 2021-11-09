package model

import "github.com/jinzhu/gorm"

type Course struct {
	Model
	Title         string `json:"title" `
	Introduction  string `json:"introduction"`
	Provider      string `json:"provider"`
	ImgFileName   string `json:"img_file_name"`
	StorePath     string `json:"store_path"`
	TypeId        int    `json:"type_id"`
	UpdateVersion int    `json:"update_version"  gorm:"default:0"`
	CateLevel     int    `json:"cate_level"` // 1：一级目录， 点中之后， 直接进文件， 2：存在二级目录， 点中之后进下一级
}

func (Course) TableName() string {
	return "course"
}

func init() {
	autoMigrateModels = append(autoMigrateModels, &Course{})
}

type CourseTitleItem struct {
	Id     int    `gorm:"json:"id"`
	TypeId int    `json:"type_id"`
	Title  string `json:"title"`
}

type CourseTypeItem struct {
	TypeId int    `json:"type_id"`
	Name   string `json:"name"`
}

func GetAllCourseIds() (a []*CourseTitleItem, err error) {
	err = db.Debug().Raw("select id, type_id, title from course").Find(&a).Error
	return
}

func GetAllCourseGroup() (courseTypeItems []*CourseTypeItem, err error) {
	sql := "select  c.type_id, t.name from course c inner join course_type t on c.type_id = t.id group by c.type_id"
	err = db.Debug().Raw(sql).Find(&courseTypeItems).Error
	return
}

func InsertCourse(c *Course) (err error) {
	err = db.Debug().Create(c).Error
	return
}

func InsertCourseWithTransaction(tx *gorm.DB, c *Course) (err error) {
	err = tx.Debug().Create(c).Error
	return
}

func UpdateCourse1(c *Course) (err error) {
	err = db.Debug().Table("course ").Update(c).Error
	return
}

func UpdateCourse(title string, introduction string, id int) (err error) {
	err = db.Debug().Exec(" update course set title= ? , introduction = ?  where id = ? ", title, introduction, id).Error
	return
}

func UpdateCourseFileCount(id, fileCount int) (err error) {
	err = db.Debug().Exec(" update course set file_count = ?  where id = ? ", fileCount, id).Error
	return
}

func UpdateCourseUpdateVersion(id int) (err error) {
	err = db.Debug().Exec(" update course set update_version=update_version+1  where id = ? ", id).Error
	return
}

func DeleteCourse(id int) (err error) {
	err = db.Debug().Exec("delete from course where id = ? ", id).Error
	return
}

func DeleteCourseType(id int) (err error) {
	err = db.Debug().Exec("delete from course_type where id = ? ", id).Error
	return
}

func FindCourseByTypeId(typeId int) (courses []*Course, err error) {
	err = db.Debug().Model(&Course{}).Select("*").Where("type_id = ? ", typeId).Find(&courses).Error
	return
}

func FindCourseById(courseId int) (course *Course, err error) {
	course = new(Course)
	err = db.Debug().Model(&Course{}).Select("*").Where("id = ? ", courseId).First(course).Error
	return
}
