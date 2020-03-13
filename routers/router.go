package routers

import (
	"feelingliu/common"
	"feelingliu/middleware"
	auth "feelingliu/routers/api"
	v1 "feelingliu/routers/api/v1"
	"github.com/gin-gonic/gin"
)

func init() {
	r := gin.Default()

	gin.SetMode(common.RUNMODE)
	r.GET("/auth", auth.Getauth)

	api := r.Group("api/v1")
	api.Use(middleware.JWT())
	{
		//  获取标签列表
		api.GET("/tags", v1.GetTags)

		//  新建标签
		api.POST("tags", v1.AddTag)

		//  更新指定标签
		api.PUT("/tags/:id", v1.EditTag)

		//  删除指定标签
		api.DELETE("tags/:id", v1.DeleteTag)

		//  获取所有文章
		api.GET("/articles", v1.GetArticles)

		//  获取分类文章
		api.GET("/articles/:id", v1.GetArticles)

		//  新建文章
		api.POST("/articles", v1.AddArticle)

		//  修改文章
		api.PUT("/articles/:id", v1.Editarticle)

		//  删除文章
		api.DELETE("/articles/:id", v1.DeleteArticle)

	}
	r.Run()
}
