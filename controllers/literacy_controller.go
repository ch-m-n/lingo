package controllers

import (
	"lingo/async"
	"lingo/database"
	"lingo/models"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func GetWordLevel(c *gin.Context) {
	words := new(models.InputParagraph)
	e := c.BindJSON(&words)
	if e != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": e.Error()})
		return
	}
	var word_level models.OutputParagraph
	var word models.Literacy
	for i := 0; i < len(words.Words); i++ {
		future := async.Exec(func() interface{} {
			return database.ConnDB().Table("literacy").Exec("SELECT * FROM literacy"+
															"WHERE word='" + words.Words[i].Word + "'"+
															"AND lang_iso='" + words.Words[i].Lang_iso + "'"+
															"AND user_id="+words.User_id).Scan(&word)
		})
		err := future.Await()
		if err != nil {
			c.JSON(http.StatusNotAcceptable, gin.H{"error": err})
			break
		} else {
			word_level.Level = append(word_level.Level, word)
		}
	}
	c.JSON(http.StatusOK, gin.H{"data": word_level, "status": http.StatusOK})
}

func AddWordLevel(c *gin.Context) {
	words := new(models.InputParagraph)
	e := c.BindJSON(&words)
	if e != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": e.Error()})
		return
	}
	
	for i := 0; i < len(words.Words); i++ {
		future := async.Exec(func() interface{} {
			known_level := 0
			if (words.Words[i].Known_level!=0){
				known_level = words.Words[i].Known_level
			}
			AddWord(c, words.Words[i].Word, words.Words[i].Lang_iso)
			return database.ConnDB().Table("literacy").Exec("UPDATE literacy SET known_level="+strconv.Itoa(known_level)+" WHERE user_id='"+words.User_id+"' AND word='"+words.Words[i].Word+"';"+
															"INSERT INTO literacy(user_id, word, lang_iso, known_level)"+
															"SELECT '"+words.User_id+"','"+words.Words[i].Word+"','"+words.Words[i].Lang_iso+"',"+strconv.Itoa(known_level)+
															" WHERE NOT EXISTS(SELECT 1 FROM literacy WHERE user_id='"+words.User_id+"'AND word='"+words.Words[i].Word+"');"+
															"INSERT INTO note(user_id, word, note)"+
															"VALUES('"+words.User_id+"', '"+words.Words[i].Word+"', '')"+
															"ON CONFLICT DO NOTHING").Error
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
