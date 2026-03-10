package utils

import "github.com/shopspring/decimal"

func SafeDecimal(d *decimal.Decimal) decimal.Decimal {
	if d == nil {
		return decimal.Zero
	}
	return *d
}
