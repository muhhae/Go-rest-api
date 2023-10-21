package middleware

import (
	"context"
	"os"
	"rest-api/connection"
	"rest-api/models"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func Auth(ctx *gin.Context) {
	tokenString := ctx.GetHeader("Token")
	if tokenString == "" {
		ctx.AbortWithStatusJSON(401, gin.H{
			"error": "Token is required",
		})
		return
	}
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		_, ok := token.Method.(*jwt.SigningMethodHMAC)
		if !ok {
			ctx.AbortWithStatusJSON(401, gin.H{
				"error": "Unauthorized",
			})
			return nil, nil
		}
		return []byte(os.Getenv("SECRET_KEY")), nil
	})
	if err != nil {
		ctx.AbortWithStatusJSON(401, gin.H{
			"error": err.Error(),
		})
		return
	}
	claims, ok := token.Claims.(jwt.MapClaims)
	// log.Println("claim", claims, "ok", ok, "valid", token.Valid)
	if !ok || !token.Valid {
		ctx.AbortWithStatusJSON(401, gin.H{
			"error": "Unauthorized",
		})
		return
	}
	userID, err := primitive.ObjectIDFromHex(claims["id"].(string))
	if err != nil {
		ctx.AbortWithStatusJSON(401, gin.H{
			"error": err.Error(),
		})
		return
	}
	user_data := connection.User().FindOne(context.TODO(), bson.M{"_id": userID})
	if user_data.Err() != nil {
		ctx.AbortWithStatusJSON(401, gin.H{
			"error": user_data.Err().Error(),
		})
		return
	}
	user := models.User{}
	err = user_data.Decode(&user)
	if err != nil {
		ctx.AbortWithStatusJSON(401, gin.H{
			"error": err.Error(),
		})
		return
	}
	ctx.Set("User", user)
	ctx.Next()
}
