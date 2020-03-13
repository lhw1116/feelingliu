package v1

import (
	"feelingliu/common"
	"feelingliu/modles"
	"github.com/gin-gonic/gin"
	"net/http"
)

//获取多个文章标签
func GetTags(c *gin.Context) {
	data := make([]string,0)
	db := common.MysqlDB
	var res []modles.Tag
	db.Select("tag_name").Find(&res)
	for _, v := range res {
		data = append(data, v.TagName)
	}
	c.JSON(http.StatusOK,gin.H{
		"code":"200",
		"msg": "ok",
		"data":data,
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
