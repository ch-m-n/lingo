package controllers

import (
	"lingo/auth"
	"lingo/database"
	"lingo/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

func CreateUser(c *gin.Context) {
	user := new(models.CreateUser)
	e := c.BindJSON(&user)
	if e != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": e.Error()})
		return
	}

	err := database.ConnDB().Table("users_profile").Exec("INSERT INTO users_profile(id, username, email, pwd, created_at, edited_at) " +
		"VALUES(gen_random_uuid(), '" + user.Username + "', '" + user.Email + "', '" + user.Pwd + "', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP) ").Error
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
	err := database.ConnDB().Table("users_profile").Raw("SELECT * FROM users_profile WHERE email='" + getUser.Email + "'").Scan(&user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err})
	}
	if user.Email != "" {
		if models.VerifyHash(getUser.Pwd, user.Pwd) {
			tokenString, err := auth.GenerateJWT(user.Email, user.Username)
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
