package config

// HDMI

func IsHDMISyncEnabled() bool {
	return GetBool("signal.enabled")
}

func IsSignalSourceTime() bool {
	return GetString("signal.source") == "time"
}

// BEQ

func IsBeqEnabled() bool {
	return GetBool("ezbeq.enabled")
}

func IsBeqTVEnabled() bool {
	return GetBool("ezbeq.enableTvBeq")
}

func IsBeqNotifyOnLoadEnabled() bool {
	return GetBool("ezbeq.notifyOnLoad")
}

func IsBeqDryRun() bool {
	return GetBool("ezbeq.dryRun")
}

// MQTT
func IsMQTTEnabled() bool {
	return GetBool("mqtt.enabled")
}

// Jellyfin

func IsJellyfinEnabled() bool {
	return GetBool("jellyfin.enabled")
}

func IsJellyfinSkipTMDB() bool {
	return GetBool("jellyfin.skiptmdb")
}

// Home Assistant

func IsHomeAssistantEnabled() bool {
	return GetBool("homeAssistant.enabled")
}

func IsHomeAssistantTriggerAVRMasterVolumeChangeOnEvent() bool {
	return GetBool("homeassistant.triggeravrmastervolumechangeonevent") && IsHomeAssistantEnabled()
}

func IsHomeAssistantTriggerLightsOnEvent() bool {
	return GetBool("homeassistant.triggerlightsonevent") && IsHomeAssistantEnabled()
}

