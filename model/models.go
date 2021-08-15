package model

import (
	"fmt"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"thxy/logger"
	"thxy/setting"
	"time"
)

var db *gorm.DB

var autoMigrateModels = make([]interface{}, 0)

type Model struct {
	Id          int       `gorm:"primary_key" json:"id"`
	GmtCreate   time.Time `json:"gmt_create"`
	GmtModified time.Time `json:"gmt_modified"`
}

func Setup() {
	conf := setting.TomlConfig
	var err error
	dbConnect := fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8&parseTime=True&loc=Local",
		conf.Test.Database.User,
		setting.TomlConfig.Test.Database.Password,
		setting.TomlConfig.Test.Database.Host,
		setting.TomlConfig.Test.Database.Name)
	db, err = gorm.Open(setting.TomlConfig.Test.Database.Type, dbConnect)
	if err != nil {
		logger.Error.Println("InitDb error:", err)
	} else {
		logger.Debug.Println("InitDb 成功，dbConnect=", dbConnect)
	}

	db.SingularTable(true)
	db.Callback().Create().Replace("gorm:update_time_stamp", updateTimeStampForCreateCallback)
	db.Callback().Update().Replace("gorm:update_time_stamp", updateTimeStampForUpdateCallback)
	db.DB().SetMaxIdleConns(10)
	db.DB().SetMaxOpenConns(100)
	autoMigrate()

}

// CloseDB closes database connection (unnecessary)
func CloseDB() {
	defer func() {
		db.Close()
	}()
}

// updateTimeStampForCreateCallback will set `CreatedOn`, `ModifiedOn` when creating
func updateTimeStampForCreateCallback(scope *gorm.Scope) {
	if !scope.HasError() {
		if createTimeField, ok := scope.FieldByName("GmtCreate"); ok {
			if createTimeField.IsBlank {
				createTimeField.Set(time.Now())
			}
		}

		if modifyTimeField, ok := scope.FieldByName("GmtModified"); ok {
			if modifyTimeField.IsBlank {
				modifyTimeField.Set(time.Now())
			}
		}
	}
}

// updateTimeStampForUpdateCallback will set `ModifiedOn` when updating
func updateTimeStampForUpdateCallback(scope *gorm.Scope) {
	if _, ok := scope.Get("gorm:update_column"); !ok {
		scope.SetColumn("GmtModified", time.Now())
	}
}

// addExtraSpaceIfExist adds a separator
func addExtraSpaceIfExist(str string) string {
	if str != "" {
		return " " + str
	}
	return ""
}

func autoMigrate() {
	if len(autoMigrateModels) > 0 {
		logger.Info.Println("auto migrate 开始")
		// todo redo
		db.Debug().AutoMigrate(autoMigrateModels...)
		logger.Info.Println("auto migrate 结束")
	}
}
