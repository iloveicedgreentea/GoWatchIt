package avr

import (
	"github.com/reiver/go-telnet"
	"github.com/iloveicedgreentea/go-plex/internal/config"
)
// AVRClient is an interface for interacting with any AVR
type AVRClient interface {
    GetCodec() (string, error)
}

// GetAVRClient returns a new instance of an AVRClient based on a brand like denon
func GetAVRClient(url string) AVRClient {
    
    log.Debug(config.GetString("ezbeq.avrbrand"))
    switch config.GetString("ezbeq.avrbrand") {
    case "denon":
        log.Debug("Creating Denon AVR client")
        return &DenonClient{ServerURL: url, Port: "23", TelClient: telnet.StandardCaller}
    // Add cases for other brands
    default:
		log.Error("No AVR brand set in config")
        return nil
    }
}