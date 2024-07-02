package logger

import (
	"bytes"
	"encoding/json"
	"go.uber.org/zap"
	"net/http"
)

type (
	responseData struct {
		status int
		size   int
	}
	loggerResponseWriter struct {
		http.ResponseWriter
		responseData *responseData
	}
)

// Log будет доступен всему коду как синглтон.
var Log *zap.Logger = zap.NewNop()

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

func LogData(message string, data interface{}) {
	var buf bytes.Buffer

	jsonEncoder := json.NewEncoder(&buf)
	if err := jsonEncoder.Encode(data); err != nil {
		Log.Debug("error encoding response", zap.Error(err))
		return
	}

	Log.Info(message, zap.Any("data", json.RawMessage(buf.Bytes())))

}
