package database

import (
	"database/sql"
	"fmt"

	_ "github.com/microsoft/go-mssqldb"
)

// InitDB initializes database connection and returns *sql.DB
func InitDB(dsn string) (*sql.DB, error) {
	db, err := sql.Open("sqlserver", dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	// Test connection
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	// Set connection pool settings
	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(5)

	return db, nil
}

// RunMigration creates necessary tables
func RunMigration(db *sql.DB) error {
	// Create QueueCounter table
	if _, err := db.Exec(`
		IF NOT EXISTS (SELECT * FROM sysobjects WHERE name='QueueCounter' and xtype='U')
		BEGIN
			CREATE TABLE QueueCounter (
				ID INT PRIMARY KEY DEFAULT 1,
				CurrentNumber NVARCHAR(3) NOT NULL DEFAULT 'A0',
				LastUpdatedAt DATETIME2 NOT NULL DEFAULT GETUTCDATE(),
				VersionLock INT NOT NULL DEFAULT 0
			)
			
			-- Insert default row
			INSERT INTO QueueCounter (ID, CurrentNumber, LastUpdatedAt, VersionLock)
			VALUES (1, 'A0', GETUTCDATE(), 0)
		END
	`); err != nil {
		return fmt.Errorf("failed to create QueueCounter table: %w", err)
	}

	// Create QueueTickets table
	if _, err := db.Exec(`
		IF NOT EXISTS (SELECT * FROM sysobjects WHERE name='QueueTickets' and xtype='U')
		BEGIN
			CREATE TABLE QueueTickets (
				ID BIGINT PRIMARY KEY IDENTITY(1,1),
				TicketNumber NVARCHAR(3) UNIQUE NOT NULL,
				CreatedAt DATETIME2 NOT NULL DEFAULT GETUTCDATE(),
				IssuedAt DATETIME2 NULL,
				Status NVARCHAR(10) NOT NULL DEFAULT 'PENDING',
				INDEX idx_ticket_number (TicketNumber),
				INDEX idx_issued_at (IssuedAt)
			)
		END
	`); err != nil {
		return fmt.Errorf("failed to create QueueTickets table: %w", err)
	}

	return nil
}

// CloseDB closes database connection
func CloseDB(db *sql.DB) error {
	return db.Close()
}
