package database

import (
	"fmt"
	"github.com/ghjan/gin-blog/pkg/setting"
	"github.com/jinzhu/gorm"
	"os"
	"sync"
)

var instance *gorm.DB
var once sync.Once

func Connect(env string) *gorm.DB {
	db, err := gorm.Open(setting.DatabaseSetting.Type, fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8&parseTime=True&loc=Local",
		setting.DatabaseSetting.User,
		setting.DatabaseSetting.Password,
		setting.DatabaseSetting.Host,
		setting.DatabaseSetting.Name))

	if err != nil {
		panic(err.Error())
	}

	return db
}

func GetInstance() *gorm.DB {
	env := os.Getenv("ENV")
	once.Do(func() {
		instance = Connect(env)
		if env != "production" {
			instance.LogMode(true)
		}
	})
	return instance
}

func Close() {
	if instance == nil {
		return
	}
	instance.Close()
}
