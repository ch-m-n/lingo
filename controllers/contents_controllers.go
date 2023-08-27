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
		tx.MustExec(`INSERT INTO inventory(user_id, head_id, lang_iso)
					VALUES($1,$2,$3)`,content_info.User_id, content_info.Id, content_info.Lang_iso)
		tx.Commit()
		return database.ConnDB().Get(&content,"SELECT * FROM contents WHERE title=$1 AND id =$2",content_info.Title, content_info.Id)
	})
	future.Await()
	c.JSON(http.StatusOK, gin.H{"data": content})
}

func GetAllContents(c *gin.Context) {
	user := new(models.RequestAllContent)
	e := c.BindJSON(&user)
	if e != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": e.Error()})
		return
	}
	
	content_titles := []models.Content{}
	future := async.Exec(func() interface{} {
		return database.ConnDB().Select(&content_titles, 
										`SELECT DISTINCT ON(title) id, user_id, title , lang_iso, created_at , edited_at , img  
										FROM contents 
										WHERE user_id=$1 AND lang_iso=$2`,
										user.User_id, user.Lang_iso)
	})
	future.Await()
	c.JSON(http.StatusOK, gin.H{"data": content_titles} )
}

func AddContents(c *gin.Context) {
	content_input := new(models.Content)
	e := c.BindJSON(&content_input)
	if e != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": e.Error()})
		return
	}
	
	future := async.Exec(func() interface{} {
		tx.MustExec(`INSERT INTO contents(id, user_id, title, lang_iso, body, created_at, edited_at, img)
					VALUES(gen_random_uuid(),$1,$2,$3,$4,CURRENT_TIMESTAMP,CURRENT_TIMESTAMP,$5)`, 
					content_input.User_id,content_input.Title,content_input.Lang_iso,content_input.Body,content_input.Img)
		return tx.Commit()
	})
	err := future.Await()
	if err != nil {
		c.JSON(http.StatusNotAcceptable, gin.H{"error": err})
	} else {
		c.JSON(http.StatusOK, gin.H{"status": http.StatusOK})
	}
}

func EditContent(c *gin.Context) {
	content_input := new(models.Content)
	e := c.BindJSON(&content_input)
	if e != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": e.Error()})
		return
	}
	future := async.Exec(func() interface{} {
		_, err := database.ConnDB().Exec(`UPDATE contents
					SET title=$1, body=$2, edited_at=CURRENT_TIMESTAMP
					WHERE user_id=$3 AND id=$4 AND lang_iso=$5`,
					content_input.Title, content_input.Body, content_input.User_id, content_input.Id, content_input.Lang_iso)
		return err	
	})
	err := future.Await()
	if err != nil {
		c.JSON(http.StatusNotAcceptable, gin.H{"error": err})
	}
}
