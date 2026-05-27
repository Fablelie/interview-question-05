package database

import (
	"context"
	"database/sql"
	"fmt"

	"example.com/internal/domain"
)

// QueueRepository implements domain.QueueTicketRepository
type QueueRepository struct {
	db        *sql.DB
	generator *TicketNumberGenerator
}

// NewQueueRepository creates new queue repository
func NewQueueRepository(db *sql.DB) domain.QueueTicketRepository {
	return &QueueRepository{
		db:        db,
		generator: &TicketNumberGenerator{},
	}
}

// GetNextTicketNumber retrieves and increments the next ticket number with optimistic locking
func (r *QueueRepository) GetNextTicketNumber(ctx context.Context) (string, error) {
	tx, err := r.db.BeginTx(ctx, &sql.TxOptions{
		Isolation: sql.LevelSerializable, // Highest isolation level to prevent conflicts
	})
	if err != nil {
		return "", fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	var currentNumber string
	var versionLock int

	// Get current queue with lock
	err = tx.QueryRowContext(
		ctx,
		"SELECT CurrentNumber, VersionLock FROM QueueCounter WHERE ID = 1",
	).Scan(&currentNumber, &versionLock)

	if err != nil {
		return "", fmt.Errorf("failed to get current queue: %w", err)
	}

	// Generate next number
	nextNumber, err := r.generator.GetNextNumber(currentNumber)
	if err != nil {
		return "", fmt.Errorf("failed to generate next number: %w", err)
	}

	// Update counter with optimistic lock (check version before update)
	result, err := tx.ExecContext(
		ctx,
		`UPDATE QueueCounter 
		 SET CurrentNumber = @NextNumber, 
		     VersionLock = VersionLock + 1,
		     LastUpdatedAt = GETUTCDATE()
		 WHERE ID = 1 AND VersionLock = @OldVersion`,
		sql.Named("NextNumber", nextNumber),
		sql.Named("OldVersion", versionLock),
	)

	if err != nil {
		return "", fmt.Errorf("failed to update queue counter: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return "", fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		// Version mismatch - try again (concurrent request detected)
		return "", fmt.Errorf("concurrent request detected - version conflict")
	}

	if err := tx.Commit(); err != nil {
		return "", fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nextNumber, nil
}

// SaveTicket saves a new ticket to database
func (r *QueueRepository) SaveTicket(ctx context.Context, ticket *domain.QueueTicket) error {
	query := `
		MERGE INTO QueueTickets AS Target
		USING (SELECT @Number AS TicketNumber) AS Source
		ON (Target.TicketNumber = Source.TicketNumber)
		WHEN MATCHED THEN
			UPDATE SET 
				IssuedAt = @IssuedAt,
				Status = @Status
		WHEN NOT MATCHED THEN
			INSERT (TicketNumber, CreatedAt, IssuedAt, Status)
			VALUES (@Number, @CreatedAt, @IssuedAt, @Status);
	`

	_, err := r.db.ExecContext(
		ctx,
		query,
		sql.Named("Number", ticket.Number),
		sql.Named("CreatedAt", ticket.CreatedAt),
		sql.Named("IssuedAt", ticket.IssuedAt),
		sql.Named("Status", ticket.Status),
	)

	if err != nil {
		return fmt.Errorf("failed to save ticket: %w", err)
	}

	return nil
}

// GetCurrentQueue returns the current queue counter
func (r *QueueRepository) GetCurrentQueue(ctx context.Context) (*domain.QueueCounter, error) {
	var counter domain.QueueCounter

	err := r.db.QueryRowContext(
		ctx,
		"SELECT ID, CurrentNumber, LastUpdatedAt, VersionLock FROM QueueCounter WHERE ID = 1",
	).Scan(&counter.ID, &counter.CurrentNumber, &counter.LastUpdatedAt, &counter.VersionLock)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("queue counter not found")
		}
		return nil, fmt.Errorf("failed to get current queue: %w", err)
	}

	return &counter, nil
}

// ClearQueue resets the queue to A0
func (r *QueueRepository) ClearQueue(ctx context.Context) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to begin clear transction: %w", err)
	}
	defer tx.Rollback()

	_, err = tx.ExecContext(
		ctx,
		`UPDATE QueueCounter 
		 SET CurrentNumber = '00', 
		     VersionLock = VersionLock + 1,
		     LastUpdatedAt = GETUTCDATE()
		 WHERE ID = 1`,
	)
	if err != nil {
		return fmt.Errorf("failed to reset queue counter: %w", err)
	}

	_, err = tx.ExecContext(
		ctx,
		`UPDATE QueueTickets
		 SET Status = 'CLEARED'
		 WHERE Status <> 'CLEARED'`,
	)
	if err != nil {
		return fmt.Errorf("failed to update ticket status to cleared: %w", err)
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit clear transaction: %w", err)
	}

	return nil
}

// GetTicketByNumber retrieves ticket info by number
func (r *QueueRepository) GetTicketByNumber(ctx context.Context, number string) (*domain.QueueTicket, error) {
	var ticket domain.QueueTicket

	err := r.db.QueryRowContext(
		ctx,
		"SELECT ID, TicketNumber, CreatedAt, IssuedAt, Status FROM QueueTickets WHERE TicketNumber = @Number",
		sql.Named("Number", number),
	).Scan(&ticket.ID, &ticket.Number, &ticket.CreatedAt, &ticket.IssuedAt, &ticket.Status)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("ticket not found")
		}
		return nil, fmt.Errorf("failed to get ticket: %w", err)
	}

	return &ticket, nil
}

// GetLatestTicket gets the most recently issued ticket
func (r *QueueRepository) GetLatestTicket(ctx context.Context) (*domain.QueueTicket, error) {
	var ticket domain.QueueTicket

	err := r.db.QueryRowContext(
		ctx,
		`SELECT TOP 1 ID, TicketNumber, CreatedAt, IssuedAt, Status 
		 FROM QueueTickets 
		 WHERE Status = 'ISSUED'
		 ORDER BY IssuedAt DESC`,
	).Scan(&ticket.ID, &ticket.Number, &ticket.CreatedAt, &ticket.IssuedAt, &ticket.Status)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("no tickets found")
		}
		return nil, fmt.Errorf("failed to get latest ticket: %w", err)
	}

	return &ticket, nil
}
