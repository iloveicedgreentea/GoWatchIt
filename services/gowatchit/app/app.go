package app

import (
	"context"
	"database/sql"
	"errors"

	"github.com/iloveicedgreentea/gowatchit/pkg/beq"
	"github.com/iloveicedgreentea/gowatchit/services/gowatchit/app/command"
	"github.com/iloveicedgreentea/gowatchit/services/gowatchit/app/query"
	domainbeq "github.com/iloveicedgreentea/gowatchit/services/gowatchit/domain/beq"
)

type App struct {
	Commands Commands
	Queries  Queries
}

type Commands struct {
	ProcessWebhook *command.ProcessWebhookHandler
	SaveConfig     *command.SaveConfigHandler
	ScrapeAuthors  *command.ScrapeAuthorsHandler
}

type Queries struct {
	GetBeq     *query.GetBeqHandler
	GetConfig  *query.GetConfigHandler
	GetLogs    *query.GetLogsHandler
	GetAuthors *query.GetAuthorsHandler
}

func NewApplication(ctx context.Context, beqClient *beq.BeqClient, db *sql.DB) (*App, error) {
	if beqClient == nil {
		return nil, errors.New("beqClient is nil")
	}
	if db == nil {
		return nil, errors.New("db is nil")
	}

	processWebhookHandler, err := command.NewProcessWebhookHandler(beqClient)
	if err != nil {
		return nil, err
	}

	saveConfigHandler, err := command.NewSaveConfigHandler()
	if err != nil {
		return nil, err
	}

	// Initialize BEQ author scraping infrastructure
	scraper := domainbeq.NewAuthorScraper()
	repo := domainbeq.NewAuthorRepository(db)

	scrapeAuthorsHandler, err := command.NewScrapeAuthorsHandler(scraper, repo)
	if err != nil {
		return nil, err
	}

	beqGetHandler, err := query.NewGetBeqHandler(beqClient)
	if err != nil {
		return nil, err
	}

	getConfigHandler, err := query.NewGetConfigHandler()
	if err != nil {
		return nil, err
	}

	getLogsHandler, err := query.NewGetLogsHandler()
	if err != nil {
		return nil, err
	}

	getAuthorsHandler, err := query.NewGetAuthorsHandler(repo)
	if err != nil {
		return nil, err
	}

	return &App{
		Commands: Commands{
			ProcessWebhook: processWebhookHandler,
			SaveConfig:     saveConfigHandler,
			ScrapeAuthors:  scrapeAuthorsHandler,
		},
		Queries: Queries{
			GetBeq:     beqGetHandler,
			GetConfig:  getConfigHandler,
			GetLogs:    getLogsHandler,
			GetAuthors: getAuthorsHandler,
		},
	}, nil
}
