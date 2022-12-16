package models

type MinidspRequest struct {
	Command string `json:"command"`
}

type MinidspCommandRequest struct {
	Overwrite   bool   `json:"overwrite"`
	Slot        string `json:"slot"`
	Inputs      []int  `json:"inputs"`
	Outputs     []int  `json:"outputs"`
	CommandType string `json:"commandType"`
	Commands    string `json:"commands"`
}