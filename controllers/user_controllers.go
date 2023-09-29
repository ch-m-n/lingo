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
		_, err := database.ConnDB().Exec(`INSERT INTO users_profile(id, username, email, pwd, created_at, edited_at, verified) 
										VALUES(gen_random_uuid(), $1, $2, $3, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP, $4)`, 
					user.Username, user.Email, string(models.PassHash(user.Pwd)), "false")
		return err
	})
	database.ConnDB().Close()
	err := future.Await()
	if err != nil {
		c.JSON(http.StatusNotAcceptable, gin.H{"error": err})
	}else{
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
	user := models.User{}
	future := async.Exec(func() interface{} {
		return database.ConnDB().Get(&user,"SELECT * FROM users_profile WHERE email=$1", getUser.Email)
	})
	database.ConnDB().Close()
	err := future.Await()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err})
	}
	if user.Email != "" {
		if models.VerifyHash(getUser.Pwd, user.Pwd) {
			tokenString, err := auth.GenerateJWT(user.Id.String(), user.Email, user.Username)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				c.Abort()
				return
			}
			c.JSON(http.StatusOK, gin.H{"result": "Correct", "token": tokenString, "status": http.StatusOK})
		} else {
			c.JSON(http.StatusOK, gin.H{"result": "Incorrect Email or Password", "user": user, "status": http.StatusOK})
			c.Abort()
		}
	} else {
		c.JSON(http.StatusOK, gin.H{"result": user, "status": http.StatusOK})
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

		user := models.User{}
		future := async.Exec(func() interface{} {
			return database.ConnDB().Get(&user,"SELECT * FROM users_profile WHERE email=$1", fmt.Sprint(email))
		})
		database.ConnDB().Close()
		err := future.Await()

		if err != nil {
			context.JSON(http.StatusInternalServerError, gin.H{"error": err})
		}

		context.JSON(http.StatusOK, &user)
	}

}

func EditUser(c *gin.Context) {
	user_info := new(models.EditUser)
	e := c.BindJSON(&user_info)
	if e != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": e.Error()})
		return
	}

	future := async.Exec(func() interface{} {
		_, err := database.ConnDB().Exec(`UPDATE users_profile SET username=$1, pwd=$2, edited_at=CURRENT_TIMESTAMP 
					WHERE id=$3`, 
					user_info.Username, string(models.PassHash(user_info.Pwd)), user_info.Id)
		return err
	})
	database.ConnDB().Close()
	err := future.Await()
	if err != nil {
		c.JSON(http.StatusNotAcceptable, gin.H{"error": err})
	}else{
		c.JSON(http.StatusOK, gin.H{"status": http.StatusOK})
	}
}