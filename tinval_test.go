package tinval_test

import (
	"testing"

	"github.com/creativefabrica/tinval"
	"github.com/creativefabrica/tinval/tinvaltest"
	"github.com/stretchr/testify/assert"
)

func TestValidator_Validate(t *testing.T) {
	validationClientMock := tinvaltest.NewMockValidationClient(t)
	validator := tinval.NewValidator(
		tinval.WithUKVATClient(validationClientMock),
		tinval.WithEUVATClient(validationClientMock),
	)

	t.Run("valid TIN", func(t *testing.T) {
		ctx := t.Context()
		id := tinval.MustParse("NL822010690B01")
		validationClientMock.EXPECT().Validate(ctx, id).Return(nil)
		err := validator.Validate(ctx, id.String())
		assert.NoError(t, err)
	})

	t.Run("invalid TIN", func(t *testing.T) {
		ctx := t.Context()
		id := tinval.MustParse("NL822010690B02")
		validationClientMock.EXPECT().Validate(ctx, id).Return(tinval.ErrInvalidFormat)
		err := validator.Validate(ctx, id.String())
		assert.ErrorIs(t, err, tinval.ErrInvalidFormat)
	})

	t.Run("service unavailable", func(t *testing.T) {
		ctx := t.Context()
		id := tinval.MustParse("NL822010690B03")
		validationClientMock.EXPECT().Validate(ctx, id).Return(tinval.ErrServiceUnavailable)
		err := validator.Validate(ctx, id.String())
		assert.ErrorIs(t, err, tinval.ErrServiceUnavailable)
	})

	t.Run("not found", func(t *testing.T) {
		ctx := t.Context()
		id := tinval.MustParse("NL822010690B04")
		validationClientMock.EXPECT().Validate(ctx, id).Return(tinval.ErrNotFound)
		err := validator.Validate(ctx, id.String())
		assert.ErrorIs(t, err, tinval.ErrNotFound)
	})

	t.Run("invalid country code", func(t *testing.T) {
		ctx := t.Context()
		err := validator.Validate(ctx, "AR822010690B05")
		assert.ErrorIs(t, err, tinval.ErrInvalidCountryCode)
	})

	t.Run("invalid TIN length", func(t *testing.T) {
		ctx := t.Context()
		err := validator.Validate(ctx, "NL")
		assert.ErrorIs(t, err, tinval.ErrInvalidFormat)
	})
}
