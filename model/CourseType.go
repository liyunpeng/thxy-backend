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

func FindAllCourseTypes() (a []*CourseType, err error) {
	//a = new(CourseFile)
	err = db.Debug().Model(&CourseType{}).Select("*").Find(&a).Error
	return
}