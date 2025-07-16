package ukvat_test

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/creativefabrica/tinval"
	"github.com/creativefabrica/tinval/ukvat"
)

func TestClient_Validate(t *testing.T) {
	tests := []struct {
		name      string
		creds     ukvat.ClientCredentials
		vatNumber tinval.TIN
		wantErr   error
	}{
		{
			name:      "Missing credentials",
			creds:     ukvat.ClientCredentials{},
			vatNumber: tinval.MustParse("GB123456789"),
			wantErr:   tinval.ErrServiceUnavailable,
		},
		{
			name: "Valid TIN",
			creds: ukvat.ClientCredentials{
				ID:     os.Getenv("UKVAT_API_CLIENT_ID"),
				Secret: os.Getenv("UKVAT_API_CLIENT_SECRET"),
			},
			vatNumber: tinval.MustParse("GB146295999727"),
			wantErr:   nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := ukvat.NewClient(
				tt.creds,
				ukvat.WithBaseURL(ukvat.TestServiceBaseURL),
			)
			err := c.Validate(t.Context(), tt.vatNumber)
			assert.ErrorIs(t, err, tt.wantErr)
		})
	}
}
