package types

type ImgBase64 struct {
	File string `json:"file"`
}

type WXLoginParams struct {
	AvatarUrl string `json:"avatar_url"`
	Country   string `json:"country"`
	Province  string `json:"province"`
	City      string `json:"city"`
	Language  string `json:"language"`
	Nickname  string `json:"nickname"`
	JSCode    string `json:"js_code"`
	Gender    int    `json:"gender"`
}

type WXBindParams struct {
	Key    string `json:"key"`    // 绑定的手机号或邮箱
	Vcode  string `json:"vcode"`  // 验证码
	Openid string `json:"openid"` // 微信 openid
}

type WxToken struct {
	AccessToken string `json:"access_token"`
	ExpiresIn   int    `json:"expires_in"`
	Errcode     int    `json:"errcode"`
	Errmsg      string `json:"errmsg"`
}

type WxOpenid struct {
	Openid     string `json:"openid"`
	SessionKey string `json:"session_key"`
	Unionid    string `json:"unionid"`
	Errcode    int    `json:"errcode"`
	Errmsg     string `json:"errmsg"`
}
