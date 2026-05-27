package dto

import "time"

// GetNextTicketRequest represents the request for getting next ticket
type GetNextTicketRequest struct{}

// GetNextTicketResponse represents the response for next ticket
type GetNextTicketResponse struct {
	TicketNumber string    `json:"ticket_number"`
	IssuedAt     time.Time `json:"issued_at"`
	Status       string    `json:"status"`
}

// GetCurrentQueueResponse represents current queue status
type GetCurrentQueueResponse struct {
	CurrentNumber string    `json:"current_number"`
	LastUpdated   time.Time `json:"last_updated"`
}

// ClearQueueRequest represents clear queue request
type ClearQueueRequest struct{}

// ClearQueueResponse represents clear queue response
type ClearQueueResponse struct {
	Message       string    `json:"message"`
	CurrentNumber string    `json:"current_number"`
	ClearedAt     time.Time `json:"cleared_at"`
}

// ErrorResponse represents error response
type ErrorResponse struct {
	Error   string `json:"error"`
	Message string `json:"message"`
}
