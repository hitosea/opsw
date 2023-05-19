package database

import (
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
		Token:    utils.GenerateString(32),
		Avatar:   "",
	}
	err = db.Create(&user).Error
	if err != nil {
		return nil, err
	}
	return user, nil
}
