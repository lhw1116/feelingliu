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

		api.POST("/user/logout",v1.Logout)

		api.GET("/user/info", v1.GetUserInfo)

		api.GET("/user/about", v1.GetUserAbout)

		api.GET("/tags", v1.GetAllTags)

		api.GET("/articles", v1.GetArticles)

		api.GET("/articles/:id", v1.GetArticle)

		api.Use(middleware.JWt())

		api.PATCH("/user/edit", v1.EditUser)

		api.POST("/tags", v1.CreateTag)

		api.PUT("/tags/:id", v1.EditTag)

		api.DELETE("/tags/:id", v1.DeleteTag)

		api.POST("/articles", v1.CreateArticle)
		/*
		Error 1452: Cannot add or update a child row: a foreign key constraint fails (`blog`.`article`, CONSTRAINT `tags` FOREIGN KEY (`tag_id`) REFERENCES `tag` (`id`) ON DELETE CASCADE ON UPDATE CASCADE)
		 */
	}
	return r

}