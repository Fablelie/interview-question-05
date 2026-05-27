package main

import (
	"fmt"
	"log"

	"github.com/gofiber/fiber/v3"

	"example.com/config"
	"example.com/internal/application/usecase"
	"example.com/internal/infrastructure/database"
	"example.com/internal/infrastructure/http"
)

func main() {
	// Load configuration
	cfg := config.LoadConfig()

	// Initialize database
	db, err := database.InitDB(cfg.GetDSN())
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer database.CloseDB(db)

	// Run migrations
	if err := database.RunMigration(db); err != nil {
		log.Fatalf("Failed to run migrations: %v", err)
	}

	log.Println("✓ Database initialized successfully")

	// Initialize repository
	queueRepo := database.NewQueueRepository(db)

	// Initialize use case
	queueUseCase := usecase.NewQueueUseCase(queueRepo)

	// Create Fiber app
	app := fiber.New(fiber.Config{
		AppName: "Queue Ticket System v1.0",
	})

	// Setup CORS middleware
	app.Use(func(c fiber.Ctx) error {
		c.Set("Access-Control-Allow-Origin", "*")
		c.Set("Access-Control-Allow-Headers", "Origin, Content-Type, Accept")
		c.Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")

		if c.Method() == "OPTIONS" {
			return c.SendStatus(200)
		}
		return c.Next()
	})

	// Setup routes
	http.SetupRoutes(app, queueUseCase)

	// Start server
	addr := fmt.Sprintf("%s:%d", cfg.ServerHost, cfg.ServerPort)
	log.Printf("🚀 Starting server on %s", addr)
	if err := app.Listen(addr); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
