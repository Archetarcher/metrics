package config

import (
	"encoding/json"
	"github.com/Archetarcher/metrics.git/internal/server/logger"
	"os"
)

func (c *AppConfig) parseJson() {
	configFile, err := os.Open(c.ConfigPath)
	defer configFile.Close()
	if err != nil {
		logger.Log.Info("failed to read json config")
	}
	jsonParser := json.NewDecoder(configFile)
	jErr := jsonParser.Decode(c)

	if jErr != nil {
		logger.Log.Info("failed to read json config")
	}
}
