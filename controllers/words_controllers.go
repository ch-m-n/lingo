package controllers

import (
	"lingo/async"
	"lingo/database"
	"lingo/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

func GetWord(c *gin.Context){
	word := new(models.Word)
	e := c.BindJSON(&word)
	if e != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": e.Error()})
		return
	}
	var words models.Word
	future := async.Exec(func() interface{} {
			return database.ConnDB().Table("words").Exec("SELECT * FROM words WHERE word='"+word.Word+"'").Scan(&words)
	})
	err := future.Await()
	if err != nil {
		c.JSON(http.StatusNotAcceptable, gin.H{"error": err})
	} else {
		c.JSON(http.StatusOK, &words)
		// c.JSON(http.StatusOK, gin.H{"status": http.StatusOK})
	}
}

func AddWord(c *gin.Context, word string, lang_iso string) {
// func AddWord(c *gin.Context) {
	// word := new(models.Word)
	// e := c.BindJSON(&word)
	// if e != nil {
	// 	c.JSON(http.StatusBadRequest, gin.H{"message": e.Error()})
	// 	return
	// }

	future := async.Exec(func() interface{} {
			return database.ConnDB().Table("words").Exec("INSERT INTO words(word, lang_iso) " +
																"VALUES('"+word+"','"+lang_iso+"') ON CONFLICT DO NOTHING").Error
	})
	err := future.Await()
	if err != nil {
		c.JSON(http.StatusNotAcceptable, gin.H{"error": err})
	} else {
		// c.JSON(http.StatusOK, gin.H{"data": user, "status": http.StatusOK})
		c.JSON(http.StatusOK, gin.H{"status": http.StatusOK})
		return
	}
}