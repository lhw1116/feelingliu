package modles

import (
	"fmt"
	"github.com/go-redis/redis/v7"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
)

var DB *gorm.DB
var RedisPool *redis.Client
var err error

func init() {

	//  Init database connect...
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8&parseTime=True&loc=Local",
		DBInfo.User,
		DBInfo.Password,
		DBInfo.Host,
		DBInfo.Port,
		DBInfo.DBName)
	DB, err = gorm.Open(DBInfo.Mode, dsn)
	if err != nil {
		panic(err)
	}
	DB.LogMode(true)
	DB.SingularTable(true)
	DB.DB().SetMaxIdleConns(10)
	DB.DB().SetMaxOpenConns(100)

	//  Init redispool connect...
	RedisPool = redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", RedisInfo.Host, RedisInfo.Port),
		Password: RedisInfo.Password,
		DB:       RedisInfo.DB,
	})
	_, err := RedisPool.Ping().Result()
	if err != nil {
		panic(err)
	}
}
