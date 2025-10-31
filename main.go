package main

import (
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/visitha2001/go-jwt-auth/database"
	"github.com/visitha2001/go-jwt-auth/routes"
)

func main() {
	// Connect to database (this also runs migrations)
	database.ConnectDB()

	// Create new Fiber app
	app := fiber.New()

	// --- Middleware ---
	app.Use(logger.New()) // Simple request logger

	// CORS for frontend interaction
	app.Use(cors.New(cors.Config{
		AllowOrigins:     "http://localhost:3000", // Your frontend URL
		AllowCredentials: true,
	}))

	// --- Setup Routes ---
	routes.SetupAuthRoutes(app)
	routes.ItemRoutes(app)

	// Handle not found
	app.Use(func(c *fiber.Ctx) error {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"status": "error", "message": "Route not found",
		})
	})
	// root route
	app.Get("/", func(c *fiber.Ctx) error {
		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"status": "success", "message": "Welcome to the go-fiber server",
		})
	})

	// Start server
	port := ":" + os.Getenv("PORT")
	if port == ":" {
		port = ":8082"
	}
	log.Println("Server running on port", port)
	if err := app.Listen(port); err != nil {
		log.Fatal(err)
	}
}
