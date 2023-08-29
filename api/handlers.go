package api

import (
	"github.com/gin-gonic/gin"
    "os"
)

func GetConfig(c *gin.Context) {
	data, err := os.ReadFile("../config.json")
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