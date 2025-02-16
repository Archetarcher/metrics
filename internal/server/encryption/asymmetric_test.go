package encryption

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestAsymmetric_DecryptAsymmetric(t *testing.T) {

	type args struct {
		key  string
		text []byte
	}

	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "negative test #1",
			args: args{
				text: []byte("teststring"),
				key:  "../../private.pem",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			_, eErr := NewAsymmetric(tt.args.key).Decrypt(tt.args.text)

			assert.Equal(t, tt.wantErr, eErr != nil)
		})
	}
}
