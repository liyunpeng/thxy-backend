package model

import (
	"github.com/pkg/errors"
	"thxy/types"
	"thxy/utils"
)

type Admin struct {
	Model
	Account           string           `json:"account" gorm:"size:20;comment:'账户'"`
	Phone             string           `json:"phone" gorm:"size:20;comment:'手机号'"`
	Pwd               string           `json:"pwd" gorm:"size:255;comment:'密码'"`
	Marks             string           `json:"marks" gorm:"size:255;comment:'备注'"`
	LogTime           types.NormalTime `json:"log_time" gorm:"comment:'最近登录时间'"`
	LastChangePwdTime types.NormalTime `json:"last_change_pwd_time" gorm:"type:datetime;comment:'上一次修改password的时间'"`
	Name              string           `json:"name" gorm:"size:255;comment:'管理员姓名'"`
	Roles             []string         `json:"roles" gorm:"-"`
	IsDeleted         bool             `json:"is_deleted" gorm:"default:false"`
}

type AdminUser struct {
	AdminSessionId string   `json:"admin_session_id"`
	AdminId        int      `json:"admin_id"`
	Account        string   `json:"account"`
	Jurisdiction   int      `json:"jurisdiction"`
	LastTime       int64    `json:"last_time"`
	Phone          string   `json:"phone"`
	Roles          []string `json:"roles"`
	IsExpire       bool     `json:"is_expire"`
}

type AdminListItem struct {
	Id        int      `json:"id"`
	Account   string   `json:"account"`
	Name      string   `json:"name"`
	Marks     string   `json:"marks"`
	Phone     string   `json:"phone"`
	GmtCreate string   `json:"gmt_create"`
	Roles     []string `json:"roles"`
}

func (Admin) TableName() string {
	return "admin"
}

func init() {
	autoMigrateModels = append(autoMigrateModels, &Admin{})
}

func GetAdminList() (res types.CommonList, err error) {
	//offset, size := utils.PaginationHelper(params.Page, params.PageSize, params.Begidx, params.Count)
	list := make([]*AdminListItem, 0)
	var total int
	stmt := db.Table("admin").Select("id, account, name, marks, phone, gmt_create")
	err = stmt.Count(&total).Error
	if err != nil {
		err = errors.Wrap(err, "count admin table failed.")
		return
	}
	err = stmt.Order("gmt_create desc").Find(&list).Error
	if err != nil {
		err = errors.Wrap(err, "query admin table failed.")
		return
	}
	res.Total = total
	res.List = list
	return
}

func AddAdmin(admin *Admin) (err error) {
	err = db.Create(admin).Error
	return
}

func UpdatePwd(pwd string) (err error) {
	err = db.Debug().Exec("update admin set pwd = ? ", pwd).Error
	return
}

func UpdateAdmin(admin *Admin) (err error) {
	err = db.Model(admin).Update(admin).Error
	return
}

func DeleteAdminById(id int) (err error) {
	a := new(Admin)
	a.Id = id
	err = db.Delete(a).Error
	return
}

func GetAdminListByIds(ids []int) (al []*Admin, err error) {
	al = make([]*Admin, 0)
	err = db.Model(&Admin{}).Where("id in (?) and is_deleted = false", ids).Find(&al).Error
	return
}

func GetAdminByPhone(phone string) (admin *Admin, err error) {
	admin = new(Admin)
	var list []*Admin
	err = db.Model(&Admin{}).Where("phone = ? and is_deleted = false", phone).Find(&list).Error
	if err == nil && len(list) > 0 {
		admin = list[0]
	}
	return
}

func AdminExistsByPhone(phone string) bool {
	admin, err := GetAdminByPhone(phone)
	return err != nil || admin.Id > 0
}

func AdminExistsByAccount(account string) (r bool) {
	admin, err := GetAdminByAccount(account)
	return err != nil || admin.Id > 0
}

func GetAdminByAccount(account string) (admin *Admin, err error) {
	admin = new(Admin)
	var list []*Admin
	err = db.Model(&Admin{}).Where("account = ? and is_deleted = false", account).Find(&list).Error
	if err == nil && len(list) > 0 {
		admin = list[0]
	}
	return
}

func GetAdminListBySess(asList []AdminSess) ([]*Admin, error) {
	ids := make([]int, len(asList))
	for i, item := range asList {
		ids[i] = item.AdminId
	}
	return GetAdminListByIds(ids)
}

func UpdateAdminPwdByAccount(account, pwd string) (err error) {
	t := utils.CurrentTimestr()
	err = db.Exec("update `admin` set pwd=?, last_change_pwd_time=? where account =? ", pwd, t, account).Error
	return
}

func GetAdminById(id int) (res *Admin, err error) {
	res = new(Admin)
	err = db.Model(res).Where("id = ? and is_deleted = false", id).First(res).Error
	return
}
