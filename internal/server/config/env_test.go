package config

import (
	"sync"
	"testing"
)

func TestAppConfig_parseEnv(t *testing.T) {
	type fields struct {
		RunAddr         string
		GRPCRunAddr     string
		Key             string
		Session         string
		LogLevel        string
		MigrationsPath  string
		TrustedSubnet   string
		PrivateKeyPath  string
		FileStoragePath string
		DatabaseDsn     string
		ConfigPath      string
		StoreInterval   int
		Restore         bool
		EnableGRPC      bool
	}
	tests := []struct {
		name   string
		fields fields
	}{
		{
			name: "positive test 1",
			fields: fields{
				RunAddr:         "",
				GRPCRunAddr:     "",
				Key:             "",
				Session:         "",
				LogLevel:        "",
				MigrationsPath:  "",
				TrustedSubnet:   "",
				PrivateKeyPath:  "",
				FileStoragePath: "",
				DatabaseDsn:     "",
				ConfigPath:      "",
				StoreInterval:   0,
				Restore:         false,
				EnableGRPC:      false,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &AppConfig{
				RunAddr:         tt.fields.RunAddr,
				GRPCRunAddr:     tt.fields.GRPCRunAddr,
				Key:             tt.fields.Key,
				Session:         tt.fields.Session,
				LogLevel:        tt.fields.LogLevel,
				MigrationsPath:  tt.fields.MigrationsPath,
				TrustedSubnet:   tt.fields.TrustedSubnet,
				PrivateKeyPath:  tt.fields.PrivateKeyPath,
				FileStoragePath: tt.fields.FileStoragePath,
				DatabaseDsn:     tt.fields.DatabaseDsn,
				ConfigPath:      tt.fields.ConfigPath,
				StoreInterval:   tt.fields.StoreInterval,
				Restore:         tt.fields.Restore,
				EnableGRPC:      tt.fields.EnableGRPC,
				mux:             sync.Mutex{},
			}
			c.parseEnv()
		})
	}
}
