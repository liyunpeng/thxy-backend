package session

import (
	"errors"
	"github.com/gin-gonic/gin"
	"strings"
	"thxy/api"
	"thxy/api/admin"
	"thxy/logger"

	"thxy/redisclient"
	"thxy/setting"
	"thxy/utils"
)

func CheckVerCode(account string, vcode string, c *gin.Context) (bool, error) {
	vCode, _ := redisclient.GetVcode(account)
	if vCode == nil {
		return false, errors.New(" GetVcode is null ")
	}
	now := utils.CurrentTimestamp()
	if strings.Contains(account, "@") {
		if vCode.Lasttime > 0 &&
			vCode.Lasttime+setting.EMAIL_CODE_EXPIRE > now &&
			vCode.Code == vcode {
			redisclient.DeleteVcode(account)
			return true, nil
		}
	} else {
		if vCode.Lasttime > 0 &&
			vCode.Lasttime+setting.SMS_CODE_EXPIRE > now &&
			vCode.Code == vcode {
			redisclient.DeleteVcode(account)
			return true, nil
		}
	}
	return false, errors.New(" CheckVerCode error")
}

func CheckAdminSession() gin.HandlerFunc {
	return func(c *gin.Context) {
		sid := c.GetHeader(setting.HTTPSidHeader)
		logger.Info.Printf("check admin session get sid: %v", sid)
		if sid == "" {
			sid = c.PostForm("sid")
		}
		if sid == "" {
			sid = c.Query("sid")
		}
		as, ok := admin.Sessions.QueryloginS(sid)
		if !ok {
			c.Abort()
			api.JSONExpire(c, "登录过期", nil)
		} else {
			c.Set(setting.AdminSessKey, as)
			c.Next()
		}
	}
}

