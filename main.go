package main

import (
	"main/controllers"

	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()

	router.POST("/code/java", controllers.JavaExecuter)

	router.Run(":3090")
}
