package controllers

import (
	"lingo/database"
	"lingo/models"
	"net/http"
	"strings"

	"regexp"

	"github.com/gin-gonic/gin"
	"github.com/lib/pq"
)

func stringProcessor(s string) []string {
	var words_list []string
	// var nonAlphanumericRegex = regexp.MustCompile(`[^a-zA-z0-9 ]+`)
	var nonAlphanumericRegex = regexp.MustCompile(`([^{?=\S*'-}{\p{L}'-}{0-9}]+)`)
	var punctuationRegex = regexp.MustCompile(`[^'\P{P}']+`)
	words := nonAlphanumericRegex.ReplaceAllString(s, " ")
	words = strings.TrimSpace(words)
	words = punctuationRegex.ReplaceAllString(words, " ")
	words = strings.ToLower(words)
	list := strings.Split(words, " ")
	for i := 0; i < len(list); i++ {
		if len(list[i]) > 0 {
			words_list = append(words_list, list[i])
		}
	}
	return words_list
}

func removeDuplicate[T string](sliceList []T) []T {
	allKeys := make(map[T]bool)
	list := []T{}
	for _, item := range sliceList {
		if _, value := allKeys[item]; !value {
			allKeys[item] = true
			list = append(list, item)
		}
	}
	return list
}

func GetContents(c *gin.Context) {
	content_info := new(models.RequestContent)
	e := c.BindJSON(&content_info)
	if e != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": e.Error()})
		return
	}

	head := models.Head{}
	content := []models.Content{}
	var list []string
	literacy := []models.Literacy{}
	notes := []models.Note{}

	database.ConnDB().Get(&head, `SELECT * FROM head WHERE id=$1 AND lang_iso=$2`, content_info.Head_id, content_info.Lang_iso)
	Add2Inventory(c, content_info.My_id, content_info.Head_id, content_info.Lang_iso)
	database.ConnDB().Select(&content, "SELECT * FROM contents WHERE head_id=$1", content_info.Head_id)
	list = append(list, stringProcessor(head.Title)...)
	for i := 0; i < len(content); i++ {
		words := content[i].Body
		list = append(list, stringProcessor(words)...)
	}
	list = removeDuplicate[string](list)
	AddWord(c, list, head.Lang_iso)
	database.ConnDB().MustExec(`INSERT INTO literacy(user_id, word, lang_iso, known_level)
									SELECT $1, word, $3, 0
									FROM UNNEST(CAST($2 as text[])) T (word)
									WHERE NOT EXISTS (SELECT * FROM literacy WHERE user_id=$1 AND word = T.word AND lang_iso=CAST($3 AS VARCHAR))`, content_info.My_id, pq.Array(list), content_info.Lang_iso)
	database.ConnDB().MustExec(`INSERT INTO note(user_id, word, note, lang_iso)
									SELECT $1, word, '', $3
									FROM UNNEST(CAST($2 as text[])) T (word)
									WHERE NOT EXISTS (SELECT * FROM note WHERE user_id=$1 AND word = T.word AND lang_iso=CAST($3 AS VARCHAR))`, content_info.My_id, pq.Array(list), content_info.Lang_iso)
	database.ConnDB().Select(&notes, `SELECT * FROM note WHERE word=ANY($1) AND user_id=$2 AND lang_iso=CAST($3 AS VARCHAR)`, pq.Array(list), content_info.My_id, content_info.Lang_iso)
	database.ConnDB().Select(&literacy, `SELECT * FROM literacy WHERE word=ANY($1) AND user_id=$2 AND lang_iso=CAST($3 AS VARCHAR)`, pq.Array(list), content_info.My_id, content_info.Lang_iso)
	c.JSON(http.StatusOK, gin.H{"head": head, "content": content, "literacy": literacy, "notes": notes})
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

	database.ConnDB().Select(&inventory, `SELECT head_id FROM inventory WHERE user_id=$1 AND lang_iso= $2`,
		content_info.User_id, content_info.Lang_iso)
	database.ConnDB().Select(&head, `SELECT * FROM head WHERE user_id = $1 AND id=ANY($2)`,
		content_info.User_id, pq.Array(inventory))

	c.JSON(http.StatusOK, gin.H{"inventory": inventory, "head": head})

}

func GetAllContents(c *gin.Context) {
	user := new(models.RequestAllContent)
	e := c.BindJSON(&user)
	if e != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": e.Error()})
		return
	}

	content_titles := []models.Head{}
	var inventory []string
	inventory_head := []models.Head{}
	database.ConnDB().Select(&inventory, `SELECT head_id FROM inventory WHERE user_id=$1 AND lang_iso=$2`, user.My_id, user.Lang_iso)
	database.ConnDB().Select(&inventory_head, `SELECT * FROM head WHERE id=ANY($1) AND lang_iso=$2`, pq.Array(inventory), user.Lang_iso)
	database.ConnDB().Select(&content_titles, `SELECT * FROM head WHERE lang_iso=$1`, user.Lang_iso)

	c.JSON(http.StatusOK, gin.H{"data": content_titles, "head": inventory_head})
}

func AddContents(c *gin.Context) {
	content_input := new(models.CreateContent)
	e := c.BindJSON(&content_input)
	if e != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": e.Error()})
		return
	}

	var head_id string
	tx1 := database.ConnDB().MustBegin()
	err1 := tx1.QueryRow(`INSERT INTO head(id, user_id, title, lang_iso, img)
										VALUES(gen_random_uuid(),$1,$2,$3,$4)RETURNING id`,
		content_input.User_id, content_input.Title, content_input.Lang_iso, content_input.Img).Scan(&head_id)
	if err1 != nil {
		tx1.Rollback()
		c.JSON(http.StatusNotAcceptable, gin.H{"error1": err1})
	} else {
		err := tx1.Commit()
		if err != nil {
			c.JSON(http.StatusNotAcceptable, gin.H{"error2": err})
		}
	}

	tx2 := database.ConnDB().MustBegin()
	_, err2 := tx2.Exec(`INSERT INTO contents(id, user_id, head_id, lang_iso, body, created_at, edited_at)
					VALUES(gen_random_uuid(),$1,$2,$3,$4,CURRENT_TIMESTAMP,CURRENT_TIMESTAMP)`,
		content_input.User_id, head_id, content_input.Lang_iso, content_input.Body)
	if err2 != nil {
		tx2.Rollback()
		c.JSON(http.StatusNotAcceptable, gin.H{"error3": err2.Error()})
	} else {
		err := tx2.Commit()
		if err != nil {
			c.JSON(http.StatusNotAcceptable, gin.H{"error4": err})
		}
	}

	c.JSON(http.StatusOK, gin.H{"status": http.StatusOK})
}

func EditContent(c *gin.Context) {
	content_input := new(models.CreateContent)
	e := c.BindJSON(&content_input)
	if e != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": e.Error()})
		return
	}
	book := models.Head{}

	database.ConnDB().Get(&book, `SELECT * FROM head WHERE user_id=$1 AND title=$2`, content_input.Title, content_input.Title)
	tx.QueryRow(`UPDATE head SET title=$1 AND img=$2 WHERE id=$3`, content_input.Title, content_input.Img, book.Id)
	_, err := database.ConnDB().Exec(`
					UPDATE contents
					SET body=$1, edited_at=CURRENT_TIMESTAMP
					WHERE user_id=$2 AND head_id=$3`,
		content_input.Body, content_input.User_id, book.Id)

	if err != nil {
		c.JSON(http.StatusNotAcceptable, gin.H{"error": err})
	} else {
		c.JSON(http.StatusOK, gin.H{"status": http.StatusOK})
	}
}
