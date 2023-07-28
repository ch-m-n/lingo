package middlewares

import (
	"lingo/auth"
	"net/http"

	"github.com/gin-gonic/gin"
)

func Auth() gin.HandlerFunc{
	return func(context *gin.Context) {
		tokenString := context.GetHeader("Authorization")
		if tokenString == "" {
			context.JSON(http.StatusOK, gin.H{"error": "request does not contain an access token"})
			context.Abort()
			return
		}
		err:= auth.ValidateToken(tokenString)
		if err != nil {
			context.JSON(http.StatusOK, gin.H{"error": err.Error()})
			context.Abort()
			return
		}
		context.Next()
	}
}