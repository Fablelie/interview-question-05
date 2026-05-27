package handler

import (
	"context"
	"errors"
	"net/http"

	"github.com/gofiber/fiber/v3"

	"example.com/internal/application/dto"
	"example.com/internal/application/usecase"
)

// QueueHandler handles queue-related HTTP requests
type QueueHandler struct {
	queueUseCase *usecase.QueueUseCase
}

// NewQueueHandler creates new queue handler
func NewQueueHandler(queueUseCase *usecase.QueueUseCase) *QueueHandler {
	return &QueueHandler{
		queueUseCase: queueUseCase,
	}
}

// GetNextTicket handles request to get next ticket
// POST /api/tickets/next
func (h *QueueHandler) GetNextTicket(c fiber.Ctx) error {
	ctx := context.Background()

	response, err := h.queueUseCase.GetNextTicket(ctx)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(dto.ErrorResponse{
			Error:   "GET_TICKET_ERROR",
			Message: err.Error(),
		})
	}

	return c.Status(http.StatusOK).JSON(response)
}

// GetCurrentQueue handles request to get current queue status
// GET /api/queue/current
func (h *QueueHandler) GetCurrentQueue(c fiber.Ctx) error {
	ctx := context.Background()

	response, err := h.queueUseCase.GetCurrentQueue(ctx)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(dto.ErrorResponse{
			Error:   "GET_QUEUE_ERROR",
			Message: err.Error(),
		})
	}

	return c.Status(http.StatusOK).JSON(response)
}

// ClearQueue handles request to clear the queue
// POST /api/queue/clear
func (h *QueueHandler) ClearQueue(c fiber.Ctx) error {
	ctx := context.Background()

	response, err := h.queueUseCase.ClearQueue(ctx)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(dto.ErrorResponse{
			Error:   "CLEAR_QUEUE_ERROR",
			Message: err.Error(),
		})
	}

	return c.Status(http.StatusOK).JSON(response)
}

// GetLatestTicket handles request to get latest ticket
// GET /api/tickets/latest
func (h *QueueHandler) GetLatestTicket(c fiber.Ctx) error {
	ctx := context.Background()

	ticket, err := h.queueUseCase.GetLatestTicket(ctx)
	if err != nil {
		if errors.Is(err, errors.New("no tickets found")) {
			return c.Status(http.StatusNotFound).JSON(dto.ErrorResponse{
				Error:   "NOT_FOUND",
				Message: "No tickets available",
			})
		}

		return c.Status(http.StatusInternalServerError).JSON(dto.ErrorResponse{
			Error:   "GET_TICKET_ERROR",
			Message: err.Error(),
		})
	}

	return c.Status(http.StatusOK).JSON(fiber.Map{
		"ticket_number": ticket.Number,
		"issued_at":     ticket.IssuedAt,
		"status":        ticket.Status,
	})
}
