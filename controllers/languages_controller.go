package controllers

import (
	"lingo/async"
	"lingo/database"
	"lingo/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

func GetLangs(c *gin.Context) {
	var langs models.Languages
	future := async.Exec(func() interface{} {
		return database.ConnDB().Table("users_profile").Exec("SELECT * FROM LANGUAGES").Scan(&langs)
	})
	err := future.Await()
	if err != nil {
		c.JSON(http.StatusNotAcceptable, gin.H{"error": err})
	} else {
		// c.JSON(http.StatusOK, gin.H{"data": user, "status": http.StatusOK})
		c.JSON(http.StatusOK, gin.H{"status": http.StatusOK})
	}
}
