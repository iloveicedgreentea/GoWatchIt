package command

import (
	"context"
	"fmt"

	"github.com/iloveicedgreentea/gowatchit/pkg/logger"
	"github.com/iloveicedgreentea/gowatchit/services/gowatchit/domain/beq"
	"go.uber.org/zap"
)

// ScrapeAuthorsCommand triggers BEQ author scraping
type ScrapeAuthorsCommand struct{}

// ScrapeAuthorsHandler handles the scrape authors command
type ScrapeAuthorsHandler struct {
	scraper *beq.AuthorScraper
	repo    *beq.AuthorRepository
}

// NewScrapeAuthorsHandler creates a new ScrapeAuthorsHandler
func NewScrapeAuthorsHandler(scraper *beq.AuthorScraper, repo *beq.AuthorRepository) (*ScrapeAuthorsHandler, error) {
	if scraper == nil {
		return nil, fmt.Errorf("scraper is nil")
	}
	if repo == nil {
		return nil, fmt.Errorf("repository is nil")
	}

	return &ScrapeAuthorsHandler{
		scraper: scraper,
		repo:    repo,
	}, nil
}

// Handle executes the scrape authors command
func (h *ScrapeAuthorsHandler) Handle(ctx context.Context, cmd *ScrapeAuthorsCommand) error {
	log := logger.GetLoggerFromContext(ctx)
	log.Info("Starting BEQ author scrape")

	// Scrape authors from catalogue
	authors, err := h.scraper.ScrapeAuthors(ctx)
	if err != nil {
		log.Error("Failed to scrape authors", zap.Error(err))
		return fmt.Errorf("failed to scrape authors: %w", err)
	}

	// Upsert authors to database
	if err := h.repo.UpsertAuthors(ctx, authors); err != nil {
		log.Error("Failed to upsert authors", zap.Error(err))
		return fmt.Errorf("failed to upsert authors: %w", err)
	}

	log.Info("Successfully scraped and stored BEQ authors", zap.Int("count", len(authors)))
	return nil
}
