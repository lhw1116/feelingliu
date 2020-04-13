package middleware

import (
	"feelingliu/modles"
	"fmt"
	"github.com/gin-gonic/gin"
	"io"
	"os"
)

type Level int

const (
	INFO Level = iota
	WARNING
	ERROR
	FATAL
)

var (
	file *os.File
	e    error
)


func CustomLogger() gin.HandlerFunc {
	gin.DisableConsoleColor()
	file, err := os.OpenFile("logs/gin.log", os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		panic(err)
	}
	gin.DefaultWriter = io.MultiWriter(file,os.Stdout)

	g :=gin.LoggerWithFormatter(func(params gin.LogFormatterParams) string {
		levelFlags := []string{"INFO","WRAN","ERROR","FATAL"}
		var level string
		status := params.StatusCode

		switch {
		case status > 499:
			level = levelFlags[FATAL]
		case status > 399:
			level = levelFlags[ERROR]
		case status > 299:
			level = levelFlags[WARNING]
		default:
			level = levelFlags[INFO]
		}

		return fmt.Sprintf("[%s] - %s - [%s] \"%s %s %s %d %s \"%s\" %s\"\n",
			level,
			params.ClientIP,
			params.TimeStamp.Format(modles.AppInfo.TimeFormat),
			params.Method,
			params.Path,
			params.Request.Proto,
			status,
			params.Latency,
			params.Request.UserAgent(),
			params.ErrorMessage,
		)


	})

	return g
}

func CloseLogFile() {
	if err := file.Close(); err != nil {
		return
	}
}