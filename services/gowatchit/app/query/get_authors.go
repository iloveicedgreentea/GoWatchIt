package query

import (
	"context"
	"fmt"

	"github.com/iloveicedgreentea/gowatchit/services/gowatchit/domain/beq"
)

// GetAuthorsQuery queries for BEQ authors
type GetAuthorsQuery struct{}

// GetAuthorsHandler handles the get authors query
type GetAuthorsHandler struct {
	repo *beq.AuthorRepository
}

// NewGetAuthorsHandler creates a new GetAuthorsHandler
func NewGetAuthorsHandler(repo *beq.AuthorRepository) (*GetAuthorsHandler, error) {
	if repo == nil {
		return nil, fmt.Errorf("repository is nil")
	}

	return &GetAuthorsHandler{
		repo: repo,
	}, nil
}

// Handle executes the get authors query
func (h *GetAuthorsHandler) Handle(ctx context.Context, query *GetAuthorsQuery) ([]string, error) {
	return h.repo.GetAllAuthors(ctx)
}
