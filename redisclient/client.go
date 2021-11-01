package redisclient

import (
	"context"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"sync"
	"thxy/model"
	"thxy/setting"
	"thxy/utils"
	"time"

	//"gitlab.forceup.in/forcepool/userbackend/model"
	//"gitlab.forceup.in/forcepool/userbackend/setting"
	//"gitlab.forceup.in/forcepool/userbackend/utils"

	"github.com/go-redis/redis/v8"
)

var redClient *redis.Client
var cliMutex sync.RWMutex
var env string

func GetClient() *redis.Client {
	cliMutex.RLock()
	if redClient != nil {
		defer cliMutex.RUnlock()
		return redClient
	}
	cliMutex.RUnlock()
	cliMutex.Lock()
	defer cliMutex.Unlock()
	if redClient != nil {
		return redClient
	}

	cfg := setting.TomlConfig.Test.Redis
	//env = setting.AppConfig.Redis.Env
	redClient = redis.NewClient(&redis.Options{
		Addr:     cfg.Host,
		Password: cfg.Password, // no password set
		//DB:       cfg.DB,       // use default DB
	})

	return redClient
}

func SetFileData(key string, data []byte) (err error) {
	client := GetClient()
	if client == nil {
		err = fmt.Errorf("failed connect to redis")
		return
	}

	cliMutex.Lock()
	defer cliMutex.Unlock()

	err = client.Set(context.Background(), key, data, 0).Err()
	if err != nil {
		return fmt.Errorf("设置SIDUSERCODE失败: %v", err)
	}
	return
}

func GetFileData(key string) (data []byte, err error) {
	client := GetClient()
	if client == nil {
		err = fmt.Errorf("failed connect to redis")
		return
	}

	cliMutex.Lock()
	defer cliMutex.Unlock()

	data, err = client.Get(context.Background(), key).Bytes()
	if err != nil {
		return nil, fmt.Errorf("设置SIDUSERCODE失败: %v", err)
	}
	return
}

func SetUserSession(key string, value *model.Loginuser) (err error) {
	client := GetClient()
	if client == nil {
		err = fmt.Errorf("failed connect to redis")
		return
	}

	cliMutex.Lock()
	defer cliMutex.Unlock()

	singleton := fmt.Sprintf("%s%s:%s", env, setting.SIDUSERCODE, value.UserCode)
	oldsid, err := client.Get(context.Background(), singleton).Result()
	if err != redis.Nil && len(oldsid) > 0 {
		if err = client.Del(context.Background(), oldsid).Err(); err != nil {
			return fmt.Errorf("删除同账号已登录用户sid失败: %v", err)
		}
	} else if err != redis.Nil && err != nil {
		return fmt.Errorf("查询redis失败! %v", err)
	}

	key = fmt.Sprintf("%s%s:%s", env, setting.SID_PREFIX, key)
	vbytes, err := json.Marshal(value)
	if err != nil {
		return
	}
	payload, err := utils.Encrypt([]byte(utils.SecretKey), vbytes)
	if err != nil {
		return
	}
	// 设置30分钟失效时间
	//expire := time.Minute * 30
	//// todo  暂时设置为24小时
	expire := time.Hour * 24

	err = client.Set(context.Background(), singleton, key, 0).Err()
	if err != nil {
		return fmt.Errorf("设置SIDUSERCODE失败: %v", err)
	}
	err = client.Set(context.Background(), key, hex.EncodeToString(payload), expire).Err()
	return
}

func GetUserSession(key string) (res *model.Loginuser, err error) {
	client := GetClient()
	if client == nil {
		err = fmt.Errorf("failed connect to redis")
		return
	}

	cliMutex.RLock()
	defer cliMutex.RUnlock()

	key = fmt.Sprintf("%s%s:%s", env, setting.SID_PREFIX, key)
	session, err := client.Get(context.Background(), key).Result()
	if err != nil {
		return
	}
	payload, err := hex.DecodeString(session)
	if err != nil {
		return
	}

	data, err := utils.Decrypt([]byte(utils.SecretKey), payload)
	if err != nil {
		return
	}
	res = new(model.Loginuser)
	err = json.Unmarshal(data, res)
	if err != nil {
		return
	}
	return
}

func DeleteUserSessionByCode(code string) error {
	client := GetClient()
	if client == nil {
		return fmt.Errorf("failed connect to redis")
	}

	cliMutex.Lock()
	defer cliMutex.Unlock()

	singleton := fmt.Sprintf("%s%s:%s", env, setting.SIDUSERCODE, code)
	oldsid, err := client.Get(context.Background(), singleton).Result()
	if err != redis.Nil && len(oldsid) > 0 {
		if err = client.Del(context.Background(), oldsid).Err(); err != nil {
			return fmt.Errorf("删除同账号已登录用户sid失败: %v", err)
		}
	} else if err != redis.Nil && err != nil {
		return fmt.Errorf("查询redis失败! %v", err)
	}

	return nil
}

func DeleteUserSession(key string) (err error) {
	client := GetClient()
	if client == nil {
		err = fmt.Errorf("failed connect to redis")
		return
	}

	cliMutex.Lock()
	defer cliMutex.Unlock()

	key = fmt.Sprintf("%s%s:%s", env, setting.SID_PREFIX, key)
	client.Del(context.Background(), key)
	return
}

func CloseClient() error {
	cliMutex.RLock()
	defer cliMutex.RUnlock()

	if redClient != nil {
		return redClient.Close()
	}
	return nil
}

func GetVcode(key string) (res *model.Vcode_t, err error) {
	client := GetClient()
	if client == nil {
		err = fmt.Errorf("failed connect to redis")
		return
	}
	key = fmt.Sprintf("%s%s:%s", env, setting.VERCODE_PREFIX, key)
	limit, err := client.Get(context.Background(), key).Result()
	if err != nil {
		return
	}
	payload, err := hex.DecodeString(limit)
	if err != nil {
		return
	}

	data, err := utils.Decrypt([]byte(utils.SecretKey), payload)
	if err != nil {
		return
	}
	res = new(model.Vcode_t)
	err = json.Unmarshal(data, res)
	return
}

func SetVcode(key string, value *model.Vcode_t) (err error) {
	client := GetClient()
	if client == nil {
		err = fmt.Errorf("failed connect to redis")
		return
	}
	key = fmt.Sprintf("%s%s:%s", env, setting.VERCODE_PREFIX, key)
	vBytes, err := json.Marshal(value)
	if err != nil {
		return
	}
	payload, err := utils.Encrypt([]byte(utils.SecretKey), vBytes)
	if err != nil {
		return
	}
	// 设置5分钟失效
	expire := time.Minute * 5
	err = client.Set(context.Background(), key, hex.EncodeToString(payload), expire).Err()
	return
}

//func GetUserLimited(key string) (res *model.UserLimited, err error) {
//	client := GetClient()
//	if client == nil {
//		err = fmt.Errorf("failed connect to redis")
//		return
//	}
//	key = fmt.Sprintf("%s%s:%s", env, setting.USERCODE_PREFIX, key)
//	limit, err := client.Get(context.Background(), key).Result()
//	if err != nil {
//		res = new(model.UserLimited)
//		return
//	}
//
//	//for testing
//	//limit = "053fdcd533d7fd183a6b9e9b954ee96f100df185057c2871b90da990d11a360baa13fed5e51950f3087c2e543d6282d251ab2ebca289eb6333c3b1dcdd4e28505f89edbdc5a13db3098a246eca89f20d487f386a16bf47f142e351487e0d1d60f2476878e5a8e957ff94fffa7dbe24710e4dff79c4c5f339a14de4b0d535443e57056f70684ef70772529e50ae3443d4b5eab449c96786b0ec56d6b29f65bab5e5fe2d23ec987232ddb93b1fb097293d6e44097c309909e37bd94b126f50bb778449b4771b5a27b8efa0ded2d9ecb01980"
//
//	payload, err := hex.DecodeString(limit)
//	if err != nil {
//		return
//	}
//
//	data, err := utils.Decrypt([]byte(utils.SecretKey), payload)
//	if err != nil {
//		return
//	}
//	res = new(model.UserLimited)
//	err = json.Unmarshal(data, res)
//	return
//}
//
//func SetUserLimited(key string, value *model.UserLimited) (err error) {
//	client := GetClient()
//	if client == nil {
//		err = fmt.Errorf("failed connect to redis")
//		return
//	}
//	key = fmt.Sprintf("%s%s:%s", env, setting.USERCODE_PREFIX, key)
//	vBytes, err := json.Marshal(value)
//	if err != nil {
//		return
//	}
//	payload, err := utils.Encrypt([]byte(utils.SecretKey), vBytes)
//	if err != nil {
//		return
//	}
//	// 设置24小时失效时间
//	expire := time.Hour * 24
//	err = client.Set(context.Background(), key, hex.EncodeToString(payload), expire).Err()
//	return
//}
//
//func GetVcode(key string) (res *model.Vcode_t, err error) {
//	client := GetClient()
//	if client == nil {
//		err = fmt.Errorf("failed connect to redis")
//		return
//	}
//	key = fmt.Sprintf("%s%s:%s", env, setting.VERCODE_PREFIX, key)
//	limit, err := client.Get(context.Background(), key).Result()
//	if err != nil {
//		return
//	}
//	payload, err := hex.DecodeString(limit)
//	if err != nil {
//		return
//	}
//
//	data, err := utils.Decrypt([]byte(utils.SecretKey), payload)
//	if err != nil {
//		return
//	}
//	res = new(model.Vcode_t)
//	err = json.Unmarshal(data, res)
//	return
//}
//
//func SetVcode(key string, value *model.Vcode_t) (err error) {
//	client := GetClient()
//	if client == nil {
//		err = fmt.Errorf("failed connect to redis")
//		return
//	}
//	key = fmt.Sprintf("%s%s:%s", env, setting.VERCODE_PREFIX, key)
//	vBytes, err := json.Marshal(value)
//	if err != nil {
//		return
//	}
//	payload, err := utils.Encrypt([]byte(utils.SecretKey), vBytes)
//	if err != nil {
//		return
//	}
//	// 设置5分钟失效
//	expire := time.Minute * 5
//	err = client.Set(context.Background(), key, hex.EncodeToString(payload), expire).Err()
//	return
//}

func DeleteVcode(key string) (err error) {
	client := GetClient()
	if client == nil {
		err = fmt.Errorf("failed connect to redis")
		return
	}
	key = fmt.Sprintf("%s%s:%s", env, setting.VERCODE_PREFIX, key)
	client.Del(context.Background(), key)
	return
}

func SetWXToken(token string) (err error) {
	client := GetClient()
	if client == nil {
		err = fmt.Errorf("failed connect to redis")
		return
	}
	// 目前只有一个小程序appid, 所以不用区分预发和生产，可以使用同一个key
	key := fmt.Sprintf("%s", setting.WXTOKEN_KEY)

	// 设置24小时失效时间
	expire := time.Second * 7000
	err = client.Set(context.Background(), key, token, expire).Err()
	return
}

func GetWXToken() (token string, err error) {
	client := GetClient()
	if client == nil {
		err = fmt.Errorf("failed connect to redis")
		return
	}
	key := fmt.Sprintf("%s", setting.WXTOKEN_KEY)
	token, err = client.Get(context.Background(), key).Result()
	return
}
