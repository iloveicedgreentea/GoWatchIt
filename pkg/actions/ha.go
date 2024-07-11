package actions

import (
    "context"
    "fmt"
    "github.com/iloveicedgreentea/go-plex/internal/homeassistant"
)

type VolumeAction struct {
    HAClient *homeassistant.HomeAssistantClient
}

func (va *VolumeAction) Execute(ctx context.Context, params map[string]interface{}) error {
    volume, ok := params["volume"].(float64)
    if !ok {
        return fmt.Errorf("invalid volume parameter")
    }
    // Implement volume change logic using HAClient
    return va.HAClient.SetVolume(volume)
}

type LightAction struct {
    HAClient *homeassistant.HomeAssistantClient
}

func (la *LightAction) Execute(ctx context.Context, params map[string]interface{}) error {
    state, ok := params["state"].(string)
    if !ok {
        return fmt.Errorf("invalid state parameter")
    }
    // Implement light change logic using HAClient
    return la.HAClient.SetLight(state)
}

type SyncAction struct {
    HAClient *homeassistant.HomeAssistantClient
}

func (sa *SyncAction) Execute(ctx context.Context, params map[string]interface{}) error {
    timeout, ok := params["timeout"].(int)
    if !ok {
        timeout = 30 // default timeout
    }
    // Implement sync wait logic using HAClient
    return sa.HAClient.WaitForSync(ctx, timeout)
}

type HAEventAction struct {
    HAClient *homeassistant.HomeAssistantClient
}

func (ha *HAEventAction) Execute(ctx context.Context, params map[string]interface{}) error {
    eventType, ok := params["event_type"].(string)
    if !ok {
        return fmt.Errorf("invalid event_type parameter")
    }

    eventData, ok := params["event_data"].(map[string]interface{})
    if !ok {
        eventData = make(map[string]interface{})
    }

    return ha.HAClient.SendEvent(eventType, eventData)
}