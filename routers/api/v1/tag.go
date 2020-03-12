package v1

import (
	"feelingliu/common"
	"feelingliu/modles"
	"fmt"
	"github.com/gin-gonic/gin"
)

//获取多个文章标签
func GetTags(c *gin.Context) {
	db := common.MysqlDB
	var users modles.User
	db.Where("id = ?",1).Find(&users)
	fmt.Println(users)
	c.JSON(200,gin.H{
		"message":users.Password,
	})
}

//新增文章标签
func AddTag(c *gin.Context) {
}

//修改文章标签
func EditTag(c *gin.Context) {
}

//删除文章标签
func DeleteTag(c *gin.Context) {
}
