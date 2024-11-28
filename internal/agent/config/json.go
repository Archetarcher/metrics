package config

import (
	"encoding/json"
	"github.com/Archetarcher/metrics.git/internal/server/logger"
	"log"
	"os"
)

func (c *AppConfig) parseJSON() {
	configFile, err := os.Open(c.ConfigPath)
	defer func() {
		cErr := configFile.Close()
		if cErr != nil {
			log.Fatal("failed to close file")
		}
	}()
	if err != nil {
		logger.Log.Info("failed to read json config")
	}
	jsonParser := json.NewDecoder(configFile)
	jErr := jsonParser.Decode(c)

	if jErr != nil {
		logger.Log.Info("failed to read json config")
	}
}
