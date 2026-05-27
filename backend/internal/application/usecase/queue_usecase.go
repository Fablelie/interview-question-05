package usecase

import (
	"context"
	"fmt"
	"sync"
	"time"

	"example.com/internal/application/dto"
	"example.com/internal/domain"
)

// QueueUseCase handles all queue ticket operations
type QueueUseCase struct {
	queueRepo domain.QueueTicketRepository
	mu        sync.Mutex // Distributed lock at application level
}

// NewQueueUseCase creates new queue use case
func NewQueueUseCase(queueRepo domain.QueueTicketRepository) *QueueUseCase {
	return &QueueUseCase{
		queueRepo: queueRepo,
	}
}

// GetNextTicket retrieves the next ticket number
func (uc *QueueUseCase) GetNextTicket(ctx context.Context) (*dto.GetNextTicketResponse, error) {
	// Double-lock mechanism: application level + database level
	uc.mu.Lock()
	defer uc.mu.Unlock()

	// Get next ticket number from repository
	ticketNumber, err := uc.queueRepo.GetNextTicketNumber(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get next ticket number: %w", err)
	}

	// Create ticket domain
	now := time.Now()
	ticket := &domain.QueueTicket{
		Number:    ticketNumber,
		CreatedAt: now,
		IssuedAt:  &now,
		Status:    "ISSUED",
	}

	// Save ticket to database
	if err := uc.queueRepo.SaveTicket(ctx, ticket); err != nil {
		return nil, fmt.Errorf("failed to save ticket: %w", err)
	}

	return &dto.GetNextTicketResponse{
		TicketNumber: ticketNumber,
		IssuedAt:     now,
		Status:       "ISSUED",
	}, nil
}

// GetCurrentQueue retrieves current queue status
func (uc *QueueUseCase) GetCurrentQueue(ctx context.Context) (*dto.GetCurrentQueueResponse, error) {
	queueCounter, err := uc.queueRepo.GetCurrentQueue(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get current queue: %w", err)
	}

	return &dto.GetCurrentQueueResponse{
		CurrentNumber: queueCounter.CurrentNumber,
		LastUpdated:   queueCounter.LastUpdatedAt,
	}, nil
}

// ClearQueue resets the queue counter
func (uc *QueueUseCase) ClearQueue(ctx context.Context) (*dto.ClearQueueResponse, error) {
	// Use lock to prevent concurrent clear operations
	uc.mu.Lock()
	defer uc.mu.Unlock()

	if err := uc.queueRepo.ClearQueue(ctx); err != nil {
		return nil, fmt.Errorf("failed to clear queue: %w", err)
	}

	return &dto.ClearQueueResponse{
		Message:       "Queue cleared successfully",
		CurrentNumber: "A0",
		ClearedAt:     time.Now(),
	}, nil
}

// GetLatestTicket retrieves the last issued ticket
func (uc *QueueUseCase) GetLatestTicket(ctx context.Context) (*domain.QueueTicket, error) {
	return uc.queueRepo.GetLatestTicket(ctx)
}
