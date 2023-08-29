package config

import (
	"log"

	"github.com/spf13/viper"
	"github.com/fsnotify/fsnotify"
	"fmt"
)

var v *viper.Viper

func init() {
	v = viper.New()
	v.SetConfigFile("../../config.json")
	if err := v.ReadInConfig(); err != nil {
		log.Fatalf("error reading in config file: %v", err)
	}
	v.OnConfigChange(func(e fsnotify.Event) {
		fmt.Println("Config file changed:", e.Name)
	})
	v.WatchConfig()
}

func GetString(key string) string {
	return v.GetString(key)
}

func GetBool(key string) bool {
	return v.GetBool(key)
}

func GetInt(key string) int {
	return v.GetInt(key)
}

func GetIntSlice(key string) []int {
	return v.GetIntSlice(key)
}

