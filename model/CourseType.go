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

func UpdateCourseTypeById( name string , id int ) (err error){
	if err != nil {
		return
	}
	err = db.Debug().Exec("update course_type set name = ? where id = ? ", name, id).Error
	return err
}

func AddCourseType ( c *CourseType ) (err error){
	if err != nil {
		return
	}
	err = db.Debug().Create(c).Error
	return err
}