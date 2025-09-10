package tinval_test

import (
	"testing"

	"github.com/creativefabrica/tinval"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestParse(t *testing.T) {
	tests := []struct {
		name string // description of this test case
		// Named input parameters for target function.
		s           string
		countryCode string
		want        tinval.TIN
		wantErr     error
	}{
		{
			name:        "invalid ID number (empty string)",
			s:           "    ",
			want:        tinval.TIN{},
			countryCode: "NL",
			wantErr:     tinval.ErrInvalidFormat,
		},
		{
			name:        "invalid ID number (too short)",
			s:           "NL",
			countryCode: "NL",
			want:        tinval.TIN{},
			wantErr:     tinval.ErrInvalidFormat,
		},
		{
			name:        "valid AU Tax ID number",
			s:           "51824753556",
			countryCode: "AU",
			want: tinval.TIN{
				CountryCode: "AU",
				Number:      "51824753556",
			},
			wantErr: nil,
		},
		{
			name:        "invalid AU Tax ID number format (bad length)",
			s:           "1234567891",
			countryCode: "AU",
			want:        tinval.TIN{},
			wantErr:     tinval.ErrInvalidFormat,
		},
		{
			name:        "invalid AU Tax ID number format (bad check digits)",
			s:           "41824753556",
			countryCode: "AU",
			want:        tinval.TIN{},
			wantErr:     tinval.ErrInvalidFormat,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, gotErr := tinval.Parse(tt.s, tt.countryCode)
			assert.Equal(t, tt.want, got)
			require.ErrorIs(t, tt.wantErr, gotErr)
		})
	}
}
