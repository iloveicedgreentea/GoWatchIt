package avr

import (
	"bytes"
	"fmt"
	"log/slog"
	"strings"
	"time"

	"github.com/iloveicedgreentea/go-plex/internal/common"
	"github.com/iloveicedgreentea/go-plex/internal/logger"
	"github.com/reiver/go-telnet"
)

var log = logger.GetLogger()

// DenonClient is an AVRClient for Denon AVRs
type DenonClient struct {
	ServerURL string
	Port      string
	TelClient telnet.Caller
}

// make a request to denon via telnet
func (c *DenonClient) makeReq(command string) (string, error) {
	// TODO: redo to add context timeout here
	conn, err := telnet.DialTo(fmt.Sprintf("%s:%s", c.ServerURL, c.Port))
	if err != nil {
		return "", err
	}
	cmd := fmt.Sprintf("%s\r", command)

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
			// store val in final result
			result = append(result, data[0])
			// read response one at a time
			// if char is 13 (CR) then break
			if bytes.Equal(data, []byte{13}) {
				break
			}
		}

		if err != nil {
			break
		}
	}

	log.Debug("result", slog.String("result", string(result)))

	return string(result), nil
}

// MapDenonToBeq takes denon codec names and Client codecs and normalizes it to BEQ mappings
func MapDenonToBeq(denonCodec, clientCodec string) string {
	// take plex data also to get channel info. AVR can only say truehd, not 7.1 or 5.1
	// TODO: check if atmos, otherwise check dtsx, , etc, truehd and compare to plex data, etc
	switch {
	// TODO: what about DD+ Atmos? AVR will return Atmos but what does file name say
	// TODO: if denon says atmos and plex says e-ac3 or eac3 5.1 etc then DD+ Atmos
	// TODO: if plex says truehd and avr says atmos, return atmos
	case common.InsensitiveContains(denonCodec, "dolby atmos"):
		return "Atmos"
	//  what is this codec?
	// case common.InsensitiveContains(denonCodec, "dolby hd"):
	// 	return "AtmosMaybe"
	case common.InsensitiveContains(denonCodec, "DOLBY DIGITAL +"):
		return "DD+"
	// TODO: test
	case common.InsensitiveContains(denonCodec, "DTS:X"):
		return "DTS-X"

	// TODO: to implement
	// "DOLBY DIGITAL",
	//     "DOLBY DIGITAL +",
	//     "STANDARD(DOLBY)",
	//     "DOLBY SURROUND",
	//     "DOLBY HD",
	//     "DOLBY ATMOS",
	//     "DOLBY AUDIO - DOLBY SURROUND",
	//     "DOLBY TRUEHD",
	//     "DOLBY AUDIO - DOLBY DIGITAL PLUS",
	//     "DOLBY AUDIO - TRUEHD + DSUR",
	//     "DOLBY AUDIO - DOLBY TRUEHD",
	//     "DOLBY AUDIO - TRUEHD + NEURAL:X",
	//     "DOLBY AUDIO - DD + NEURAL:X",
	//     "DOLBY AUDIO - DD + DSUR",
	//     "DOLBY AUDIO - DD+   + NEURAL:X",
	//     "DOLBY AUDIO - DD+   + DSUR",
	//     "DOLBY AUDIO-DD+ +DSUR",
	//     "DOLBY AUDIO - DOLBY DIGITAL",
	//     "DOLBY AUDIO-DSUR",
	//     "DOLBY AUDIO-DD+DSUR",
	//     "DOLBY PRO LOGIC",
	// DTS-HD MSTR
	// DTS-HD
	// DTS
	// DTS SURROUND
	case common.InsensitiveContains(denonCodec, "DTS-HD MA 7.1") && !common.InsensitiveContains(denonCodec, "DTS:X") && !common.InsensitiveContains(denonCodec, "DTS-X"):
		return "DTS-HD MA 7.1"
	// DTS HA MA 5.1
	case common.InsensitiveContains(denonCodec, "DTS-HD MA 5.1"):
		return "DTS-HD MA 5.1"
	// DTS 5.1
	case common.InsensitiveContains(denonCodec, "DTS 5.1"):
		return "DTS 5.1"
	// TrueHD 5.1
	case common.InsensitiveContains(denonCodec, "TRUEHD 5.1"):
		return "TrueHD 5.1"
	// TrueHD 6.1
	case common.InsensitiveContains(denonCodec, "TRUEHD 6.1"):
		return "TrueHD 6.1"
	// DTS HRA
	case common.InsensitiveContains(denonCodec, "DTS-HD HRA 7.1"):
		return "DTS-HD HR 7.1"
	case common.InsensitiveContains(denonCodec, "DTS-HD HRA 5.1"):
		return "DTS-HD HR 5.1"
	// LPCM
	case common.InsensitiveContains(denonCodec, "LPCM 5.1"):
		return "LPCM 5.1"
	case common.InsensitiveContains(denonCodec, "LPCM 7.1"):
		return "LPCM 7.1"
	case common.InsensitiveContains(denonCodec, "LPCM 2.0"):
		return "LPCM 2.0"
	case common.InsensitiveContains(denonCodec, "AAC Stereo"):
		return "AAC 2.0"
	case common.InsensitiveContains(denonCodec, "AC3 5.1") || common.InsensitiveContains(denonCodec, "EAC3 5.1"):
		return "AC3 5.1"
	case common.InsensitiveContains(denonCodec, "stereo"):
		return "Stereo"
	default:
		return "Empty"
	}
}

// GetAudioMode returns the current audio mode like dolby atmos, stereo, etc. Best used to detect atmos from truehd
func (c *DenonClient) GetCodec() (string, error) {
	mode, err := c.makeReq("MS?")
	res := strings.ToLower(mode[1:])
	log.Debug("res", slog.Any("resp", res))
	return res, err
}
