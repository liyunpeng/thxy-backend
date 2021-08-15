package setting

import (
	"fmt"
	"github.com/jinzhu/configor"
	"thxy/logger"
	"thxy/types"
)

var TomlConfig *types.AppConfig

func InitConfig(fileName string) {
	TomlConfig = new(types.AppConfig)
	err := configor.Load(TomlConfig, fileName)
	if err != nil {
		panic(fmt.Sprintf("fail to load app config:\n %v\n", err))
	}
	logger.Info.Println(" 配置加载完成")
}
