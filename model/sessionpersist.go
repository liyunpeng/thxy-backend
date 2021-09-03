package model

import (
	"database/sql"
	"fmt"
	"thxy/logger"
	"thxy/setting"
	"thxy/types"
	"thxy/utils"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

const ticker_interval = 60 * 60

var conf *types.AppConfig

func SessionPersisInit() {
	conf = setting.TomlConfig
	initTables()
	go checkExpire()
}

type AdminSess struct {
	Id      int
	AdminId int
	LogTime int64
	Sid     string
}

type UserSess struct {
	Id      int
	UserId  int
	LogTime int64
	Sid     string
}

func SessionDB() (*sql.DB, error) {
	return sql.Open(conf.Session.Type, conf.Session.Path)
}

func PersisAdminsess(adminid int, sid string) (err error) {
	_, err = GetAdminsessByID(adminid)
	if err != nil {
		fmt.Println(err)
		err = AddToAdminsess(adminid, sid)
	} else {
		err = UpdateAdminsessLogTime(adminid, sid)
	}
	return
}

func UpdateAdminsessLogTime(adminid int, sid string) (err error) {
	db, err := SessionDB()
	defer db.Close()
	if err != nil {
		fmt.Println(err)
		return
	}
	stmt := fmt.Sprintf("update adminsess set log_time=%d, sid=\"%s\" where admin_id=%d", utils.CurrentTimestamp(), sid, adminid)
	//fmt.Println(stmt)
	_, err = db.Exec(stmt)
	if err != nil {
		fmt.Println(err)
	}
	return
}

func GetAdminsessByID(adminid int) (as AdminSess, err error) {
	db, err := SessionDB()
	defer db.Close()
	if err != nil {
		fmt.Println(err)
		return
	}
	rows, err := db.Query(fmt.Sprintf("select * from adminsess where admin_id=%d", adminid))
	if err != nil {
		fmt.Println(err)
		return
	}
	defer rows.Close()
	len := 0
	for rows.Next() {
		as = AdminSess{}
		err = rows.Scan(&as.Id, &as.AdminId, &as.Sid, &as.LogTime)
		if err != nil {
			fmt.Println(err)
			return
		}
		len++
	}
	err = rows.Err()
	if err != nil {
		fmt.Println(err)
		return
	}
	if len == 0 {
		err = fmt.Errorf("no record")
		return
	}
	return
}

func AddToAdminsess(adminid int, sid string) (err error) {
	db, err := SessionDB()
	defer db.Close()
	if err != nil {
		fmt.Println(err)
		return
	}
	stmt := fmt.Sprintf("insert into adminsess (admin_id, log_time, sid) values (%d, %d, \"%s\")", adminid, utils.CurrentTimestamp(), sid)
	//fmt.Println(stmt)
	_, err = db.Exec(stmt)
	if err != nil {
		fmt.Println(err)
	}
	return
}

func RemoveAdminsess(adminid int) {
	db, err := SessionDB()
	defer db.Close()
	if err != nil {
		fmt.Println(err)
		return
	}
	stmt := fmt.Sprintf("delete from adminsess where admin_id=%d", adminid)
	//fmt.Println(stmt)
	_, err = db.Exec(stmt)
	if err != nil {
		fmt.Println(err)
	}
	return
}

func SessAdminList() (asList []AdminSess, err error) {
	db, err := SessionDB()
	defer db.Close()
	if err != nil {
		fmt.Println(err)
		return
	}
	curTime := utils.CurrentTimestamp()

	stmt := fmt.Sprintf("select * from adminsess where log_time > %d", curTime-int64(conf.Session.Life))
	//fmt.Println(stmt)
	rows, err := db.Query(stmt)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer rows.Close()
	asList = make([]AdminSess, 0, 50)
	for rows.Next() {
		as := AdminSess{}
		err = rows.Scan(&as.Id, &as.AdminId, &as.Sid, &as.LogTime)
		if err != nil {
			fmt.Println(err)
			return
		}
		asList = append(asList, as)
	}
	err = rows.Err()
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(asList)
	return
}

func PersisUsersess(userid int, sid string) (err error) {
	_, err = GetUsersessByID(userid)
	if err != nil {
		fmt.Println(err)
		err = AddToUsersess(userid, sid)
	} else {
		err = UpdateUsersessLogTime(userid, sid)
	}
	return
}

func UpdateUsersessLogTime(userid int, sid string) (err error) {
	db, err := SessionDB()
	defer db.Close()
	if err != nil {
		fmt.Println(err)
		return
	}
	stmt := fmt.Sprintf("update usersess set log_time=%d, sid=\"%s\" where user_id=%d", utils.CurrentTimestamp(), sid, userid)
	//fmt.Println(stmt)
	_, err = db.Exec(stmt)
	if err != nil {
		fmt.Println(err)
	}
	return
}

func GetUsersessByID(userid int) (us UserSess, err error) {
	db, err := SessionDB()
	defer db.Close()
	if err != nil {
		fmt.Println(err)
		return
	}
	rows, err := db.Query(fmt.Sprintf("select * from usersess where user_id=%d", userid))
	if err != nil {
		fmt.Println(err)
		return
	}
	defer rows.Close()
	len := 0
	for rows.Next() {
		us = UserSess{}
		err = rows.Scan(&us.Id, &us.UserId, &us.Sid, &us.LogTime)
		if err != nil {
			fmt.Println(err)
			return
		}
		len++
	}
	err = rows.Err()
	if err != nil {
		fmt.Println(err)
		return
	}
	if len == 0 {
		err = fmt.Errorf("no record")
		return
	}
	return
}

func AddToUsersess(userid int, sid string) (err error) {
	db, err := SessionDB()
	defer db.Close()
	if err != nil {
		fmt.Println(err)
		return
	}
	stmt := fmt.Sprintf("insert into usersess (user_id, log_time, sid) values (%d, %d, \"%s\")", userid, utils.CurrentTimestamp(), sid)
	//fmt.Println(stmt)
	_, err = db.Exec(stmt)
	if err != nil {
		fmt.Println(err)
	}
	return
}

func RemoveUsersess(userid int64) {
	db, err := SessionDB()
	defer db.Close()
	if err != nil {
		fmt.Println(err)
		return
	}
	stmt := fmt.Sprintf("delete from usersess where user_id=%d", userid)
	//fmt.Println(stmt)
	_, err = db.Exec(stmt)
	if err != nil {
		fmt.Println(err)
	}
	return
}

func SessUserList() (usList []UserSess, err error) {
	db, err := SessionDB()
	defer db.Close()
	if err != nil {
		fmt.Println(err)
		return
	}
	curTime := utils.CurrentTimestamp()
	stmt := fmt.Sprintf("select * from usersess where log_time > %d", curTime-int64(conf.Session.Life))
	//fmt.Println(stmt)
	rows, err := db.Query(stmt)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer rows.Close()
	usList = make([]UserSess, 0, 50)
	for rows.Next() {
		us := UserSess{}
		err = rows.Scan(&us.Id, &us.UserId, &us.Sid, &us.LogTime)
		if err != nil {
			fmt.Println(err)
			return
		}
		usList = append(usList, us)
	}
	err = rows.Err()
	if err != nil {
		fmt.Println(err)
	}
	return
}

func cleanExpiredAdminsess() {
	db, err := SessionDB()
	defer db.Close()
	if err != nil {
		fmt.Println(err)
		return
	}
	curTime := utils.CurrentTimestamp()
	expired := curTime - int64(conf.Session.Life)
	stmt := fmt.Sprintf("delete from adminsess where log_time < %d", expired)
	//fmt.Println(stmt)
	_, err = db.Exec(stmt)
	if err != nil {
		fmt.Println(err)
	}
	return
}

func cleanExpiredUsersess() {
	db, err := SessionDB()
	defer db.Close()
	if err != nil {
		fmt.Println(err)
		return
	}
	curTime := utils.CurrentTimestamp()

	expired := curTime - int64(conf.Session.Life)
	stmt := fmt.Sprintf("delete from usersess where log_time < %d", expired)
	//fmt.Println(stmt)
	_, err = db.Exec(stmt)
	if err != nil {
		fmt.Println(err)
	}
	return
}

func checkExpire() {
	ticker := time.NewTicker(time.Duration(ticker_interval * 1e9))
	for {
		select {
		case t := <-ticker.C:
			logger.Info.Println("adminsess expire check task - %v", t)
			cleanExpiredAdminsess()
			cleanExpiredUsersess()
		}
	}
}

func initTables() {
	db, err := SessionDB()
	defer db.Close()
	if err != nil {
		fmt.Println(err)
		return
	}
	// check adminsess table
	sqlStmt := `CREATE TABLE IF NOT EXISTS adminsess (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		admin_id int(11) NOT NULL,
		sid varchar(255) NOT NULL,
		log_time int(11) NOT NULL
	);`
	_, err = db.Exec(sqlStmt)
	if err != nil {
		fmt.Printf("%q: %s\n", err, sqlStmt)
		return
	}

	// check usersess table

	sqlStmt = `CREATE TABLE IF NOT EXISTS usersess (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		user_id int(11) NOT NULL,
		sid varchar(255) NOT NULL,
		log_time int(11) NOT NULL
	);`
	_, err = db.Exec(sqlStmt)
	if err != nil {
		fmt.Printf("%q: %s\n", err, sqlStmt)
		return
	}
}
