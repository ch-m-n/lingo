package controllers

import (
	"fmt"
	"lingo/async"
	"lingo/auth"
	"lingo/database"
	"lingo/models"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
)

func CreateUser(c *gin.Context) {
	user := new(models.CreateUser)
	e := c.BindJSON(&user)
	if e != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": e.Error()})
		return
	}

	future := async.Exec(func() interface{} {
			return database.ConnDB().Table("users_profile").Exec("INSERT INTO users_profile(id, username, email, pwd, created_at, edited_at, verified) " +
			"VALUES(gen_random_uuid(), '" + user.Username + "', '" + user.Email + "', '" + string(models.PassHash(user.Pwd)) + "', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP, false) ").Error
	})
	err := future.Await()
	if err != nil {
		c.JSON(http.StatusNotAcceptable, gin.H{"error": err})
	} else {
		// c.JSON(http.StatusOK, gin.H{"data": user, "status": http.StatusOK})
		c.JSON(http.StatusOK, gin.H{"status": http.StatusOK})
	}
}

func VerifyUser(c *gin.Context) {
	getUser := new(models.GetUser)
	e := c.BindJSON(&getUser)
	if e != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": e.Error()})
		return 
	}
	var user models.User
	future := async.Exec(func() interface{} {
		return database.ConnDB().Table("users_profile").Raw("SELECT * FROM users_profile WHERE email='" + getUser.Email + "'").Scan(&user)
	})
	err := future.Await()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err})
	}
	if user.Email != "" {
		if models.VerifyHash(getUser.Pwd, user.Pwd) {
			tokenString, err := auth.GenerateJWT(user.Id.String(),user.Email, user.Username)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				c.Abort()
				return
			}
			c.JSON(http.StatusOK, gin.H{"result": "Correct", "token":tokenString, "status": http.StatusOK})
		} else {
			c.JSON(http.StatusOK, gin.H{"result": "Incorrect Email", "status": http.StatusOK})
			c.Abort()
		}
	} else {
		c.JSON(http.StatusOK, gin.H{"result": "Incorrect Email", "status": http.StatusOK})
		c.Abort()
	}
}



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
		future := async.Exec(func() interface{} {
			return database.ConnDB().Table("users_profile").Raw("SELECT * FROM users_profile WHERE email='" + fmt.Sprint(email) + "'").Scan(&user)
		})
		err := future.Await()
		
		if err != nil {
			context.JSON(http.StatusInternalServerError, gin.H{"error": err})
		}

		context.JSON(http.StatusOK, &user)
	}

}
