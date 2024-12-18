package logger

import (
	"go.uber.org/zap"
	"log"
)

// Log will be available as singleton.
var Log *zap.Logger

func init() {
	Log = zap.NewNop()

	defer func() {
		err := Log.Sync()
		if err != nil {
			log.Fatal("failed to init logger")
		}
	}()
}

// Initialize initiates singleton of Log with appropriate log level.
func Initialize(level string) error {

	lvl, err := zap.ParseAtomicLevel(level)
	if err != nil {
		return err
	}
	cfg := zap.NewProductionConfig()
	cfg.Level = lvl

	zl, err := cfg.Build()

	if err != nil {
		return err
	}
	Log = zl

	return nil
}
