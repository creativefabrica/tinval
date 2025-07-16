package abn_test

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/creativefabrica/tinval"
	"github.com/creativefabrica/tinval/abn"
)

func TestClient_Validate(t *testing.T) {
	tests := []struct {
		name      string
		guid      string
		vatNumber tinval.TIN
		wantErr   error
	}{
		{
			name:      "Missing credentials",
			guid:      "",
			vatNumber: tinval.MustParse("AU51824753556"),
			wantErr:   tinval.ErrServiceUnavailable,
		},
		{
			name:      "Valid TIN number",
			guid:      os.Getenv("ABN_API_AUTH_GUID"),
			vatNumber: tinval.MustParse("AU51824753556"),
			wantErr:   nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := abn.NewClient(tt.guid)
			err := c.Validate(t.Context(), tt.vatNumber)
			assert.ErrorIs(t, err, tt.wantErr)
		})
	}
}
