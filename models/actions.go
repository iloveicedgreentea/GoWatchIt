package models

type Action string

const (
	ActionPlay Action = "play"
	ActionPause Action = "pause"
	ActionStop  Action = "stop"
)