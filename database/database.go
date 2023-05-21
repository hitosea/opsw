package database

import (
	"encoding/json"
	"errors"
	"gorm.io/driver/mysql"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"opsw/utils"
	"opsw/vars"
	"strings"
)

func InDB(str string) (*gorm.DB, error) {
	sp := strings.Split(str, "://")
	dbType := "sqlite"
	dbPath := str
	if len(sp) == 2 {
		dbType = strings.ToLower(sp[0])
		dbPath = sp[1]
	}
	if dbType == "mysql" {
		return gorm.Open(mysql.Open(dbPath), &gorm.Config{})
	} else {
		return gorm.Open(sqlite.Open(dbPath), &gorm.Config{})
	}
}

func CloseDB(db *gorm.DB) {
	if sqlDB, err := db.DB(); err == nil {
		_ = sqlDB.Close()
	}
}

func Init() (*gorm.DB, error) {
	db, err := InDB(vars.Config.DB)
	if err != nil {
		return nil, err
	}
	defer CloseDB(db)
	//
	autoIncrement := "AUTOINCREMENT"
	if strings.HasPrefix(strings.ToLower(vars.Config.DB), "mysql") {
		autoIncrement = "AUTO_INCREMENT"
	}
	sqls := utils.Sql("/install.sql", autoIncrement)
	for _, s := range sqls {
		err = db.Exec(s).Error
		if err != nil {
			return nil, err
		}
	}
	return db, err
}

func UserGet(query any) (*User, error) {
	db, err := InDB(vars.Config.DB)
	if err != nil {
		return nil, err
	}
	defer CloseDB(db)
	//
	var user *User
	db.Where(query).Last(&user)
	if user.Id == 0 {
		return nil, errors.New("用户不存在")
	}
	return user, nil
}

func UserCheck(email, password string) (*User, error) {
	user, err := UserGet(map[string]any{
		"email": email,
	})
	if err != nil {
		return nil, err
	}
	if user.Password != utils.StringMd52(password, user.Encrypt) {
		return nil, errors.New("邮箱或密码错误")
	}
	return user, nil
}

func UserCreate(email, name, password string) (*User, error) {
	db, err := InDB(vars.Config.DB)
	if err != nil {
		return nil, err
	}
	defer CloseDB(db)
	//
	var user *User
	db.Where(map[string]any{
		"email": email,
	}).Last(&user)
	if user.Id > 0 {
		return nil, errors.New("邮箱地址已存在")
	}
	//
	encrypt := utils.GenerateString(6)
	user = &User{
		Email:    email,
		Name:     name,
		Encrypt:  encrypt,
		Password: utils.StringMd52(password, encrypt),
		Token:    utils.Base64Encode("u:%s", utils.GenerateString(22)),
		Avatar:   "",
	}
	err = db.Create(&user).Error
	if err != nil {
		return nil, err
	}
	return user, nil
}

func ServerCreate(server *Server, serverUser *ServerUser) error {
	db, err := InDB(vars.Config.DB)
	if err != nil {
		return err
	}
	defer CloseDB(db)
	//
	err = db.Transaction(func(tx *gorm.DB) error {
		err = tx.Create(server).Error
		if err != nil {
			return err
		}
		serverUser.ServerId = server.Id
		err = tx.Create(serverUser).Error
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return err
	}
	return nil
}

func ServerGet(query any, userId int32, owner bool) (*Server, error) {
	db, err := InDB(vars.Config.DB)
	if err != nil {
		return nil, err
	}
	defer CloseDB(db)
	//
	var server *Server
	db.Where(query).Last(&server)
	if server.Id == 0 {
		return nil, errors.New("服务器不存在")
	}
	if userId > -1 || owner {
		if userId == 0 {
			return nil, errors.New("没有权限-1")
		}
		var serverUser *ServerUser
		db.Where(map[string]any{
			"server_id": server.Id,
			"user_id":   userId,
		}).Last(&serverUser)
		if serverUser.Id == 0 {
			return nil, errors.New("没有权限-2")
		}
		if owner {
			if serverUser.OwnerId != userId {
				return nil, errors.New("没有权限-3")
			}
		}
	}
	return server, nil
}

func ServerUpdate(server *Server) error {
	db, err := InDB(vars.Config.DB)
	if err != nil {
		return err
	}
	defer CloseDB(db)
	//
	err = db.Save(server).Error
	if err != nil {
		return err
	}
	return nil
}

func ServerDelete(query any, userId int32) (*Server, error) {
	db, err := InDB(vars.Config.DB)
	if err != nil {
		return nil, err
	}
	defer CloseDB(db)
	//
	server, err := ServerGet(query, userId, true)
	if err != nil {
		return nil, err
	}
	err = db.Transaction(func(tx *gorm.DB) error {
		err = tx.Delete(&server).Error
		if err != nil {
			return err
		}
		err = tx.Delete(&ServerUser{
			ServerId: server.Id,
		}).Error
		return nil
	})
	if err != nil {
		return nil, err
	}
	return server, nil
}

func ServerInfoGet(serverId int32) (*ServerInfo, error) {
	db, err := InDB(vars.Config.DB)
	if err != nil {
		return nil, err
	}
	defer CloseDB(db)
	//
	var serverInfo *ServerInfo
	db.Where(map[string]any{
		"server_id": serverId,
	}).Last(&serverInfo)
	if serverInfo.Id == 0 {
		return nil, errors.New("服务器信息不存在")
	}
	return serverInfo, nil
}

func ServerInfoUpdate(serverId int32, data any) error {
	ss, err := json.Marshal(data)
	if err != nil {
		return err
	}
	var info *ServerInfo
	if err = json.Unmarshal(ss, &info); err != nil {
		return err
	}
	info.ServerId = serverId
	//
	db, err := InDB(vars.Config.DB)
	if err != nil {
		return err
	}
	defer CloseDB(db)
	//
	var serverInfo *ServerInfo
	db.Where(map[string]any{
		"server_id": serverId,
	}).Last(&serverInfo)
	if serverInfo.Id == 0 {
		err = db.Create(info).Error
		if err != nil {
			return err
		}
	} else {
		info.Id = serverInfo.Id
		info.CreatedAt = serverInfo.CreatedAt
		err = db.Save(info).Error
		if err != nil {
			return err
		}
	}
	return nil
}
