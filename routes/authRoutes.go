package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/visitha2001/go-jwt-auth/handlers"
	"github.com/visitha2001/go-jwt-auth/middleware"
)

func SetupAuthRoutes(app *fiber.App) {
	auth := app.Group("/api/auth")

	// Public routes
	auth.Post("/signup", handlers.SignUp)
	auth.Post("/signin", handlers.SignIn)
	auth.Post("/signout", handlers.SignOut)

	// Protected route
	// This route will first run the AuthRequired middleware
	auth.Get("/profile", middleware.AuthRequired, handlers.GetProfile)
}
