package api

import (
	"feelingliu/public"
	"feelingliu/utils"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
)

type Auth struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func Getauth(c *gin.Context) {
	var a Auth

	data := make(map[string]interface{})
	code := public.INVALID_PARAMS

	err := c.BindJSON(&a)
	if err != nil {
		fmt.Println(err)
	} else {
		isExist := utils.CheckAuth(a.Username, a.Password)
		if isExist {
			token, err := utils.GenerateToken(a.Username, a.Password)
			if err != nil {
				code = public.ERROR_AUTH_TOKEN
			} else {
				data["token"] = token
				code = public.SUCCESS
			}
		} else {
			code = public.ERROR_USERPASS
		}
	}
	c.JSON(http.StatusOK, gin.H{
		"code": code,
		"msg":  public.GetMsg(code),
		"data": data,
	})

}
