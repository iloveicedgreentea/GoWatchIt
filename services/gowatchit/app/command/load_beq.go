package command

import (
	"context"
	"errors"

	"github.com/iloveicedgreentea/gowatchit/services/gowatchit/domain/beq"
)

type LoadBeqCommand struct {
	Payload *beq.BEQPayload
}

type LoadBeqHandler struct {
	BeqClient *beq.BeqClient
}

func NewLoadBeqHandler(beqClient *beq.BeqClient) (*LoadBeqHandler, error) {
	if beqClient == nil {
		return nil, errors.New("beqClient is nil")
	}
	return &LoadBeqHandler{
		BeqClient: beqClient,
	}, nil
}

// Handle loads BEQ to device
func (h *LoadBeqHandler) Handle(ctx context.Context, cmd *LoadBeqCommand) error {
	return h.BeqClient.LoadBeqProfile(ctx, cmd.Payload)
}
