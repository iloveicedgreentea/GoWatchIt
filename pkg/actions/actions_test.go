package actions

func main() {
    // ... other setup code ...

    haClient := homeassistant.NewClient(/* ... */)
    actionHandler := control.NewActionHandler()

    // Register actions
    actionHandler.RegisterAction("changeVolume", &control.VolumeAction{HAClient: haClient})
    actionHandler.RegisterAction("changeLight", &control.LightAction{HAClient: haClient})
    actionHandler.RegisterAction("waitForSync", &control.SyncAction{HAClient: haClient})
    actionHandler.RegisterAction("sendHAEvent", &control.HAEventAction{HAClient: haClient})

    // Use the new action to send a Home Assistant event
    err := actionHandler.ExecuteAction(context.Background(), "sendHAEvent", map[string]interface{}{
        "event_type": "change_lights",
        "event_data": map[string]interface{}{
            "entity_id": "light.living_room",
            "brightness": 255,
        },
    })
    if err != nil {
        log.Printf("Error sending Home Assistant event: %v", err)
    }

    // ... rest of your application logic ...
}