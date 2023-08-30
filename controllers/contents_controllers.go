package controllers

import (
	"lingo/async"
	"lingo/database"
	"lingo/models"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/lib/pq"
)

func GetContents(c *gin.Context) {
	content_info := new(models.RequestContent)
	e := c.BindJSON(&content_info)
	if e != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": e.Error()})
		return
	}
	head := models.Head{}
	content := models.Content{}
	future := async.Exec(func() interface{} {
		database.ConnDB().Get(&head, `SELECT * FROM head WHERE user_id=$1 AND title = $2 and lang_iso=$3)`,
									content_info.Author_id, content_info.Title, content_info.Lang_iso)
		Add2Inventory(c, content_info.My_id, content_info.Head_id, content_info.Lang_iso)
		return database.ConnDB().Get(&content, "SELECT * FROM contents WHERE title=$1 AND head_id=$2", content_info.Title, content_info.Head_id)
	})
	future.Await()
	c.JSON(http.StatusOK, gin.H{"head": head, "content": content})
}

func GetInventoryContents(c *gin.Context) {
	content_info := new(models.RequestInventoryContent)
	e := c.BindJSON(&content_info)
	if e != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": e.Error()})
		return
	}
	var inventory []string
	head := []models.Head{}
	future := async.Exec(func() interface{} {

		database.ConnDB().Select(&inventory,`SELECT head_id FROM inventory WHERE user_id=$1 AND lang_iso= $2`, 
								content_info.User_id, content_info.Lang_iso)			

		return database.ConnDB().Select(&head, `SELECT * FROM head WHERE user_id = $1 AND id=ANY($2)`,
			content_info.User_id, pq.Array(inventory))
	})
	err := future.Await()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"errors": err})
	} else {
		c.JSON(http.StatusOK, gin.H{"inventory": inventory, "head": head})
	}

}

func GetAllContents(c *gin.Context) {
	user := new(models.RequestAllContent)
	e := c.BindJSON(&user)
	if e != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": e.Error()})
		return
	}
	
	content_titles := []models.Head{}
	future := async.Exec(func() interface{} {
		return database.ConnDB().Select(&content_titles,
										//`SELECT DISTINCT ON(title) id, user_id, title , lang_iso, created_at , edited_at , img  
										`SELECT *
										FROM head 
										WHERE lang_iso=$1`,
			user.Lang_iso)
	})
	future.Await()
	c.JSON(http.StatusOK, gin.H{"data": content_titles})
}

func AddContents(c *gin.Context) {
	content_input := new(models.CreateContent)
	e := c.BindJSON(&content_input)
	if e != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": e.Error()})
		return
	}

	var head_id string
	future := async.Exec(func() interface{} {

		tx.QueryRow(`INSERT INTO head(id, user_id, title, img)
										VALUES(gen_random_uuid(),$1,$2,$3)RETURNING id`, content_input.User_id, content_input.Title, content_input.Img).Scan(head_id)
		
		tx.QueryRow(`INSERT INTO contents(id, user_id, head_id, lang_iso, body, created_at, edited_at)
					VALUES(gen_random_uuid(),$1,$2,$3,$4,CURRENT_TIMESTAMP,CURRENT_TIMESTAMP)`,
			content_input.User_id, head_id, content_input.Lang_iso, content_input.Body)
		return tx.Commit()
	})
	err := future.Await()
	if err != nil {
		c.JSON(http.StatusNotAcceptable, gin.H{"error": err})
	}else{
		c.JSON(http.StatusOK, gin.H{"status": http.StatusOK})
	}
}

func EditContent(c *gin.Context) {
	content_input := new(models.CreateContent)
	e := c.BindJSON(&content_input)
	if e != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": e.Error()})
		return
	}
	book := models.Head{}
	future := async.Exec(func() interface{} {

		database.ConnDB().Get(&book,`SELECT * FROM head WHERE user_id=$1 AND title=$2`, content_input.Title, content_input.Title)

		tx.QueryRow(`UPDATE head SET title=$1 AND img=$2 WHERE id=$3`, content_input.Title, content_input.Img, book.Id)

		_, err := database.ConnDB().Exec(`
					UPDATE contents
					SET body=$1, edited_at=CURRENT_TIMESTAMP
					WHERE user_id=$2 AND head_id=$3`,
				content_input.Body, content_input.User_id, book.Id)
		return err
	})
	err := future.Await()
	if err != nil {
		c.JSON(http.StatusNotAcceptable, gin.H{"error": err})
	}else{
		c.JSON(http.StatusOK, gin.H{"status": http.StatusOK})
	}
}
