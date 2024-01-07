package avr

import (
	"github.com/reiver/go-telnet"
)
// AVRClient is an interface for interacting with any AVR
type AVRClient interface {
    GetCodec() (string, error)
}

// GetAVRClient returns a new instance of an AVRClient based on a brand like denon
func GetAVRClient(brand, url string) AVRClient {
    switch brand {
    case "denon":
        return &DenonClient{ServerURL: url, Port: "23", TelClient: telnet.StandardCaller}
    // Add cases for other brands
    default:
        return nil
    }
}