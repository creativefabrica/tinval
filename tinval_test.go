package tinval_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/creativefabrica/tinval"
	"github.com/creativefabrica/tinval/tinvaltest"
)

func TestValidator_Validate(t *testing.T) {
	validationClientMock := tinvaltest.NewMockValidationClient(t)
	validator := tinval.NewValidator(
		tinval.WithUKVATClient(validationClientMock),
		tinval.WithEUVATClient(validationClientMock),
		tinval.WithANBClient(validationClientMock),
	)

	tests := []struct {
		name      string
		tin       string
		country   string
		setupMock func(*tinvaltest.MockValidationClient, *testing.T)
		wantErr   error
	}{
		{
			name:    "valid VAT",
			tin:     "NL822010690B01",
			country: "NL",
			setupMock: func(mock *tinvaltest.MockValidationClient, t *testing.T) {
				ctx := t.Context()
				id := tinval.MustParse("NL822010690B01", "NL")
				mock.EXPECT().Validate(ctx, id).Return(nil)
			},
			wantErr: nil,
		},
		{
			name:    "valid ABN",
			tin:     "51824753556",
			country: "AU",
			setupMock: func(mock *tinvaltest.MockValidationClient, t *testing.T) {
				ctx := t.Context()
				id := tinval.MustParse("51824753556", "AU")
				mock.EXPECT().Validate(ctx, id).Return(nil)
			},
			wantErr: nil,
		},
		{
			name:    "valid UK VAT",
			tin:     "GB146295999727",
			country: "GB",
			setupMock: func(mock *tinvaltest.MockValidationClient, t *testing.T) {
				ctx := t.Context()
				id := tinval.MustParse("GB146295999727", "GB")
				mock.EXPECT().Validate(ctx, id).Return(nil)
			},
			wantErr: nil,
		},
		{
			name:    "invalid TIN",
			tin:     "NL822010690B02",
			country: "NL",
			setupMock: func(mock *tinvaltest.MockValidationClient, t *testing.T) {
				ctx := t.Context()
				id := tinval.MustParse("NL822010690B02", "NL")
				mock.EXPECT().Validate(ctx, id).Return(tinval.ErrInvalidFormat)
			},
			wantErr: tinval.ErrInvalidFormat,
		},
		{
			name:    "service unavailable",
			tin:     "NL822010690B03",
			country: "NL",
			setupMock: func(mock *tinvaltest.MockValidationClient, t *testing.T) {
				ctx := t.Context()
				id := tinval.MustParse("NL822010690B03", "NL")
				mock.EXPECT().Validate(ctx, id).Return(tinval.ErrServiceUnavailable)
			},
			wantErr: tinval.ErrServiceUnavailable,
		},
		{
			name:    "not found",
			tin:     "NL822010690B04",
			country: "NL",
			setupMock: func(mock *tinvaltest.MockValidationClient, t *testing.T) {
				ctx := t.Context()
				id := tinval.MustParse("NL822010690B04", "NL")
				mock.EXPECT().Validate(ctx, id).Return(tinval.ErrNotFound)
			},
			wantErr: tinval.ErrNotFound,
		},
		{
			name:    "invalid country code",
			tin:     "822010690B05",
			country: "AR",
			setupMock: func(mock *tinvaltest.MockValidationClient, t *testing.T) {
				// No mock setup needed for this case as it fails before reaching the validation client
			},
			wantErr: tinval.ErrUnsupportedTaxType,
		},
		{
			name:    "invalid TIN length",
			tin:     "NL",
			country: "NL",
			setupMock: func(mock *tinvaltest.MockValidationClient, t *testing.T) {
				// No mock setup needed for this case as it fails before reaching the validation client
			},
			wantErr: tinval.ErrInvalidFormat,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setupMock(validationClientMock, t)

			ctx := t.Context()
			err := validator.Validate(ctx, tt.tin, tt.country)

			if tt.wantErr == nil {
				assert.NoError(t, err)
			} else {
				assert.ErrorIs(t, err, tt.wantErr)
			}
		})
	}
}
