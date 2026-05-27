package http

import (
	"github.com/gofiber/fiber/v3"

	"example.com/internal/application/usecase"
	"example.com/internal/infrastructure/http/handler"
)

// Router sets up all API routes
func SetupRoutes(app *fiber.App, queueUseCase *usecase.QueueUseCase) {
	queueHandler := handler.NewQueueHandler(queueUseCase)

	// API routes
	api := app.Group("/api")

	// Queue routes
	tickets := api.Group("/tickets")
	tickets.Post("/next", queueHandler.GetNextTicket)
	tickets.Get("/latest", queueHandler.GetLatestTicket)

	// Queue counter routes
	queue := api.Group("/queue")
	queue.Get("/current", queueHandler.GetCurrentQueue)
	queue.Post("/clear", queueHandler.ClearQueue)

	// Health check
	app.Get("/health", func(c fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"status":  "ok",
			"service": "queue-ticket-system",
		})
	})
}
