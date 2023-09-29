package controllers

import (
	"lingo/async"
	"lingo/database"
	"lingo/models"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/lib/pq"
)

func GetAllWordLevel(c *gin.Context) {
	words := new(models.InputGetLiteracy)
	e := c.BindJSON(&words)
	if e != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": e.Error()})
		return
	}
	word_level := []models.Literacy{}
	future := async.Exec(func() interface{} {
		return database.ConnDB().Select(&word_level,
			`SELECT * FROM literacy WHERE user_id=$1 AND lang_iso=$2`,
			words.User_id, words.Lang_iso)
	})
	database.ConnDB().Close()
	future.Await()
	c.JSON(http.StatusOK, gin.H{"words": word_level})
}

func GetWordLevel(c *gin.Context) {
	words := new(models.InputGetLiteracy)
	e := c.BindJSON(&words)
	if e != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": e.Error()})
		return
	}
	word_level := []models.Literacy{}
	future := async.Exec(func() interface{} {
		return database.ConnDB().Select(&word_level,
			`SELECT * FROM literacy 
			WHERE user_id=$1 AND lang_iso=$2 AND word=ANY($3)`,
			words.User_id, words.Lang_iso, pq.Array(words.Words))
	})
	database.ConnDB().Close()
	future.Await()
	c.JSON(http.StatusOK, gin.H{"words": word_level})
}

func AddWordLevel(c *gin.Context) {
	words := new(models.Literacy)
	e := c.BindJSON(&words)
	if e != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": e.Error()})
	}
	// AddWord(c, words.Words)
	future := async.Exec(func() interface{} {
		existed := 0
		database.ConnDB().Get(&existed, `SELECT COUNT(*) FROM literacy WHERE user_id = $1 AND word = $2`, words.User_id, words.Word)
		if existed == 0 {
			AddNote(c, words.User_id, words.Word, words.Lang_iso)
			_, err := database.ConnDB().Exec(`INSERT INTO literacy(user_id, word, lang_iso, known_level)
												VALUES($1,$2,$3,0)`, words.User_id, words.Word, words.Lang_iso)
			return err
		} else {
			_, err := database.ConnDB().Exec(`UPDATE literacy SET known_level=$3
												WHERE user_id=$1 AND word=$2`, words.User_id, words.Word, words.Known_level)
			return err
		}
	})
	database.ConnDB().Close()
	err := future.Await()
	if err != nil {
		c.JSON(http.StatusNotAcceptable, gin.H{"error": err})
	} else {
		c.JSON(http.StatusOK, gin.H{"status": http.StatusOK})
	}

}
