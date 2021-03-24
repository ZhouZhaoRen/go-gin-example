package models

import (
	"fmt"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"go-gin-example/pkg/setting"
	"log"
	"time"
)

var db *gorm.DB

type Model struct {
	ID         int `gorm:"primary_key" json:"id"`
	CreatedOn  int `json:"created_on"`
	ModifiedOn int `json:"modified_on"`
}

func Setup() {
	var err error

	// 进行数据库连接
	db, err = gorm.Open(setting.DatabaseSetting.Type, fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8&parseTime=True&loc=Local",
		setting.DatabaseSetting.User,
		setting.DatabaseSetting.Password,
		setting.DatabaseSetting.Host,
		setting.DatabaseSetting.Name))

	if err != nil {
		log.Println("连接数据库出错：", err)
	}
	// 数据库表名的命名规则
	gorm.DefaultTableNameHandler = func(db *gorm.DB, defaultTableName string) string {
		return setting.DatabaseSetting.TablePrefix + defaultTableName
	}
	// 其他设置
	db.SingularTable(true) // 全局设置表名不可以为复数形式
	db.LogMode(true)
	db.DB().SetMaxIdleConns(10)  // 最大空闲连接
	db.DB().SetMaxOpenConns(100) // 最大连接数

	// 注册callback
	db.Callback().Create().Replace("gorm:update_time_stamp", updateTimeStampForCreateCallback)
	db.Callback().Update().Replace("gorm:update_time_stamp", updateTimeStampForUpdateCallback)
}

// 关闭连接的函数
func CloseDB() {
	defer db.Close()
}

// 实现Callbacks 创建的时候触发这个函数  创建的时候创建和修改时间都得改变
func updateTimeStampForCreateCallback(scope *gorm.Scope) {
	if !scope.HasError() {
		nowTime := time.Now().Unix()
		// 创建时间
		if createTimeField, ok := scope.FieldByName("CreatedOn"); ok {
			if createTimeField.IsBlank {
				createTimeField.Set(nowTime)
			}
		}
		// 修改时间
		if modifiedTimeField, ok := scope.FieldByName("ModifiedOn"); ok {
			if modifiedTimeField.IsBlank {
				modifiedTimeField.Set(nowTime)
			}
		}
	}
}

// 修改的时候触发这个函数  修改的时候只需要改变修改时间
func updateTimeStampForUpdateCallback(scope *gorm.Scope){
	if _, ok := scope.Get("gorm:update_column"); !ok {
		scope.SetColumn("ModifiedOn", time.Now().Unix())
	}
}
