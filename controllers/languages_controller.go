package controllers

import (
	"lingo/async"
	"lingo/database"
	"lingo/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

func GetLangs(c *gin.Context) {
	var langs []models.Languages
	future := async.Exec(func() interface{} {
		return database.ConnDB().Get(&langs,"SELECT * FROM LANGUAGES")
	})
	err := future.Await()
	if err != nil {
		c.JSON(http.StatusNotAcceptable, gin.H{"error": err})
	} else {
		c.JSON(http.StatusOK, gin.H{"data": langs})
	}
}
