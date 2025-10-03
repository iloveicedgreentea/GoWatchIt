package beq

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/iloveicedgreentea/gowatchit/pkg/logger"
	"go.uber.org/zap"
)

// AuthorRepository handles database operations for BEQ authors
type AuthorRepository struct {
	db *sql.DB
}

// NewAuthorRepository creates a new AuthorRepository instance
func NewAuthorRepository(db *sql.DB) *AuthorRepository {
	return &AuthorRepository{
		db: db,
	}
}

// UpsertAuthors inserts or updates authors in the database
func (r *AuthorRepository) UpsertAuthors(ctx context.Context, authors []string) error {
	log := logger.GetLoggerFromContext(ctx)
	log.Debug("Upserting authors to database", zap.Int("count", len(authors)))

	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer func() {
		if err := tx.Rollback(); err != nil {
			log.Error("Failed to rollback transaction", zap.Error(err))
		}
	}()

	now := time.Now().UTC().Format(time.RFC3339)

	// Upsert each author using INSERT OR IGNORE followed by UPDATE
	// This ensures we don't create duplicates with NULL ids
	for _, author := range authors {
		// First, try to insert if it doesn't exist
		_, err := tx.ExecContext(ctx, `
			INSERT INTO beqauthor (name, updated_at)
			SELECT ?, ?
			WHERE NOT EXISTS (SELECT 1 FROM beqauthor WHERE name = ?)
		`, author, now, author)
		if err != nil {
			return fmt.Errorf("failed to insert author %s: %w", author, err)
		}

		// Then update the timestamp if it already exists
		_, err = tx.ExecContext(ctx, `
			UPDATE beqauthor SET updated_at = ? WHERE name = ?
		`, now, author)
		if err != nil {
			return fmt.Errorf("failed to update author %s: %w", author, err)
		}
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	log.Info("Successfully upserted authors", zap.Int("count", len(authors)))
	return nil
}

// GetAllAuthors retrieves all author names from the database
func (r *AuthorRepository) GetAllAuthors(ctx context.Context) ([]string, error) {
	log := logger.GetLoggerFromContext(ctx)
	log.Debug("Fetching all authors from database")

	rows, err := r.db.QueryContext(ctx, "SELECT name FROM beqauthor ORDER BY name")
	if err != nil {
		return nil, fmt.Errorf("failed to query authors: %w", err)
	}
	defer func() {
		if err := rows.Close(); err != nil {
			log.Error("Failed to close rows", zap.Error(err))
		}
	}()

	authors := make([]string, 0)
	for rows.Next() {
		var name string
		if err := rows.Scan(&name); err != nil {
			return nil, fmt.Errorf("failed to scan author: %w", err)
		}
		authors = append(authors, name)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating rows: %w", err)
	}

	log.Debug("Fetched authors from database", zap.Int("count", len(authors)))
	return authors, nil
}
