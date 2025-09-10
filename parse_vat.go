package tinval

import (
	"regexp"
	"strings"
)

const vatMinLength = 3

//nolint:gochecknoglobals // This is a constant map of country codes to their VAT ID number regex patterns.
var patternsVAT = map[string]*regexp.Regexp{
	"AT": regexp.MustCompile(`U[A-Z0-9]{8}`),
	"BE": regexp.MustCompile(`(0[0-9]{9}|[0-9]{10})`),
	"BG": regexp.MustCompile(`[0-9]{9,10}`),
	"CH": regexp.MustCompile(
		`(?:E(?:-| )[0-9]{3}(?:\.| )[0-9]{3}(?:\.| )[0-9]{3}( MWST)?|E[0-9]{9}(?:MWST)?)`,
	),
	"CY": regexp.MustCompile(`[0-9]{8}[A-Z]`),
	"CZ": regexp.MustCompile(`[0-9]{8,10}`),
	"DE": regexp.MustCompile(`[0-9]{9}`),
	"DK": regexp.MustCompile(`[0-9]{8}`),
	"EE": regexp.MustCompile(`[0-9]{9}`),
	"EL": regexp.MustCompile(`[0-9]{9}`),
	"ES": regexp.MustCompile(`[A-Z][0-9]{7}[A-Z]|[0-9]{8}[A-Z]|[A-Z][0-9]{8}`),
	"FI": regexp.MustCompile(`[0-9]{8}`),
	"FR": regexp.MustCompile(`([A-Z]{2}|[0-9]{2})[0-9]{9}`),
	// Supposedly the regex for GB numbers is `[0-9]{9}|[0-9]{12}|(GD|HA)[0-9]{3}`,
	// but our validator service only accepts numbers with 9 or 12 digits following the country code.
	// Seems like the official site only accepts 9 digits... https://www.gov.uk/check-uk-vat-number
	"GB": regexp.MustCompile(`([0-9]{9}|[0-9]{12})`),
	"HR": regexp.MustCompile(`[0-9]{11}`),
	"HU": regexp.MustCompile(`[0-9]{8}`),
	"IE": regexp.MustCompile(`[A-Z0-9]{7}[A-Z]|[A-Z0-9]{7}[A-W][A-I]`),
	"IT": regexp.MustCompile(`[0-9]{11}`),
	"LT": regexp.MustCompile(`([0-9]{9}|[0-9]{12})`),
	"LU": regexp.MustCompile(`[0-9]{8}`),
	"LV": regexp.MustCompile(`[0-9]{11}`),
	"MT": regexp.MustCompile(`[0-9]{8}`),
	"NL": regexp.MustCompile(`[0-9]{9}B[0-9]{2}`),
	"PL": regexp.MustCompile(`[0-9]{10}`),
	"PT": regexp.MustCompile(`[0-9]{9}`),
	"RO": regexp.MustCompile(`[0-9]{2,10}`),
	"SE": regexp.MustCompile(`[0-9]{12}`),
	"SI": regexp.MustCompile(`[0-9]{8}`),
	"SK": regexp.MustCompile(`[0-9]{10}`),
	"XI": regexp.MustCompile(`([0-9]{9}|[0-9]{12})`), // Northern Ireland, same format as GB
}

func ParseVAT(s string) (TIN, error) {
	s = strings.ReplaceAll(s, " ", "")

	if len(s) < vatMinLength {
		return TIN{}, ErrInvalidFormat
	}

	s = strings.ToUpper(s)
	num := TIN{
		CountryCode: s[:2],
		Number:      s[2:],
	}

	pattern, ok := patternsVAT[num.CountryCode]
	if !ok {
		return TIN{}, ErrInvalidCountryCode
	}

	if !pattern.MatchString(num.Number) {
		return TIN{}, ErrInvalidFormat
	}

	return num, nil
}
