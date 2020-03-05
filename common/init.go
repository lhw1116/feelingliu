package common

//  程序初始化入口
func init() {
	viper()
	mysqlconn()
	redisconn()
}