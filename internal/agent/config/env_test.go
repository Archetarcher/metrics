package config

import (
	"reflect"
	"testing"
)

func TestAppConfig_parseEnv(t *testing.T) {
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
			fields: fields{ServerRunAddr: "", GRPCRunAddr: "", LogLevel: "", Key: "",
				PublicKeyPath: "", ConfigPath: "", Session: Session{}, ReportInterval: 0, PollInterval: 0, RateLimit: 0, EnableGRPC: false},
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
			c.parseEnv()
		})
	}
}

func Test_getEnvOrDefault(t *testing.T) {
	type args struct {
		env string
		def any
		t   int
	}
	c := &AppConfig{}
	tests := []struct {
		name string
		args args
		want any
	}{
		{
			name: "positive test 1",
			args: args{
				env: envReportIntervalName,
				def: c.ReportInterval,
				t:   3,
			},
			want: 0,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := getEnvOrDefault(tt.args.env, tt.args.def, tt.args.t); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("getEnvOrDefault() = %v, want %v", got, tt.want)
			}
		})
	}
}
