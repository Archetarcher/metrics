package logger

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestInitialize(t *testing.T) {

	tests := []struct {
		name    string
		level   string
		wantErr bool
	}{
		{
			name:    "positive test #1",
			level:   "info",
			wantErr: false,
		},
		{
			name:    "negative test #2",
			level:   "infos",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := Initialize(tt.level)

			assert.Equal(t, tt.wantErr, err != nil)
		})
	}
}
