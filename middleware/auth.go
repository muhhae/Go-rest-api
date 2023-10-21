package middleware

import (
	"os"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
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
	ctx.Set("userID", claims["id"])
	ctx.Next()
}
