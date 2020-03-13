package middleware

import (
	"feelingliu/modles"
	"feelingliu/public"
	"feelingliu/utils"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

func JWT() gin.HandlerFunc {
	return func(c *gin.Context) {
		var getauth modles.GetAuth
		var code int
		var data interface{}

		code = public.SUCCESS
		err := c.BindJSON(&getauth)
		if err != nil {
			fmt.Println(err)
		}
		if getauth.Token == "" {
			code = public.INVALID_PARAMS
		} else {
			claims, err := utils.ParseToken(getauth.Token)
			if err != nil {
				code = public.ERROR_AUTH_CHECK_TOKEN_FAIL
			} else if time.Now().Unix() > claims.ExpiresAt {
				code = public.ERROR_AUTH_CHECK_TOKEN_TIMEOUT
			}
		}
		if code != public.SUCCESS {
			c.JSON(http.StatusUnauthorized, gin.H{
				"code": code,
				"msg":  public.GetMsg(code),
				"data": data,
			})

			c.Abort()
			return
		}
		c.Next()
	}
}
