package session

import (
	"encoding/hex"
	"encoding/json"
	"strings"

	"gitlab.forceup.in/forcepool/userbackend/api"
	"gitlab.forceup.in/forcepool/userbackend/log"
	"gitlab.forceup.in/forcepool/userbackend/mes"
	"gitlab.forceup.in/forcepool/userbackend/models"
	"gitlab.forceup.in/forcepool/userbackend/redisclient"
	"gitlab.forceup.in/forcepool/userbackend/settings"
	"gitlab.forceup.in/forcepool/userbackend/utils"

	"errors"

	"github.com/gin-gonic/gin"
)

func CheckUserSession() gin.HandlerFunc {
	return func(c *gin.Context) {
		// sid := c.GetHeader(settings.HTTPSidHeader)
		// //logging.Infof("check admin session get sid: %v", sid)
		// if sid == "" {
		// 	sid = c.PostForm("sid")
		// }
		// if sid == "" {
		// 	sid = c.Query("sid")
		// }
		// us, ok := user.UserSessions.QueryloginS(sid)
		logger := log.GetLogger()
		session := c.GetHeader(settings.XUSER)
		valid := true
		userInfo := new(models.Loginuser)
		if session == "" {
			logger.Warn("### empty X-forcepool-user abort ### ")
			valid = false
			// 荣来测试环境也是在nginx中拦截session并设置到redis中
			if settings.AppConfig.App.Station == "ronglai" {
				sid := c.GetHeader(settings.HTTPSidHeader)
				if sid != "" {
					if res, err := redisclient.GetUserSession(sid); err == nil {
						valid = true
						//fmt.Printf("%#v\n", res)
						userInfo = res
					}
				}
			} else {
				// 开发或测试环境 鉴权失败可能是没配置nginx 尝试直接从redis拿用户信息
				if settings.AppConfig.App.Runmode != settings.RunmodeProd {
					sid := c.GetHeader(settings.HTTPSidHeader)
					if sid != "" {
						if res, err := redisclient.GetUserSession(sid); err == nil {
							valid = true
							//fmt.Printf("%#v\n", res)
							userInfo = res
						}
					}
				}
			}
		} else {
			payload, err := hex.DecodeString(session)
			if err != nil {
				logger.Warnf("### can not decode session, cause: %v ###", err)
				valid = false
			}

			data, err := utils.Decrypt([]byte(utils.SecretKey), payload)
			if err != nil {
				logger.Warnf("### can not decrypt session payload, cause: %v ###", err)
				valid = false
			}

			if err := json.Unmarshal(data, &userInfo); err != nil {
				logger.Warnf("### can not unmarshal session payload, cause: %v", err)
				valid = false
			}
		}
		if valid == false {
			c.Abort()
			api.JSONExpire(c, mes.ByCtx(c, mes.Expire), nil)
		} else {
			// 判断用户的状态，只有等于0的时候才能设置session
			if user, err := models.GetUserbyCode(userInfo.UserCode); err == nil && user.Status == 0 {
				c.Set(settings.UserSessKey, userInfo)
				c.Next()
			} else {
				c.Abort()
				api.JSONExpire(c, mes.ByCtx(c, mes.Expire), nil)
			}
		}
	}
}

func CheckUserLimited(usercode string, account string) (bool, string) {
	userLimited, err := redisclient.GetUserLimited(usercode)
	if err != nil && userLimited == nil {
		return false, "TryAgain"
	}
	now := utils.CurrentTimestamp()
	if strings.Contains(account, "@") {
		if userLimited.LastSentEmailTime > 0 &&
			now <= userLimited.LastSentEmailTime+60 {
			return false, "TryAgain"
		} else if userLimited.FirstSentEmailTime <= 0 {
			userLimited.FirstSentEmailTime = now
			userLimited.SendCountEmail = 0
		}
		//redis中有值，但第一次发送的时间距现在已经超过24小时
		if userLimited.FirstSentEmailTime+3600*24 <= now {
			userLimited.FirstSentEmailTime = now
			userLimited.SendCountEmail = 1
		} else if userLimited.SendCountEmail >= settings.EMAIL_COUNT_LIMITED {
			return false, "TooMany"
		} else {
			userLimited.SendCountEmail++
		}
		userLimited.LastSentEmailTime = now
	} else {
		if userLimited.LastSentSmsTime > 0 &&
			now <= userLimited.LastSentSmsTime+60 &&
			userLimited.LastSentSmsTime != userLimited.FirstSentSmsTime { //用户若仅修改手机号码一次，不做60秒的限制
			return false, "TryAgain"
		} else if userLimited.FirstSentSmsTime <= 0 {
			userLimited.FirstSentSmsTime = now
			userLimited.SendCountSms = 0
		}
		//redis中有值，但第一次发送的时间距现在已经超过24小时
		if userLimited.FirstSentSmsTime+3600*24 <= now {
			userLimited.FirstSentSmsTime = now
			userLimited.SendCountSms = 1
		} else if userLimited.SendCountSms >= settings.SMS_COUNT_LIMITED {
			return false, "TooMany"
		} else {
			userLimited.SendCountSms++
		}
		userLimited.LastSentSmsTime = now
	}

	err = redisclient.SetUserLimited(usercode, userLimited)
	if err != nil {
		return false, "TryAgain"
	}

	return true, ""
}

func CheckVerCode(account string, vcode string, c *gin.Context) (bool, error) {
	vCode, _ := redisclient.GetVcode(account)
	if vCode == nil {
		return false, errors.New(mes.ByCtx(c, mes.CodeNotExist))
	}
	now := utils.CurrentTimestamp()
	if strings.Contains(account, "@") {
		if vCode.Lasttime > 0 &&
			vCode.Lasttime+settings.EMAIL_CODE_EXPIRE > now &&
			vCode.Code == vcode {
			redisclient.DeleteVcode(account)
			return true, nil
		}
	} else {
		if vCode.Lasttime > 0 &&
			vCode.Lasttime+settings.SMS_CODE_EXPIRE > now &&
			vCode.Code == vcode {
			redisclient.DeleteVcode(account)
			return true, nil
		}
	}
	return false, errors.New(mes.ByCtx(c, mes.CodeError))
}

func CheckUserAuthLimited(usercode string) (bool, string) {
	userLimited, err := redisclient.GetUserLimited(usercode)
	if err != nil && userLimited == nil {
		return false, "TryAgain"
	}

	now := utils.CurrentTimestamp()

	if userLimited.LastSentAuthTime > 0 &&
		now <= userLimited.LastSentAuthTime+60 {
		return false, "TryAgain"
	} else if userLimited.FirstSentAuthTime <= 0 {
		userLimited.FirstSentAuthTime = now
		userLimited.SendCountAuth = 0
	}
	//redis中有值，但第一次发送的时间距现在已经超过24小时
	if userLimited.FirstSentAuthTime+3600*24 <= now {
		userLimited.FirstSentAuthTime = now
		userLimited.SendCountAuth = 1
	} else if userLimited.SendCountAuth >= settings.AUTH_COUNT_LIMITED {
		return false, "TooMany"
	} else {
		userLimited.SendCountAuth++
	}
	userLimited.LastSentAuthTime = now

	err = redisclient.SetUserLimited(usercode, userLimited)
	if err != nil {
		return false, "TryAgain"
	}

	return true, ""
}
