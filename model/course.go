package model

type Course struct {
	Model
	Title     string `json:"title"`
	Introduce string `json:"introduce"`
	Provider  string `json:"provider"`
	TypeId    int    `json:"type_id"`
	CataLevel int    `json:"cata_level"` // 1：一级目录， 点中之后， 直接进文件， 2：存在二级目录， 点中之后进下一级
}

func (Course) TableName() string {
	return "course"
}

func init() {
	autoMigrateModels = append(autoMigrateModels, &Course{})
}

func FindCourseByTypeId(typeId int) (a []*Course, err error) {
	//a = new(Course)
	err = db.Debug().Model(&Course{}).Select("*").Where("type_id = ? ", typeId).Find(&a).Error
	return
}
