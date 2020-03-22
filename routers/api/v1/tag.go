package v1

import (
	"feelingliu/common"
	"feelingliu/modles"
	"github.com/gin-gonic/gin"
)

//获取多个文章标签
func GetTags(c *gin.Context) {
	//data := make([]string,0)
	db := common.MysqlDB
	var res []modles.Tag
	db.Find(&res)
	var arr []string
	for _, v := range res {
		arr = append(arr,v.TagName)
	}


	r := make(map[string][]string)
	r["tags"] = arr
	c.JSON(200,gin.H{
		"code":"200",
		"msg": "ok",
		"data":r,
	})

}
//  新增文章
func AddTag(c *gin.Context) {
	//db := common.MysqlDB
}

//  修改文章标签
func EditTag(c *gin.Context) {
}

//  删除文章标签
func DeleteTag(c *gin.Context) {
}
