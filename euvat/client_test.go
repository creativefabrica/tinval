package euvat_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/creativefabrica/tinval"
	"github.com/creativefabrica/tinval/euvat"
)

func Test_Client_Validate(t *testing.T) {
	tests := []struct {
		name      string
		vatNumber tinval.TIN
		wantErr   error
	}{
		{
			name:      "valid TIN",
			vatNumber: tinval.MustParse("NL822010690B01"),
			wantErr:   nil,
		},
		{
			name:      "non existing TIN",
			vatNumber: tinval.MustParse("NL822010690B02"),
			wantErr:   tinval.ErrNotFound,
		},
		{
			name:      "invalid format",
			vatNumber: tinval.TIN{CountryCode: "XX", Number: "822010690B01"},
			wantErr:   tinval.ErrInvalidFormat,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := euvat.NewClient()
			err := c.Validate(t.Context(), tt.vatNumber)
			assert.ErrorIs(t, err, tt.wantErr)
		})
	}
}
