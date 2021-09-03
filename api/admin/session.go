package admin

import (
	"sync"
	"thxy/logger"
	"thxy/model"
	"thxy/setting"
	"thxy/utils"
	"time"

)

var Sessions *SeesionBuckets

type SeesionBuckets struct {
	bLock   sync.RWMutex
	bucket  map[string]*model.AdminUser
	account map[string]*model.AdminUser
	adminId map[int]*model.AdminUser
}

func SessionInit() {
	Sessions = new(SeesionBuckets)
	Sessions.bucket = make(map[string]*model.AdminUser)
	Sessions.account = make(map[string]*model.AdminUser)
	Sessions.adminId = make(map[int]*model.AdminUser)

	ticker := time.NewTicker(time.Hour * 1)
	go func() {
		for range ticker.C {
			if time.Now().Format("15") == "00" {
				Sessions.bLock.Lock()
				now := utils.CurrentTimestamp()
				for k, v := range Sessions.bucket {
					if now-v.LastTime > int64(setting.TomlConfig.Session.Life) {
						delete(Sessions.bucket, k)
						delete(Sessions.adminId, v.AdminId)
						delete(Sessions.account, v.Account)
					}
				}
				Sessions.bLock.Unlock()
				logger.Info.Println("session clear 24 hour.")
			}
		}
	}()

	if asList, err := model.SessAdminList(); err == nil && len(asList) > 0 {
		if adminList, err := model.GetAdminListBySess(asList); err == nil && len(adminList) > 0 {
			for _, item := range adminList {
				id := item.Id
				account := item.Account
				var admin *model.AdminUser
				sid := getSid(asList, id)
				if sid == "" {
					continue
				}
				logger.Info.Println(id, sid)
				admin = &model.AdminUser{
					AdminSessionId: sid,
					AdminId:        id,
					Account:        account,
					LastTime:       utils.CurrentTimestamp(),
					Phone:          item.Phone,
				}
				Sessions.AddAdmin(admin)
				//admin.Roles = e.GetRolesForUser(admin.Account)
				logger.Info.Println("admin=", *admin)
			}
		}
	}
}
func getSid(list []model.AdminSess, adminid int) string {
	for i, l := 0, len(list); i < l; i++ {
		if list[i].AdminId == adminid {
			return list[i].Sid
		}
	}
	return ""
}

func (b *SeesionBuckets) AddAdmin(admin *model.AdminUser) {
	b.bLock.Lock()
	if uu, ok := b.account[admin.Account]; ok {
		logger.Info.Printf("Account:[%s]已经在线,另一端将会下线", admin.Account)
		delete(b.bucket, uu.AdminSessionId)
		delete(b.account, uu.Account)
		delete(b.adminId, uu.AdminId)
	}
	logger.Info.Printf("Account:[%s] add in Sessions", admin.Account)
	b.bucket[admin.AdminSessionId] = admin
	b.account[admin.Account] = admin
	b.adminId[admin.AdminId] = admin
	b.bLock.Unlock()
	return
}

func (b *SeesionBuckets) DelAdmin(sion string) {
	b.bLock.Lock()
	if uu, ok := b.bucket[sion]; ok {
		delete(b.bucket, sion)
		delete(b.account, uu.Account)
		delete(b.adminId, uu.AdminId)
	} else {
		logger.Info.Printf("AdminSessionId:[%s]用户不在线", sion)
		//logger.Info.Printlnf("AdminSessionId:[%s]用户不在线", sion)
	}
	b.bLock.Unlock()
	return
}

func (b *SeesionBuckets) DelAdminById(admin_id int) {
	b.bLock.Lock()
	defer b.bLock.Unlock()
	if uu, exist := b.adminId[admin_id]; exist {
		delete(b.bucket, uu.AdminSessionId)
		delete(b.account, uu.Account)
		delete(b.adminId, uu.AdminId)
	}
	return
}

func (b *SeesionBuckets) QueryloginS(sion string) (admin *model.AdminUser, ok bool) {
	b.bLock.Lock()
	if admin, ok = b.bucket[sion]; ok {
		logger.Info.Printf("AdminSessionId:[%s]用户在线", sion)
	} else {
		logger.Info.Printf("AdminSessionId:[%s]用户不在线", sion)
	}
	b.bLock.Unlock()
	return
}

func (b *SeesionBuckets) QueryloginA(account string) (admin *model.AdminUser, ok bool) {
	b.bLock.Lock()
	if admin, ok = b.account[account]; ok {
		logger.Info.Printf("AdminSessionId:[%s]用户在线", account)
	} else {
		logger.Info.Printf("AdminSessionId:[%s]用户不在线", account)
	}
	b.bLock.Unlock()
	return
}

func (b *SeesionBuckets) QueryloginB(sion string) (ok bool) {
	b.bLock.Lock()
	_, ok = b.bucket[sion]
	b.bLock.Unlock()
	return
}
