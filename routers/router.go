package routers

import "github.com/gin-gonic/gin"

func init() {
	router := gin.Default()
	router.Group("/v1")
	router.Run()
}
