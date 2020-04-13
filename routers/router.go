package routers

import (
	v1 "feelingliu/api/v1"
	"feelingliu/middleware"
	"feelingliu/modles"
	"github.com/gin-gonic/gin"
	"net/http"
)

func InitRouter() *gin.Engine {
	r := gin.New()
	r.HandleMethodNotAllowed = true

	r.NoMethod(func(c *gin.Context) {
		c.JSON(http.StatusMethodNotAllowed, gin.H{"result": false, "error": "Method Not Allowed"})
		return
	})

	r.Use(middleware.CustomLogger()) //  write logs middleware

	r.Use(middleware.CorsMiddleware()) //  intercept special ip address middleware

	r.Use(gin.Recovery())

	api := r.Group(modles.AppInfo.ApiBaseUrl)
	{
		api.Static(modles.AppInfo.StaticBasePath, modles.AppInfo.UploadBasePath)
		api.POST("/user/login",v1.Login)
	}
	return r

}