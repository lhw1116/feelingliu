package utils

import "os"
import "fmt"

func WriteErrorLog(s string) {
	//s := err.Error()
	file, _ := os.OpenFile("logs/error.log", os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0666)
	_, e := file.WriteString(s)
	if e != nil {
		fmt.Println(e)
	}
}
