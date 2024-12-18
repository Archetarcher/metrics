package encryption

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestAsymmetric_Encrypt(t *testing.T) {

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
				key:  "../../../public.pem",
			},
			wantErr: false,
		},
		{
			name: "negative test #2",
			args: args{
				text: []byte("teststring"),
				key:  "../../public.pem",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			_, eErr := NewAsymmetric(tt.args.key).Encrypt(tt.args.text)

			assert.Equal(t, tt.wantErr, eErr != nil)
		})
	}
}

func TestGenKey(t *testing.T) {
	tests := []struct {
		name    string
		n       int
		wantErr bool
	}{
		{
			name:    "positive test #1",
			n:       16,
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := GenKey(tt.n)

			assert.Equal(t, tt.wantErr, err != nil)
		})
	}
}
