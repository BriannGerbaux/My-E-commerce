package main

import (
	"go-backend/src/handlers"
	"go-backend/src/middleware"

	"github.com/gin-gonic/gin"
)

//func GetCart(c *gin.Context) {
//	query := "SELECT Products.name, Products.description, Products.price_in_dollar, Products.thumbnail_url FROM Products JOIN Cart ON "
//}

func main() {
	router := gin.Default()
	dbService := middleware.DbConnection()
	router.Use(middleware.InjectDatabaseService(&dbService))

	apiRouter := router.Group("/api")
	apiRouter.Use(middleware.UserAuthMiddleware())

	apiRouter.GET("/users/:id", handlers.GetUserById)
	apiRouter.GET("/products", handlers.GetProducts)

	adminRouter := router.Group("/admin")
	adminRouter.Use(middleware.AdminAuthMiddleware())
	adminRouter.POST("/product", handlers.PostProduct)
	router.Run(":8181")
}
