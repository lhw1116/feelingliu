package v1

import (
	"feelingliu/service"
	"feelingliu/utils"
	"github.com/gin-gonic/gin"
	"net/http"
)

func GetAllTags(c *gin.Context) {
	tags, e := service.Tag{}.GetAll()
	if e != nil {
		c.JSON(http.StatusInternalServerError, utils.GenResponse(40008, nil, e))
	}
	c.JSON(http.StatusOK, utils.GenResponse(20000, tags, nil))
}
