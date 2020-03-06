package main

import (
	"feelingliu/common"
	"fmt"
)

//  程序启动入口
func main() {
	getString := common.Viper.GetString("liuhanwen")
	fmt.Println(getString)
}
