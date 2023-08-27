package controllers

import (
	"lingo/async"
	"lingo/database"
	"net/http"

	"github.com/gin-gonic/gin"
)

func add2Inventory(c *gin.Context, user_id string, head_id string, lang_iso string){

	future := async.Exec(func() interface{} {
		tx := database.ConnDB().MustBegin()
		tx.MustExec(`INSERT INTO inventory(user_id, head_id, lang_iso)
					VALUES ($1,$2,$3)`,user_id, head_id, lang_iso)
		return tx.Commit()
	})
	err := future.Await()
	if err != nil {
		c.JSON(http.StatusNotAcceptable, gin.H{"error": err})
	} else {
		c.JSON(http.StatusOK, gin.H{"status": http.StatusOK})
	}
}