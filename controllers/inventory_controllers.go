package controllers

import (
	"lingo/async"
	"lingo/database"
	"lingo/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

func Add2Inventory(c *gin.Context, user_id string, head_id string, lang_iso string){

	future := async.Exec(func() interface{} {
		existed := 0
		database.ConnDB().Get(&existed,`SELECT COUNT(*) FROM inventory WHERE user_id = $1 AND head_id = $2`,user_id, head_id)
		if existed == 0{
			_, err := database.ConnDB().Exec(`INSERT INTO inventory(user_id, head_id, lang_iso)
						VALUES ($1,$2,$3)`,user_id, head_id, lang_iso)
			return err
		}
		return nil
	})
	err := future.Await()
	if err != nil {
		c.JSON(http.StatusNotAcceptable, gin.H{"error": err})
	}else{
		c.JSON(http.StatusOK, gin.H{"status": http.StatusOK})
	}
}

func GetInventory(c *gin.Context) {
	user_info := new(models.GetInventory)
	e := c.BindJSON(&user_info)
	if e != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": e.Error()})
		return
	}
	books := []string{}
	future := async.Exec(func() interface{} {
		return database.ConnDB().Select(&books,`SELECT head_id FROM inventory WHERE user_id=$1 AND lang_iso=$2`,user_info.User_id, user_info.Lang_iso)
	})
	err := future.Await()
	if err != nil {
		c.JSON(http.StatusNotAcceptable, gin.H{"error": err})
	} else {
		c.JSON(http.StatusOK, gin.H{"data": books})
	}
}