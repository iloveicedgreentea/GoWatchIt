package app

import (
	"context"
	"errors"

	"github.com/iloveicedgreentea/gowatchit/pkg/beq"
	"github.com/iloveicedgreentea/gowatchit/services/gowatchit/app/command"
	"github.com/iloveicedgreentea/gowatchit/services/gowatchit/app/query"
)

type App struct {
	Commands Commands
	Queries  Queries
}

type Commands struct {
	ProcessWebhook *command.ProcessWebhookHandler
}

type Queries struct {
	GetBeq *query.GetBeqHandler
}

func NewApplication(ctx context.Context, beqClient *beq.BeqClient) (*App, error) {
	if beqClient == nil {
		return nil, errors.New("beqClient is nil")
	}

	processWebhookHandler, err := command.NewProcessWebhookHandler(beqClient)
	if err != nil {
		return nil, err
	}

	beqGetHandler, err := query.NewGetBeqHandler(beqClient)
	if err != nil {
		return nil, err
	}

	return &App{
		Commands: Commands{
			ProcessWebhook: processWebhookHandler,
		},
		Queries: Queries{
			GetBeq: beqGetHandler,
		},
	}, nil
}
