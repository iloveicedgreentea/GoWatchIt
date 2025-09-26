package avr

import (
	"github.com/iloveicedgreentea/go-plex/internal/config"
	"github.com/reiver/go-telnet"
)

// AVRClient is an interface for interacting with any AVR
type AVRClient interface {
	GetCodec() (string, error)
}

// GetAVRClient returns a new instance of an AVRClient based on a brand like denon
func GetAVRClient() AVRClient {
	log.Debug(config.GetEZBeqAvrBrand())
	switch config.GetEZBeqAvrBrand() {
	case "denon":
		log.Debug("Creating Denon AVR client")
		return &DenonClient{ServerURL: config.GetEZBeqAvrURL(), Port: "23", TelClient: telnet.StandardCaller}
	// Add cases for other brands
	default:
		log.Error("unsupported AVR brand set in config")
		return nil
	}
}
