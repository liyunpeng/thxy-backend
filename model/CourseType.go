package model

type CourseType struct {
	Model
	Name     string `json:"name"`
	//Introduce string `json:"introduce"`
	//Provider  string `json:"provider"`
}

func (CourseType) TableName() string {
	return "course_type"
}

func init() {
	autoMigrateModels = append(autoMigrateModels, &CourseType{})
}

func FindAllCourseTypes() (courseTypes []*CourseType, err error) {
	err = db.Debug().Model(&CourseType{}).Select("*").Find(&courseTypes).Error
	return
}