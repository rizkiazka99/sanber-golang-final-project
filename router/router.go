package router

import (
	"golang-final-project/controllers"

	"github.com/gin-gonic/gin"
)

func StartServer() *gin.Engine {
	router := gin.Default()

	router.POST("/api/register", controllers.Register)
	router.POST("/api/login", controllers.Login)

	return router
}
