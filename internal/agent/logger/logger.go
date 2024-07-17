package logger

import (
	"go.uber.org/zap"
)

// Log будет доступен всему коду как синглтон.
var Log *zap.Logger

func init() {
	Log = zap.NewNop()

	defer Log.Sync()
}

// Initialize инициализирует синглтон логера с необходимым уровнем логирования.
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
