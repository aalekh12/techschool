package util

const (
	USD = "USD"
	INR = "INR"
	EUR = "EUR"
)

func IsCurrencySupport(currecny string) bool {
	switch currecny {
	case USD, INR, EUR:
		return true
	}
	return false
}
