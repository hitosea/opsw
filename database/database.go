package database

import (
	"errors"
	"gorm.io/driver/mysql"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"opsw/utils"
	"opsw/vars"
	"strings"
	"time"
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

func closeDB(db *gorm.DB) {
	if sqlDB, err := db.DB(); err == nil {
		_ = sqlDB.Close()
	}
}

func Init() (*gorm.DB, error) {
	db, err := InDB(vars.Config.DB)
	if err != nil {
		return nil, err
	}
	defer closeDB(db)
	//
	autoIncrement := "AUTOINCREMENT"
	if strings.HasPrefix(strings.ToLower(vars.Config.DB), "mysql") {
		autoIncrement = "AUTO_INCREMENT"
	}
	sqls := utils.Sql("/database.sql", autoIncrement)
	for _, s := range sqls {
		err = db.Exec(s).Error
		if err != nil {
			return nil, err
		}
	}
	return db, err
}

func CreateUser(email, name, password string) (*vars.UserModel, error) {
	db, err := InDB(vars.Config.DB)
	if err != nil {
		return nil, err
	}
	defer closeDB(db)
	//
	var userData *vars.UserModel
	db.Table("users").Where(map[string]any{
		"email": email,
	}).Last(&userData)
	if userData.ID > 0 {
		return nil, errors.New("邮箱地址已存在")
	}
	//
	encrypt := utils.GenerateString(6)
	user := &vars.UserModel{
		Email:     email,
		Name:      name,
		Encrypt:   encrypt,
		Password:  utils.StringMd52(password, encrypt),
		Token:     utils.GenerateString(32),
		CreatedAt: uint32(time.Now().Unix()),
		UpdatedAt: uint32(time.Now().Unix()),
	}
	result := db.Table("users").Create(&user)
	if result.Error != nil {
		return nil, result.Error
	}
	return user, nil
}
