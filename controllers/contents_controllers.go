package controllers

import (
	"lingo/async"
	"lingo/database"
	"lingo/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

func GetContents(c *gin.Context) {
	content_info := new(models.RequestContent)
	e := c.BindJSON(&content_info)
	if e != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": e.Error()})
		return
	}
	var content models.Content
	future := async.Exec(func() interface{} {
		return database.ConnDB().Table("contents").Raw("SELECT * FROM contents WHERE title='" + content_info.Title +
			"' AND lang_iso='" + content_info.Lang_iso + "'").Scan(&content)
	})
	err := future.Await()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err})
	}
}

func AddContents(c *gin.Context) {
	content := new(models.Content)
	e := c.BindJSON(&content)
	if e != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": e.Error()})
		return
	}
	future := async.Exec(func() interface{} {
		return database.ConnDB().Table("contents").Exec("INSERT INTO contents(id, user_id, title, lang_iso, body, created_at, edited_at) " +
			"VALUES(gen_random_uuid(), '" + content.User_id.String() + "', '" + content.Title + "', '" + content.Lang_iso + "', '" + content.Body + "', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP) ").Error
	})
	err := future.Await()
	if err != nil {
		c.JSON(http.StatusNotAcceptable, gin.H{"error": err})
	} else {
		// c.JSON(http.StatusOK, gin.H{"data": user, "status": http.StatusOK})
		c.JSON(http.StatusOK, gin.H{"status": http.StatusOK})
	}
}

func EditContent(c *gin.Context) {

}
