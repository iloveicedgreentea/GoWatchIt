package actions

import (
    "context"
)

type ActionHandler struct {
    Actions map[string]Action
}

type Action interface {
    Execute(ctx context.Context, params map[string]interface{}) error
}

func NewActionHandler() *ActionHandler {
    return &ActionHandler{
        Actions: make(map[string]Action),
    }
}

func (ah *ActionHandler) RegisterAction(name string, action Action) {
    ah.Actions[name] = action
}

func (ah *ActionHandler) ExecuteAction(ctx context.Context, name string, params map[string]interface{}) error {
    action, exists := ah.Actions[name]
    if !exists {
        return fmt.Errorf("unknown action: %s", name)
    }
    return action.Execute(ctx, params)
}