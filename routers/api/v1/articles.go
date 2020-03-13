package v1

import "github.com/gin-gonic/gin"


//  获取所有文章
func GetArticles(c *gin.Context)  {
	c.JSON(200,gin.H{
		"liuhanwen":"qianjiayu",
	})
}

//  根据分类获取文章
func GetArticlesByTag(c *gin.Context)  {

}

//  添加文章
func AddArticle(c *gin.Context)  {

}

//  修改文章
func Editarticle(c *gin.Context)  {

}

//  删除文章
func DeleteArticle(c *gin.Context) {

}

