package session

import (
	"errors"
	"github.com/gin-gonic/gin"
	"strings"
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
