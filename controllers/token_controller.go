package controllers

import (
	"fmt"
	"lingo/database"
	"lingo/models"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
)

func GetUser(context *gin.Context) {
	tokenString := context.Request.Header["Authorization"][0]
	token, _, err := new(jwt.Parser).ParseUnverified(tokenString, jwt.MapClaims{})
	if err != nil {
		context.JSON(http.StatusOK, gin.H{"Error": err})
	}
	if claims, ok := token.Claims.(jwt.MapClaims); ok {
		// obtains claims
		// username := claims["username"]
		email := claims["email"]

		var user models.User
		err := database.ConnDB().Table("users_profile").Raw("SELECT * FROM users_profile WHERE email='" + fmt.Sprint(email) + "'").Scan(&user)
		if err != nil {
			context.JSON(http.StatusInternalServerError, gin.H{"error": err})
		}

		context.JSON(http.StatusOK, &user)
	}

}
