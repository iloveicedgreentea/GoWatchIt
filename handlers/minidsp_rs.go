package handlers

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/iloveicedgreentea/go-plex/ezbeq"
	"github.com/iloveicedgreentea/go-plex/models"
	"github.com/spf13/viper"
)

// https://minidsp-rs.pages.dev/cli/master/mute

// based on event type, determine what to do
func minidspRouter(payload models.MinidspRequest, vip *viper.Viper, beqClient *ezbeq.BeqClient) {
	log.Debug(payload.Command)

	switch {
	case strings.Contains(payload.Command, "off"):
		muteOff(beqClient)
	case strings.Contains(payload.Command, "on"):
		muteOn(beqClient)
	}
}

// send minidsp command via ezbeq
func doMinidspCommand(action string, beqClient *ezbeq.BeqClient) {
	r := models.MinidspCommandRequest{
		Overwrite: true,
		Slot: "1",
		Inputs: []int{1, 2},
		Outputs: []int{1, 2, 3, 4},
		CommandType: "rs",
		Commands: action,
	}

	j, err := json.Marshal(r)
	if err != nil {
		log.Error(err)
	}

	beqClient.MakeCommand(j)

}
// TODO: test this
// TODO: add this to home assistant via automation -> harmony button -> HA automation -> this path
func muteOn(beqClient *ezbeq.BeqClient) {
	log.Debug("running mute on")
	doMinidspCommand("mute on", beqClient)
}

func muteOff(beqClient *ezbeq.BeqClient) {
	log.Debug("running mute off")
	doMinidspCommand("mute off", beqClient)
}

// process webhook 
func ProcessMinidspWebhook(miniDsp chan<- models.MinidspRequest, vip *viper.Viper) http.Handler {
	log.Debug("minidsp triggered")
	fn := func(w http.ResponseWriter, r *http.Request) {
		var payload models.MinidspRequest

		err := json.NewDecoder(r.Body).Decode(&payload)
		if err != nil {
			log.Error(err)
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		
		miniDsp <- payload
	}

	return http.HandlerFunc(fn)
}

// entry point for background tasks
func MiniDspWorker(minidspChan <-chan models.MinidspRequest, vip *viper.Viper) {
	log.Info("Minidsp started")

	var beqClient *ezbeq.BeqClient
	var err error

	if vip.GetBool("ezbeq.enabled") {
		log.Info("Started with ezbeq enabled")
		beqClient, err = ezbeq.NewClient(vip.GetString("ezbeq.url"), vip.GetString("ezbeq.port"))
		if err != nil {
			log.Error(err)
		}
	}

	// block forever until closed so it will wait in background for work
	for i := range minidspChan {
		// determine what to do
		minidspRouter(i, vip, beqClient)
	}
}