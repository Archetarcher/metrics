package config

import "testing"

func TestAppConfig_parseJSON(t *testing.T) {
	type fields struct {
		ServerRunAddr  string
		GRPCRunAddr    string
		LogLevel       string
		Key            string
		PublicKeyPath  string
		ConfigPath     string
		Session        Session
		ReportInterval int
		PollInterval   int
		RateLimit      int
		EnableGRPC     bool
	}
	tests := []struct {
		name   string
		fields fields
	}{
		{
			fields: fields{
				ServerRunAddr:  "",
				GRPCRunAddr:    "",
				LogLevel:       "",
				Key:            "",
				PublicKeyPath:  "",
				ConfigPath:     "../../../agent-config.json",
				Session:        Session{},
				ReportInterval: 0,
				PollInterval:   0,
				RateLimit:      0,
				EnableGRPC:     false,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &AppConfig{
				ServerRunAddr:  tt.fields.ServerRunAddr,
				GRPCRunAddr:    tt.fields.GRPCRunAddr,
				LogLevel:       tt.fields.LogLevel,
				Key:            tt.fields.Key,
				PublicKeyPath:  tt.fields.PublicKeyPath,
				ConfigPath:     tt.fields.ConfigPath,
				Session:        tt.fields.Session,
				ReportInterval: tt.fields.ReportInterval,
				PollInterval:   tt.fields.PollInterval,
				RateLimit:      tt.fields.RateLimit,
				EnableGRPC:     tt.fields.EnableGRPC,
			}
			c.parseJSON()
		})
	}
}
