package controllers

import (
	"lingo/async"
	"lingo/database"
	"lingo/models"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/lib/pq"
)

func AddNote(c *gin.Context, user_id string, word string, lang_iso string) {
	future := async.Exec(func() interface{} {

		_, err := database.ConnDB().Exec(`
							INSERT INTO note(user_id, word, note, lang_iso)
							VALUES($1, $2, '', $3)`, user_id, word, lang_iso)
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

func GetAllNotes(c *gin.Context) {
	input := new(models.InputGetAllNote)
	e := c.BindJSON(&input)
	if e != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": e.Error()})
		return
	}
	notes := []models.Note{}
	future := async.Exec(func() interface{} {
		database.ConnDB().Select(&notes, `SELECT * FROM note WHERE user_id=$1 AND lang_iso=$2`, input.User_id, input.Lang_iso)
		return database.ConnDB().Close()
	})
	future.Await()

	c.JSON(http.StatusOK, gin.H{"data": &notes})
}

func GetNote(c *gin.Context) {
	words := new(models.InputGetNote)
	e := c.BindJSON(&words)
	if e != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": e.Error()})
		return
	}
	note := []models.Note{}
	future := async.Exec(func() interface{} {
		database.ConnDB().Select(&note, `SELECT * FROM note WHERE user_id=$1 AND word=ANY($2)`, words.User_id, pq.Array(words.Words))
		return database.ConnDB().Close()
	})
	future.Await()

	c.JSON(http.StatusOK, gin.H{"data": &note} )
}

func EditNote(c *gin.Context) {

	note_input := new(models.Note)
	e := c.BindJSON(&note_input)
	if e != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": e.Error()})
		return
	}
	future := async.Exec(func() interface{} {
		tx := database.ConnDB().MustBegin()
		tx.MustExec(`UPDATE note 
					SET note=$1
					WHERE user_id=$2 AND word=$3 AND lang_iso=$4`, note_input.Note, note_input.User_id, note_input.Word, note_input.Lang_iso)
		tx.Commit()
		return database.ConnDB().Close()
	})
	err := future.Await()
	if err != nil {
		c.JSON(http.StatusNotAcceptable, gin.H{"error": err})
	} else{
		c.JSON(http.StatusOK, gin.H{"status": http.StatusOK})
	}
}
