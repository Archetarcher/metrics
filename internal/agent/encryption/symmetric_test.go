package encryption

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestEncryptSymmetric(t *testing.T) {

	type args struct {
		text []byte
		key  string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "positive test #1",
			args: args{
				text: []byte("teststring"),
				key:  "secretkey",
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := NewSymmetric(tt.args.key).Encrypt(tt.args.text)

			assert.Equal(t, tt.wantErr, err != nil)
		})
	}
}
