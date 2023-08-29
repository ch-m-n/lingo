package controllers

import (
	"lingo/async"
	"lingo/database"
	"lingo/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

func GetWord(c *gin.Context) {
	word := new(models.Word)
	e := c.BindJSON(&word)
	if e != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": e.Error()})
	}
	var words models.Word
	future := async.Exec(func() interface{} {
		return database.ConnDB().Get(&words, "SELECT * FROM words WHERE word=$1", word.Word)
	})
	err := future.Await()
	if err != nil {
		c.JSON(http.StatusNotAcceptable, gin.H{"error": err})
	} else {
		c.JSON(http.StatusOK, &words)
	}
}

func AddWord(c *gin.Context, word_level []models.Literacy) {

	future := async.Exec(func() interface{} {
		_, err := database.ConnDB().NamedExec(`INSERT INTO words(word, lang_iso) 
					VALUES(:word, :lang_iso) ON CONFLICT DO NOTHING`, word_level)
		return err
	})
	err := future.Await()
	if err != nil {
		c.JSON(http.StatusNotAcceptable, gin.H{"error": err})
	}else{
		c.JSON(http.StatusOK, gin.H{"status": http.StatusOK})
	}
}
