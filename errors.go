package tinval

import "errors"

var (
	ErrInvalidFormat      = errors.New("invalid format")
	ErrNotFound           = errors.New("not found")
	ErrServiceUnavailable = errors.New("validation service unavailable")
	ErrInvalidCountryCode = errors.New("invalid country code")
)
