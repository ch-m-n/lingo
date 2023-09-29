package controllers

import (
	"lingo/async"
	"lingo/database"
	"lingo/models"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/lib/pq"
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

func AddWord(c *gin.Context, words []string, lang string) {

	future := async.Exec(func() interface{} {

		return database.ConnDB().MustExec(`INSERT INTO words(word, lang_iso)
										SELECT word, $2
										FROM UNNEST(CAST($1 as text[])) T (word)
										WHERE NOT EXISTS (SELECT * FROM words WHERE word = T.word)`, pq.Array(words), &lang)
	})
	database.ConnDB().Close()
	future.Await()
	// if err != nil {
	// 	c.JSON(http.StatusNotAcceptable, gin.H{"error": err})
	// } else {
	// 	c.JSON(http.StatusOK, gin.H{"status": http.StatusOK})
	// }
}
