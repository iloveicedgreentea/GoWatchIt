package webhooks

import (
    "net/http"
    "context"
)

type WebhookHandler struct {
    Processors map[string]WebhookProcessor
}

type WebhookProcessor interface {
    Process(r *http.Request) error
    StartWorker(ctx context.Context, ready chan<- bool)
}

func NewWebhookHandler() *WebhookHandler {
    return &WebhookHandler{
        Processors: make(map[string]WebhookProcessor),
    }
}

func (wh *WebhookHandler) RegisterProcessor(name string, processor WebhookProcessor) {
    wh.Processors[name] = processor
}

func (wh *WebhookHandler) HandleWebhook(w http.ResponseWriter, r *http.Request) {
    processorName := r.URL.Path[len("/webhook/"):]
    processor, exists := wh.Processors[processorName]
    if !exists {
        http.Error(w, "Unknown webhook type", http.StatusBadRequest)
        return
    }

    err := processor.Process(r)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    w.WriteHeader(http.StatusOK)
}