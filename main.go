package main

import (
	"lingo/controllers"
	"lingo/database"
	"lingo/middlewares"
	"time"

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
	database.ConnDB().MustExec(database.Schema())
	router := initRouter()
	router.Run(":5000")
}

func initRouter() *gin.Engine {
	router := gin.Default()
	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"PUT", "PATCH", "GET", "POST"},
		AllowHeaders:     []string{"Authorization", "Referer", "User-Agent"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		// AllowOriginFunc: func(origin string) bool {
		//   return origin == "https://github.com"
		// },
		MaxAge: 12 * time.Hour,
	  }))
	path := router.Group("/")
	{
		path.GET("/",controllers.Home)
	}
	api := router.Group("/api")
	{
		api.POST("/user/register", controllers.CreateUser)
		api.POST("/user/login", controllers.VerifyUser)
		secured := api.Group("/secured").Use(middlewares.Auth())
		{
			secured.GET("/user", controllers.GetUser)
			secured.GET("/user/word/get", controllers.GetWord)
			secured.POST("/user/content/get", controllers.GetContents)
			secured.GET("/user/content/get_inventory", controllers.GetInventoryContents)
			secured.POST("/user/content/get_all", controllers.GetAllContents)
			secured.GET("/user/literacy/get", controllers.GetWordLevel)
			secured.GET("/user/literacy/get_all", controllers.GetAllWordLevel)
			secured.GET("/user/note/get", controllers.GetNote)
			secured.GET("/user/note/get_all", controllers.GetAllNotes)
			secured.GET("/user/inventory/get", controllers.GetInventory)
			// secured.POST("/user/word/add", controllers.AddWord)
			secured.POST("/user/edit", controllers.EditUser)
			secured.POST("/user/content/add", controllers.AddContents)
			secured.POST("/user/content/edit", controllers.EditContent)
			secured.POST("/user/literacy/add", controllers.AddWordLevel)
			secured.POST("/user/note/edit", controllers.EditNote)
		}
	}

	return router
}
