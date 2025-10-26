package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/visitha2001/go-jwt-auth/handlers"
	"github.com/visitha2001/go-jwt-auth/middleware"
)

func ItemRoutes(app fiber.Router) {
	itemGroup := app.Group("/api/items")
	itemGroup.Get("/", handlers.GetItems)
	itemGroup.Get("/:id", handlers.GetItem)
	itemGroup.Post("/", handlers.CreateItem)
	itemGroup.Put("/update/:id", middleware.AuthRequired, handlers.UpdateItem)
	itemGroup.Delete("/delete/:id", middleware.AuthRequired, handlers.DeleteItem)
	itemGroup.Get("/user/:id", middleware.AuthRequired, handlers.GetItemsByUserID)
}
