package config

import (
	"github.com/iloveicedgreentea/go-plex/internal/logger"

	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"

	"os"
	"path/filepath"
)

var v *viper.Viper
var log = logger.GetLogger()

func init() {
	v = viper.New()
	// Try the directory of the executable first
	ex, err := os.Executable()
	if err != nil {
		log.Fatalf("Could not find executable path: %v", err)
	}

	exPath := filepath.Dir(ex)
	// docker path
	configPath1 := "/data/config.json"
	// local
	configPath2 := filepath.Join(exPath, "../config.json")

	var found bool

	// Try the first path
	v.SetConfigFile(configPath1)

	err = v.ReadInConfig()
	if err == nil {
		found = true
	}

	// If not found, try the fallback path
	if !found {
		v.SetConfigFile(configPath2)
		err = v.ReadInConfig()
		if err == nil {
			found = true
		}
	}
	// if in a testing env
	if !found {

		// Try the directory of the executable first
		cwd, err := os.Getwd()
		if err != nil {
			log.Fatalf("Could not find executable path: %v", err)
		}
		configpathTest := filepath.Join(cwd, "../../config.json")
		v.SetConfigFile(configpathTest)
		err = v.ReadInConfig()
		if err == nil {
			found = true
		}
	}

	// If still not found, log an error
	if !found {
		log.Error("no suitable config file found. Please set one up in the UI")
		return
	}

	v.OnConfigChange(func(e fsnotify.Event) {
		log.Debugf("Config file changed: %s", e.Name)
		log.Info("Config reloaded")
	})
	// hot reload for config
	v.WatchConfig()
}

func Set(key string, value interface{}) {
	v.Set(key, value)
}

func SaveConfigFile(configPath string) error {
	return v.WriteConfigAs(configPath)
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
