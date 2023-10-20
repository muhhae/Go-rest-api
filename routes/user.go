package routes

import (
	"context"
	"rest-api/connection"
	"rest-api/models"

	"github.com/gin-gonic/gin"
)

func InitUserRoute(router *gin.RouterGroup) {
	user := router.Group("/user")
	{
		user.GET("/", func(context *gin.Context) {
			context.JSON(200, gin.H{
				"message": "User API is working",
			})
		})
		user.GET("/sign-in", signIn)
		user.POST("/sign-up", signUp)
	}
}

func signIn(context *gin.Context) {
	if context.Request.Body == nil {
		context.JSON(400, gin.H{
			"message": "Bad request",
		})
		return
	}
	data := make(map[string]interface{})
	err := context.BindJSON(&data)

	if err != nil {
		context.AbortWithStatusJSON(400, gin.H{
			"error": err.Error(),
		})
		return
	}

	if data["username"] == nil || data["password"] == nil {
		context.AbortWithStatusJSON(400, gin.H{
			"error": "Username or password is missing",
		})
		return
	}

	if data["username"] != "admin" || data["password"] != "admin" {
		context.AbortWithStatusJSON(401, gin.H{
			"error": "Username or password is incorrect",
		})
		return
	}

	context.JSON(200, gin.H{
		"message": "Sign in successfully",
	})
}

func signUp(ctx *gin.Context) {
	if ctx.Request.Body == nil {
		ctx.JSON(400, gin.H{
			"message": "Bad request",
		})
		return
	}
	data := make(map[string]interface{})
	err := ctx.BindJSON(&data)

	if err != nil {
		ctx.AbortWithStatusJSON(400, gin.H{
			"error": err.Error(),
		})
		return
	}

	if data["username"] == nil || data["password"] == nil {
		ctx.AbortWithStatusJSON(400, gin.H{
			"error": "Username or password is missing",
		})
		return
	}

	user_col := connection.User()
	new_user := models.User{Username: data["username"].(string), Password: data["password"].(string)}

	res, err := user_col.InsertOne(context.TODO(), new_user)
	if err != nil {
		ctx.AbortWithStatusJSON(400, gin.H{
			"message": "Sign up failed",
			"error":   err.Error(),
		})
	}

	ctx.JSON(200, gin.H{
		"message": "Sign up successfully",
		"newUserID":  res.InsertedID,
	})
}
