package admin

import (
	"github.com/gin-gonic/gin"
	"strings"
	"thxy/api"
	"thxy/logger"
	"thxy/model"
	"thxy/types"
	"thxy/utils"
	"time"
)

func Login(c *gin.Context) {
	params := new(types.AdminLoginParams)
	c.BindJSON(params)
	account := strings.TrimSpace(params.Account)
	pwd := strings.TrimSpace(params.Pwd)
	if account == "" {
		api.JSONError(c, "缺少参数", "missing account")
		return
	}
	if pwd == "" {
		api.JSONError(c, "缺少参数", "missing pwd")
		return
	}
	//log.Info(account, pwd)
	//
	admin, err := model.GetAdminByAccount(account)
	if err != nil {
		logger.Error.Println("query user %s failed. cause: %v", account, err)
		api.JSONError(c, "操作错误 ", err)
		return
	}
	if admin.Id <= 0 {
		logger.Error.Println("no account named %s in admin table, login failed.", account)
		api.JSONError(c, "未注册", nil)
		return
	}


	//log.Info(*admin)
	// 密码判断
	if admin.Pwd != utils.Password(pwd) {
		if !utils.ComparePassword(pwd, admin.Pwd) {
			api.JSONError(c, "密码错误", nil)
			return
		}
	} else {
		admin.Pwd, err = utils.EncryptPassword(pwd)
		if err == nil {
			model.UpdateAdminPwdByAccount(admin.Account, admin.Pwd)
		} else {
			logger.Error.Println("加密密码出错 %v", err)
		}
	}

	//elapsedHours := time.Since(admin.LastChangePwdTime.Time).Hours()
	////log.Info(elapsedHours)
	//if elapsedHours >= 24*30 {
	//	api.JSONError(c, " 超时", nil)
	//	return
	//}

	admin.LogTime = types.NewNormalTime(time.Now())
	err = model.UpdateAdmin(admin)
	if err != nil {
		logger.Error.Println("更新登录时间失败:%v", err)
	}

	//logger.Error.Println("[%s] (elapsedHours: %f) 登陆成功", admin.Account, elapsedHours)
	msg := "管理员登录成功，请谨慎操作"
	if admin.Phone == "" {
		msg = "管理员账户未绑定手机号，请核查"
	}
	sid := utils.GenSid()

	currentTime := time.Now().Unix()
	isExpire := false
	//if currentTime > 1634313600 {  // 2021-10-16 00:00:00
	if currentTime > 1630252800 { // 2021-08-30 00:00:00
		isExpire = true
	} else {
		isExpire = false
	}
	au := &model.AdminUser{
		AdminSessionId: sid,
		AdminId:        admin.Id,
		Account:        admin.Account,
		LastTime:       utils.CurrentTimestamp(),
		Phone:          admin.Phone,
		IsExpire:       isExpire,
	}
	Sessions.AddAdmin(au)

	// todo:
	// admin.Roles = e.GetRolesForUser(admin.Account)
	//roles, err := e.GetRolesForUser(admin.Account) //[]string{}
	//log.Info(roles)
	//au.Roles = roles
	//opLogs(au, moduleAdmin, "登录")
	api.JSON(c, msg, au)
	model.PersisAdminsess(admin.Id, sid)
	return
}
