package controllers


import (
	"net/http"
	"time"

	"github.com/ch-m-n/lingo/db"
	"github.com/ch-m-n/lingo/models"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func CreateUser(c *gin.Context) {
	//Input example
	// {
	// 	"username": "Bizet",
	// 	"email": "bizet@",
	// 	"phone_num": "12345"
	// }
	user := new(models.User)
	e := c.BindJSON(&user)
	switch e != nil {
	case true:
		c.JSON(http.StatusBadRequest, gin.H{"message": e})
		return
	}

	userPayload := models.User{
		Id:         uuid.New().String(),
		Username:   user.Username,
		Email:      user.Email,
		Pwd:  		user.Pwd,
		Created_At: time.Now().Local(),
	}

	err := db.Dbconn().Table("users").Create(&userPayload).Error
	switch err != nil {
	case true:
		c.JSON(http.StatusInternalServerError, gin.H{"error": err})
	}

	c.JSON(http.StatusOK, gin.H{"data": userPayload, "status": http.StatusOK})
}

func GetUser(c *gin.Context){
	var result models.User
	switch err := db.Dbconn().Where("username = ?", c.Param("username")).First(&result).Error; err != nil {
	case true:
		c.JSON(http.StatusInternalServerError, gin.H{"error": err})
	}
	c.JSON(http.StatusOK, gin.H{"data":result})
}

