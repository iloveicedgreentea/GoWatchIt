package avr

import (
	"bytes"
	"fmt"
	"strings"
	"time"

	"github.com/iloveicedgreentea/go-plex/internal/logger"
	"github.com/reiver/go-telnet"
)

var log = logger.GetLogger()


// DenonClient is a client for Denon AVRs
type DenonClient struct {
	ServerURL string
	Port      string
	TelClient telnet.Caller
}


// make a request to denon via telnet
func (c *DenonClient) makeReq(command string) (string, error) {
	conn, err := telnet.DialTo(fmt.Sprintf("%s:%s", c.ServerURL, c.Port))
	if err != nil {
		return "", err
	}
	cmd := fmt.Sprintf("%s\r", command)
	log.Debugf("Sending command: %s", cmd)

	// send cmd
	_, err = conn.Write([]byte(cmd))
	if err != nil {
		return "", err
	}
	time.Sleep(500 * time.Millisecond)
	// receive
	var buffer [1]byte
	p := buffer[:]
	// final result
	var result []byte
	log.Debug("Receiving data")
	// Read until carriage return
	for {
		// this function is weird, it will just read the length of the byte[] given, but block
		// so you need to give it 1 length array and read 1 byte at a time
		n, err := conn.Read(p) // will block if nothing else to send
		if n > 0 {
			data := p[:n]
			// read response one at a time
			// if char is 13 (CR) then break
			if bytes.Equal(data, []byte{13}) {
				break
			}
			// store val in final result
			result = append(result, data[0])
		}

		if err != nil {
			break
		}
	}

	if err != nil {
		return "", err
	}

	log.Debugf("Got result: %s", string(result))

	return string(result), nil
}

// GetAudioMode returns the current audio mode like dolby atmos, stereo, etc
// Getting some weird garbage out of an AVR-X4700H I'm testing against, so
// I've added a 10 retry loop for if it gets any of that noise
func (c *DenonClient) GetCodec() (string, error) {
    var mode string
    var err error
    for i:=0; i<10; i++ {
        mode, err = c.makeReq("MS?")
        if strings.Contains(mode, "OPSMLALL") || strings.Contains(mode, "OPINFINS") {
            time.Sleep(750 * time.Millisecond)
            continue
        } else {
            return strings.ToLower(mode[2:]), err
        }
    }
    return strings.ToLower(mode[2:]), err
}
