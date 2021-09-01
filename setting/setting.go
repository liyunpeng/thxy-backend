package setting

import (
	"fmt"
	"github.com/jinzhu/configor"
	"thxy/logger"
	"thxy/types"
)

const SID_PREFIX = "sid"

const USERCODE_PREFIX = "usercode"
const (
	EMAIL_COUNT_LIMITED = 20
	SMS_COUNT_LIMITED   = 20
	AUTH_COUNT_LIMITED  = 3
	EMAIL_CODE_EXPIRE   = 60 * 5
	SMS_CODE_EXPIRE     = 60 * 5
)
const VERCODE_PREFIX = "vcode"

const SIDUSERCODE = "sucode"
const AUTHIMG_PREFIX = "at"

const WXTOKEN_KEY = "wxtoken"

var TomlConfig *types.AppConfig

func InitConfig(fileName string) {
	TomlConfig = new(types.AppConfig)
	err := configor.Load(TomlConfig, fileName)
	if err != nil {
		panic(fmt.Sprintf("fail to load app config:\n %v\n", err))
	}
	logger.Info.Println(" 配置加载完成")
}
