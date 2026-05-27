package domain

import (
	"context"
	"time"
)

// QueueTicket represents a ticket in the queue system
type QueueTicket struct {
	ID        int64
	Number    string // Format: A0 to Z9
	CreatedAt time.Time
	IssuedAt  *time.Time
	Status    string // PENDING, ISSUED, CLEARED
}

// QueueCounter tracks the current queue number
type QueueCounter struct {
	ID            int
	CurrentNumber string
	LastUpdatedAt time.Time
	VersionLock   int // For optimistic locking
}

// QueueTicketRepository defines methods for queue ticket data access
type QueueTicketRepository interface {
	// GetNextTicketNumber retrieves and increments the next ticket number
	// Returns the new ticket number or error
	GetNextTicketNumber(ctx context.Context) (string, error)

	// SaveTicket saves a new ticket to database
	SaveTicket(ctx context.Context, ticket *QueueTicket) error

	// GetCurrentQueue returns the current queue counter
	GetCurrentQueue(ctx context.Context) (*QueueCounter, error)

	// ClearQueue resets the queue to 00
	ClearQueue(ctx context.Context) error

	// GetTicketByNumber retrieves ticket info by number
	GetTicketByNumber(ctx context.Context, number string) (*QueueTicket, error)

	// GetLatestTicket gets the most recently issued ticket
	GetLatestTicket(ctx context.Context) (*QueueTicket, error)
}
