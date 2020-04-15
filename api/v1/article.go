package v1

import (
	"feelingliu/service"
	"feelingliu/utils"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

//  Find all articles
func GetArticles(c *gin.Context) {
	//  Get receive param...
	limit := c.DefaultQuery("limit", "")
	page := c.DefaultQuery("page", "")
	//q := c.DefaultQuery("q", "")
	tag := c.DefaultQuery("tag", "")
	key := c.DefaultQuery("key", "")
	status := c.DefaultQuery("status", "")
	admin := c.DefaultQuery("admin", "")

	if tag != "" {
		fmt.Println("进来了")
		data, e := service.GetArticlesByTag(service.SetLimitPage(limit, page), service.SetAdmin(admin), service.SetTag(tag))
		if e != nil {
			c.JSON(http.StatusInternalServerError, utils.GenResponse(40022, nil, e))
			return
		}
		c.JSON(http.StatusOK, utils.GenResponse(20000, data, nil))
		return
	}

	if key != "" || status != "" {
		data, e := service.SearchArticle(key, status, service.SetLimitPage(limit, page), service.SetAdmin(admin), service.SetSearch(true))
		if e != nil {
			c.JSON(http.StatusInternalServerError, utils.GenResponse(40022, nil, e))
			return
		}
		c.JSON(http.StatusOK, utils.GenResponse(20000, data, nil))
		return
	}
	//if q != "" {
	//	articles, e := service.SearchFromES(service.SetQ(q), service.SetLimitPage(limit, page))
	//	if e != nil {
	//		c.JSON(http.StatusInternalServerError, utils.GenResponse(40033, nil, e))
	//		return
	//	}
	//	c.JSON(http.StatusOK, utils.GenResponse(20000, articles, nil))
	//	return
	//}

	//  no by tag or elasticsearch...
	a := service.Article{}
	data, e := a.GetAll(service.SetLimitPage(limit, page), service.SetAdmin(admin))
	if e != nil {
		c.JSON(http.StatusInternalServerError, utils.GenResponse(40022, nil, e))
		return
	}
	c.JSON(http.StatusOK, utils.GenResponse(20000, data, nil))
	return
}

//  Find article by id
func GetArticle(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	admin := c.DefaultQuery("admin", "")
	r := service.Article{ID: id}

	articleDetail, e := r.GetOne(service.SetAdmin(admin))
	fmt.Println(articleDetail)
	if e != nil {
		c.JSON(http.StatusNotFound, utils.GenResponse(40020, nil, e))
		return
	}
	c.JSON(http.StatusOK, utils.GenResponse(20000, articleDetail, nil))
	return
}

//  json: cannot unmarshal array into Go struct field Article.tag_id of type int
func CreateArticle(c *gin.Context) {
	article := &service.Article{}
	if e := c.ShouldBindJSON(article); e != nil {
		fmt.Println("e:",e)
		c.JSON(http.StatusBadRequest, utils.GenResponse(40024, nil, e))
		return
	}

	a, e := article.Create()
	if e != nil {
		c.JSON(http.StatusInternalServerError, utils.GenResponse(40024, nil, e))
		return
	}
	c.JSON(http.StatusOK, utils.GenResponse(20000, a, nil))
	return
}


