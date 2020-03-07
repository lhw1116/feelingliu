package common

import (
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
)

//  创建mysql连接池
func mysqlconn() {

	//  定义mysql连接串
	user := Viper.GetString("mysql.user")
	pass := Viper.GetString("mysql.pass")
	ip := Viper.GetString("mysql.ip")
	database := Viper.GetString("mysql.database")

	dsn := fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8",user,pass,ip,database)

	// open connect
	db, err := gorm.Open("mysql", dsn)
	if err != nil {
		fmt.Println(err)
	}
	defer db.Close()
}