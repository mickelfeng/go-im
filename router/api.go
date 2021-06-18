/**
  @author:panliang
  @data:2021/6/18
  @note
**/
package router

import "github.com/gin-gonic/gin"

var router *gin.Engine

func RegisterApiRoutes(router *gin.Engine)  {
	router.GET("/", func(context *gin.Context) {
		context.JSON(200, gin.H{
			"message": "hello world!!!",
		})
	})
}