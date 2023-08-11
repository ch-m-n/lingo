package main

import (
	_ "io/ioutil"
	"lingo/controllers"
	_ "lingo/database"
	"lingo/middlewares"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	// query, err := ioutil.ReadFile("./database/schema.sql")
	// if err != nil {
	// 	panic(err)
	// }
	// if err := database.ConnDB().Exec(string(query)); err != nil {
	// 	panic(err)
	// }

	router := initRouter()
	router.Run(":5000")
}

func initRouter() *gin.Engine {
	router := gin.Default()
	router.Use(cors.Default())
	api := router.Group("/api")
	{
		api.POST("/user/register", controllers.CreateUser)
		api.POST("/user/login", controllers.VerifyUser)
		secured := api.Group("/secured").Use(middlewares.Auth())
		{
			secured.GET("/user", controllers.GetUser)
			secured.GET("/user/word/get", controllers.GetWord)
			secured.GET("/user/content/get", controllers.GetContents)
			secured.GET("/user/literacy/get", controllers.GetWordLevel)
			secured.GET("/user/note/get", controllers.GetNote)
			// secured.POST("/user/word/add", controllers.AddWord)
			secured.POST("/user/content/add", controllers.AddContents)
			secured.POST("/user/content/edit", controllers.EditContent)
			secured.POST("/user/literacy/add", controllers.AddWordLevel)
			secured.POST("/user/note/edit", controllers.EditNote)
		}
	}

	return router
}
