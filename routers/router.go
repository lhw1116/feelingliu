package routers

import (
	v1 "feelingliu/api/v1"
	"feelingliu/middleware"
	"feelingliu/modles"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	_ "github.com/swaggo/gin-swagger/example/basic/docs"
	"net/http"
)

func InitRouter() *gin.Engine {
	r := gin.New()
	url := ginSwagger.URL("http://localhost:8888/swagger/doc.json") // The url pointing to API definition
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler, url))
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
		api.POST("/user/login", v1.Login)

		api.POST("/user/logout", v1.Logout)

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

		api.DELETE("/articles/:id", v1.DeleteArticle)

		api.PUT("/articles/:id", v1.EditArticle)
	}
	return r

}
