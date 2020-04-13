package modles

import (
	"github.com/spf13/viper"
	"time"
)

var AppInfo = &App{}
var ServerInfo = &Server{}
var DBInfo = &DataBase{}
var RedisInfo = &Redis{}

func init() {
	//设置配置文件的名字
	viper.SetConfigName("conf")
	viper.AddConfigPath("./config")

	viper.SetConfigType("yaml")
	if err := viper.ReadInConfig(); err != nil {
		panic(err)
	}

	//  Init projrct base information
	AppInfo.TimeFormat = viper.GetString("app.timeFormat")
	AppInfo.JwtSecret = viper.GetString("app.jwtSecret")
	AppInfo.TokenTimeout = viper.GetInt64("app.tokenTimeout")
	AppInfo.StaticBasePath = viper.GetString("app.staticBasePath")
	AppInfo.UploadBasePath = viper.GetString("app.uploadBasePath")
	AppInfo.ImageRelPath = viper.GetString("app.imageRelPath")
	AppInfo.AvatarRelPath = viper.GetString("app.avatarRelPath")
	AppInfo.ApiBaseUrl = viper.GetString("app.apiBaseUrl")

	//  Init projrct run env base information
	ServerInfo.RunMode = viper.GetString("server.runMode")
	ServerInfo.ServerAddr = viper.GetString("server.serverAddr")
	ServerInfo.ReadTimeout = time.Duration(viper.GetInt("server.readTimeout")) * time.Second
	ServerInfo.WriteTimeout = time.Duration(viper.GetInt("server.writeTimeout")) * time.Second

	//  Init projrct database information
	DBInfo.Mode = viper.GetString("database.mode")
	DBInfo.Host = viper.GetString("database.host")
	DBInfo.Port = viper.GetString("database.port")
	DBInfo.User = viper.GetString("database.user")
	DBInfo.Password = viper.GetString("database.password")
	DBInfo.DBName = viper.GetString("database.dbName")

	//  Init projrct redis information
	RedisInfo.Host = viper.GetString("redis.host")
	RedisInfo.Port = viper.GetString("redis.port")
	RedisInfo.Password = viper.GetString("redis.password")
	RedisInfo.DB = viper.GetInt("redis.db")
	RedisInfo.CacheTime = viper.GetInt("redis.cacheTime")
}
