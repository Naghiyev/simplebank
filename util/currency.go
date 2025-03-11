package util

const (
	EUR = "EUR"
	GBP = "GBP"
	USD = "USD"
	JPY = "JPY"
)

func IsSupportedCurrency(currency string) bool {
	switch currency {
	case EUR, GBP, USD, JPY:
		return true
	}
	return false
}
