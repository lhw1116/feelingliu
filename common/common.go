package common

import (
	"feelingliu/utils"
	"github.com/go-redis/redis/v7"
	"github.com/jinzhu/gorm"
	"fmt"
)

var DB *gorm.DB
var RedisPool *redis.Client
var err error

func init() {

	//  Init database connect...
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8&parseTime=True&loc=Local",
		utils.DBInfo.User,
		utils.DBInfo.Password,
		utils.DBInfo.Host,
		utils.DBInfo.Port,
		utils.DBInfo.DBName,)
	DB, err = gorm.Open(utils.DBInfo.Mode,dsn)
	if err != nil {
		panic(err)
	}
	DB.DB().SetMaxIdleConns(10)
	DB.DB().SetMaxOpenConns(100)

	//  Init redispool connect...
	RedisPool = redis.NewClient(&redis.Options{
		Addr: fmt.Sprint("%s:%s",utils.RedisInfo.Host,utils.RedisInfo.Port),
		Password:utils.RedisInfo.Password,
		DB:utils.RedisInfo.DB,


	})
}
