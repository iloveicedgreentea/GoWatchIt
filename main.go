package main

import (
	"fmt"
	"net/http"

	"github.com/fsnotify/fsnotify"
	"github.com/iloveicedgreentea/go-plex/handlers"
	"github.com/iloveicedgreentea/go-plex/logger"
	"github.com/iloveicedgreentea/go-plex/models"
	"github.com/spf13/viper"
)

func readConfig() (*viper.Viper, error) {
	v := viper.New()
	v.SetConfigFile("config.json")
	err := v.ReadInConfig()

	// Hot reload config
	v.OnConfigChange(func(e fsnotify.Event) {
		fmt.Println("Config file changed:", e.Name)
	})
	v.WatchConfig()

	return v, err
}

func main() {
	log := logger.GetLogger()
	vip, err := readConfig()
	if err != nil {
		log.Fatal(err)
	}
	log.Debug("Started in debug mode...")
	mux := http.NewServeMux()

	// you can copy this schema to create event handlers for any service
	// create channel to receive jobs
	var plexChan = make(chan models.PlexWebhookPayload, 5)
	var minidspChan = make(chan models.MinidspRequest, 5)

	// run worker forever in background
	go handlers.PlexWorker(plexChan, vip)
	go handlers.MiniDspWorker(minidspChan, vip)

	// pass the chan to the handlers
	plexWh := handlers.ProcessWebhook(plexChan, vip)
	dspWh := handlers.ProcessMinidspWebhook(minidspChan, vip)
	// jellyfin might one day be supported
	// jfWh := handlers.ProcessJellyfinWebhook(plexChan, vip)
	// healthcheck
	health := handlers.ProcessHealthcheckWebhook()

	// Add plex webhook handler
	mux.Handle("/plexwebhook", plexWh)
	// minidsp
	mux.Handle("/minidspwebhook", dspWh)
	// jellyfin
	// mux.Handle("/jellyfinwebhook", jfWh)

	// healthcheck
	mux.Handle("/health", health)


	log.Info("Starting server")
	err = http.ListenAndServe(fmt.Sprintf(":%s", vip.GetString("main.listenPort")), mux)
	log.Fatal(err)
}
