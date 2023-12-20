package api

import (
	"fmt"
	"os"

	"path/filepath"

	"github.com/gin-gonic/gin"
	"github.com/iloveicedgreentea/go-plex/internal/config"
	"github.com/iloveicedgreentea/go-plex/internal/logger"
)

var log = logger.GetLogger()

func GenConfigPaths() (string, string){
	ex, err := os.Executable()
	if err != nil {
		log.Error(err)
	}

	exPath := filepath.Dir(ex)
	configPath1 := "/config.json" // docker
	configPath2 := filepath.Join(exPath, "../config.json") // Fallback path (for local)

	log.Debugf("Config paths: %s, %s", configPath1, configPath2)

	return configPath1, configPath2
}

// GetConfigPath returns the path to the config file
func GetConfigPath() (string, error) {
	configPath1, configPath2 := GenConfigPaths()

	if _, err := os.Stat(configPath1); err == nil {
		return configPath1, nil
	} else if _, err := os.Stat(configPath2); err == nil {
		return configPath2, nil
	}
	return "", os.ErrNotExist
}

// ConfigExists checks if the config exists for the API
func ConfigExists(c *gin.Context) {
	configPath, err := GetConfigPath()
	if err != nil {
		c.JSON(500, gin.H{"exists": false})
		return
	}
	_, err = os.Stat(configPath)
	c.JSON(200, gin.H{"exists": err == nil})
}

// GetConfig returns the config for the API
func GetConfig(c *gin.Context) {
	path, err := GetConfigPath()
	// if not found, create it
	if err != nil {
		log.Debugf("Didn't get config: %v", err)
		err = CreateConfig(c)
		if err != nil {
			log.Debugf("Didn't create config: %v", err)
			c.JSON(500, gin.H{"error": "unable to create config"})
			return
		}
	}
	data, err := os.ReadFile(path)
	if err != nil {
		log.Debugf("Didn't read config: %v", err)
		c.JSON(500, gin.H{"error": "unable to read config"})
		return
	}
	c.Data(200, "application/json", data)
}

// CreateConfig creates a new config file
func CreateConfig(c *gin.Context) error {
	log.Debug("Creating new config")
	configPath1, configPath2 := GenConfigPaths()

	// try to create config in the first path
	file, err := os.Create(configPath1)
	if err != nil {
		log.Debugf("Unable to create config in %s: %v", configPath1, err)
		// try to create config in the second path
		file, err = os.Create(configPath2)
		if err != nil {
			// if we can't create it in either path, return the error
			log.Errorf("Unable to create config in %s: %v", configPath2, err)
			return fmt.Errorf("unable to create config in %s or %s", configPath1, configPath2)
		}
	}
	defer file.Close()

	log.Debug("Successfully created config file")
	return nil
}

// SaveConfig saves the config for the API
func SaveConfig(c *gin.Context) {
	var jsonData map[string]interface{}

	if err := c.ShouldBindJSON(&jsonData); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		fmt.Println(c.Request.Body)
		return
	}

	path, err := GetConfigPath()
	if err != nil {
		log.Error("unable to get config")
		c.JSON(500, gin.H{"error": "unable to get config"})
		return
	}

	// Loop through the incoming JSON map to set keys in Viper
	for key, value := range jsonData {
		switch v := value.(type) {
		case map[string]interface{}:
			for subKey, subValue := range v {
				config.Set(fmt.Sprintf("%s.%s", key, subKey), subValue)
			}
		default:
			config.Set(key, value)
		}
	}

	// Use your SaveConfigFile function to save the updated configuration
	if err := config.SaveConfigFile(path); err != nil {
		log.Error("unable to save config")
		c.JSON(500, gin.H{"error": "Unable to save config"})
		return
	}

	c.JSON(200, gin.H{"message": "Config saved successfully"})
}
