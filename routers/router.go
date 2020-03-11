package routers

import (
	"feelingliu/common"
	v1 "feelingliu/routers/api/v1"
	"github.com/gin-gonic/gin"
)

func init() {
	r := gin.Default()
	gin.SetMode(common.RUNMODE)
	api := r.Group("api/v1")
	{
		//  获取标签列表
		api.GET("/tags",v1.GetTags)

		//  新建标签
		api.POST("tags",v1.AddTag)

		//  更新指定标签
		api.PUT("/tags/:id",v1.EditTag)

		//  删除指定标签
		api.DELETE("tags/:id",v1.DeleteTag)
	}
	r.Run()
}
