package routes

import (
	"github.com/gin-gonic/gin"
)

func SetRoutes(router *gin.Engine) {
	home := router.Group("/")
	{
		home.GET("/", func(context *gin.Context) {
			context.JSON(200, gin.H{
				"message": "API is working",
			})
		})

		InitUserRoute(home)
		InitPostRoute(home)
	}
}
