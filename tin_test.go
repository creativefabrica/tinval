package tinval_test

import (
	"testing"

	"github.com/creativefabrica/tinval"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestParse(t *testing.T) {
	type args struct {
		s string
	}
	tests := []struct {
		name    string
		args    args
		want    tinval.TIN
		wantErr error
	}{
		{
			name: "invalid ID number (empty string)",
			args: args{
				s: "    ",
			},
			want:    tinval.TIN{},
			wantErr: tinval.ErrInvalidFormat,
		},
		{
			name: "invalid ID number (too short)",
			args: args{
				s: "NL",
			},
			want:    tinval.TIN{},
			wantErr: tinval.ErrInvalidFormat,
		},
		{
			name: "valid AU Tax ID number",
			args: args{
				s: "AU51824753556",
			},
			want: tinval.TIN{
				CountryCode: "AU",
				Number:      "51824753556",
			},
			wantErr: nil,
		},
		{
			name: "invalid AU Tax ID number format (bad length)",
			args: args{
				s: "AU5182475355",
			},
			want:    tinval.TIN{},
			wantErr: tinval.ErrInvalidFormat,
		},
		{
			name: "invalid AU Tax ID number format (bad check digits)",
			args: args{
				s: "AU41824753556",
			},
			want:    tinval.TIN{},
			wantErr: tinval.ErrInvalidFormat,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tinval.Parse(tt.args.s)
			assert.Equal(t, tt.want, got)
			require.ErrorIs(t, tt.wantErr, err)
		})
	}
}
