package controllers

import (
	"lingo/async"
	"lingo/database"
	"lingo/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

func GetNote(c *gin.Context) {
	words := new(models.InputGetNote)
	e := c.BindJSON(&words)
	if e != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": e.Error()})
		return
	}
	var notes models.OutputGetNote
	var note models.Note
	for i := 0; i < len(words.Words); i++ {
		future := async.Exec(func() interface{} {
			return database.ConnDB().Table("note").Raw("SELECT * FROM note "+
														"WHERE word='" + words.Words[i] + "' "+
														"AND user_id='"+words.User_id+"'").Scan(&note)
		})
		future.Await()
		notes.Notes = append(notes.Notes, note)
	}
	c.JSON(http.StatusOK, gin.H{"data": notes, "status": http.StatusOK})
}

func EditNote(c *gin.Context) {

	notes := new(models.InputUpdateNote)
	e := c.BindJSON(&notes)
	if e != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": e.Error()})
		return
	}
	for i := 0; i < len(notes.Notes); i++ {
		future := async.Exec(func() interface{} {
			return database.ConnDB().Table("note").Exec("UPDATE note SET note='"+notes.Notes[i].Note+"' "+
														"WHERE user_id='"+notes.Notes[i].User_id.String()+"' "+
														"AND word='"+notes.Notes[i].Word+"'").Error
		})
		err := future.Await()
		if err != nil {
			c.JSON(http.StatusNotAcceptable, gin.H{"error": err})
			break
		} else {
			c.JSON(http.StatusOK, gin.H{"status": http.StatusOK})
		}
	}
}