package router

import (
	"golang-final-project/controllers"

	"github.com/gin-gonic/gin"
)

func StartServer() *gin.Engine {
	router := gin.Default()

	router.Static("/uploads", "./uploads")

	router.POST("/api/register", controllers.Register)
	router.POST("/api/login", controllers.Login)

	router.POST("/api/items", controllers.PostItem)
	router.GET("/api/items", controllers.GetItems)
	router.GET("/api/items/:id", controllers.GetItemById)
	router.PUT("/api/items/:id", controllers.UpdateItem)
	router.DELETE("/api/items/:id", controllers.DeleteItem)

	router.POST("/api/carts", controllers.PostCart)
	router.GET("/api/carts", controllers.GetCarts)
	router.GET("/api/carts/:id", controllers.GetCartById)
	router.GET("/api/carts/:id/users", controllers.GetCartsByUserId)
	// router.PUT("/api/carts/:id", controllers.UpdateCart)
	// router.DELETE("/api/carts/:id/cart_items", controllers.DeleteCartItems)
	router.DELETE("/api/carts/:id", controllers.DeleteCart)

	router.PUT("/api/pay/:cart_id", controllers.PayCart)

	return router
}
