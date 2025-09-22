package query

import (
	"context"

	"github.com/iloveicedgreentea/gowatchit/services/gowatchit/domain/beq"
)

type GetBeqQuery struct {
	Payload *beq.BEQPayload
}

type GetBeqHandler struct {
	BeqClient *beq.BeqClient
}

func NewGetBeqHandler(beqClient *beq.BeqClient) (*GetBeqHandler, error) {
	if beqClient == nil {
		return nil, nil
	}
	return &GetBeqHandler{
		BeqClient: beqClient,
	}, nil
}

// Handle gets BEQ profile from device
func (h *GetBeqHandler) Handle(ctx context.Context, cmd *GetBeqQuery) (map[string]string, error) {
	err := h.BeqClient.GetStatus(ctx)
	if err != nil {
		return nil, err
	}
	return h.BeqClient.GetCurrentProfile(ctx)
}
