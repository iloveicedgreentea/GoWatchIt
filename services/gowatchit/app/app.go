package app

import (
	"context"
	"errors"

	"github.com/iloveicedgreentea/gowatchit/services/gowatchit/app/command"
	"github.com/iloveicedgreentea/gowatchit/services/gowatchit/app/query"
	"github.com/iloveicedgreentea/gowatchit/services/gowatchit/domain/beq"
)

type App struct{
	Commands Commands
	Queries  Queries
}

type Commands struct {
	LoadBeq *command.LoadBeqHandler
}

type Queries struct {
	GetBeq *query.GetBeqHandler
}

func NewApplication(ctx context.Context, beqClient *beq.BeqClient) (*App, error) {
	if beqClient == nil {
		return nil, errors.New("beqClient is nil")
	}

	beqHandler, err := command.NewLoadBeqHandler(beqClient)
	if err != nil {
		return nil, err
	}

	beqGetHandler, err := query.NewGetBeqHandler(beqClient)
	if err != nil {
		return nil, err
	}
	return &App{
		Commands: Commands{
			LoadBeq: beqHandler,
		},
		Queries: Queries{
			GetBeq: beqGetHandler,
		},
	}, nil
}
