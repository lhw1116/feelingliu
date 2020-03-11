package common

import (
	"fmt"
	"github.com/spf13/viper"
)
var (
	Viper *viper.Viper
	RUNMODE string
)


//  viper 加载配置文件初始化
func viperLeader() {
	Viper = viper.New()

	//设置配置文件的名字
	Viper.SetConfigName("conf")
	Viper.AddConfigPath("./config")
	//设置配置文件类型
	Viper.SetConfigType("yaml");
	if err := Viper.ReadInConfig(); err != nil{
		fmt.Println(err)
	}
	RUNMODE = viper.GetString("RUNMODE")
}