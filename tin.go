package tinval

import (
	"fmt"
)

type TIN struct {
	CountryCode string
	Number      string
}

func (id TIN) String() string {
	taxType, _ := TaxTypeFor(id.CountryCode)
	switch taxType {
	case "au_abn":
		return id.Number
	default:
		return fmt.Sprintf("%s%s", id.CountryCode, id.Number)
	}
}
