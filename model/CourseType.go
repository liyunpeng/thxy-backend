package model

type CourseType struct {
	Model
	Name                string `json:"name"`
	CourseUpdateVersion int    `json:"course_update_version"`
	//Introduce string `json:"introduce"`
	//Provider  string `json:"provider"`
}

func (CourseType) TableName() string {
	return "course_type"
}

func init() {
	autoMigrateModels = append(autoMigrateModels, &CourseType{})
}

func FindCourseTypeById(typeId int) (courseType *CourseType, err error) {
	courseType = new(CourseType)
	err = db.Debug().Model(&CourseType{}).Select("*").Where("id = ? ", typeId).First(courseType).Error
	return
}

func FindAllCourseTypes() (courseTypes []*CourseType, err error) {
	err = db.Debug().Model(&CourseType{}).Select("*").Find(&courseTypes).Error
	return
}

func UpdateCourseTypeById(name string, id int) (err error) {
	err = db.Debug().Exec("update course_type set name = ? where id = ? ", name, id).Error
	if err != nil {
		return
	}
	return err
}

func UpdateCourseTypeUpdateVersionById(id int) (err error) {
	err = db.Debug().Exec("update course_type set course_update_version=course_update_version+1 where id = ? ", id).Error
	if err != nil {
		return
	}
	return err
}

func AddCourseType(c *CourseType) (err error) {
	if err != nil {
		return
	}
	err = db.Debug().Create(c).Error
	return err
}
