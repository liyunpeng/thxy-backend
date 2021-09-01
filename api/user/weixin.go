package user

import (
	"encoding/json"
	"fmt"
	"github.com/parnurzeal/gorequest"
	"strings"
	"thxy/api"
	"thxy/logger"
	"thxy/middleware/session"
	"thxy/model"
	"thxy/redisclient"
	"thxy/setting"
	"thxy/types"
	"thxy/utils"
	"time"

	"github.com/gin-gonic/gin"
	//"github.com/parnurzeal/gorequest"
	//"gitlab.forceup.in/forcepool/userbackend/api"
	//"gitlab.forceup.in/forcepool/userbackend/api/types"
	//"gitlab.forceup.in/forcepool/userbackend/log"
	//"gitlab.forceup.in/forcepool/userbackend/mes"
	//"gitlab.forceup.in/forcepool/userbackend/models"
	//"gitlab.forceup.in/forcepool/userbackend/redisclient"
	//"gitlab.forceup.in/forcepool/userbackend/settings"
	//"gitlab.forceup.in/forcepool/userbackend/utils"
)

//检测邮箱/手机  与 验证码是否对应，注意：验证之后未删除
func CheckVcode(key string, vcode string, c *gin.Context) (bool, error) {
	//如果线上配置文件中不小心写成了off，用户只是收不到验证码，但不至于任何验证码都能随意登录
	//因此这里还是不要加"off"的判断比较安全
	//conf := settings.AppConfig
	//if conf.App.Runmode != settings.RunmodeProd && conf.App.Vcode == "off" {
	//	return true, nil
	//}
	////curTime := utils.CurrentTimestamp()
	//if val, ok := VerCodes[key]; ok {
	//	if val.code == strings.ToUpper(vcode) {
	//		if val.effectivetime < curTime {
	//			return false, errors.New(mes.ByCtx(c, mes.CodeExpire))
	//		}
	//		return true, nil
	//
	//	}
	//	return false, errors.New(mes.ByCtx(c, mes.CodeError))
	//
	//}
	//return false, errors.New(mes.ByCtx(c, mes.CodeNotExist))
	return session.CheckVerCode(key, vcode, c)
}

// 微信用户绑定
func WXBind(c *gin.Context) {
	params := new(types.WXBindParams)
	c.ShouldBindJSON(params)

	key := strings.TrimSpace(params.Key)
	vcode := strings.TrimSpace(params.Vcode)
	openid := strings.TrimSpace(params.Openid)
	if key == "" {
		api.JSONError(c, "缺少参数", "missing key")
		return
	}
	if vcode == "" {
		api.JSONError(c, "缺少参数", "missing vcode")
		return
	}
	if openid == "" {
		api.JSONError(c, "缺少参数", "missing openid")
		return
	}

	// check if vcode match
	if ok, err := CheckVcode(key, vcode, c); !ok {
		api.JSONError(c, err.Error(), err)
		return
	}

	// 检查是否已绑定
	wxUser, err := model.GetWXUserByOpenid(openid)
	if err != nil {
		api.JSONError(c, err.Error(), err)
		return
	}
	if wxUser.UserCode != "" {
		api.JSONError(c, "用户已绑定", "")
		return
	}

	// 检查 forcepool 是否 存在 该用户
	var user *model.User
	//var err error
	var registeredBy string
	// 通过邮箱查询用户
	//if utils.IsValidEmail(key) {
	//	user, err = model.GetUserByEmail(key)
	//	if err != nil {
	//		api.JSONError(c, "搜索用户异常", err)
	//		return
	//	}
	//	registeredBy = "email"
	//} else { // 尝试通过手机查询用户
	//	// 检查虚拟号段支持
	//	if utils.FindMobileSuplier(key) == "" {
	//		api.JSONError(c, "不支持虚拟运营商号段", nil)
	//		return
	//	}
	//	user, err = model.GetUserByPhone(key)
	//	if err != nil {
	//		api.JSONError(c, "搜索用户异常", err)
	//		return
	//	}
	//	registeredBy = "phone"
	//}
	// 不存在 则 创建用户
	if user == nil || user.Id == 0 {
		p := &types.RegisterParams{}
		p.Key = key
		wx_init_pass := utils.GenRandStr(8)
		p.Pwd = wx_init_pass
		u, err := model.RegisterUser(p, registeredBy)
		if err != nil {
			api.JSONError(c, "注册错误", err)
			return
		}
		user = u
	}

	// 建立绑定关系 返回用户 session
	wxUser.UserCode = user.Code
	wxUser.BindTime = types.NewNormalTime(time.Now())
	err = model.UpdateWXUser(wxUser)
	if err != nil {
		api.JSONError(c, err.Error(), nil)
		return
	}
	//err = model.UpdateUserLogTime(user)
	//if err != nil {
	//	logger.Info.Println("更新登录时间失败:%v", err)
	//}
	//logger.Info.Println("[%s][%s]登陆成功", user.Phone, user.Email)

	loginedUser, err := AddUserSessionBySid(user, utils.GenSid())
	if err != nil {
		api.JSONError(c, err.Error(), err)
		return
	}
	redisclient.SetUserSession(loginedUser.SessionId, loginedUser)
	api.JSON(c, "登录成功", loginedUser)
}

// 微信登录接口
func WXLogin(c *gin.Context) {

	params := new(types.WXLoginParams)
	c.ShouldBindJSON(params)

	jscode := strings.TrimSpace(params.JSCode)

	if jscode == "" {
		api.JSONError(c, "缺少参数", "missing jscode")
		return
	}
	// 通过js_code 获取用户 openid
	wxOpenid, err := getOpenid(jscode)
	if err != nil {
		api.JSONError(c, err.Error(), nil)
		return
	}
	// 通过 openid 查询 wx_user 是否存在 是否和 用户绑定
	wxUser, err := model.GetWXUserByOpenid(wxOpenid.Openid)
	if err != nil {
		logger.Warning.Println("GetWXUserByOpenid err=", err)
	}
	if wxUser == nil || wxUser.Id == 0 {
		// 创建 wx_user entry
		wxUser = &model.WXUser{
			Openid:    wxOpenid.Openid,
			Unionid:   wxOpenid.Unionid,
			AvatarUrl: params.AvatarUrl,
			Country:   params.Country,
			Province:  params.Province,
			City:      params.City,
			Language:  params.Language,
			Nickname:  params.Nickname,
			Gender:    params.Gender,
		}
		err := model.CreateWXUser(wxUser)
		if err != nil {
			api.JSONError(c, err.Error(), nil)
			return
		}
	}
	type WxLoginRes struct {
		WxInfo      types.WxOpenid   `json:"wx_info"`
		UserSession *model.Loginuser `json:"user_session"`
	}
	res := WxLoginRes{
		WxInfo: wxOpenid,
	}
	// 如果未绑定 返回成功 但没有 user session 信息
	if wxUser.UserCode == "" {
		api.JSON(c, "登录成功", res)
		return
	}

	// 如果已绑定 返回 用户 session
	user, err := model.GetUserbyCode(wxUser.UserCode)
	if err != nil {
		api.JSONError(c, "登录失败", err)
		return
	}
	if user == nil || user.Id == 0 {
		api.JSONError(c, "搜索用户异常", nil)
		return
	}

	err = model.UpdateUserLogTime(user)
	if err != nil {
		logger.Info.Println("更新登录时间失败:%v", err)
	}

	loginedUser, err := AddUserSessionBySid(user, utils.GenSid())
	if err != nil {
		api.JSONError(c, err.Error(), err)
		return
	}

	redisclient.SetUserSession(loginedUser.SessionId, loginedUser)
	res.UserSession = loginedUser
	api.JSON(c, "登录成功", res)
	return
}

func AddUserSessionBySid(us *model.User, sid string) (u *model.Loginuser, err error) {
	loginUser := &model.Loginuser{
		SessionId: sid,
		UserId:    us.Id,
		UserCode:  us.Code,
		LastTime:  utils.CurrentTimestamp(),
		LogTime:   us.LogTime.String(),
		//Headurl:         headurl,
	}

	return loginUser, nil
}

func WXToken(c *gin.Context) {
	var res types.WxToken

	token, err := getWXToken()
	if err != nil {
		api.JSONError(c, "WXToken 失败", res)
		return
	}
	res.AccessToken = token
	api.JSON(c, "操作成功", res)
}

func getOpenid(jscode string) (types.WxOpenid, error) {

	config := setting.TomlConfig
	var res types.WxOpenid

	var wx_appid = config.WxAppid
	var wx_screct = config.WxSecret
	var wx_gateway = config.WxGateway
	const wx_method = "/sns/jscode2session"
	const grant_type = "authorization_code"

	request := gorequest.New()

	requestUrl := fmt.Sprintf("%v%v?grant_type=%v&appid=%v&secret=%v&js_code=%v",
		wx_gateway, wx_method, grant_type,
		wx_appid, wx_screct, jscode)

	logger.Info.Println(" 微信登录 requestUrl =",  requestUrl)

	resp, body, errs := request.Get(requestUrl).End()
	if len(errs) > 0 {
		logger.Warning.Println(errs)
		return res, fmt.Errorf("%v", errs)
	}
	if resp.StatusCode != 200 {
		return res, fmt.Errorf("%v", body)
	}
	logger.Info.Println(body)

	err := json.Unmarshal([]byte(body), &res)
	if err != nil {
		return res, err
	}
	if res.Errcode != 0 {
		return res, fmt.Errorf("%v", res)
	}

	return res, nil
}

func getWXToken() (string, error) {

	config := setting.TomlConfig
	var res types.WxToken
	if token, e := redisclient.GetWXToken(); e == nil && token != "" {
		res.AccessToken = token
		logger.Info.Println("get token from redis: ", token)
		return token, nil
	} else {
		logger.Warning.Println(token, e)
	}

	var wx_appid = config.WxAppid
	var wx_screct = config.WxSecret
	var wx_gateway = config.WxGateway
	const wx_method = "/cgi-bin/token"
	const grant_type = "client_credential"

	request := gorequest.New()

	requestUrl := fmt.Sprintf("%v%v?grant_type=%v&appid=%v&secret=%v",
		wx_gateway, wx_method, grant_type,
		wx_appid, wx_screct)
	logger.Info.Println("获取微信token的url= ", requestUrl)

	resp, body, errs := request.Get(requestUrl).End()
	if len(errs) > 0 {
		logger.Warning.Println(errs)
		return "", fmt.Errorf("%v", errs)
	}
	if resp.StatusCode != 200 {
		return "", fmt.Errorf("%v", body)
	}
	logger.Info.Println(body)

	err := json.Unmarshal([]byte(body), &res)
	if err != nil {
		return "", err
	}
	if res.Errcode != 0 {
		return "", fmt.Errorf("%v", res)
	}
	err = redisclient.SetWXToken(res.AccessToken)
	if err != nil {
		logger.Warning.Println(err)
	}
	return res.AccessToken, nil
}
