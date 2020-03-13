package common

import (
	"feelingliu/modles"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
)

var MysqlDB *gorm.DB

//  创建mysql连接池
func mysqlconn() {

	//  定义错误类型(因为需要赋值给全局变量的原因)
	var err error
	//  定义mysql连接串
	dbType := Viper.GetString("mysql.type")
	user := Viper.GetString("mysql.user")
	pass := Viper.GetString("mysql.pass")
	ip := Viper.GetString("mysql.ip")
	database := Viper.GetString("mysql.database")

	dsn := fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8",user,pass,ip,database)
	// open connect
	MysqlDB, err = gorm.Open(dbType, dsn)
	if err != nil {
		fmt.Println(err)
	}
	//  全局禁用表名复数
	MysqlDB.SingularTable(true)
	// 启用Logger，显示详细日志

	MysqlDB.LogMode(true)

	MysqlDB.DB().SetMaxOpenConns(100)
	MysqlDB.DB().SetConnMaxLifetime(100)
	MysqlDB.DB().SetMaxIdleConns(100)
	MysqlDB.AutoMigrate(&modles.User{})
	MysqlDB.AutoMigrate(&modles.Auth{})

}

func CloseDB() {
	defer MysqlDB.Close()
}