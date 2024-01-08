package main

import (
	"livecode-3-arvisy/config"
	"livecode-3-arvisy/handler"
	"livecode-3-arvisy/middleware"

	"github.com/labstack/echo"
)

func main() {
	e := echo.New()

	db := config.InitDB()

	userHandler := handler.NewUserHandler(db)
	user := e.Group("/users")
	{
		user.POST("/register", userHandler.Register)
		user.POST("/login", userHandler.Login)
		user.GET("/carts", userHandler.GetCart, middleware.Authentication)
		user.POST("/carts", userHandler.CreateCart, middleware.Authentication)
		user.DELETE("/carts/:id", userHandler.DeleteCart, middleware.Authentication)
		user.GET("/orders", userHandler.GetOrder, middleware.Authentication)
		user.POST("/orders", userHandler.CreateOrder, middleware.Authentication)
	}

	productHandler := handler.NewProductHandler(db)
	product := e.Group("/products")
	{
		product.GET("", productHandler.GetAllProduct)
		product.GET("/:id", productHandler.GetProductByID)
	}

	e.Start(":8080")
}
