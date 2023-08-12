package controllers

import (
	"lingo/async"
	"lingo/database"
	"lingo/models"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

func GetContents(c *gin.Context) {
	content_info := new(models.RequestContent)
	e := c.BindJSON(&content_info)
	if e != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": e.Error()})
		return
	}
	var content []models.Content
	future := async.Exec(func() interface{} {
		return database.ConnDB().Table("contents").Where("title=?",content_info.Title).Find(&content)
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
	
	var content_titles []models.Content
	future := async.Exec(func() interface{} {
		return database.ConnDB().Select("title").Where("user_id=?",user.User_id.String()).Where("lang_iso=?",user.Lang_iso).Distinct("title").Find(&content_titles)
	})
	future.Await()
	c.JSON(http.StatusOK, &content_titles)
}

func AddContents(c *gin.Context) {
	content_input := new(models.Content)
	e := c.BindJSON(&content_input)
	if e != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": e.Error()})
		return
	}
	
	future := async.Exec(func() interface{} {
		content := models.Content{
			Id: content_input.Id, 
			User_id: content_input.User_id,
			Title: content_input.Title,
			Lang_iso: content_input.Lang_iso,
			Body: content_input.Body,
			Created_at: content_input.Created_at,
			Edited_at: content_input.Edited_at,
			Img: content_input.Img,
		}
		return database.ConnDB().Create(&content)
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
	content_input := new(models.Content)
	e := c.BindJSON(&content_input)
	if e != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": e.Error()})
		return
	}
	future := async.Exec(func() interface{} {
		content := models.Content{
			Id: content_input.Id, 
			User_id: content_input.User_id,
			Title: content_input.Title,
			Lang_iso: content_input.Lang_iso,
			Body: content_input.Body,
			Created_at: content_input.Created_at,
			Edited_at: content_input.Edited_at,
			Img: content_input.Img,
		}

		return database.ConnDB().Updates(&content).Where("id=?",content_input.Id).Where("lang_iso=?",content_input.Lang_iso).Update("edited_at", time.Now()).Error
	
	})
	err := future.Await()
	if err != nil {
		c.JSON(http.StatusNotAcceptable, gin.H{"error": err})
	} else {
		// c.JSON(http.StatusOK, gin.H{"data": user, "status": http.StatusOK})
		c.JSON(http.StatusOK, gin.H{"status": http.StatusOK})
	}
}
