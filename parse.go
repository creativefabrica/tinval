package tinval

func Parse(s string, countryCode string) (TIN, error) {
	taxType, ok := TaxTypeFor(countryCode)
	if !ok {
		return TIN{}, ErrInvalidCountryCode
	}

	switch taxType {
	case "au_abn":
		return ParseABN(s)
	case "eu_vat", "gb_vat":
		return ParseVAT(s)
	default:
		return TIN{}, ErrUnsupportedTaxType
	}
}

func MustParse(s string, countryCode string) TIN {
	id, err := Parse(s, countryCode)
	if err != nil {
		panic(err)
	}

	return id
}
