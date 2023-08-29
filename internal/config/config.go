package config

import (
	"log"

	"github.com/spf13/viper"
	"github.com/fsnotify/fsnotify"

	"os"
	"path/filepath"
)

var v *viper.Viper

func init() {
	v = viper.New()
	// Try the directory of the executable first
	ex, err := os.Executable()
	if err != nil {
		log.Fatalf("Could not find executable path: %v", err)
	}

	exPath := filepath.Dir(ex)
	configPath1 := filepath.Join(exPath, "../config.json")

	// Fallback path (for Docker)
	configPath2 := "/config.json"

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

	// If still not found, log an error
	if !found {
		log.Fatalf("Error: No suitable config file found.")
		return
	}

	v.OnConfigChange(func(e fsnotify.Event) {
		log.Println("Config file changed:", e.Name)
	})

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

