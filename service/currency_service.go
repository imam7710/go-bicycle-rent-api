package service

type CurrencyService interface {
	Convert(amount float64, from string, to string) (float64, error)
}
