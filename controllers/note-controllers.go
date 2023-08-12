package controllers

import (
	"lingo/async"
	"lingo/database"
	"lingo/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

func GetAllNotes(c *gin.Context) {
	input := new(models.InputGetAllNote)
	e := c.BindJSON(&input)
	if e != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": e.Error()})
		return
	}
	var notes []models.Note
		future := async.Exec(func() interface{} {
			return database.ConnDB().Table("note").Where("user_id=?",input.User_id).Where("lang_iso=?",input.Lang_iso).Find(&notes)
		})
		future.Await()
	
	c.JSON(http.StatusOK, gin.H{"data": notes, "status": http.StatusOK})
}

func GetNote(c *gin.Context) {
	words := new(models.InputGetNote)
	e := c.BindJSON(&words)
	if e != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": e.Error()})
		return
	}
	var notes []models.Note
	future := async.Exec(func() interface{} {
		return database.ConnDB().Table("note").Where("word IN?", words.Words).Where("user_id=?", words.User_id).Find(&notes)
	})
	future.Await()

	c.JSON(http.StatusOK, gin.H{"data": notes, "status": http.StatusOK})
}

func EditNote(c *gin.Context) {

	note_input := new(models.Note)
	e := c.BindJSON(&note_input)
	if e != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": e.Error()})
		return
	}
	future := async.Exec(func() interface{} {
		note := models.Note{
			User_id: note_input.User_id,
			Word: note_input.Word,
			Note: note_input.Note,
		}
		return database.ConnDB().Table("note").Where("word=?",note_input.Word).Updates(&note).Error
	})
	err := future.Await()
	if err != nil {
		c.JSON(http.StatusNotAcceptable, gin.H{"error": err})
	} else {
		c.JSON(http.StatusOK, gin.H{"status": http.StatusOK})
	}

}
