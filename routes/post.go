package routes

import (
	"context"
	"rest-api/connection"
	"rest-api/middleware"
	"rest-api/models"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func InitPostRoute(router *gin.RouterGroup) {
	post := router.Group("/post")
	{
		post.GET("/", func(context *gin.Context) {
			context.JSON(200, gin.H{
				"message": "Post API is working",
			})
		})
		post.POST("/create", middleware.Auth, createPost)
		post.GET("/all", getAllPost)
		post.GET("/detail/:id", getPostDetail)
		post.PUT("/update/:id", updatePost)
		post.DELETE("/delete/:id", deletePost)
	}
}

func createPost(ctx *gin.Context) {
	user, ok := ctx.Get("User")
	if !ok {
		ctx.AbortWithStatusJSON(401, gin.H{
			"error": "Unauthorized",
		})
		return
	}
	post := models.Post{}
	err := ctx.BindJSON(&post)
	if err != nil {
		ctx.AbortWithStatusJSON(400, gin.H{
			"error": err.Error(),
		})
		return
	}
	author := user.(models.User)
	post.AuthorID = author.ID
	post.Author = author.Username
	post.CreatedAt = time.Now()
	post.UpdatedAt = time.Now()
	res, err := connection.Post().InsertOne(context.TODO(), post)
	if err != nil {
		ctx.AbortWithStatusJSON(400, gin.H{
			"error": err.Error(),
		})
		return
	}
	post.ID = res.InsertedID.(primitive.ObjectID)
	ctx.JSON(200, gin.H{
		"message": "Post created successfully",
		"post":    post,
	})
}
func getAllPost(ctx *gin.Context) {

}
func getPostDetail(ctx *gin.Context) {

}
func updatePost(ctx *gin.Context) {

}
func deletePost(ctx *gin.Context) {

}
