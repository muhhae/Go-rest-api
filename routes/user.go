package routes

import (
	"context"
	"os"
	"rest-api/connection"
	"rest-api/middleware"
	"rest-api/models"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
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
		user.GET("/profile", middleware.Auth, profile)
	}
}

func signIn(context *gin.Context) {
	if context.Request.Body == nil {
		context.JSON(400, gin.H{
			"message": "Bad request",
		})
		return
	}
	user_input := models.User{}
	err := context.BindJSON(&user_input)
	if err != nil {
		context.AbortWithStatusJSON(400, gin.H{
			"error": err.Error(),
		})
		return
	}
	userData := connection.User().FindOne(context, bson.M{"email": user_input.Email})
	if userData.Err() == mongo.ErrNoDocuments {
		context.AbortWithStatusJSON(400, gin.H{
			"error": "Email is not registered",
		})
		return
	}
	user := models.User{}
	err = userData.Decode(&user)
	if err != nil {
		context.AbortWithStatusJSON(400, gin.H{
			"error": err.Error(),
		})
	}
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(user_input.Password)); err != nil {
		context.AbortWithStatusJSON(400, gin.H{
			"error": "Wrong password",
		})
		return
	}
	token := jwt.New(jwt.SigningMethodHS256)
	token_claim := token.Claims.(jwt.MapClaims)
	token_claim["id"] = user.ID
	token_claim["exp"] = time.Now().Add(time.Hour * 24).Unix()
	token_string, err := token.SignedString([]byte(os.Getenv("SECRET_KEY")))
	if err != nil {
		context.AbortWithStatusJSON(400, gin.H{
			"error": err.Error(),
		})
	}

	context.JSON(200, gin.H{
		"message": "Sign in successfully",
		"token":   token_string,
		"user":    user,
	})
}

func signUp(ctx *gin.Context) {
	if ctx.Request.Body == nil {
		ctx.JSON(400, gin.H{
			"message": "Bad request",
		})
		return
	}
	new_user := models.User{}
	err := ctx.BindJSON(&new_user)

	if err != nil {
		ctx.AbortWithStatusJSON(400, gin.H{
			"error": err.Error(),
		})
		return
	}

	if new_user.Username == "" || new_user.Password == "" || new_user.Email == "" {
		ctx.AbortWithStatusJSON(400, gin.H{
			"error": "Email or username or password is missing",
		})
		return
	}

	if count, err := connection.User().CountDocuments(context.TODO(), bson.M{"email": new_user.Email}); err == nil && count > 0 {
		ctx.AbortWithStatusJSON(400, gin.H{
			"error": "Email is already taken",
		})
		return
	} else if err != nil {
		ctx.AbortWithStatusJSON(400, gin.H{
			"error": err.Error(),
		})
		return
	}

	if count, err := connection.User().CountDocuments(context.TODO(), bson.M{"username": new_user.Username}); err == nil && count > 0 {
		ctx.AbortWithStatusJSON(400, gin.H{
			"error": "Username is already taken",
		})
		return
	} else if err != nil {
		ctx.AbortWithStatusJSON(400, gin.H{
			"error": err.Error(),
		})
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(new_user.Password), 10)
	if err != nil {
		ctx.AbortWithStatusJSON(400, gin.H{
			"error": err.Error(),
		})
		return
	}

	new_user.Password = string(hashedPassword)
	new_user.Role = "user"
	new_user.Verified = false

	res, err := connection.User().InsertOne(context.TODO(), new_user)
	new_user.ID = res.InsertedID.(primitive.ObjectID)

	if err != nil {
		ctx.AbortWithStatusJSON(400, gin.H{
			"message": "Sign up failed",
			"error":   err.Error(),
		})
	}

	ctx.JSON(200, gin.H{
		"message":       "Sign up successfully",
		"new_user_data": new_user,
	})
}

func profile(ctx *gin.Context) {
	userID, ok := ctx.Get("userID")
	if !ok {
		ctx.AbortWithStatusJSON(400, gin.H{
			"error": "User ID is not found",
		})
		return
	}
	id, err := primitive.ObjectIDFromHex(userID.(string))
	if err != nil {
		ctx.AbortWithStatusJSON(400, gin.H{
			"error": err.Error(),
		})
	}
	userData := connection.User().FindOne(context.TODO(), bson.M{"_id": id})
	if userData.Err() == mongo.ErrNoDocuments {
		ctx.AbortWithStatusJSON(400, gin.H{
			"error": "User not found",
		})
		return
	}
	user := models.User{}
	err = userData.Decode(&user)
	if err != nil {
		ctx.AbortWithStatusJSON(400, gin.H{
			"error": err.Error(),
		})
	}
	ctx.JSON(200, gin.H{
		"user": user,
	})
}
