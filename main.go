package main

import (
	"github.com/ch-m-n/lingo/controllers"
	"github.com/gin-gonic/gin"
)

func main(){
	r := gin.Default()
	r.POST("/crtU", controllers.CreateUser)
	r.GET("/getU/:username", controllers.GetUser)
	r.Run(":8888")
}