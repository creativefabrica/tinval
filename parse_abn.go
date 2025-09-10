package tinval

import (
	"regexp"
	"strings"
)

//nolint:gochecknoglobals // This is a constant map of country codes to their ABN number regex patterns.
var patternsABN = regexp.MustCompile(`[0-9]{11}`)

func ParseABN(s string) (TIN, error) {
	s = strings.ReplaceAll(s, " ", "")

	if !patternsABN.MatchString(s) {
		return TIN{}, ErrInvalidFormat
	}

	if !validateABN(s) {
		return TIN{}, ErrInvalidFormat
	}

	return TIN{
		CountryCode: "AU",
		Number:      s,
	}, nil
}

// validateABN will check if an ABN is valid.
// For more information on how this works you can
// refer to: https://abr.business.gov.au/Help/AbnFormat
func validateABN(abn string) bool {
	// If the first check digit is a 0 then it's not a valid ABN
	if abn[:1] == "0" {
		return false
	}

	// Subtract 1 from the first check digit of the abn
	abnByte := []byte(abn)
	abnByte[0]--
	abn = string(abnByte)

	if abn == "" || len(abn) != 11 {
		return false
	}

	abnWeights := []int{10, 1, 3, 5, 7, 9, 11, 13, 15, 17, 19}
	var weightingSum int

	for i := range abnWeights {
		num := int(abn[i]) - 48
		weightingSum += num * abnWeights[i]
	}

	// If the weightedSum is a multiple of 89 then it's a valid ABN
	return (weightingSum % 89) == 0
}
