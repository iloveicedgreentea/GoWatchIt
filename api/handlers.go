package api

import (
	"github.com/gin-gonic/gin"
	"os"
	"path/filepath"
)

func GetConfigPath() (string, error) {
	ex, err := os.Executable()
	if err != nil {
		return "", err
	}

	exPath := filepath.Dir(ex)
	configPath1 := filepath.Join(exPath, "../config.json")
	configPath2 := "/config.json" // Fallback path (for Docker)

	if _, err := os.Stat(configPath1); err == nil {
		return configPath1, nil
	} else if _, err := os.Stat(configPath2); err == nil {
		return configPath2, nil
	}
	return "", os.ErrNotExist
}


func GetConfig(c *gin.Context) {
	path, err := GetConfigPath()
	if err != nil {
		c.JSON(500, gin.H{"error": "unable to find config"})
		return
	}
	data, err := os.ReadFile(path)
	if err != nil {
		c.JSON(500, gin.H{"error": "unable to read config"})
		return
	}
	c.Data(200, "application/json", data)
}

func SaveConfig(c *gin.Context) {
	var jsonData map[string]interface{}
	if err := c.ShouldBindJSON(&jsonData); err != nil {
		c.JSON(400, gin.H{"error": "Invalid Payload"})
		return
	}

	// TODO: save config to viper or something
	c.JSON(200, gin.H{"message": "Config saved successfully"})
}